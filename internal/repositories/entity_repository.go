package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"n1h41/zolaris-backend-app/internal/domain"
)

type EntityRepository struct {
	db *pgxpool.Pool
}

func NewEntityRepository(dbPool *pgxpool.Pool) EntityRepository {
	return EntityRepository{
		db: dbPool,
	}
}

func (r *EntityRepository) CheckEntityPresence(ctx context.Context, userId string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM z_entity WHERE user_id = $1)`
	err := r.db.QueryRow(ctx, query, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	return exists, nil
}

func (r *EntityRepository) CreateSubEntity(ctx context.Context, categoryId string, user domain.User) error {
	query := `INSERT INTO z_entity (user_id, name, category_id, parent_id) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, user.ID, user.FirstName, categoryId, user.ParentID)
	if err != nil {
		return fmt.Errorf("failed to create sub entity: %w", err)
	}

	return nil
}
