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
		ProductName   string  `json:"product_name"`
		StorageName   *string `json:"storage_name"`
		AccountName   *string `json:"account_name"`
		LotNumber     *string `json:"lot_number"`
		CustomerPhone *string `json:"customer_phone"`
		CustomerEmail *string `json:"customer_email"`
		Category      *string `json:"category"`
		Subcategory   *string `json:"subcategory"`
		QuantityUnit  string  `json:"quantity_unit"`
		Quantity      int     `json:"quantity"`
		Amount        float64 `json:"amount"`
		Deposit       float64 `json:"deposit"`
		CustomerPaid  float64 `json:"customer_paid"`
		Notes         *string `json:"notes"`
		ImageURL      *string `json:"image_url"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		log.Printf("❌ Failed to parse request body: %v", err)
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Only product_name is required
	if req.ProductName == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "product_name is required")
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
		utils.ErrorJson(w, http.StatusBadRequest, "invalid quantity_unit")
		return
	}

	if req.Quantity == 0 {
		req.Quantity = 1
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	item := &models.Item{
		UserID:        userID,
		Notes:         req.Notes,
		ProductName:   req.ProductName,
		StorageName:   req.StorageName,
		AccountName:   req.AccountName,
		LotNumber:     req.LotNumber,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		Category:      req.Category,
		Subcategory:   req.Subcategory,
		QuantityUnit:  req.QuantityUnit,
		Quantity:      req.Quantity,
		Amount:        req.Amount,
		Deposit:       req.Deposit,
		CustomerPaid:  req.CustomerPaid,
		ImageURL:      req.ImageURL,
	}

	id, err := h.app.Models.Item.Insert(r.Context(), item)
	if err != nil {
		log.Printf("❌ Failed to insert item: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	item.ID = id

	ip := r.RemoteAddr
	ua := r.UserAgent()
	description := fmt.Sprintf("Created item: %s", req.ProductName)
	if req.CustomerPhone != nil && *req.CustomerPhone != "" {
		description += fmt.Sprintf(" (%s)", *req.CustomerPhone)
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

func (h *Handler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ProductName   *string  `json:"product_name"`
		StorageName   *string  `json:"storage_name"`
		AccountName   *string  `json:"account_name"`
		LotNumber     *string  `json:"lot_number"`
		CustomerPhone *string  `json:"customer_phone"`
		CustomerEmail *string  `json:"customer_email"`
		Category      *string  `json:"category"`
		Subcategory   *string  `json:"subcategory"`
		QuantityUnit  *string  `json:"quantity_unit"`
		Quantity      *int     `json:"quantity"`
		Amount        *float64 `json:"amount"`
		Deposit       *float64 `json:"deposit"`
		CustomerPaid  *float64 `json:"customer_paid"`
		Notes         *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ProductName != nil && *req.ProductName == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "product_name cannot be empty")
		return
	}

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

	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.ProductName != nil {
		updates["product_name"] = *req.ProductName
	}
	if req.StorageName != nil {
		updates["storage_name"] = *req.StorageName
	}
	if req.AccountName != nil {
		updates["account_name"] = *req.AccountName
	}
	if req.LotNumber != nil {
		updates["lot_number"] = *req.LotNumber
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
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Deposit != nil {
		updates["deposit"] = *req.Deposit
	}
	if req.CustomerPaid != nil {
		updates["customer_paid"] = *req.CustomerPaid
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	if len(updates) > 1 {
		err = h.app.Models.Item.UpdateFields(r.Context(), itemID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item")
			return
		}
	}

	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: fmt.Sprintf("Updated item: %s", existing.ProductName),
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

	ip := r.RemoteAddr
	ua := r.UserAgent()
	description := fmt.Sprintf("Deleted item: %s", existing.ProductName)
	if existing.CustomerPhone != nil && *existing.CustomerPhone != "" {
		description += fmt.Sprintf(" (%s)", *existing.CustomerPhone)
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: description,
		EntityType:  "items",
		EntityID:    itemID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "item deleted successfully", nil)
}
