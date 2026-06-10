package auth

import (
	"avito-tech/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type authRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

// Login implements [Repository].
func (r *authRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query, args, err := r.builder.
		Select("id", "email", "password_hash", "role", "created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u entity.User
	
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &u, nil
}


// Register implements [Repository].
func (r *authRepository) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	query, args, err := r.builder.Insert("users").Columns("email", "password_hash", "role").
		Values(user.Email, user.PasswordHash, user.Role).Suffix("RETURNING id, created_at").ToSql()

	if err != nil {
		return nil, err
	}

	var id uint
	var createdAt time.Time

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	user.ID = id
	user.CreatedAt = createdAt

	return user, nil
}

func NewAuthRepository(db *sql.DB) Repository {
	return &authRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
