package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"time"
)

func (h *Handler) CreateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ShipmentType   string   `json:"shipment_type"`
		ItemID         int64    `json:"item_id"`
		Quantity       int      `json:"quantity"`
		VehicleNumber  *string  `json:"vehicle_number"`
		DriverName     *string  `json:"driver_name"`
		DriverPhone    *string  `json:"driver_phone"`
		FromLocation   *string  `json:"from_location"`
		CompanyName    *string  `json:"company_name"`
		ToLocation     *string  `json:"to_location"`
		ReceiverName   *string  `json:"receiver_name"`
		ReceiverPhone  *string  `json:"receiver_phone"`
		CustomerCharge *float64 `json:"customer_charge"`
		CompanyCost    *float64 `json:"company_cost"`
		CompanyPaid    *float64 `json:"company_paid"`
		CustomerPaid   *float64 `json:"customer_paid"`
		Notes          *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ShipmentType == "" || req.ItemID == 0 || req.Quantity <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validTypes := map[string]bool{"pickup": true, "delivery": true}
	if !validTypes[req.ShipmentType] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid shipment_type")
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

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	shipment := &models.Shipment{
		UserID:         &userID,
		ShipmentType:   req.ShipmentType,
		ItemID:         req.ItemID,
		Quantity:       req.Quantity,
		VehicleNumber:  req.VehicleNumber,
		DriverName:     req.DriverName,
		DriverPhone:    req.DriverPhone,
		FromLocation:   req.FromLocation,
		CompanyName:    req.CompanyName,
		ToLocation:     req.ToLocation,
		ReceiverName:   req.ReceiverName,
		ReceiverPhone:  req.ReceiverPhone,
		CustomerCharge: req.CustomerCharge,
		CompanyCost:    req.CompanyCost,
		CompanyPaid:    req.CompanyPaid,
		CustomerPaid:   req.CustomerPaid,
		Status:         "pending",
		Notes:          req.Notes,
	}

	id, err := h.app.Models.Shipment.Insert(r.Context(), shipment)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create shipment")
		return
	}

	if err := h.app.Models.Item.AdjustQuantity(r.Context(), req.ItemID, -req.Quantity); err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update item quantity")
		return
	}

	shipment.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created shipment #" + formatInt64(id),
		EntityType:  "shipments",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "shipment created successfully", shipment)
}

