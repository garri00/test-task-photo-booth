package postgresql

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"

	"test-task-photo-booth/src/config"
)

const MigrationFilesDestination = "file://./migrations"

// MigrateUp
func MigrateUp(configs config.PostgresDBConf, l zerolog.Logger) error {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", configs.Username, configs.Password, configs.Host, configs.Port, configs.Database)

	mg, err := migrate.New(MigrationFilesDestination, connectionString)
	if err != nil {
		err = fmt.Errorf("migrate.New() failed: %w", err)
		l.Err(err).Send()

		return err
	}

	if err := mg.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		err = fmt.Errorf("migrate.Up() failed: %w", err)
		l.Err(err).Send()

		return err
	}

	version, _, err := mg.Version()
	if err != nil {
		err = fmt.Errorf("migrate.Up() to version %v failed: %w", version, err)
		l.Err(err).Send()

		return err
	}

	l.Info().Msgf("database migrated to ver=%v", version)

	return nil
}
