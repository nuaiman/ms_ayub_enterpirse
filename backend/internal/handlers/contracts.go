package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"time"
)

func (h *Handler) CreateContractHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CustomerID      int64    `json:"customer_id"`
		DurationType    string   `json:"duration_type"`
		Duration        *int     `json:"duration"`
		Price           *float64 `json:"price"`
		SecurityDeposit *float64 `json:"security_deposit"`
		EstimatedValue  *float64 `json:"estimated_value"`
		Notes           *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 || req.DurationType == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validTypes := map[string]bool{
		"day": true, "week": true, "month": true, "year": true,
	}
	if !validTypes[req.DurationType] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
		return
	}

	customer, err := h.app.Models.Customer.GetByID(r.Context(), req.CustomerID)
	if err != nil || customer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	contract := &models.Contract{
		UserID:          &userID,
		CustomerID:      req.CustomerID,
		DurationType:    req.DurationType,
		Duration:        req.Duration,
		Price:           req.Price,
		SecurityDeposit: req.SecurityDeposit,
		EstimatedValue:  req.EstimatedValue,
		Notes:           req.Notes,
		Status:          "active",
	}

	id, err := h.app.Models.Contract.Insert(r.Context(), contract)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create contract")
		return
	}

	contract.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created contract #" + formatInt64(id),
		EntityType:  "contracts",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "contract created successfully", contract)
}

func (h *Handler) GetContractHandler(w http.ResponseWriter, r *http.Request) {
	contractID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	contract, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || contract == nil {
		utils.ErrorJson(w, http.StatusNotFound, "contract not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "contract details fetched", contract)
}

func (h *Handler) GetAllContractsHandler(w http.ResponseWriter, r *http.Request) {
	contracts, err := h.app.Models.Contract.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch contracts")
		return
	}

	if contracts == nil {
		contracts = []models.Contract{}
	}

	utils.SuccessJson(w, http.StatusOK, "contracts fetched successfully", contracts)
}

func (h *Handler) UpdateContractHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CustomerID      int64    `json:"customer_id"`
		DurationType    string   `json:"duration_type"`
		Duration        *int     `json:"duration"`
		Price           *float64 `json:"price"`
		SecurityDeposit *float64 `json:"security_deposit"`
		EstimatedValue  *float64 `json:"estimated_value"`
		Status          string   `json:"status"`
		Notes           *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CustomerID == 0 || req.DurationType == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validTypes := map[string]bool{
		"day": true, "week": true, "month": true, "year": true,
	}
	if !validTypes[req.DurationType] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
		return
	}

	if req.Status != "" {
		validStatuses := map[string]bool{
			"active": true, "completed": true, "cancelled": true,
		}
		if !validStatuses[req.Status] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid status")
			return
		}
	}

	contractID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid contract id")
		return
	}

	existing, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "contract not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	customer, err := h.app.Models.Customer.GetByID(r.Context(), req.CustomerID)
	if err != nil || customer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	var endedAt *time.Time
	if (req.Status == "completed" || req.Status == "cancelled") && existing.Status == "active" {
		now := time.Now()
		endedAt = &now
	}

	contract := &models.Contract{
		CustomerID:      req.CustomerID,
		DurationType:    req.DurationType,
		Duration:        req.Duration,
		Price:           req.Price,
		SecurityDeposit: req.SecurityDeposit,
		EstimatedValue:  req.EstimatedValue,
		Status:          req.Status,
		Notes:           req.Notes,
	}

	err = h.app.Models.Contract.Update(r.Context(), contractID, contract)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update contract")
		return
	}

	if req.Status != existing.Status {
		err = h.app.Models.Contract.UpdateStatus(r.Context(), contractID, req.Status, endedAt)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update contract status")
			return
		}
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated contract #" + formatInt64(contractID),
		EntityType:  "contracts",
		EntityID:    contractID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedContract, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || updatedContract == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated contract")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "contract updated successfully", updatedContract)
}

func (h *Handler) ChangeContractStatusHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validStatuses := map[string]bool{
		"active": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid status")
		return
	}

	contractID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid contract id")
		return
	}

	existing, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "contract not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	var endedAt *time.Time
	if (req.Status == "completed" || req.Status == "cancelled") && existing.Status == "active" {
		now := time.Now()
		endedAt = &now
	}

	err = h.app.Models.Contract.UpdateStatus(r.Context(), contractID, req.Status, endedAt)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update contract status")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Changed contract #" + formatInt64(contractID) + " status to " + req.Status,
		EntityType:  "contracts",
		EntityID:    contractID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedContract, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || updatedContract == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated contract")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "contract status updated successfully", updatedContract)
}

func (h *Handler) DeleteContractHandler(w http.ResponseWriter, r *http.Request) {
	contractID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid contract id")
		return
	}

	existing, err := h.app.Models.Contract.GetByID(r.Context(), contractID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "contract not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.app.Models.Contract.Delete(r.Context(), contractID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete contract")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted contract #" + formatInt64(contractID),
		EntityType:  "contracts",
		EntityID:    contractID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "contract deleted successfully", nil)
}
