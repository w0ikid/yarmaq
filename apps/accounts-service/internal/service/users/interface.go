package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type UsersRepo interface {
	Create(ctx context.Context, user models.User) (*models.User, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByZitadelID(ctx context.Context, zitadelID string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)

	Update(ctx context.Context, user models.User) (*models.User, error)
	UpdateRoles(ctx context.Context, zitadelID string, roles []string) (*models.User, error)

	Delete(ctx context.Context, id uuid.UUID) error
}
