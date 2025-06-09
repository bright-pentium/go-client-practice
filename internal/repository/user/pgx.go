package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPgxUserRepository(dbpool *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{
		dbpool: dbpool,
	}
}

func (repo *PgxUserRepository) CreateUser(ctx context.Context, ID uuid.UUID, name string, account string, passwordHash []byte) (*domain.User, error) {
	var user domain.User
	errfmt := "%w: %s"
	query := `INSERT INTO users (id, name, account, password_hash) VALUES ($1, $2, $3, $4) RETURNING id, name, account, password_hash`
	err := repo.dbpool.QueryRow(ctx, query, ID, name, account, passwordHash).Scan(&user.ID, &user.Name, &user.Account, &user.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				// unique_violation (e.g., ID or name already exists)
				return nil, fmt.Errorf(errfmt, domain.ErrUserAlreadyExists, pgErr.Error())
			case "23514":
				// check_violation (e.g., name too short, invalid input)
				return nil, fmt.Errorf(errfmt, domain.ErrInvalidUserData, pgErr.Error())
			}
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err.Error())
	}
	return &user, nil
}

func (repo *PgxUserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*domain.User, error) {
	var user domain.User
	errfmt := "%w: %s"
	query := `SELECT id, name, account, password_hash FROM users WHERE id = $1`
	err := repo.dbpool.QueryRow(ctx, query, ID).Scan(&user.ID, &user.Name, &user.Account, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(errfmt, domain.ErrUserNotFound, err.Error())
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err.Error())
	}
	return &user, nil
}

func (repo *PgxUserRepository) GetUserByAccount(ctx context.Context, account string) (*domain.User, error) {
	var user domain.User
	errfmt := "%w: %s"
	query := `SELECT id, name, account, password_hash FROM users WHERE account = $1`
	err := repo.dbpool.QueryRow(ctx, query, account).Scan(&user.ID, &user.Name, &user.Account, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(errfmt, domain.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err)
	}
	return &user, nil
}

func (repo *PgxUserRepository) UpdateUserByID(
	ctx context.Context,
	ID uuid.UUID,
	name string,
	passwordHash []byte,
) (*domain.User, error) {
	var user domain.User
	errfmt := "%w: %s"

	// Start building query and args
	query := `UPDATE users SET `
	args := []interface{}{}
	argIndex := 1
	updates := []string{}

	if name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if passwordHash != nil {
		updates = append(updates, fmt.Sprintf("password_hash = $%d", argIndex))
		args = append(args, passwordHash)
		argIndex++
	}

	if len(updates) == 0 {
		// Nothing to update
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, "no fields to update")
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, account, password_hash", argIndex)
	args = append(args, ID)

	err := repo.dbpool.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Account, &user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(errfmt, domain.ErrUserNotFound, err.Error())
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err.Error())
	}

	return &user, nil
}

func (repo *PgxUserRepository) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := repo.dbpool.Exec(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrGeneralUser, err.Error())
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
