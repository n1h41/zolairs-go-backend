package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"n1h41/zolaris-backend-app/internal/domain"
)

// UserRepository handles all user-related database operations with PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(dbPool *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{
		db: dbPool,
	}
}

func (r *UserRepository) GetUserIdByCognitoId(ctx context.Context, cId string) (string, error) {
	var userId string

	query := `select user_id from z_users where cognito_id = $1`

	if err := r.db.QueryRow(ctx, query, cId).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil // User not found
		}
		return "", fmt.Errorf("failed to get user ID by Cognito ID: %w", err)
	}

	return userId, nil
}

// GetUserByID retrieves a user by ID from PostgreSQL
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT user_id, email, first_name, last_name, phone, 
		       address, parent_id, created_at, updated_at
		FROM z_users 
		WHERE user_id = $1
	`

	row := r.db.QueryRow(ctx, query, userID)

	user := &domain.User{}
	var addressJSON []byte
	var parentID *string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&addressJSON,
		&parentID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found, return nil without error
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse address from JSON
	if len(addressJSON) > 0 && string(addressJSON) != "null" {
		if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
			return nil, fmt.Errorf("failed to parse address JSON: %w", err)
		}
	}

	user.ParentID = parentID // May be nil
	return user, nil
}

// CreateUser creates a new user in PostgreSQL
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Convert address struct to JSON
	addressJSON, err := json.Marshal(user.Address)
	if err != nil {
		return fmt.Errorf("failed to convert address to JSON: %w", err)
	}

	query := `
		INSERT INTO z_users (
			user_id, email, first_name, last_name, phone,
			address, parent_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		addressJSON,
		user.ParentID,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// UpdateUser updates user in PostgreSQL
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Convert address struct to JSON
	addressJSON, err := json.Marshal(user.Address)
	if err != nil {
		return fmt.Errorf("failed to convert address to JSON: %w", err)
	}

	query := `
		UPDATE z_users SET
			first_name = $1,
			last_name = $2,
			phone = $3,
			address = $4,
			updated_at = $5
		WHERE user_id = $6
	`

	// Update the timestamp
	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		addressJSON,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", user.ID)
	}

	return nil
}

// CheckHasParentID checks if a user has a parent ID in PostgreSQL
func (r *UserRepository) CheckHasParentID(ctx context.Context, userID string) (bool, error) {
	query := `
		SELECT 
			CASE WHEN parent_id IS NULL THEN false ELSE true END 
		FROM z_users 
		WHERE user_id = $1
	`

	var hasParent bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&hasParent)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("user not found: %w", err)
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	return hasParent, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT user_id, email, first_name, last_name, phone, 
		       address, parent_id, created_at, updated_at
		FROM z_users 
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	user := &domain.User{}
	var addressJSON []byte
	var parentID *string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&addressJSON,
		&parentID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found, return nil without error
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse address from JSON
	if len(addressJSON) > 0 && string(addressJSON) != "null" {
		if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
			return nil, fmt.Errorf("failed to parse address JSON: %w", err)
		}
	}

	user.ParentID = parentID
	return user, nil
}

// GetChildUsers gets all child users for a parent user
func (r *UserRepository) GetChildUsers(ctx context.Context, parentID string) ([]*domain.User, error) {
	query := `
		SELECT user_id, email, first_name, last_name, phone, 
		       address, parent_id, created_at, updated_at
		FROM z_users 
		WHERE parent_id = $1
	`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var addressJSON []byte
		var parentID *string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&addressJSON,
			&parentID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		// Parse address from JSON
		if len(addressJSON) > 0 && string(addressJSON) != "null" {
			if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
				return nil, fmt.Errorf("failed to parse address JSON: %w", err)
			}
		}

		user.ParentID = parentID
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
