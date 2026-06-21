package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	Address      *string   `json:"address"`
	IDType       *string   `json:"id_type"`
	IDNumber     *string   `json:"id_number"`
	ImageURL     *string   `json:"image_url"`
	Password     string    `json:"-"`
	Role         string    `json:"role"`
	RefreshToken *string   `json:"-"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(ctx context.Context, user *User) (int64, error) {
	query := `
		INSERT INTO users (
			name,
			username,
			email,
			phone,
			address,
			id_type,
			id_number,
			image_url,
			password,
			role
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		user.Name,
		user.Username,
		user.Email,
		user.Phone,
		user.Address,
		user.IDType,
		user.IDNumber,
		user.ImageURL,
		user.Password,
		user.Role,
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

func (m *UserModel) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, id_type, id_number, image_url, 
		       password, role, refresh_token, is_active, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (m *UserModel) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, id_type, id_number, image_url,
		       password, role, refresh_token, is_active, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	row := m.DB.QueryRowContext(ctx, query, username)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (m *UserModel) GetByRefreshToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, id_type, id_number, image_url,
		       password, role, refresh_token, is_active, created_at, updated_at
		FROM users
		WHERE refresh_token = ?
	`

	row := m.DB.QueryRowContext(ctx, query, token)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, id_type, id_number, image_url,
		       password, role, refresh_token, is_active, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	row := m.DB.QueryRowContext(ctx, query, email)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (m *UserModel) GetAll(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, username, email, phone, address, id_type, id_number, image_url,
		       password, role, refresh_token, is_active, created_at, updated_at
		FROM users
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
			&u.Email,
			&u.Phone,
			&u.Address,
			&u.IDType,
			&u.IDNumber,
			&u.ImageURL,
			&u.Password,
			&u.Role,
			&u.RefreshToken,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) UpdateRefreshToken(ctx context.Context, userID int64, token *string) error {
	query := `
		UPDATE users
		SET refresh_token = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, token, userID)
	return err
}

func (m *UserModel) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (m *UserModel) UpdateProfile(ctx context.Context, userID int64, user *User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, phone = ?, address = ?, 
		    id_type = ?, id_number = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.Address,
		user.IDType,
		user.IDNumber,
		userID,
	)
	return err
}

// UpdateImage updates only the image_url for a user
// Pass nil to remove the image, or a string pointer to set a new URL
func (m *UserModel) UpdateImage(ctx context.Context, userID int64, imageURL *string) error {
	query := `
		UPDATE users
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, userID)
	return err
}

func (m *UserModel) ResetAllPasswordsExceptRoot(ctx context.Context, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id != 1
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword)
	return err
}

func (m *UserModel) AdminUpdatePasswordAndInvalidateSessions(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (m *UserModel) UpdateRole(ctx context.Context, userID int64, role string) error {
	query := `
		UPDATE users
		SET role = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, role, userID)
	return err
}

func (m *UserModel) SetActiveStatus(ctx context.Context, userID int64, active bool) error {
	query := `
		UPDATE users
		SET is_active = ?, refresh_token = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	val := 0
	if active {
		val = 1
	}

	_, err := m.DB.ExecContext(ctx, query, val, userID)
	return err
}
