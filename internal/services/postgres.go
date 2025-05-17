package services

import (
	"context"
	"go-template/internal/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres() (*pgxpool.Pool, error) {
	// stringConnection := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v",
	// 	common.GetEnv("PGUSER"), common.GetEnv("PGPASSWORD"), common.GetEnv("PGHOST"), common.GetEnv("PGPORT"), common.GetEnv("PGDATABASE"))
	stringConnection := common.GetEnv("DB_DSN")
	config, err := pgxpool.ParseConfig(stringConnection)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
