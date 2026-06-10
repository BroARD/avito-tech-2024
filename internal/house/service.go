package house

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/house/dto"
	"context"
)

type Service interface {
	Create(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error)
}

type houseService struct {
	repo Repository
}

// Create implements [Service].
func (s *houseService) Create(ctx context.Context, createdHouse *dto.CreateHouseInput) (*entity.House, error) {
	house := &entity.House{
		Address: createdHouse.Address,
		Year: createdHouse.Year,
		Developer: createdHouse.Developer,
	}

	return s.repo.Create(ctx, house)
}

func NewHouseService(repo Repository) Service {
	return &houseService{
		repo: repo,
	}
}
