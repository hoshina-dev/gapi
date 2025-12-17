package ports

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type AdminAreaService interface {
	GetAll(ctx context.Context, admin_level *int32) ([]*domain.AdminArea, error)
	GetByID(ctx context.Context, id int) (*domain.AdminArea, error)
	GetByCode(ctx context.Context, code string, admin_level int32) (*domain.AdminArea, error)
}
