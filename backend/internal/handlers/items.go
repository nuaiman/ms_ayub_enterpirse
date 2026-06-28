package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"time"
)

func (h *Handler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		// Optional fields
		CustomerID      *int64   `json:"customer_id"`
		DurationType    *string  `json:"duration_type"`
		Duration        *int     `json:"duration"`
		Price           *float64 `json:"price"`
		SecurityDeposit *float64 `json:"security_deposit"`
		EstimatedValue  *float64 `json:"estimated_value"`

		// Required field
		Name string `json:"name"`

		// Optional fields with defaults
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

	// Validate required field
	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	// Validate optional customer if provided
	if req.CustomerID != nil {
		customer, err := h.app.Models.Customer.GetByID(r.Context(), *req.CustomerID)
		if err != nil || customer == nil {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
	}

	// Validate optional contract fields if provided
	if req.DurationType != nil {
		validTypes := map[string]bool{
			"day": true, "week": true, "month": true, "year": true,
		}
		if !validTypes[*req.DurationType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
			return
		}
	}

	// Validate quantity unit if provided
	if req.QuantityUnit != "" {
		validUnits := map[string]bool{
			"pcs": true, "g": true, "kg": true, "ton": true, "ml": true, "liter": true,
		}
		if !validUnits[req.QuantityUnit] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
			return
		}
	}

	// Set defaults
	if req.QuantityUnit == "" {
		req.QuantityUnit = "pcs"
	}
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Set default status if contract fields are provided
	var status *string
	if req.DurationType != nil {
		defaultStatus := "active"
		status = &defaultStatus
	}

	item := &models.Item{
		UserID:          &userID,
		Notes:           req.Notes,
		CustomerID:      req.CustomerID,
		DurationType:    req.DurationType,
		Duration:        req.Duration,
		Price:           req.Price,
		SecurityDeposit: req.SecurityDeposit,
		EstimatedValue:  req.EstimatedValue,
		Status:          status,
		Name:            req.Name,
		QuantityUnit:    req.QuantityUnit,
		Quantity:        req.Quantity,
		Weight:          req.Weight,
		Width:           req.Width,
		Height:          req.Height,
		Length:          req.Length,
		ImageURL:        req.ImageURL,
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

	description := "Created item: " + req.Name
	if item.HasCustomer() {
		description += " for customer #" + formatInt64(*req.CustomerID)
	}
	if item.HasContract() {
		description += " (with contract)"
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: description,
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
		CustomerID      *int64   `json:"customer_id"`
		Name            *string  `json:"name"`
		DurationType    *string  `json:"duration_type"`
		Duration        *int     `json:"duration"`
		Price           *float64 `json:"price"`
		SecurityDeposit *float64 `json:"security_deposit"`
		EstimatedValue  *float64 `json:"estimated_value"`
		Status          *string  `json:"status"`
		QuantityUnit    *string  `json:"quantity_unit"`
		Quantity        *int     `json:"quantity"`
		Weight          *float64 `json:"weight"`
		Width           *float64 `json:"width"`
		Height          *float64 `json:"height"`
		Length          *float64 `json:"length"`
		Notes           *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate name if provided
	if req.Name != nil && *req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name cannot be empty")
		return
	}

	// Validate optional customer if provided
	if req.CustomerID != nil {
		customer, err := h.app.Models.Customer.GetByID(r.Context(), *req.CustomerID)
		if err != nil || customer == nil {
			utils.ErrorJson(w, http.StatusNotFound, "customer not found")
			return
		}
	}

	// Validate optional contract fields if provided
	if req.DurationType != nil {
		validTypes := map[string]bool{
			"day": true, "week": true, "month": true, "year": true,
		}
		if !validTypes[*req.DurationType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
			return
		}
	}

	if req.Status != nil {
		validStatuses := map[string]bool{
			"active": true, "completed": true, "cancelled": true,
		}
		if !validStatuses[*req.Status] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid status")
			return
		}
	}

	// Validate quantity unit if provided
	if req.QuantityUnit != nil {
		validUnits := map[string]bool{
			"pcs": true, "g": true, "kg": true, "ton": true, "ml": true, "liter": true,
		}
		if !validUnits[*req.QuantityUnit] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
			return
		}
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

	// Build update map with only fields that were provided
	updates := make(map[string]interface{})

	// Add fields that were provided
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CustomerID != nil {
		updates["customer_id"] = *req.CustomerID
	} else if req.CustomerID == nil && existing.CustomerID != nil {
		// Explicitly set to NULL if customer_id was sent as null
		updates["customer_id"] = nil
	}
	if req.DurationType != nil {
		updates["duration_type"] = *req.DurationType
	} else if req.DurationType == nil && existing.DurationType != nil {
		// Explicitly set to NULL if duration_type was sent as null
		updates["duration_type"] = nil
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	} else if req.Duration == nil && existing.Duration != nil {
		updates["duration"] = nil
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	} else if req.Price == nil && existing.Price != nil {
		updates["price"] = nil
	}
	if req.SecurityDeposit != nil {
		updates["security_deposit"] = *req.SecurityDeposit
	} else if req.SecurityDeposit == nil && existing.SecurityDeposit != nil {
		updates["security_deposit"] = nil
	}
	if req.EstimatedValue != nil {
		updates["estimated_value"] = *req.EstimatedValue
	} else if req.EstimatedValue == nil && existing.EstimatedValue != nil {
		updates["estimated_value"] = nil
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	} else if req.Status == nil && existing.Status != nil {
		updates["status"] = nil
	}
	if req.QuantityUnit != nil {
		updates["quantity_unit"] = *req.QuantityUnit
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	} else if req.Weight == nil && existing.Weight != nil {
		updates["weight"] = nil
	}
	if req.Width != nil {
		updates["width"] = *req.Width
	} else if req.Width == nil && existing.Width != nil {
		updates["width"] = nil
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	} else if req.Height == nil && existing.Height != nil {
		updates["height"] = nil
	}
	if req.Length != nil {
		updates["length"] = *req.Length
	} else if req.Length == nil && existing.Length != nil {
		updates["length"] = nil
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	} else if req.Notes == nil && existing.Notes != nil {
		updates["notes"] = nil
	}

	// Always update updated_at
	updates["updated_at"] = time.Now()

	// Only update if there are changes (more than just updated_at)
	if len(updates) > 1 {
		err = h.app.Models.Item.UpdateFields(r.Context(), itemID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item")
			return
		}
	}

	// Handle status change if status was provided and item has contract
	if req.Status != nil && existing.Status != nil && *req.Status != *existing.Status {
		var endedAt *time.Time
		if (*req.Status == "completed" || *req.Status == "cancelled") && *existing.Status == "active" {
			now := time.Now()
			endedAt = &now
		}
		err = h.app.Models.Item.UpdateStatus(r.Context(), itemID, *req.Status, endedAt)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item status")
			return
		}
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated item: " + existing.Name,
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

func (h *Handler) ChangeItemStatusHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if item has a contract
	if existing.Status == nil {
		utils.ErrorJson(w, http.StatusBadRequest, "item does not have a contract")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	var endedAt *time.Time
	if (req.Status == "completed" || req.Status == "cancelled") && *existing.Status == "active" {
		now := time.Now()
		endedAt = &now
	}

	err = h.app.Models.Item.UpdateStatus(r.Context(), itemID, req.Status, endedAt)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item status")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Changed item #" + formatInt64(itemID) + " status to " + req.Status,
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

	utils.SuccessJson(w, http.StatusOK, "item status updated successfully", updatedItem)
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
