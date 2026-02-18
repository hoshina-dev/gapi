package graph

import "github.com/hoshina-dev/gapi/internal/core/ports"

//go:generate go tool gqlgen generate

type Resolver struct {
	adminAreaService ports.AdminAreaService
	osmLineService   ports.OSMLineService
}

func NewResolver(adminAreaService ports.AdminAreaService, osmLineService ports.OSMLineService) *Resolver {
	return &Resolver{
		adminAreaService: adminAreaService,
		osmLineService:   osmLineService,
	}
}
