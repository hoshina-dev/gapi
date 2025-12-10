package ports

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type CountryRepository interface {
	List(ctx context.Context) ([]domain.Country, error)
	GetByID(ctx context.Context, id int) (*domain.Country, error)
}
