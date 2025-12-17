package ports

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type CountryRepository interface {
	List(ctx context.Context, admin_level *int) ([]*domain.Country, error)
	GetByID(ctx context.Context, id int) (*domain.Country, error)
	GetByCode(ctx context.Context, code string, admin_level int) (*domain.Country, error)
}
