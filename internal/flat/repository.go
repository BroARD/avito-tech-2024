package flat

import (
	"avito-tech/internal/entity"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	Create(ctx context.Context, house *entity.Flat) (*entity.Flat, error)
	GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error)
	GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error)
	ChangeStatus(ctx context.Context, id uint, newStatus string) (*entity.Flat, error)
}

type flatRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

// ChangeStatus implements [Repository].
func (r *flatRepository) ChangeStatus(ctx context.Context, id uint, newStatus string) (*entity.Flat, error) {
	cb := r.builder.
		Update("flats").
		Set("status", newStatus).
		Where(sq.Eq{"id": id})

	switch newStatus {
	case "on moderation":
		cb = cb.Where(sq.Eq{"status": []string{"created", "declined"}})
	case "approved", "declined":
		cb = cb.Where(sq.Eq{"status": "on moderation"})
	}

	query, args, err := cb.Suffix("RETURNING id, house_id, number, price, rooms, status").ToSql()
	if err != nil {
		return nil, err
	}

	var f entity.Flat
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&f.ID, &f.HouseID, &f.Number, &f.Price, &f.Rooms, &f.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &f, nil
}


// GetAllByHouseID implements [Repository].
func (r *flatRepository) GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	query, args, err := r.builder.
		Select("id", "house_id", "number", "price", "rooms", "status").
		From("flats").
		Where(sq.Eq{"house_id": houseID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flats []entity.Flat

	for rows.Next() {
		var f entity.Flat
		err := rows.Scan(&f.ID, &f.HouseID, &f.Number, &f.Price, &f.Rooms, &f.Status)
		if err != nil {
			return nil, err
		}

		flats = append(flats, f)
	}

	return flats, nil
}

// GetApprovedByHouseID implements [Repository].
func (r *flatRepository) GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	query, args, err := r.builder.
		Select("id", "house_id", "number", "price", "rooms", "status").
		From("flats").
		Where(sq.Eq{"house_id": houseID, "status": "approved"}).
		ToSql()
	if err != nil {
		return nil, err
	}


	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flats []entity.Flat

	for rows.Next() {
		var f entity.Flat
		err := rows.Scan(&f.ID, &f.HouseID, &f.Number, &f.Price, &f.Rooms, &f.Status)
		if err != nil {
			return nil, err
		}

		flats = append(flats, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return flats, nil
}

// Create implements [Repository].
func (r *flatRepository) Create(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
	query, args, err := r.builder.
		Insert("flats").
		Columns("house_id", "number", "price", "rooms").
		Values(flat.HouseID, flat.Number, flat.Price, flat.Rooms).
		Suffix("RETURNING id, status").
		ToSql()
	if err != nil {
		return nil, err
	}

	var id uint
	var status string

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &status)

	if err != nil {
		return nil, err
	}

	flat.ID = id
	flat.Status = status

	return flat, nil
}

func NewFlatRepository(db *sql.DB) Repository {
	return &flatRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
