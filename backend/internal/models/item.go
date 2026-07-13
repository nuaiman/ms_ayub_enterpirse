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
	UserID int64   `json:"user_id"`
	Notes  *string `json:"notes"`

	ProductName string  `json:"product_name"`
	StorageName *string `json:"storage_name"`
	AccountName *string `json:"account_name"`
	LotNumber   *string `json:"lot_number"`

	CustomerPhone *string `json:"customer_phone"`
	CustomerEmail *string `json:"customer_email"`

	Category    *string `json:"category"`
	Subcategory *string `json:"subcategory"`

	QuantityUnit string   `json:"quantity_unit"`
	Quantity     int      `json:"quantity"`
	Weight       *float64 `json:"weight"`
	WeightUnit   *string  `json:"weight_unit"`

	Amount       float64 `json:"amount"`
	Deposit      float64 `json:"deposit"`
	CustomerPaid float64 `json:"customer_paid"`

	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemModel struct {
	DB *sql.DB
}

func (m *ItemModel) Insert(ctx context.Context, item *Item) (int64, error) {
	query := `
		INSERT INTO items (
			user_id, notes,
			product_name, storage_name, account_name, lot_number,
			customer_phone, customer_email,
			category, subcategory,
			quantity_unit, quantity, weight, weight_unit,
			amount, deposit, customer_paid,
			image_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		item.UserID, item.Notes,
		item.ProductName, item.StorageName, item.AccountName, item.LotNumber,
		item.CustomerPhone, item.CustomerEmail,
		item.Category, item.Subcategory,
		item.QuantityUnit, item.Quantity, item.Weight, item.WeightUnit,
		item.Amount, item.Deposit, item.CustomerPaid,
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
		SELECT id, user_id, notes,
		       product_name, storage_name, account_name, lot_number,
		       customer_phone, customer_email,
		       category, subcategory,
		       quantity_unit, quantity, weight, weight_unit,
		       amount, deposit, customer_paid,
		       image_url, created_at, updated_at
		FROM items WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	item := &Item{}

	err := row.Scan(
		&item.ID, &item.UserID, &item.Notes,
		&item.ProductName, &item.StorageName, &item.AccountName, &item.LotNumber,
		&item.CustomerPhone, &item.CustomerEmail,
		&item.Category, &item.Subcategory,
		&item.QuantityUnit, &item.Quantity, &item.Weight, &item.WeightUnit,
		&item.Amount, &item.Deposit, &item.CustomerPaid,
		&item.ImageURL, &item.CreatedAt, &item.UpdatedAt,
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
		SELECT id, user_id, notes,
		       product_name, storage_name, account_name, lot_number,
		       customer_phone, customer_email,
		       category, subcategory,
		       quantity_unit, quantity, weight, weight_unit,
		       amount, deposit, customer_paid,
		       image_url, created_at, updated_at
		FROM items ORDER BY id DESC
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
			&i.ID, &i.UserID, &i.Notes,
			&i.ProductName, &i.StorageName, &i.AccountName, &i.LotNumber,
			&i.CustomerPhone, &i.CustomerEmail,
			&i.Category, &i.Subcategory,
			&i.QuantityUnit, &i.Quantity, &i.Weight, &i.WeightUnit,
			&i.Amount, &i.Deposit, &i.CustomerPaid,
			&i.ImageURL, &i.CreatedAt, &i.UpdatedAt,
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

func (m *ItemModel) UpdateFields(ctx context.Context, id int64, updates map[string]any) error {
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
	query := "UPDATE items SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *ItemModel) Update(ctx context.Context, id int64, item *Item) error {
	query := `
		UPDATE items SET notes = ?,
		    product_name = ?, storage_name = ?, account_name = ?, lot_number = ?,
		    customer_phone = ?, customer_email = ?,
		    category = ?, subcategory = ?,
		    quantity_unit = ?, quantity = ?, weight = ?, weight_unit = ?,
		    amount = ?, deposit = ?, customer_paid = ?,
		    image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		item.Notes,
		item.ProductName, item.StorageName, item.AccountName, item.LotNumber,
		item.CustomerPhone, item.CustomerEmail,
		item.Category, item.Subcategory,
		item.QuantityUnit, item.Quantity, item.Weight, item.WeightUnit,
		item.Amount, item.Deposit, item.CustomerPaid,
		item.ImageURL, id,
	)
	return err
}

func (m *ItemModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	_, err := m.DB.ExecContext(ctx, `UPDATE items SET image_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, imageURL, id)
	return err
}

func (m *ItemModel) Delete(ctx context.Context, id int64) error {
	_, err := m.DB.ExecContext(ctx, `DELETE FROM items WHERE id = ?`, id)
	return err
}

func (m *ItemModel) AdjustQuantity(ctx context.Context, id int64, delta int) error {
	_, err := m.DB.ExecContext(ctx, `UPDATE items SET quantity = quantity + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, delta, id)
	return err
}
