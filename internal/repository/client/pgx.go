package client

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

type PgxClientRepository struct {
	dbpool *pgxpool.Pool
}

func NewPgxClientRepository(dbpool *pgxpool.Pool) *PgxClientRepository {
	return &PgxClientRepository{
		dbpool: dbpool,
	}
}

func (repo *PgxClientRepository) ListClientsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Client, error) {
	clients := make([]domain.Client, 0)
	errfmt := "%w: %s"
	query := `SELECT id, user_id, scope, secret_hash FROM clients WHERE user_id=$1`
	rows, err := repo.dbpool.Query(ctx, query, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, fmt.Errorf("%w: %s", domain.ErrGeneralClient, pgErr.Error())
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralClient, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var client domain.Client
		if err := rows.Scan(&client.ID, &client.UserID, &client.Scope, &client.SecretHash); err != nil {
			return nil, fmt.Errorf(errfmt, domain.ErrGeneralClient, err.Error())
		}
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralClient, err.Error())
	}

	return clients, nil

}

func (repo *PgxClientRepository) CreateClient(ctx context.Context, ID uuid.UUID, userID uuid.UUID, scope []domain.Permission, secretHash []byte) (*domain.Client, error) {
	errfmt := "%w: %s"
	query := `INSERT INTO clients (id, user_id, scope, secret_hash) VALUES ($1, $2, $3, $4) RETURNING id, user_id, scope, secret_hash`

	var client domain.Client
	err := repo.dbpool.QueryRow(ctx, query, ID.String(), userID.String(), scope, secretHash).
		Scan(&client.ID, &client.UserID, &client.Scope, &client.SecretHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				// unique_violation
				return nil, fmt.Errorf(errfmt, domain.ErrClientAlreadyExists, pgErr.Error())
			case "23514":
				// check_violation (e.g. constraints)
				return nil, fmt.Errorf(errfmt, domain.ErrInvalidClientData, pgErr.Error())
			}
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralClient, err.Error())
	}

	return &client, nil
}

func (repo *PgxClientRepository) GetClientByID(ctx context.Context, ID uuid.UUID) (*domain.Client, error) {
	errfmt := "%w: %s"
	query := `SELECT id, user_id, scope, secret_hash FROM clients WHERE id = $1`

	var client domain.Client
	err := repo.dbpool.QueryRow(ctx, query, ID.String()).
		Scan(&client.ID, &client.UserID, &client.Scope, &client.SecretHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(errfmt, domain.ErrUserNotFound, err.Error())
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err.Error())
	}
	return &client, nil
}

func (repo *PgxClientRepository) UpdateClientByIDandUser(
	ctx context.Context,
	ID uuid.UUID,
	userId uuid.UUID,
	scope []domain.Permission,
	secretHash []byte,
) (*domain.Client, error) {
	errfmt := "%w: %s"

	query := `UPDATE clients SET `
	args := []interface{}{}
	argIndex := 1
	updates := []string{}

	if scope != nil {
		updates = append(updates, fmt.Sprintf("scope = $%d", argIndex))
		args = append(args, scope)
		argIndex++
	}

	if secretHash != nil {
		updates = append(updates, fmt.Sprintf("secret_hash = $%d", argIndex))
		args = append(args, secretHash)
		argIndex++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, "no fields to update")
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d RETURNING id, user_id, scope, secret_hash", argIndex, argIndex+1)
	args = append(args, ID, userId)

	var client domain.Client
	err := repo.dbpool.QueryRow(ctx, query, args...).Scan(
		&client.ID, &client.UserID, &client.Scope, &client.SecretHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(errfmt, domain.ErrUserNotFound, err.Error())
		}
		return nil, fmt.Errorf(errfmt, domain.ErrGeneralUser, err.Error())
	}

	return &client, nil
}

func (repo *PgxClientRepository) DeleteClientByIDandUser(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1 and userID = $2`
	cmdTag, err := repo.dbpool.Exec(ctx, query, ID, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrGeneralUser, err.Error())
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
