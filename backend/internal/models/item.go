package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Item struct {
	ID           int64     `json:"id"`
	UserID       *int64    `json:"user_id"`
	Notes        *string   `json:"notes"`
	ContractID   *int64    `json:"contract_id"`
	Name         string    `json:"name"`
	QuantityUnit string    `json:"quantity_unit"`
	Quantity     int       `json:"quantity"`
	Weight       *float64  `json:"weight"`
	Width        *float64  `json:"width"`
	Height       *float64  `json:"height"`
	Length       *float64  `json:"length"`
	ImageURL     *string   `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ItemModel struct {
	DB *sql.DB
}

func (m *ItemModel) Insert(ctx context.Context, item *Item) (int64, error) {
	query := `
		INSERT INTO items (
			user_id, notes, contract_id, name, quantity_unit,
			quantity, weight, width, height, length, image_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		item.UserID,
		item.Notes,
		item.ContractID,
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
		SELECT id, user_id, notes, contract_id, name, quantity_unit,
		       quantity, weight, width, height, length, image_url,
		       created_at, updated_at
		FROM items
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	item := &Item{}

	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.Notes,
		&item.ContractID,
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
		SELECT id, user_id, notes, contract_id, name, quantity_unit,
		       quantity, weight, width, height, length, image_url,
		       created_at, updated_at
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
			&i.ContractID,
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

func (m *ItemModel) GetByContractID(ctx context.Context, contractID int64) ([]Item, error) {
	query := `
		SELECT id, user_id, notes, contract_id, name, quantity_unit,
		       quantity, weight, width, height, length, image_url,
		       created_at, updated_at
		FROM items
		WHERE contract_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, contractID)
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
			&i.ContractID,
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

func (m *ItemModel) Update(ctx context.Context, id int64, item *Item) error {
	query := `
		UPDATE items
		SET notes = ?, contract_id = ?, name = ?, quantity_unit = ?,
		    quantity = ?, weight = ?, width = ?, height = ?, length = ?,
		    image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		item.Notes,
		item.ContractID,
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

func (m *ItemModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE items
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
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
