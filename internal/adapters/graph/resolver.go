package graph

import "github.com/hoshina-dev/gapi/internal/core/ports"

//go:generate go tool gqlgen generate

type Resolver struct {
	countryService ports.CountryService
}

func NewResolver(countryService ports.CountryService) *Resolver {
	return &Resolver{countryService: countryService}
}
