package postgres

import (
	"context"
	"fmt"
	"go-template/internal/domain"
	"go-template/internal/infra/postgres/sqlc"
	"go-template/pkg/optional"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: sqlc.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:         uuidToPg(u.ID),
		Name:       u.Name,
		Surname:    optionalStringToPgText(u.Surname),
		TelegramID: u.TelegramID,
		CreatedAt:  timeToPgTimestamptz(u.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("create user: %w", mapError(err, "user"))
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, uuidToPg(id))
	if err != nil {
		return domain.User{}, fmt.Errorf("get user %s: %w", id, mapError(err, "user"))
	}
	return hydrateUser(row), nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	rows, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = hydrateUser(row)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:         uuidToPg(u.ID),
		Name:       u.Name,
		Surname:    optionalStringToPgText(u.Surname),
		TelegramID: u.TelegramID,
	})
	if err != nil {
		return fmt.Errorf("update user %s: %w", u.ID, mapError(err, "user"))
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUser(ctx, uuidToPg(id))
	if err != nil {
		return fmt.Errorf("delete user %s: %w", id, mapError(err, "user"))
	}
	return nil
}

// hydrateUser creates a domain.User from a DB row.
// Validation is skipped — data in the database is considered valid
// because it passed through the domain constructor on write.
func hydrateUser(row sqlc.User) domain.User {
	u := domain.User{
		BaseEntity: domain.BaseEntity{
			ID:        pgToUUID(row.ID),
			CreatedAt: row.CreatedAt.Time,
		},
		Name:       row.Name,
		TelegramID: row.TelegramID,
	}
	if row.Surname.Valid {
		u.Surname = optional.Some(row.Surname.String)
	}
	return u
}

func optionalStringToPgText(o optional.Optional[string]) pgtype.Text {
	if v, ok := o.Get(); ok {
		return pgtype.Text{String: v, Valid: true}
	}
	return pgtype.Text{}
}
