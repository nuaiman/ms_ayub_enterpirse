package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Contract struct {
	ID              int64      `json:"id"`
	UserID          *int64     `json:"user_id"`
	Notes           *string    `json:"notes"`
	CustomerID      int64      `json:"customer_id"`
	DurationType    string     `json:"duration_type"`
	Duration        *int       `json:"duration"`
	Price           *float64   `json:"price"`
	SecurityDeposit *float64   `json:"security_deposit"`
	EstimatedValue  *float64   `json:"estimated_value"`
	Status          string     `json:"status"`
	EndedAt         *time.Time `json:"ended_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ContractModel struct {
	DB *sql.DB
}

func (m *ContractModel) Insert(ctx context.Context, contract *Contract) (int64, error) {
	query := `
		INSERT INTO contracts (
			user_id, notes, customer_id, duration_type, duration,
			price, security_deposit, estimated_value, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		contract.UserID,
		contract.Notes,
		contract.CustomerID,
		contract.DurationType,
		contract.Duration,
		contract.Price,
		contract.SecurityDeposit,
		contract.EstimatedValue,
		contract.Status,
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

func (m *ContractModel) GetByID(ctx context.Context, id int64) (*Contract, error) {
	query := `
		SELECT id, user_id, notes, customer_id, duration_type, duration,
		       price, security_deposit, estimated_value, status, ended_at,
		       created_at, updated_at
		FROM contracts
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	contract := &Contract{}

	err := row.Scan(
		&contract.ID,
		&contract.UserID,
		&contract.Notes,
		&contract.CustomerID,
		&contract.DurationType,
		&contract.Duration,
		&contract.Price,
		&contract.SecurityDeposit,
		&contract.EstimatedValue,
		&contract.Status,
		&contract.EndedAt,
		&contract.CreatedAt,
		&contract.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return contract, nil
}

func (m *ContractModel) GetAll(ctx context.Context) ([]Contract, error) {
	query := `
		SELECT id, user_id, notes, customer_id, duration_type, duration,
		       price, security_deposit, estimated_value, status, ended_at,
		       created_at, updated_at
		FROM contracts
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contracts := []Contract{}

	for rows.Next() {
		var c Contract

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Notes,
			&c.CustomerID,
			&c.DurationType,
			&c.Duration,
			&c.Price,
			&c.SecurityDeposit,
			&c.EstimatedValue,
			&c.Status,
			&c.EndedAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (m *ContractModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Contract, error) {
	query := `
		SELECT id, user_id, notes, customer_id, duration_type, duration,
		       price, security_deposit, estimated_value, status, ended_at,
		       created_at, updated_at
		FROM contracts
		WHERE customer_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contracts := []Contract{}

	for rows.Next() {
		var c Contract

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Notes,
			&c.CustomerID,
			&c.DurationType,
			&c.Duration,
			&c.Price,
			&c.SecurityDeposit,
			&c.EstimatedValue,
			&c.Status,
			&c.EndedAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (m *ContractModel) Update(ctx context.Context, id int64, contract *Contract) error {
	query := `
		UPDATE contracts
		SET notes = ?, customer_id = ?, duration_type = ?, duration = ?,
		    price = ?, security_deposit = ?, estimated_value = ?,
		    status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		contract.Notes,
		contract.CustomerID,
		contract.DurationType,
		contract.Duration,
		contract.Price,
		contract.SecurityDeposit,
		contract.EstimatedValue,
		contract.Status,
		id,
	)
	return err
}

func (m *ContractModel) UpdateStatus(ctx context.Context, id int64, status string, endedAt *time.Time) error {
	query := `
		UPDATE contracts
		SET status = ?, ended_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, status, endedAt, id)
	return err
}

func (m *ContractModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM contracts WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
