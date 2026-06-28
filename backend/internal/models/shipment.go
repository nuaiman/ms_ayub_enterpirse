package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Shipment struct {
	ID             int64      `json:"id"`
	UserID         *int64     `json:"user_id"`
	ShipmentType   string     `json:"shipment_type"`
	ItemID         int64      `json:"item_id"`
	Quantity       int        `json:"quantity"`
	VehicleNumber  *string    `json:"vehicle_number"`
	DriverName     *string    `json:"driver_name"`
	DriverPhone    *string    `json:"driver_phone"`
	FromLocation   *string    `json:"from_location"`
	CompanyName    *string    `json:"company_name"`
	ToLocation     *string    `json:"to_location"`
	ReceiverName   *string    `json:"receiver_name"`
	ReceiverPhone  *string    `json:"receiver_phone"`
	CustomerCharge *float64   `json:"customer_charge"`
	CompanyCost    *float64   `json:"company_cost"`
	CompanyPaid    *float64   `json:"company_paid"`
	CustomerPaid   *float64   `json:"customer_paid"`
	Status         string     `json:"status"`
	Notes          *string    `json:"notes"`
	ShippedAt      *time.Time `json:"shipped_at"`
	ReceivedAt     *time.Time `json:"received_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ShipmentModel struct {
	DB *sql.DB
}

func (m *ShipmentModel) Insert(ctx context.Context, shipment *Shipment) (int64, error) {
	query := `
		INSERT INTO shipments (
			user_id, shipment_type, item_id, quantity,
			vehicle_number, driver_name, driver_phone,
			from_location, company_name, to_location,
			receiver_name, receiver_phone,
			customer_charge, company_cost, company_paid, customer_paid,
			status, notes
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		shipment.UserID,
		shipment.ShipmentType,
		shipment.ItemID,
		shipment.Quantity,
		shipment.VehicleNumber,
		shipment.DriverName,
		shipment.DriverPhone,
		shipment.FromLocation,
		shipment.CompanyName,
		shipment.ToLocation,
		shipment.ReceiverName,
		shipment.ReceiverPhone,
		shipment.CustomerCharge,
		shipment.CompanyCost,
		shipment.CompanyPaid,
		shipment.CustomerPaid,
		shipment.Status,
		shipment.Notes,
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

func (m *ShipmentModel) GetByID(ctx context.Context, id int64) (*Shipment, error) {
	query := `
		SELECT id, user_id, shipment_type, item_id, quantity,
		       vehicle_number, driver_name, driver_phone,
		       from_location, company_name, to_location,
		       receiver_name, receiver_phone,
		       customer_charge, company_cost, company_paid, customer_paid,
		       status, notes, shipped_at, received_at, created_at
		FROM shipments
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	s := &Shipment{}

	err := row.Scan(
		&s.ID, &s.UserID, &s.ShipmentType, &s.ItemID, &s.Quantity,
		&s.VehicleNumber, &s.DriverName, &s.DriverPhone,
		&s.FromLocation, &s.CompanyName, &s.ToLocation,
		&s.ReceiverName, &s.ReceiverPhone,
		&s.CustomerCharge, &s.CompanyCost, &s.CompanyPaid, &s.CustomerPaid,
		&s.Status, &s.Notes, &s.ShippedAt, &s.ReceivedAt, &s.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

func (m *ShipmentModel) GetAll(ctx context.Context) ([]Shipment, error) {
	query := `
		SELECT id, user_id, shipment_type, item_id, quantity,
		       vehicle_number, driver_name, driver_phone,
		       from_location, company_name, to_location,
		       receiver_name, receiver_phone,
		       customer_charge, company_cost, company_paid, customer_paid,
		       status, notes, shipped_at, received_at, created_at
		FROM shipments
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shipments := []Shipment{}

	for rows.Next() {
		var s Shipment

		err := rows.Scan(
			&s.ID, &s.UserID, &s.ShipmentType, &s.ItemID, &s.Quantity,
			&s.VehicleNumber, &s.DriverName, &s.DriverPhone,
			&s.FromLocation, &s.CompanyName, &s.ToLocation,
			&s.ReceiverName, &s.ReceiverPhone,
			&s.CustomerCharge, &s.CompanyCost, &s.CompanyPaid, &s.CustomerPaid,
			&s.Status, &s.Notes, &s.ShippedAt, &s.ReceivedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shipments = append(shipments, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shipments, nil
}

func (m *ShipmentModel) GetByItemID(ctx context.Context, itemID int64) ([]Shipment, error) {
	query := `
		SELECT id, user_id, shipment_type, item_id, quantity,
		       vehicle_number, driver_name, driver_phone,
		       from_location, company_name, to_location,
		       receiver_name, receiver_phone,
		       customer_charge, company_cost, company_paid, customer_paid,
		       status, notes, shipped_at, received_at, created_at
		FROM shipments
		WHERE item_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shipments := []Shipment{}

	for rows.Next() {
		var s Shipment

		err := rows.Scan(
			&s.ID, &s.UserID, &s.ShipmentType, &s.ItemID, &s.Quantity,
			&s.VehicleNumber, &s.DriverName, &s.DriverPhone,
			&s.FromLocation, &s.CompanyName, &s.ToLocation,
			&s.ReceiverName, &s.ReceiverPhone,
			&s.CustomerCharge, &s.CompanyCost, &s.CompanyPaid, &s.CustomerPaid,
			&s.Status, &s.Notes, &s.ShippedAt, &s.ReceivedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shipments = append(shipments, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shipments, nil
}

func (m *ShipmentModel) UpdateStatus(ctx context.Context, id int64, status string, shippedAt *time.Time, receivedAt *time.Time) error {
	query := `
		UPDATE shipments
		SET status = ?, shipped_at = ?, received_at = ?
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, status, shippedAt, receivedAt, id)
	return err
}

func (m *ShipmentModel) Update(ctx context.Context, id int64, shipment *Shipment) error {
	query := `
		UPDATE shipments
		SET vehicle_number = ?, driver_name = ?, driver_phone = ?,
		    from_location = ?, company_name = ?, to_location = ?,
		    receiver_name = ?, receiver_phone = ?,
		    customer_charge = ?, company_cost = ?, company_paid = ?, customer_paid = ?,
		    notes = ?, status = ?
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		shipment.VehicleNumber,
		shipment.DriverName,
		shipment.DriverPhone,
		shipment.FromLocation,
		shipment.CompanyName,
		shipment.ToLocation,
		shipment.ReceiverName,
		shipment.ReceiverPhone,
		shipment.CustomerCharge,
		shipment.CompanyCost,
		shipment.CompanyPaid,
		shipment.CustomerPaid,
		shipment.Notes,
		shipment.Status,
		id,
	)
	return err
}

// internal/models/shipment.go

func (m *ShipmentModel) UpdateFields(ctx context.Context, id int64, updates map[string]any) error {
	// Build the SET clause dynamically
	var setClauses []string
	var args []any

	for field, value := range updates {
		setClauses = append(setClauses, field+" = ?")
		args = append(args, value)
	}

	// Add id as the last argument for WHERE clause
	args = append(args, id)

	query := "UPDATE shipments SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *ShipmentModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM shipments WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
