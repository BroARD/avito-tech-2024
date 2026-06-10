package house

import (
	"avito-tech/internal/entity"
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	Create(ctx context.Context, house *entity.House) (*entity.House, error)
}

type houseRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}


// Create implements [Repository].
func (r *houseRepository) Create(ctx context.Context, house *entity.House) (*entity.House, error) {
	query, args, err := r.builder.Insert("houses").Columns("address", "year", "developer").
		Values(house.Address, house.Year, house.Developer).Suffix("RETURNING id, created_at").ToSql()

	if err != nil {
		return nil, err
	}

	var id uint
	var createdAt time.Time

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	house.ID = id
	house.CreatedAt = createdAt

	return house, nil
}

func NewHouseRepository(db *sql.DB) Repository {
	return &houseRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
