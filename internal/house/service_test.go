package house


import (
	"context"


	"avito-tech/internal/house/dto"
	"avito-tech/internal/entity"


)


type houseRepoMock struct {
	CreateFunc func(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error)
}

func (m *houseRepoMock) Register(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error)   { return m.CreateFunc(ctx, house) }


