package ports

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type AdminAreaRepository interface {
	List(ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error)
	GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error)
	GetByCode(ctx context.Context, code string, adminLevel int32) (*domain.AdminArea, error)
	GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error)
}
