package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
)

func (h *Handler) GetExpenseHandler(w http.ResponseWriter, r *http.Request) {
	expenseID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	expense, err := h.app.Models.Expense.GetByID(r.Context(), expenseID)
	if err != nil || expense == nil {
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "expense details fetched", expense)
}

func (h *Handler) GetAllExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.app.Models.Expense.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expenses")
		return
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	utils.SuccessJson(w, http.StatusOK, "expenses fetched successfully", expenses)
}

func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title       string  `json:"title"`
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `json:"expense_date"`
		Notes       *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.Category == "" || req.Amount <= 0 || req.ExpenseDate == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validCategories := map[string]bool{
		"rent": true, "salary": true, "bill": true, "other": true,
	}
	if !validCategories[req.Category] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid category")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	expense := &models.Expense{
		UserID:      &userID,
		Title:       req.Title,
		Category:    req.Category,
		Amount:      req.Amount,
		ExpenseDate: req.ExpenseDate,
		Notes:       req.Notes,
	}

	id, err := h.app.Models.Expense.Insert(r.Context(), expense)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create expense")
		return
	}

	expense.ID = id

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: "Created expense: " + req.Title,
		EntityType:  "expenses",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "expense created successfully", expense)
}

func (h *Handler) UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title       string  `json:"title"`
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `json:"expense_date"`
		Notes       *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.Category == "" || req.Amount <= 0 || req.ExpenseDate == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	validCategories := map[string]bool{
		"rent": true, "salary": true, "bill": true, "other": true,
	}
	if !validCategories[req.Category] {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid category")
		return
	}

	expenseID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	existing, err := h.app.Models.Expense.GetByID(r.Context(), expenseID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	expense := &models.Expense{
		Title:       req.Title,
		Category:    req.Category,
		Amount:      req.Amount,
		ExpenseDate: req.ExpenseDate,
		Notes:       req.Notes,
		ImageURL:    existing.ImageURL,
	}

	err = h.app.Models.Expense.Update(r.Context(), expenseID, expense)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update expense")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Updated expense: " + req.Title,
		EntityType:  "expenses",
		EntityID:    expenseID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	updatedExpense, err := h.app.Models.Expense.GetByID(r.Context(), expenseID)
	if err != nil || updatedExpense == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated expense")
		return
	}

	utils.SuccessJson(w, http.StatusOK, "expense updated successfully", updatedExpense)
}

func (h *Handler) DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	expenseID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	existing, err := h.app.Models.Expense.GetByID(r.Context(), expenseID)
	if err != nil || existing == nil {
		utils.ErrorJson(w, http.StatusNotFound, "expense not found")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.app.Models.Expense.Delete(r.Context(), expenseID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to delete expense")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: "Deleted expense: " + existing.Title,
		EntityType:  "expenses",
		EntityID:    expenseID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "expense deleted successfully", nil)
}
