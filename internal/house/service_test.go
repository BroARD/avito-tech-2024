package house

import (
	"context"
	"fmt"
	"testing"
	"time"

	"avito-tech/internal/entity"
	"avito-tech/internal/house/dto"
)

type houseRepoMock struct {
	CreateFunc func(ctx context.Context, house *entity.House) (*entity.House, error)
}

func (m *houseRepoMock) Create(ctx context.Context, house *entity.House) (*entity.House, error) {
	return m.CreateFunc(ctx, house)
}

func TestCreateHouse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name        string
		input       dto.CreateHouseInput
		setupMock   func(m *houseRepoMock)
		wantErr     bool
		expectedErr string
		checkResult func(t *testing.T, res *entity.House)
	}{
		{
			name: "Success - Дом создан",
			input: dto.CreateHouseInput{Address: "Test_Address", Year: 1999, Developer: "Someone"},
			setupMock: func(m *houseRepoMock) {
				m.CreateFunc = func(ctx context.Context, house *entity.House) (*entity.House, error) {
					if house.Address != "Test_Address"{
						t.Errorf("В репозиторий ушел неверный адрес: ожидали Address_Test, получили %s", house.Address)
					}
					if house.Year != 1999 {
						t.Errorf("В репозиторий ушел неверный год постройки: ожидали 1999, получили %d", house.Year)
					}
					if house.Developer != "Someone"{
						t.Errorf("В репозиторий ушел неверный адрес: ожидали Someone, получили %s", house.Developer)
					}
					house.ID = 77
					house.CreatedAt = now

					return house, nil
				}
			},
		},
		{
			name:  "Failure - Дом уже существует!",
			input: dto.CreateHouseInput{Address: "Address_Test", Year: 1999, Developer: "Someone"},
			setupMock: func(m *houseRepoMock) {m.CreateFunc = func(ctx context.Context, house *entity.House) (*entity.House, error) {
				return nil, fmt.Errorf("the house already exists")
			}},
			wantErr: true,
			expectedErr: "the house already exists",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &houseRepoMock{}
			if tc.setupMock != nil { tc.setupMock(mockRepo)}
			
			service := NewHouseService(mockRepo)
			result, err := service.Create(context.Background(), &tc.input)

			if tc.wantErr {
				if err == nil || err.Error() != tc.expectedErr {t.Errorf("Ожидалась ошибка %q, а получили ошибка %v", tc.expectedErr, err)}
				return
			}
			if err != nil { t.Fatalf("Не ожидали ошибку: %v", err) }
			if tc.checkResult != nil { tc.checkResult(t, result) }
			}) 
	}
}
