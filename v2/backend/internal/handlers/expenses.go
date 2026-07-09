package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
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
		// Salary Tracking
		IsSalary        bool    `json:"is_salary"`
		SalaryUserID    *int64  `json:"salary_user_id"`
		SalaryMonthYear *string `json:"salary_month_year"`

		// Expense Details
		Title       string  `json:"title"`
		Category    *string `json:"category"`
		Amount      float64 `json:"amount"`
		ExpenseDate string  `json:"expense_date"`
		Notes       *string `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" || req.Amount <= 0 || req.ExpenseDate == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "title, amount, and expense_date are required")
		return
	}

	// Validate salary fields
	if req.IsSalary {
		if req.SalaryUserID == nil {
			utils.ErrorJson(w, http.StatusBadRequest, "salary_user_id is required for salary expenses")
			return
		}
		if req.SalaryMonthYear == nil || *req.SalaryMonthYear == "" {
			utils.ErrorJson(w, http.StatusBadRequest, "salary_month_year is required for salary expenses")
			return
		}

		// Check if salary already exists for this user and month
		existing, err := h.app.Models.Expense.GetSalaryByUserAndMonth(r.Context(), *req.SalaryUserID, *req.SalaryMonthYear)
		if err == nil && existing != nil {
			utils.ErrorJson(w, http.StatusConflict, "salary already exists for this user and month")
			return
		}

		// Validate salary user exists
		salaryUser, err := h.app.Models.User.GetByID(r.Context(), *req.SalaryUserID)
		if err != nil || salaryUser == nil {
			utils.ErrorJson(w, http.StatusNotFound, "salary user not found")
			return
		}

		// Auto-fill category for salary
		if req.Category == nil {
			req.Category = new(string)
			*req.Category = "salary"
		}
	}

	// Validate expense date format
	_, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid expense_date format, use YYYY-MM-DD")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	expense := &models.Expense{
		UserID:          userID,
		IsSalary:        req.IsSalary,
		SalaryUserID:    req.SalaryUserID,
		SalaryMonthYear: req.SalaryMonthYear,
		Title:           req.Title,
		Category:        req.Category,
		Amount:          req.Amount,
		ExpenseDate:     req.ExpenseDate,
		Notes:           req.Notes,
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
	description := fmt.Sprintf("Created expense: %s ($%.2f)", req.Title, req.Amount)
	if req.IsSalary {
		description = fmt.Sprintf("Created salary payment for user #%d (%s)", *req.SalaryUserID, *req.SalaryMonthYear)
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "create",
		Description: description,
		EntityType:  "expenses",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "expense created successfully", expense)
}

func (h *Handler) UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		// Salary Tracking
		IsSalary        *bool   `json:"is_salary"`
		SalaryUserID    *int64  `json:"salary_user_id"`
		SalaryMonthYear *string `json:"salary_month_year"`

		// Expense Details
		Title       *string  `json:"title"`
		Category    *string  `json:"category"`
		Amount      *float64 `json:"amount"`
		ExpenseDate *string  `json:"expense_date"`
		Notes       *string  `json:"notes"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate fields if provided
	if req.Title != nil && *req.Title == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	if req.Amount != nil && *req.Amount <= 0 {
		utils.ErrorJson(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	if req.ExpenseDate != nil {
		_, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err != nil {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid expense_date format, use YYYY-MM-DD")
			return
		}
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

	// Build update map with only provided fields
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	if req.IsSalary != nil {
		updates["is_salary"] = *req.IsSalary
	}
	if req.SalaryUserID != nil {
		updates["salary_user_id"] = *req.SalaryUserID
	}
	if req.SalaryMonthYear != nil {
		updates["salary_month_year"] = *req.SalaryMonthYear
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.ExpenseDate != nil {
		updates["expense_date"] = *req.ExpenseDate
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	if len(updates) > 1 {
		err = h.app.Models.Expense.UpdateFields(r.Context(), expenseID, updates)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update expense")
			return
		}
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	description := fmt.Sprintf("Updated expense: %s", existing.Title)
	if existing.IsSalary {
		description = fmt.Sprintf("Updated salary payment for user #%d", *existing.SalaryUserID)
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: description,
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

	// Delete associated image file if exists
	if existing.ImageURL != nil && *existing.ImageURL != "" {
		imagePath := "." + *existing.ImageURL
		os.Remove(imagePath)
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	description := fmt.Sprintf("Deleted expense: %s ($%.2f)", existing.Title, existing.Amount)
	if existing.IsSalary {
		description = fmt.Sprintf("Deleted salary payment for user #%d (%s)", *existing.SalaryUserID, *existing.SalaryMonthYear)
	}

	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "delete",
		Description: description,
		EntityType:  "expenses",
		EntityID:    expenseID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "expense deleted successfully", nil)
}
