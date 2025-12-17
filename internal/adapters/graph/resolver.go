package graph

import "github.com/hoshina-dev/gapi/internal/core/ports"

//go:generate go tool gqlgen generate

type Resolver struct {
	adminAreaService ports.AdminAreaService
}

func NewResolver(adminAreaService ports.AdminAreaService) *Resolver {
	return &Resolver{adminAreaService: adminAreaService}
}
