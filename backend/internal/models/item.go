package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Item struct {
	ID     int64   `json:"id"`
	UserID *int64  `json:"user_id"`
	Notes  *string `json:"notes"`

	// Customer Relationship (Optional)
	CustomerID *int64 `json:"customer_id"`

	// Contract Fields (Optional)
	DurationType    *string    `json:"duration_type"`
	Duration        *int       `json:"duration"`
	Price           *float64   `json:"price"`
	SecurityDeposit *float64   `json:"security_deposit"`
	EstimatedValue  *float64   `json:"estimated_value"`
	Status          *string    `json:"status"`
	EndedAt         *time.Time `json:"ended_at"`

	// Item Fields
	Name         string   `json:"name"`          // Required
	QuantityUnit string   `json:"quantity_unit"` // Default: "pcs"
	Quantity     int      `json:"quantity"`      // Default: 1
	Weight       *float64 `json:"weight"`
	Width        *float64 `json:"width"`
	Height       *float64 `json:"height"`
	Length       *float64 `json:"length"`
	ImageURL     *string  `json:"image_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Helper method to check if item has a customer
func (i *Item) HasCustomer() bool {
	return i.CustomerID != nil
}

// Helper method to check if item has contract
func (i *Item) HasContract() bool {
	return i.DurationType != nil && i.Status != nil
}

type ItemModel struct {
	DB *sql.DB
}

func (m *ItemModel) Insert(ctx context.Context, item *Item) (int64, error) {
	query := `
		INSERT INTO items (
			user_id, notes, customer_id, 
			duration_type, duration, price, security_deposit, 
			estimated_value, status, ended_at,
			name, quantity_unit, quantity, weight, width, height, 
			length, image_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		item.UserID,
		item.Notes,
		item.CustomerID,
		item.DurationType,
		item.Duration,
		item.Price,
		item.SecurityDeposit,
		item.EstimatedValue,
		item.Status,
		item.EndedAt,
		item.Name,
		item.QuantityUnit,
		item.Quantity,
		item.Weight,
		item.Width,
		item.Height,
		item.Length,
		item.ImageURL,
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

func (m *ItemModel) GetByID(ctx context.Context, id int64) (*Item, error) {
	query := `
		SELECT id, user_id, notes, customer_id, 
		       duration_type, duration, price, security_deposit, 
		       estimated_value, status, ended_at,
		       name, quantity_unit, quantity, weight, width, height, 
		       length, image_url, created_at, updated_at
		FROM items
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	item := &Item{}

	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.Notes,
		&item.CustomerID,
		&item.DurationType,
		&item.Duration,
		&item.Price,
		&item.SecurityDeposit,
		&item.EstimatedValue,
		&item.Status,
		&item.EndedAt,
		&item.Name,
		&item.QuantityUnit,
		&item.Quantity,
		&item.Weight,
		&item.Width,
		&item.Height,
		&item.Length,
		&item.ImageURL,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

func (m *ItemModel) GetAll(ctx context.Context) ([]Item, error) {
	query := `
		SELECT id, user_id, notes, customer_id, 
		       duration_type, duration, price, security_deposit, 
		       estimated_value, status, ended_at,
		       name, quantity_unit, quantity, weight, width, height, 
		       length, image_url, created_at, updated_at
		FROM items
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}

	for rows.Next() {
		var i Item

		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Notes,
			&i.CustomerID,
			&i.DurationType,
			&i.Duration,
			&i.Price,
			&i.SecurityDeposit,
			&i.EstimatedValue,
			&i.Status,
			&i.EndedAt,
			&i.Name,
			&i.QuantityUnit,
			&i.Quantity,
			&i.Weight,
			&i.Width,
			&i.Height,
			&i.Length,
			&i.ImageURL,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (m *ItemModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Item, error) {
	query := `
		SELECT id, user_id, notes, customer_id, 
		       duration_type, duration, price, security_deposit, 
		       estimated_value, status, ended_at,
		       name, quantity_unit, quantity, weight, width, height, 
		       length, image_url, created_at, updated_at
		FROM items
		WHERE customer_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}

	for rows.Next() {
		var i Item

		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Notes,
			&i.CustomerID,
			&i.DurationType,
			&i.Duration,
			&i.Price,
			&i.SecurityDeposit,
			&i.EstimatedValue,
			&i.Status,
			&i.EndedAt,
			&i.Name,
			&i.QuantityUnit,
			&i.Quantity,
			&i.Weight,
			&i.Width,
			&i.Height,
			&i.Length,
			&i.ImageURL,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Add this method to your ItemModel
func (m *ItemModel) UpdateFields(ctx context.Context, id int64, updates map[string]any) error {
	// Build the SET clause dynamically
	var setClauses []string
	var args []any

	for field, value := range updates {
		setClauses = append(setClauses, field+" = ?")
		args = append(args, value)
	}

	// Add id as the last argument for WHERE clause
	args = append(args, id)

	query := "UPDATE items SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// Keep the original Update method for backward compatibility
func (m *ItemModel) Update(ctx context.Context, id int64, item *Item) error {
	query := `
		UPDATE items
		SET notes = ?, customer_id = ?,
		    duration_type = ?, duration = ?, price = ?, 
		    security_deposit = ?, estimated_value = ?, status = ?,
		    name = ?, quantity_unit = ?, quantity = ?,
		    weight = ?, width = ?, height = ?, length = ?,
		    image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		item.Notes,
		item.CustomerID,
		item.DurationType,
		item.Duration,
		item.Price,
		item.SecurityDeposit,
		item.EstimatedValue,
		item.Status,
		item.Name,
		item.QuantityUnit,
		item.Quantity,
		item.Weight,
		item.Width,
		item.Height,
		item.Length,
		item.ImageURL,
		id,
	)
	return err
}
func (m *ItemModel) UpdateCustomer(ctx context.Context, id int64, customerID *int64) error {
	query := `
		UPDATE items
		SET customer_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, customerID, id)
	return err
}

func (m *ItemModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE items
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *ItemModel) UpdateStatus(ctx context.Context, id int64, status string, endedAt *time.Time) error {
	query := `
		UPDATE items
		SET status = ?, ended_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, status, endedAt, id)
	return err
}

func (m *ItemModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *ItemModel) AdjustQuantity(ctx context.Context, id int64, delta int) error {
	query := `
		UPDATE items
		SET quantity = quantity + ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, delta, id)
	return err
}
