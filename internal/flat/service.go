package flat

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/flat/dto"
	"context"
)

type Service interface {
	Create(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error)
	GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error)
	GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error)
	ChangeStatus(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error)
}

type flatService struct {
	repo Repository
}

// ChangeStatus implements [Service].
func (s *flatService) ChangeStatus(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error) {
	return s.repo.ChangeStatus(ctx, flatID, newStatus)
}

// GetAllByHouseID implements [Service].
func (s *flatService) GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return s.repo.GetAllByHouseID(ctx, houseID)
}

// GetApprovedByHouseID implements [Service].
func (s *flatService) GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return s.repo.GetApprovedByHouseID(ctx, houseID)
}

// Create implements [Service].
func (s *flatService) Create(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error) {

	flat := &entity.Flat{
		HouseID: createdFlat.HouseID,
		Number:  createdFlat.Number,
		Price:   createdFlat.Price,
		Rooms:   createdFlat.Rooms,
	}

	return s.repo.Create(ctx, flat)
}

func NewFlatService(repo Repository) Service {
	return &flatService{
		repo: repo,
	}
}
