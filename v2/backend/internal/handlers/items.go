package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (h *Handler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		// Customer Fields
		Name          string  `json:"name"`
		CustomerPhone string  `json:"customer_phone"`
		CustomerEmail *string `json:"customer_email"`

		// Category
		Category    *string `json:"category"`
		Subcategory *string `json:"subcategory"`

		// Item Details
		QuantityUnit string `json:"quantity_unit"`
		Quantity     int    `json:"quantity"`

		// Dimensions
		Weight        *float64 `json:"weight"`
		WeightUnit    *string  `json:"weight_unit"`
		Width         *float64 `json:"width"`
		Height        *float64 `json:"height"`
		Length        *float64 `json:"length"`
		DimensionUnit *string  `json:"dimension_unit"`

		// Storage Contract
		DurationType *string `json:"duration_type"`
		Duration     *int    `json:"duration"`
		StartDate    *string `json:"start_date"`
		Amount       float64 `json:"amount"`
		Deposit      float64 `json:"deposit"`
		CustomerPaid float64 `json:"customer_paid"`

		// Notes and Image
		Notes    *string `json:"notes"`
		ImageURL *string `json:"image_url"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("❌ Failed to parse request body: %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("📝 CreateItem request: name=%s, phone=%s, amount=%.2f, customer_paid=%.2f",
		req.Name, req.CustomerPhone, req.Amount, req.CustomerPaid)

	// Validate required fields
	if req.Name == "" {
		log.Printf("❌ Validation failed: name is required")
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	if req.CustomerPhone == "" {
		log.Printf("❌ Validation failed: customer_phone is required")
		utils.ErrorJson(w, http.StatusBadRequest, "customer_phone is required")
		return
	}

	// Validate quantity unit
	validUnits := map[string]bool{
		"bag": true, "bottle": true, "box": true, "can": true, "carton": true,
		"cup": true, "dozen": true, "gallon": true, "pack": true, "pair": true,
		"pcs": true, "roll": true, "set": true, "sheet": true, "unit": true,
	}

	if req.QuantityUnit == "" {
		req.QuantityUnit = "pcs"
	} else if !validUnits[req.QuantityUnit] {
		log.Printf("❌ Validation failed: invalid quantity_unit '%s'", req.QuantityUnit)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
		return
	}

	// Validate weight unit if provided
	if req.WeightUnit != nil {
		validWeightUnits := map[string]bool{
			"mg": true, "g": true, "oz": true, "lb": true, "kg": true, "ton": true,
		}
		if !validWeightUnits[*req.WeightUnit] {
			log.Printf("❌ Validation failed: invalid weight_unit '%s'", *req.WeightUnit)
			utils.ErrorJson(w, http.StatusBadRequest, "invalid weight_unit")
			return
		}
	}

	// Validate dimension unit if provided
	if req.DimensionUnit != nil {
		validDimUnits := map[string]bool{
			"mm": true, "cm": true, "in": true, "ft": true, "m": true, "yd": true, "km": true,
		}
		if !validDimUnits[*req.DimensionUnit] {
			log.Printf("❌ Validation failed: invalid dimension_unit '%s'", *req.DimensionUnit)
			utils.ErrorJson(w, http.StatusBadRequest, "invalid dimension_unit")
			return
		}
	}

	// Validate storage contract fields if provided
	var status *string
	if req.DurationType != nil && *req.DurationType != "" {
		validTypes := map[string]bool{
			"day": true, "week": true, "month": true, "year": true,
		}
		if !validTypes[*req.DurationType] {
			log.Printf("❌ Validation failed: invalid duration_type '%s'", *req.DurationType)
			utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
			return
		}
		if req.Duration == nil || *req.Duration <= 0 {
			log.Printf("❌ Validation failed: duration is required when duration_type is provided")
			utils.ErrorJson(w, http.StatusBadRequest, "duration is required when duration_type is provided")
			return
		}
		if req.Amount <= 0 {
			log.Printf("❌ Validation failed: amount must be greater than 0 for storage contracts")
			utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0 for storage contracts")
			return
		}
		defaultStatus := "active"
		status = &defaultStatus
	}

	// Set defaults
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		log.Printf("❌ Failed to get userID from context")
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	log.Printf("📝 Creating item with userID: %d", userID)

	item := &models.Item{
		UserID:        userID,
		Notes:         req.Notes,
		Name:          req.Name,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		Category:      req.Category,
		Subcategory:   req.Subcategory,
		QuantityUnit:  req.QuantityUnit,
		Quantity:      req.Quantity,
		Weight:        req.Weight,
		WeightUnit:    req.WeightUnit,
		Width:         req.Width,
		Height:        req.Height,
		Length:        req.Length,
		DimensionUnit: req.DimensionUnit,
		DurationType:  req.DurationType,
		Duration:      req.Duration,
		StartDate:     req.StartDate,
		Amount:        req.Amount,
		Deposit:       req.Deposit,
		CustomerPaid:  req.CustomerPaid,
		Status:        status,
		ImageURL:      req.ImageURL,
	}

	id, err := h.app.Models.Item.Insert(r.Context(), item)
	if err != nil {
		log.Printf("❌ Failed to insert item: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	log.Printf("✅ Item created successfully with ID: %d", id)

	item.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()

	description := fmt.Sprintf("Created item for %s (%s)", req.Name, req.CustomerPhone)
	if item.HasContract() {
		description += fmt.Sprintf(" - Storage contract: %d %s(s) for $%.2f", *req.Duration, *req.DurationType, req.Amount)
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
		// Customer Fields
		Name          *string `json:"name"`
		CustomerPhone *string `json:"customer_phone"`
		CustomerEmail *string `json:"customer_email"`

		// Category
		Category    *string `json:"category"`
		Subcategory *string `json:"subcategory"`

		// Item Details
		QuantityUnit *string `json:"quantity_unit"`
		Quantity     *int    `json:"quantity"`

		// Dimensions
		Weight        *float64 `json:"weight"`
		WeightUnit    *string  `json:"weight_unit"`
		Width         *float64 `json:"width"`
		Height        *float64 `json:"height"`
		Length        *float64 `json:"length"`
		DimensionUnit *string  `json:"dimension_unit"`

		// Storage Contract
		DurationType *string  `json:"duration_type"`
		Duration     *int     `json:"duration"`
		StartDate    *string  `json:"start_date"`
		Amount       *float64 `json:"amount"`
		Deposit      *float64 `json:"deposit"`
		CustomerPaid *float64 `json:"customer_paid"`
		Status       *string  `json:"status"`

		// Notes
		Notes *string `json:"notes"`
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

	// Validate customer_phone if provided
	if req.CustomerPhone != nil && *req.CustomerPhone == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "customer_phone cannot be empty")
		return
	}

	// Validate quantity unit if provided
	if req.QuantityUnit != nil {
		validUnits := map[string]bool{
			"bag": true, "bottle": true, "box": true, "can": true, "carton": true,
			"cup": true, "dozen": true, "gallon": true, "pack": true, "pair": true,
			"pcs": true, "roll": true, "set": true, "sheet": true, "unit": true,
		}
		if !validUnits[*req.QuantityUnit] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
			return
		}
	}

	// Validate weight unit if provided
	if req.WeightUnit != nil {
		validWeightUnits := map[string]bool{
			"mg": true, "g": true, "oz": true, "lb": true, "kg": true, "ton": true,
		}
		if !validWeightUnits[*req.WeightUnit] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid weight_unit")
			return
		}
	}

	// Validate dimension unit if provided
	if req.DimensionUnit != nil {
		validDimUnits := map[string]bool{
			"mm": true, "cm": true, "in": true, "ft": true, "m": true, "yd": true, "km": true,
		}
		if !validDimUnits[*req.DimensionUnit] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid dimension_unit")
			return
		}
	}

	// Validate duration_type if provided
	if req.DurationType != nil {
		validTypes := map[string]bool{
			"day": true, "week": true, "month": true, "year": true,
		}
		if !validTypes[*req.DurationType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid duration_type")
			return
		}
	}

	// Validate status if provided
	if req.Status != nil {
		validStatuses := map[string]bool{
			"active": true, "complete": true,
		}
		if !validStatuses[*req.Status] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid status")
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
	updates["updated_at"] = time.Now()

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CustomerPhone != nil {
		updates["customer_phone"] = *req.CustomerPhone
	}
	if req.CustomerEmail != nil {
		updates["customer_email"] = *req.CustomerEmail
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Subcategory != nil {
		updates["subcategory"] = *req.Subcategory
	}
	if req.QuantityUnit != nil {
		updates["quantity_unit"] = *req.QuantityUnit
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.WeightUnit != nil {
		updates["weight_unit"] = *req.WeightUnit
	}
	if req.Width != nil {
		updates["width"] = *req.Width
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}
	if req.Length != nil {
		updates["length"] = *req.Length
	}
	if req.DimensionUnit != nil {
		updates["dimension_unit"] = *req.DimensionUnit
	}
	if req.DurationType != nil {
		updates["duration_type"] = *req.DurationType
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Deposit != nil {
		updates["deposit"] = *req.Deposit
	}
	if req.CustomerPaid != nil {
		updates["customer_paid"] = *req.CustomerPaid
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	// Only update if there are changes (more than just updated_at)
	if len(updates) > 1 {
		err = h.app.Models.Item.UpdateFields(r.Context(), itemID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item")
			return
		}
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: fmt.Sprintf("Updated item for %s", existing.Name),
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
		"active": true, "complete": true,
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

	// Check if item has a storage contract
	if existing.Status == nil {
		utils.ErrorJson(w, http.StatusBadRequest, "item does not have a storage contract")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.app.Models.Item.UpdateStatus(r.Context(), itemID, req.Status)
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
		Description: fmt.Sprintf("Changed item #%d status to %s", itemID, req.Status),
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
		Description: fmt.Sprintf("Deleted item for %s (%s)", existing.Name, existing.CustomerPhone),
		EntityType:  "items",
		EntityID:    itemID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "item deleted successfully", nil)
}
