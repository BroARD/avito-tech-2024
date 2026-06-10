package entity

import (
	"time"
)

type House struct {
	ID            uint       `json:"id" db:"id"`                                     // уникальный номер дома
	Address       string     `json:"address" db:"address"`                           // адрес
	Year          int        `json:"year" db:"year"`                                 // год постройки
	Developer     *string    `json:"developer,omitempty" db:"developer"`             // застройщик (указатель, чтобы поддерживать NULL)
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`                     // дата создания дома
	LastFlatAdded *time.Time `json:"last_flat_added,omitempty" db:"last_flat_added"` // дата добавления последней квартиры
}
