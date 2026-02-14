package mocks

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockAdminAreaService struct {
	mock.Mock
}

func (m *MockAdminAreaService) GetAll(ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	args := m.Called(ctx, adminLevel, tolerance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AdminArea), args.Error(1)
}

func (m *MockAdminAreaService) GetByID(ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	args := m.Called(ctx, id, adminLevel, tolerance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AdminArea), args.Error(1)
}

func (m *MockAdminAreaService) GetByCode(ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	args := m.Called(ctx, code, adminLevel, tolerance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AdminArea), args.Error(1)
}

func (m *MockAdminAreaService) GetChildren(ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	args := m.Called(ctx, parentCode, childLevel, tolerance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AdminArea), args.Error(1)
}

func (m *MockAdminAreaService) FilterCoordinatesByBoundary(ctx context.Context, coordinates []*domain.Coordinate, boundaryID string, adminLevel int32) ([]*domain.Coordinate, error) {
	args := m.Called(ctx, coordinates, boundaryID, adminLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Coordinate), args.Error(1)
}
