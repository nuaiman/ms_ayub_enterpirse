package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
)

func (h *Handler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContractID   *int64   `json:"contract_id"`
		Name         string   `json:"name"`
		QuantityUnit string   `json:"quantity_unit"`
		Quantity     int      `json:"quantity"`
		Weight       *float64 `json:"weight"`
		Width        *float64 `json:"width"`
		Height       *float64 `json:"height"`
		Length       *float64 `json:"length"`
		Notes        *string  `json:"notes"`
		ImageURL     *string  `json:"image_url"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validUnits := map[string]bool{
		"pcs": true, "g": true, "kg": true, "ton": true, "ml": true, "liter": true,
	}
	if req.QuantityUnit != "" && !validUnits[req.QuantityUnit] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
		return
	}
	if req.QuantityUnit == "" {
		req.QuantityUnit = "pcs"
	}
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	if req.ContractID != nil {
		contract, err := h.app.Models.Contract.GetByID(r.Context(), *req.ContractID)
		if err != nil || contract == nil {
			utils.ErrorJson(w, http.StatusNotFound, "contract not found")
			return
		}
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	item := &models.Item{
		UserID:       &userID,
		ContractID:   req.ContractID,
		Name:         req.Name,
		QuantityUnit: req.QuantityUnit,
		Quantity:     req.Quantity,
		Weight:       req.Weight,
		Width:        req.Width,
		Height:       req.Height,
		Length:       req.Length,
		Notes:        req.Notes,
	}

	id, err := h.app.Models.Item.Insert(r.Context(), item)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	item.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created item: " + req.Name,
		EntityType:  "items",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "item created successfully", item)
}

func (h *Handler) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	item, err := h.app.Models.Item.GetByID(r.Context(), itemID)
	if err != nil || item == nil {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "item details fetched", item)
}

func (h *Handler) GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := h.app.Models.Item.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch items")
		return
	}

	if items == nil {
		items = []models.Item{}
	}

	utils.SuccessJson(w, http.StatusOK, "items fetched successfully", items)
}

func (h *Handler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ContractID   *int64   `json:"contract_id"`
		Name         string   `json:"name"`
		QuantityUnit string   `json:"quantity_unit"`
		Quantity     int      `json:"quantity"`
		Weight       *float64 `json:"weight"`
		Width        *float64 `json:"width"`
		Height       *float64 `json:"height"`
		Length       *float64 `json:"length"`
		Notes        *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validUnits := map[string]bool{
		"pcs": true, "g": true, "kg": true, "ton": true, "ml": true, "liter": true,
	}
	if req.QuantityUnit != "" && !validUnits[req.QuantityUnit] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
		return
	}

	itemID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid item id")
		return
	}

	existing, err := h.app.Models.Item.GetByID(r.Context(), itemID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if req.ContractID != nil {
		contract, err := h.app.Models.Contract.GetByID(r.Context(), *req.ContractID)
		if err != nil || contract == nil {
			utils.ErrorJson(w, http.StatusNotFound, "contract not found")
			return
		}
	}

	item := &models.Item{
		ContractID:   req.ContractID,
		Name:         req.Name,
		QuantityUnit: req.QuantityUnit,
		Quantity:     req.Quantity,
		Weight:       req.Weight,
		Width:        req.Width,
		Height:       req.Height,
		Length:       req.Length,
		Notes:        req.Notes,
		ImageURL:     existing.ImageURL,
	}

	err = h.app.Models.Item.Update(r.Context(), itemID, item)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated item: " + req.Name,
		EntityType:  "items",
		EntityID:    itemID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedItem, err := h.app.Models.Item.GetByID(r.Context(), itemID)
	if err != nil || updatedItem == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated item")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "item updated successfully", updatedItem)
}

func (h *Handler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid item id")
		return
	}

	existing, err := h.app.Models.Item.GetByID(r.Context(), itemID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.app.Models.Item.Delete(r.Context(), itemID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted item: " + existing.Name,
		EntityType:  "items",
		EntityID:    itemID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "item deleted successfully", nil)
}
