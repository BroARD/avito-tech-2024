package flat // Меняем на имя тестируемого пакета

import (
	"context"
	"fmt"
	"testing"

	"avito-tech/internal/entity"
	"avito-tech/internal/flat/dto"
)

type flatRepoMock struct {
	CreateFunc               func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error)
	GetAllByHouseIDFunc      func(ctx context.Context, houseID uint) ([]entity.Flat, error)
	GetApprovedByHouseIDFunc func(ctx context.Context, houseID uint) ([]entity.Flat, error)
	ChangeStatusFunc         func(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error)
}

func (m *flatRepoMock) Create(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
	return m.CreateFunc(ctx, flat)
}

func (m *flatRepoMock) GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetAllByHouseIDFunc(ctx, houseID)
}

func (m *flatRepoMock) GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetApprovedByHouseIDFunc(ctx, houseID)
}

func (m *flatRepoMock) ChangeStatus(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error) {
	return m.ChangeStatusFunc(ctx, flatID, newStatus)
}

func TestCreateFlat(t *testing.T) {
	tests := []struct {
		name        string
		input       dto.CreateFlatInput
		setupMock   func(m *flatRepoMock)
		wantErr     bool
		expectedErr string
		checkResult func(t *testing.T, res *entity.Flat)
	}{
		{
			name:  "Success - Квартира успешно создана!",
			input: dto.CreateFlatInput{HouseID: 1, Number: 66, Price: 1000000, Rooms: 4},
			setupMock: func(m *flatRepoMock) {
				m.CreateFunc = func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
					if flat.Number != 66 {
						t.Errorf("В репозиторий ушел неверный номер квартиры: ожидали 66, получили %d", flat.Number)
					}
					if flat.Price != 1000000 {
						t.Errorf("В репозиторий ушла неверная цена")
					}
					flat.ID = 7
					flat.Status = "Created"
					return flat, nil
				}
			},
			checkResult: func(t *testing.T, res *entity.Flat) {
				if res.ID != 7 {
					t.Errorf("Ожидался ID 7, получили %d", res.ID)
				}
				if res.Status != "Created" {
					t.Errorf("Ожидался статус Created, получили %s", res.Status)
				}
			},
		},
		{
			name:  "Failure - Дома не существует!",
			input: dto.CreateFlatInput{HouseID: 13, Number: 66, Price: 1000000, Rooms: 4},
			setupMock: func(m *flatRepoMock) {
				m.CreateFunc = func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
					return nil, fmt.Errorf("house is not exists")
				}
			},
			wantErr:     true,
			expectedErr: "house is not exists",
		},
		{
			name:  "Failure - Квартира уже существует!",
			input: dto.CreateFlatInput{HouseID: 1, Number: 66, Price: 1000000, Rooms: 4},
			setupMock: func(m *flatRepoMock) {
				m.CreateFunc = func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
					return nil, fmt.Errorf("the flat already exists")
				}
			},
			wantErr:     true,
			expectedErr: "the flat already exists",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &flatRepoMock{}
			if tc.setupMock != nil {
				tc.setupMock(mockRepo)
			}

			service := NewFlatService(mockRepo)
			result, err := service.Create(context.Background(), &tc.input)

			if tc.wantErr {
				if err == nil || err.Error() != tc.expectedErr {
					t.Errorf("Ожидалась ошибка %q, получили %v", tc.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Не ожидали ошибку: %v", err)
			}
			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}
