package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Customer struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     *string   `json:"email"`
	Company   *string   `json:"company"`
	Address   *string   `json:"address"`
	IDType    *string   `json:"id_type"`
	IDNumber  *string   `json:"id_number"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerModel struct {
	DB *sql.DB
}

func (m *CustomerModel) Insert(ctx context.Context, customer *Customer) (int64, error) {
	query := `
		INSERT INTO customers (
			user_id,
			name,
			phone,
			email,
			company,
			address,
			id_type,
			id_number,
			notes
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		customer.UserID,
		customer.Name,
		customer.Phone,
		customer.Email,
		customer.Company,
		customer.Address,
		customer.IDType,
		customer.IDNumber,
		customer.Notes,
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

func (m *CustomerModel) GetByID(ctx context.Context, id int64) (*Customer, error) {
	query := `
		SELECT id, user_id, name, phone, email, company, address, 
		       id_type, id_number, notes, created_at, updated_at
		FROM customers
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	customer := &Customer{}

	err := row.Scan(
		&customer.ID,
		&customer.UserID,
		&customer.Name,
		&customer.Phone,
		&customer.Email,
		&customer.Company,
		&customer.Address,
		&customer.IDType,
		&customer.IDNumber,
		&customer.Notes,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return customer, nil
}

func (m *CustomerModel) GetByPhone(ctx context.Context, phone string) (*Customer, error) {
	query := `
		SELECT id, user_id, name, phone, email, company, address, 
		       id_type, id_number, notes, created_at, updated_at
		FROM customers
		WHERE phone = ?
	`

	row := m.DB.QueryRowContext(ctx, query, phone)

	customer := &Customer{}

	err := row.Scan(
		&customer.ID,
		&customer.UserID,
		&customer.Name,
		&customer.Phone,
		&customer.Email,
		&customer.Company,
		&customer.Address,
		&customer.IDType,
		&customer.IDNumber,
		&customer.Notes,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return customer, nil
}

func (m *CustomerModel) GetAll(ctx context.Context) ([]Customer, error) {
	query := `
		SELECT id, user_id, name, phone, email, company, address, 
		       id_type, id_number, notes, created_at, updated_at
		FROM customers
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Company,
			&c.Address,
			&c.IDType,
			&c.IDNumber,
			&c.Notes,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (m *CustomerModel) GetByUserID(ctx context.Context, userID int64) ([]Customer, error) {
	query := `
		SELECT id, user_id, name, phone, email, company, address, 
		       id_type, id_number, notes, created_at, updated_at
		FROM customers
		WHERE user_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Phone,
			&c.Email,
			&c.Company,
			&c.Address,
			&c.IDType,
			&c.IDNumber,
			&c.Notes,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (m *CustomerModel) Update(ctx context.Context, id int64, customer *Customer) error {
	query := `
		UPDATE customers
		SET name = ?, phone = ?, email = ?, company = ?, address = ?, 
		    id_type = ?, id_number = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		customer.Name,
		customer.Phone,
		customer.Email,
		customer.Company,
		customer.Address,
		customer.IDType,
		customer.IDNumber,
		customer.Notes,
		id,
	)
	return err
}

func (m *CustomerModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
