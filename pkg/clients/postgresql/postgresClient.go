package postgresql

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"test-task-photo-booth/src/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, configs config.PostgresDBConf, l *zerolog.Logger) (db *pgxpool.Pool, err error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		configs.Username,
		configs.Password,
		configs.Host,
		configs.Port,
		configs.Database,
		configs.SSLMode)

	dbPool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		err = fmt.Errorf("pgxpool.New failed: %w", err)
		l.Err(err).Send()

		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		err = fmt.Errorf("dbPool.Ping failed: %w", err)
		l.Err(err).Send()

		return nil, err
	}

	log.Info().Msg("successfully connected to PostgresDB")

	return dbPool, nil
}
