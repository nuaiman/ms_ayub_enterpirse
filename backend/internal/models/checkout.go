package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Checkout struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ItemID        int64     `json:"item_id"`
	Quantity      int       `json:"quantity"`
	CheckoutDate  time.Time `json:"checkout_date"`
	ReceiverName  *string   `json:"receiver_name"`
	ReceiverPhone *string   `json:"receiver_phone"`
	Type          string    `json:"type"` // pickup or delivery

	// Delivery Details (only for delivery type)
	VehicleNumber *string `json:"vehicle_number"`
	DriverName    *string `json:"driver_name"`
	DriverPhone   *string `json:"driver_phone"`
	FromLocation  *string `json:"from_location"`
	ToLocation    *string `json:"to_location"`

	// Financial (only for delivery)
	DeliveryCharge *float64 `json:"delivery_charge"`
	DeliveryCost   *float64 `json:"delivery_cost"`
	CustomerPaid   *float64 `json:"customer_paid"`

	// Status
	Status   string  `json:"status"` // pending, in_transit, complete
	Notes    *string `json:"notes"`
	ImageURL *string `json:"image_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CheckoutModel struct {
	DB *sql.DB
}

func (m *CheckoutModel) Insert(ctx context.Context, checkout *Checkout) (int64, error) {
	query := `
		INSERT INTO checkouts (
			user_id, item_id, quantity, checkout_date,
			receiver_name, receiver_phone, type,
			vehicle_number, driver_name, driver_phone,
			from_location, to_location,
			delivery_charge, delivery_cost, customer_paid,
			status, notes, image_url
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		checkout.UserID,
		checkout.ItemID,
		checkout.Quantity,
		checkout.CheckoutDate,
		checkout.ReceiverName,
		checkout.ReceiverPhone,
		checkout.Type,
		checkout.VehicleNumber,
		checkout.DriverName,
		checkout.DriverPhone,
		checkout.FromLocation,
		checkout.ToLocation,
		checkout.DeliveryCharge,
		checkout.DeliveryCost,
		checkout.CustomerPaid,
		checkout.Status,
		checkout.Notes,
		checkout.ImageURL,
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

func (m *CheckoutModel) GetByID(ctx context.Context, id int64) (*Checkout, error) {
	query := `
		SELECT id, user_id, item_id, quantity, checkout_date,
		       receiver_name, receiver_phone, type,
		       vehicle_number, driver_name, driver_phone,
		       from_location, to_location,
		       delivery_charge, delivery_cost, customer_paid,
		       status, notes, image_url,
		       created_at, updated_at
		FROM checkouts
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	c := &Checkout{}

	err := row.Scan(
		&c.ID, &c.UserID, &c.ItemID, &c.Quantity, &c.CheckoutDate,
		&c.ReceiverName, &c.ReceiverPhone, &c.Type,
		&c.VehicleNumber, &c.DriverName, &c.DriverPhone,
		&c.FromLocation, &c.ToLocation,
		&c.DeliveryCharge, &c.DeliveryCost, &c.CustomerPaid,
		&c.Status, &c.Notes, &c.ImageURL,
		&c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return c, nil
}

func (m *CheckoutModel) GetAll(ctx context.Context) ([]Checkout, error) {
	query := `
		SELECT id, user_id, item_id, quantity, checkout_date,
		       receiver_name, receiver_phone, type,
		       vehicle_number, driver_name, driver_phone,
		       from_location, to_location,
		       delivery_charge, delivery_cost, customer_paid,
		       status, notes, image_url,
		       created_at, updated_at
		FROM checkouts
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []Checkout{}

	for rows.Next() {
		var c Checkout

		err := rows.Scan(
			&c.ID, &c.UserID, &c.ItemID, &c.Quantity, &c.CheckoutDate,
			&c.ReceiverName, &c.ReceiverPhone, &c.Type,
			&c.VehicleNumber, &c.DriverName, &c.DriverPhone,
			&c.FromLocation, &c.ToLocation,
			&c.DeliveryCharge, &c.DeliveryCost, &c.CustomerPaid,
			&c.Status, &c.Notes, &c.ImageURL,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		checkouts = append(checkouts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkouts, nil
}

func (m *CheckoutModel) GetByItemID(ctx context.Context, itemID int64) ([]Checkout, error) {
	query := `
		SELECT id, user_id, item_id, quantity, checkout_date,
		       receiver_name, receiver_phone, type,
		       vehicle_number, driver_name, driver_phone,
		       from_location, to_location,
		       delivery_charge, delivery_cost, customer_paid,
		       status, notes, image_url,
		       created_at, updated_at
		FROM checkouts
		WHERE item_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []Checkout{}

	for rows.Next() {
		var c Checkout

		err := rows.Scan(
			&c.ID, &c.UserID, &c.ItemID, &c.Quantity, &c.CheckoutDate,
			&c.ReceiverName, &c.ReceiverPhone, &c.Type,
			&c.VehicleNumber, &c.DriverName, &c.DriverPhone,
			&c.FromLocation, &c.ToLocation,
			&c.DeliveryCharge, &c.DeliveryCost, &c.CustomerPaid,
			&c.Status, &c.Notes, &c.ImageURL,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		checkouts = append(checkouts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkouts, nil
}

func (m *CheckoutModel) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE checkouts
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, status, id)
	return err
}

func (m *CheckoutModel) UpdateFields(ctx context.Context, id int64, updates map[string]any) error {
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

	query := "UPDATE checkouts SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *CheckoutModel) Update(ctx context.Context, id int64, checkout *Checkout) error {
	query := `
		UPDATE checkouts
		SET quantity = ?, receiver_name = ?, receiver_phone = ?,
		    vehicle_number = ?, driver_name = ?, driver_phone = ?,
		    from_location = ?, to_location = ?,
		    delivery_charge = ?, delivery_cost = ?, customer_paid = ?,
		    status = ?, notes = ?, image_url = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		checkout.Quantity,
		checkout.ReceiverName,
		checkout.ReceiverPhone,
		checkout.VehicleNumber,
		checkout.DriverName,
		checkout.DriverPhone,
		checkout.FromLocation,
		checkout.ToLocation,
		checkout.DeliveryCharge,
		checkout.DeliveryCost,
		checkout.CustomerPaid,
		checkout.Status,
		checkout.Notes,
		checkout.ImageURL,
		id,
	)
	return err
}

func (m *CheckoutModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE checkouts
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *CheckoutModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM checkouts WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
