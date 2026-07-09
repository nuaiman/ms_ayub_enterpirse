package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Expense struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

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
	ImageURL    *string `json:"image_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExpenseModel struct {
	DB *sql.DB
}

func (m *ExpenseModel) Insert(ctx context.Context, expense *Expense) (int64, error) {
	query := `
		INSERT INTO expenses (
			user_id, is_salary, salary_user_id, salary_month_year,
			title, category, amount, expense_date, notes, image_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		expense.UserID,
		expense.IsSalary,
		expense.SalaryUserID,
		expense.SalaryMonthYear,
		expense.Title,
		expense.Category,
		expense.Amount,
		expense.ExpenseDate,
		expense.Notes,
		expense.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ExpenseModel) GetByID(ctx context.Context, id int64) (*Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	e := &Expense{}

	err := row.Scan(
		&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
		&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
		&e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return e, nil
}

func (m *ExpenseModel) GetAll(ctx context.Context) ([]Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		ORDER BY expense_date DESC, id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}

	for rows.Next() {
		var e Expense

		err := rows.Scan(
			&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
			&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (m *ExpenseModel) GetByCategory(ctx context.Context, category string) ([]Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		WHERE category = ?
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}

	for rows.Next() {
		var e Expense

		err := rows.Scan(
			&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
			&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (m *ExpenseModel) GetSalaries(ctx context.Context) ([]Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		WHERE is_salary = 1
		ORDER BY salary_month_year DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}

	for rows.Next() {
		var e Expense

		err := rows.Scan(
			&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
			&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (m *ExpenseModel) GetSalaryByUserAndMonth(ctx context.Context, userID int64, monthYear string) (*Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		WHERE is_salary = 1 AND salary_user_id = ? AND salary_month_year = ?
	`

	row := m.DB.QueryRowContext(ctx, query, userID, monthYear)

	e := &Expense{}

	err := row.Scan(
		&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
		&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
		&e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return e, nil
}

func (m *ExpenseModel) GetByDateRange(ctx context.Context, startDate, endDate string) ([]Expense, error) {
	query := `
		SELECT id, user_id, is_salary, salary_user_id, salary_month_year,
		       title, category, amount, expense_date, notes, image_url,
		       created_at, updated_at
		FROM expenses
		WHERE expense_date BETWEEN ? AND ?
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}

	for rows.Next() {
		var e Expense

		err := rows.Scan(
			&e.ID, &e.UserID, &e.IsSalary, &e.SalaryUserID, &e.SalaryMonthYear,
			&e.Title, &e.Category, &e.Amount, &e.ExpenseDate, &e.Notes, &e.ImageURL,
			&e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (m *ExpenseModel) Update(ctx context.Context, id int64, expense *Expense) error {
	query := `
		UPDATE expenses
		SET is_salary = ?, salary_user_id = ?, salary_month_year = ?,
		    title = ?, category = ?, amount = ?, expense_date = ?, 
		    notes = ?, image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		expense.IsSalary,
		expense.SalaryUserID,
		expense.SalaryMonthYear,
		expense.Title,
		expense.Category,
		expense.Amount,
		expense.ExpenseDate,
		expense.Notes,
		expense.ImageURL,
		id,
	)
	return err
}

func (m *ExpenseModel) UpdateFields(ctx context.Context, id int64, updates map[string]any) error {
	var setClauses []string
	var args []any

	for field, value := range updates {
		setClauses = append(setClauses, field+" = ?")
		args = append(args, value)
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)

	query := "UPDATE expenses SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *ExpenseModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE expenses
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *ExpenseModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
