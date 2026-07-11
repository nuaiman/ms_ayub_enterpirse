package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) CreateCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Type          string   `json:"type"`
		ItemID        int64    `json:"item_id"`
		Quantity      int      `json:"quantity"`
		Weight        *float64 `json:"weight"`
		ReceiverName  *string  `json:"receiver_name"`
		ReceiverPhone *string  `json:"receiver_phone"`

		// Delivery Details
		VehicleNumber *string `json:"vehicle_number"`
		DriverName    *string `json:"driver_name"`
		DriverPhone   *string `json:"driver_phone"`
		FromLocation  *string `json:"from_location"`
		ToLocation    *string `json:"to_location"`

		// Financial
		DeliveryCharge *float64 `json:"delivery_charge"`
		DeliveryCost   *float64 `json:"delivery_cost"`
		CustomerPaid   *float64 `json:"customer_paid"`

		Notes    *string `json:"notes"`
		ImageURL *string `json:"image_url"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Type == "" || req.ItemID == 0 || req.Quantity <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "type, item_id, and quantity are required")
		return
	}

	validTypes := map[string]bool{"pickup": true, "delivery": true}
	if !validTypes[req.Type] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid type, must be 'pickup' or 'delivery'")
		return
	}

	item, err := h.app.Models.Item.GetByID(r.Context(), req.ItemID)
	if err != nil || item == nil {
		utils.ErrorJson(w, http.StatusNotFound, "item not found")
		return
	}

	if item.Quantity < req.Quantity {
		utils.ErrorJson(w, http.StatusBadRequest, "insufficient item quantity")
		return
	}

	if req.Type == "delivery" {
		if req.ToLocation == nil || *req.ToLocation == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "to_location is required for delivery type")
			return
		}
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	checkout := &models.Checkout{
		UserID:         userID,
		ItemID:         req.ItemID,
		Quantity:       req.Quantity,
		Weight:         req.Weight,
		CheckoutDate:   time.Now(),
		ReceiverName:   req.ReceiverName,
		ReceiverPhone:  req.ReceiverPhone,
		Type:           req.Type,
		VehicleNumber:  req.VehicleNumber,
		DriverName:     req.DriverName,
		DriverPhone:    req.DriverPhone,
		FromLocation:   req.FromLocation,
		ToLocation:     req.ToLocation,
		DeliveryCharge: req.DeliveryCharge,
		DeliveryCost:   req.DeliveryCost,
		CustomerPaid:   req.CustomerPaid,
		Notes:          req.Notes,
		ImageURL:       req.ImageURL,
	}

	id, err := h.app.Models.Checkout.Insert(r.Context(), checkout)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create checkout")
		return
	}

	// Deduct item quantity
	if err := h.app.Models.Item.AdjustQuantity(r.Context(), req.ItemID, -req.Quantity); err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item quantity")
		return
	}

	checkout.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: fmt.Sprintf("Created %s checkout #%d for item #%d", req.Type, id, req.ItemID),
		EntityType:  "checkouts",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "checkout created successfully", checkout)
}

func (h *Handler) GetCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	checkoutID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	checkout, err := h.app.Models.Checkout.GetByID(r.Context(), checkoutID)
	if err != nil || checkout == nil {
		utils.ErrorJson(w, http.StatusNotFound, "checkout not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "checkout details fetched", checkout)
}

func (h *Handler) GetAllCheckoutsHandler(w http.ResponseWriter, r *http.Request) {
	checkouts, err := h.app.Models.Checkout.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch checkouts")
		return
	}

	if checkouts == nil {
		checkouts = []models.Checkout{}
	}

	utils.SuccessJson(w, http.StatusOK, "checkouts fetched successfully", checkouts)
}

func (h *Handler) UpdateCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Quantity      *int     `json:"quantity"`
		Weight        *float64 `json:"weight"`
		ReceiverName  *string  `json:"receiver_name"`
		ReceiverPhone *string  `json:"receiver_phone"`

		// Delivery Details
		VehicleNumber *string `json:"vehicle_number"`
		DriverName    *string `json:"driver_name"`
		DriverPhone   *string `json:"driver_phone"`
		FromLocation  *string `json:"from_location"`
		ToLocation    *string `json:"to_location"`

		// Financial
		DeliveryCharge *float64 `json:"delivery_charge"`
		DeliveryCost   *float64 `json:"delivery_cost"`
		CustomerPaid   *float64 `json:"customer_paid"`

		Notes *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	checkoutID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid checkout id")
		return
	}

	existing, err := h.app.Models.Checkout.GetByID(r.Context(), checkoutID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "checkout not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.Quantity != nil {
		if *req.Quantity < 1 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity must be at least 1")
			return
		}
		quantityDiff := existing.Quantity - *req.Quantity
		if quantityDiff != 0 {
			if err := h.app.Models.Item.AdjustQuantity(r.Context(), existing.ItemID, quantityDiff); err != nil {
				utils.ErrorJson(w, http.StatusInternalServerError, "failed to adjust item quantity")
				return
			}
		}
		updates["quantity"] = *req.Quantity
	}

	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.ReceiverName != nil {
		updates["receiver_name"] = *req.ReceiverName
	}
	if req.ReceiverPhone != nil {
		updates["receiver_phone"] = *req.ReceiverPhone
	}
	if req.VehicleNumber != nil {
		updates["vehicle_number"] = *req.VehicleNumber
	}
	if req.DriverName != nil {
		updates["driver_name"] = *req.DriverName
	}
	if req.DriverPhone != nil {
		updates["driver_phone"] = *req.DriverPhone
	}
	if req.FromLocation != nil {
		updates["from_location"] = *req.FromLocation
	}
	if req.ToLocation != nil {
		updates["to_location"] = *req.ToLocation
	}
	if req.DeliveryCharge != nil {
		updates["delivery_charge"] = *req.DeliveryCharge
	}
	if req.DeliveryCost != nil {
		updates["delivery_cost"] = *req.DeliveryCost
	}
	if req.CustomerPaid != nil {
		updates["customer_paid"] = *req.CustomerPaid
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	if len(updates) > 1 {
		err = h.app.Models.Checkout.UpdateFields(r.Context(), checkoutID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update checkout")
			return
		}
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: fmt.Sprintf("Updated checkout #%d", checkoutID),
		EntityType:  "checkouts",
		EntityID:    checkoutID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedCheckout, err := h.app.Models.Checkout.GetByID(r.Context(), checkoutID)
	if err != nil || updatedCheckout == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated checkout")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "checkout updated successfully", updatedCheckout)
}

func (h *Handler) DeleteCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	checkoutID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid checkout id")
		return
	}

	// Check if quantity should be restored (default: true)
	restoreQuantity := r.URL.Query().Get("restore_quantity")
	shouldRestore := restoreQuantity != "false"

	existing, err := h.app.Models.Checkout.GetByID(r.Context(), checkoutID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "checkout not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	// Restore item quantity if requested
	if shouldRestore {
		if err := h.app.Models.Item.AdjustQuantity(r.Context(), existing.ItemID, existing.Quantity); err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to restore item quantity")
			return
		}
	}

	err = h.app.Models.Checkout.Delete(r.Context(), checkoutID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete checkout")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	description := fmt.Sprintf("Deleted checkout #%d", checkoutID)
	if shouldRestore {
		description += " (quantity restored)"
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: description,
		EntityType:  "checkouts",
		EntityID:    checkoutID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "checkout deleted successfully", nil)
}
