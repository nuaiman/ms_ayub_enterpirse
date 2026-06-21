package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"strings"
)

func (h *Handler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string  `json:"name"`
		Phone    string  `json:"phone"`
		Email    *string `json:"email"`
		Company  *string `json:"company"`
		Address  *string `json:"address"`
		IDType   *string `json:"id_type"`
		IDNumber *string `json:"id_number"`
		Notes    *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Phone == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	if req.Email != nil && *req.Email != "" {
		if !strings.Contains(*req.Email, "@") {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	customer := &models.Customer{
		UserID:   &userID,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Address:  req.Address,
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
		Notes:    req.Notes,
	}

	id, err := h.app.Models.Customer.Insert(r.Context(), customer)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	customer.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created customer: " + req.Name,
		EntityType:  "customers",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "customer created successfully", customer)
}

func (h *Handler) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	customerID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	customer, err := h.app.Models.Customer.GetByID(r.Context(), customerID)
	if err != nil || customer == nil {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "customer details fetched", customer)
}

func (h *Handler) GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {
	customers, err := h.app.Models.Customer.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch customers")
		return
	}

	if customers == nil {
		customers = []models.Customer{}
	}

	utils.SuccessJson(w, http.StatusOK, "customers fetched successfully", customers)
}

func (h *Handler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string  `json:"name"`
		Phone    string  `json:"phone"`
		Email    *string `json:"email"`
		Company  *string `json:"company"`
		Address  *string `json:"address"`
		IDType   *string `json:"id_type"`
		IDNumber *string `json:"id_number"`
		Notes    *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Phone == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	if req.Email != nil && *req.Email != "" {
		if !strings.Contains(*req.Email, "@") {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}
	}

	customerID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	existing, err := h.app.Models.Customer.GetByID(r.Context(), customerID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	customer := &models.Customer{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Address:  req.Address,
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
		Notes:    req.Notes,
	}

	err = h.app.Models.Customer.Update(r.Context(), customerID, customer)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated customer: " + req.Name,
		EntityType:  "customers",
		EntityID:    customerID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedCustomer, err := h.app.Models.Customer.GetByID(r.Context(), customerID)
	if err != nil || updatedCustomer == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated customer")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "customer updated successfully", updatedCustomer)
}

func (h *Handler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	customerID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	existing, err := h.app.Models.Customer.GetByID(r.Context(), customerID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "customer not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.app.Models.Customer.Delete(r.Context(), customerID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted customer: " + existing.Name,
		EntityType:  "customers",
		EntityID:    customerID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "customer deleted successfully", nil)
}
