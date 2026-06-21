package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	ExpenseDate string    `json:"expense_date"`
	Notes       *string   `json:"notes"`
	ImageURL    *string   `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ExpenseModel struct {
	DB *sql.DB
}

func (m *ExpenseModel) Insert(ctx context.Context, expense *Expense) (int64, error) {
	query := `
		INSERT INTO expenses (user_id, title, category, amount, expense_date, notes, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		expense.UserID,
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
		SELECT id, user_id, title, category, amount, expense_date, notes, image_url, created_at
		FROM expenses
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	e := &Expense{}

	err := row.Scan(
		&e.ID, &e.UserID, &e.Title, &e.Category, &e.Amount,
		&e.ExpenseDate, &e.Notes, &e.ImageURL, &e.CreatedAt,
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
		SELECT id, user_id, title, category, amount, expense_date, notes, image_url, created_at
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
			&e.ID, &e.UserID, &e.Title, &e.Category, &e.Amount,
			&e.ExpenseDate, &e.Notes, &e.ImageURL, &e.CreatedAt,
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
		SET title = ?, category = ?, amount = ?, expense_date = ?, notes = ?, image_url = ?
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
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

func (m *ExpenseModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE expenses
		SET image_url = ?
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
