package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Name    string      `db:"name" json:"name"`
	Id      pgtype.UUID `db:"id" json:"id"`
	Created time.Time   `db:"created" json:"date_created"`
}

type RepositoryModel struct {
	Pool *pgxpool.Pool
}

func (r *RepositoryModel) CreateRepository(repositoryName string) (uuid.UUID, error) {

	id := uuid.New()

	query := `INSERT INTO repository (name, id) VALUES (@name, @id)`
	args := pgx.NamedArgs{
		"name": repositoryName,
		"id":   id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	commandTag, err := r.Pool.Exec(ctx, query, args)
	if err != nil {
		return uuid.Nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return uuid.Nil, ErrRecordNotFound
	}

	return id, nil
}

func (r *RepositoryModel) GetRepositoryById(repoId uuid.UUID) (Repository, error) {
	query := `SELECT * FROM repository where id = @id`
	args := pgx.NamedArgs{
		"id": repoId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.Pool.Query(ctx, query, args)
	if err != nil {
		return Repository{}, ErrQueryError
	}
	repository, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Repository])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Repository{}, ErrRecordNotFound
		default:
			return Repository{}, err
		}
	}

	return repository, nil
}

func (r *RepositoryModel) CheckIfRepositoryExist(repoId uuid.UUID) bool {
	query := `SELECT id FROM repository where id = @id`
	args := pgx.NamedArgs{
		"id": repoId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.Pool.Query(ctx, query, args)
	if err != nil {
		return false
	}
	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return false
		default:
			return false
		}
	}

	return true

}