func (h *Handler) GetShipmentHandler(w http.ResponseWriter, r *http.Request) {
	shipmentID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	shipment, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || shipment == nil {
		utils.ErrorJson(w, http.StatusNotFound, "shipment not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "shipment details fetched", shipment)
}

func (h *Handler) GetAllShipmentsHandler(w http.ResponseWriter, r *http.Request) {
	shipments, err := h.app.Models.Shipment.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch shipments")
		return
	}

	if shipments == nil {
		shipments = []models.Shipment{}
	}

	utils.SuccessJson(w, http.StatusOK, "shipments fetched successfully", shipments)
}

// internal/handlers/shipments.go

func (h *Handler) UpdateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ShipmentType   *string  `json:"shipment_type"`
		Quantity       *int     `json:"quantity"`
		VehicleNumber  *string  `json:"vehicle_number"`
		DriverName     *string  `json:"driver_name"`
		DriverPhone    *string  `json:"driver_phone"`
		FromLocation   *string  `json:"from_location"`
		CompanyName    *string  `json:"company_name"`
		ToLocation     *string  `json:"to_location"`
		ReceiverName   *string  `json:"receiver_name"`
		ReceiverPhone  *string  `json:"receiver_phone"`
		CustomerCharge *float64 `json:"customer_charge"`
		CompanyCost    *float64 `json:"company_cost"`
		CompanyPaid    *float64 `json:"company_paid"`
		CustomerPaid   *float64 `json:"customer_paid"`
		Notes          *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	shipmentID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid shipment id")
		return
	}

	// Get existing shipment
	existing, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "shipment not found")
		return
	}

	// Build update map with only fields that were provided
	updates := make(map[string]interface{})

	// Helper to add field if provided
	if req.ShipmentType != nil {
		validTypes := map[string]bool{"pickup": true, "delivery": true}
		if !validTypes[*req.ShipmentType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid shipment_type")
			return
		}
		updates["shipment_type"] = *req.ShipmentType
	}

	if req.Quantity != nil {
		if *req.Quantity < 1 {
			utils.ErrorJson(w, http.StatusBadRequest, "quantity must be at least 1")
			return
		}
		updates["quantity"] = *req.Quantity
	}

	if req.VehicleNumber != nil {
		updates["vehicle_number"] = *req.VehicleNumber
	} else if req.VehicleNumber == nil && existing.VehicleNumber != nil {
		updates["vehicle_number"] = nil
	}

	if req.DriverName != nil {
		updates["driver_name"] = *req.DriverName
	} else if req.DriverName == nil && existing.DriverName != nil {
		updates["driver_name"] = nil
	}

	if req.DriverPhone != nil {
		updates["driver_phone"] = *req.DriverPhone
	} else if req.DriverPhone == nil && existing.DriverPhone != nil {
		updates["driver_phone"] = nil
	}

	if req.FromLocation != nil {
		updates["from_location"] = *req.FromLocation
	} else if req.FromLocation == nil && existing.FromLocation != nil {
		updates["from_location"] = nil
	}

	if req.CompanyName != nil {
		updates["company_name"] = *req.CompanyName
	} else if req.CompanyName == nil && existing.CompanyName != nil {
		updates["company_name"] = nil
	}

	if req.ToLocation != nil {
		updates["to_location"] = *req.ToLocation
	} else if req.ToLocation == nil && existing.ToLocation != nil {
		updates["to_location"] = nil
	}

	if req.ReceiverName != nil {
		updates["receiver_name"] = *req.ReceiverName
	} else if req.ReceiverName == nil && existing.ReceiverName != nil {
		updates["receiver_name"] = nil
	}

	if req.ReceiverPhone != nil {
		updates["receiver_phone"] = *req.ReceiverPhone
	} else if req.ReceiverPhone == nil && existing.ReceiverPhone != nil {
		updates["receiver_phone"] = nil
	}

	if req.CustomerCharge != nil {
		updates["customer_charge"] = *req.CustomerCharge
	} else if req.CustomerCharge == nil && existing.CustomerCharge != nil {
		updates["customer_charge"] = nil
	}

	if req.CompanyCost != nil {
		updates["company_cost"] = *req.CompanyCost
	} else if req.CompanyCost == nil && existing.CompanyCost != nil {
		updates["company_cost"] = nil
	}

	if req.CompanyPaid != nil {
		updates["company_paid"] = *req.CompanyPaid
	} else if req.CompanyPaid == nil && existing.CompanyPaid != nil {
		updates["company_paid"] = nil
	}

	if req.CustomerPaid != nil {
		updates["customer_paid"] = *req.CustomerPaid
	} else if req.CustomerPaid == nil && existing.CustomerPaid != nil {
		updates["customer_paid"] = nil
	}

	if req.Notes != nil {
		updates["notes"] = *req.Notes
	} else if req.Notes == nil && existing.Notes != nil {
		updates["notes"] = nil
	}

	// Always update updated_at
	updates["updated_at"] = time.Now()

	// Only update if there are changes
	if len(updates) > 1 {
		err = h.app.Models.Shipment.UpdateFields(r.Context(), shipmentID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update shipment: "+err.Error())
			return
		}
	}

	// Audit log
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if ok {
		ip := r.RemoteAddr
		ua := r.UserAgent()
		_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
			UserID:      userID,
			Action:      "update",
			Description: "Updated shipment #" + formatInt64(shipmentID),
			EntityType:  "shipments",
			EntityID:    shipmentID,
			IPAddress:   &ip,
			UserAgent:   &ua,
		})
	}

	updatedShipment, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || updatedShipment == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated shipment")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "shipment updated successfully", updatedShipment)
}

func (h *Handler) ChangeShipmentStatusHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Status string `json:"status"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "in_transit": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid status")
		return
	}

	shipmentID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid shipment id")
		return
	}

	existing, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "shipment not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	var shippedAt *time.Time
	var receivedAt *time.Time

	if req.Status == "in_transit" && existing.Status == "pending" {
		now := time.Now()
		shippedAt = &now
	}
	if req.Status == "completed" {
		now := time.Now()
		receivedAt = &now
		if existing.ShippedAt == nil {
			shippedAt = &now
		}
	}
	if req.Status == "cancelled" && existing.Status != "cancelled" {
		if err := h.app.Models.Item.AdjustQuantity(r.Context(), existing.ItemID, existing.Quantity); err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to restore item quantity")
			return
		}
	}

	err = h.app.Models.Shipment.UpdateStatus(r.Context(), shipmentID, req.Status, shippedAt, receivedAt)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update shipment status")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Changed shipment #" + formatInt64(shipmentID) + " status to " + req.Status,
		EntityType:  "shipments",
		EntityID:    shipmentID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedShipment, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || updatedShipment == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated shipment")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "shipment status updated successfully", updatedShipment)
}

func (h *Handler) DeleteShipmentHandler(w http.ResponseWriter, r *http.Request) {
	shipmentID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid shipment id")
		return
	}

	existing, err := h.app.Models.Shipment.GetByID(r.Context(), shipmentID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "shipment not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if existing.Status != "cancelled" {
		if err := h.app.Models.Item.AdjustQuantity(r.Context(), existing.ItemID, existing.Quantity); err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to restore item quantity")
			return
		}
	}

	err = h.app.Models.Shipment.Delete(r.Context(), shipmentID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete shipment")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted shipment #" + formatInt64(shipmentID),
		EntityType:  "shipments",
		EntityID:    shipmentID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "shipment deleted successfully", nil)
}

func formatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
