package flat // Меняем на имя тестируемого пакета

import (
	"context"
	"testing"

	"avito-tech/internal/entity"
	"avito-tech/internal/flat/dto"
)

// 1. Наш мок-репозиторий. Теперь методы СТРОГО вызывают свои замыкания (функции).
type flatRepoMock struct {
	CreateFunc               func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error)
	GetAllByHouseIDFunc      func(ctx context.Context, houseID uint) ([]entity.Flat, error)
	GetApprovedByHouseIDFunc func(ctx context.Context, houseID uint) ([]entity.Flat, error)
}

// Реализуем метод интерфейса Repository
func (m *flatRepoMock) Create(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
	return m.CreateFunc(ctx, flat)
}

// ИСПРАВЛЕНО: Метод теперь реально вызывает GetAllByHouseIDFunc, а не возвращает nil
func (m *flatRepoMock) GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetAllByHouseIDFunc(ctx, houseID)
}

// Заглушки для методов, которые в данных сценариях не участвуют
func (m *flatRepoMock) GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetApprovedByHouseIDFunc(ctx, houseID)
}
func (m *flatRepoMock) ChangeStatus(ctx context.Context, id uint, status string) (*entity.Flat, error) {
	return nil, nil
}

// 2. Тест создания квартиры
func TestCreateFlat_Success(t *testing.T) {
	mockFlatRepo := &flatRepoMock{
		CreateFunc: func(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
			flat.ID = 7
			flat.Status = "created"
			return flat, nil
		},
	}

	// Так как мы находимся внутри пакета flat, вызываем функцию напрямую
	service := NewFlatService(mockFlatRepo)

	input := dto.CreateFlatInput{HouseID: 1, Number: 10, Price: 5000000, Rooms: 2}
	result, err := service.Create(context.Background(), &input)

	if err != nil {
		t.Fatalf("Ожидался успешный результат, но получена ошибка: %v", err)
	}

	if result.ID != 7 {
		t.Errorf("Ожидался ID квартиры = 7, получили = %d", result.ID)
	}

	if result.Status != "created" {
		t.Errorf("Ожидался статус 'created', получили = %s", result.Status)
	}
}

// 3. Тест получения списка квартир модератором
func TestGetApprovedFlatByHouseID_Success(t *testing.T) {
	targetHouseID := uint(3)

	mockFlatRepo := &flatRepoMock{
		GetAllByHouseIDFunc: func(ctx context.Context, houseID uint) ([]entity.Flat, error) {
			return []entity.Flat{
				{ID: 1, HouseID: houseID, Number: 10, Status: "approved"},
				{ID: 2, HouseID: houseID, Number: 11, Status: "created"},
			}, nil
		},
	}

	service := NewFlatService(mockFlatRepo)

	result, err := service.GetAllByHouseID(context.Background(), targetHouseID)

	if err != nil {
		t.Fatalf("Ожидался успешный результат, но получена ошибка: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Ожидалось получить 2 квартиры, получили = %d", len(result))
	}

	if result[0].HouseID != targetHouseID {
		t.Errorf("Ожидался ID дома = 3, получили = %d", result[0].HouseID)
	}
}

