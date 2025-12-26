package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/thealiakbari/task-pool-system/pkg/common/config"
	infraLogger "github.com/thealiakbari/task-pool-system/pkg/common/logger"
	apmpostgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

const UUIDExtension = "uuid-ossp"

var transactionTimeOut time.Duration = 60000

func NewPostgresConn(ctx context.Context, cfg config.Postgres) (*gorm.DB, error) {
	db, err := gorm.Open(apmpostgres.Open(fmt.Sprintf(`postgresql://%s:%s@%s:%d/%s?sslmode=%s&application_name=%s`,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Ssl,
		cfg.AppName)), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newGormLogger(cfg.TraceStacks),
	})
	if err != nil {
		return nil, err
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dbWrapper := DBWrapper{DB: db}
	err = dbWrapper.addExtension([]string{UUIDExtension})
	if err != nil {
		return nil, err
	}
	pdb, err := db.DB()
	if err != nil {
		return nil, err
	}

	pdb.SetConnMaxLifetime(time.Millisecond * cfg.ConnMaxLifetime)
	pdb.SetMaxIdleConns(cfg.MaxIdleConnection)
	pdb.SetMaxOpenConns(cfg.MaxOpenConnection)
	db.WithContext(ctx)

	transactionTimeOut = cfg.TransactionTimeout * time.Millisecond

	return db, nil
}

func (db *DBWrapper) addExtension(names []string) error {
	for _, name := range names {
		if err := db.DB.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%v\";", name)).Error; err != nil {
			return err
		}
	}

	return nil
}

func (db *DBWrapper) AutoMigrate(entities ...interface{}) (err error) {
	for _, model := range entities {
		err = db.DB.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBWrapper) GenerateView(views ...string) (err error) {
	for _, view := range views {
		err = db.DB.Exec(view).Error
		if err != nil {
			return fmt.Errorf("failed to execute SQL query: %w", err)
		}
	}

	return nil
}

func (db *DBWrapper) GenerateFunction(fnCode ...string) (err error) {
	for _, view := range fnCode {
		err = db.DB.Exec(view).Error
		if err != nil {
			return fmt.Errorf("failed to execute SQL query: %w", err)
		}
	}

	return nil
}

// Migrate need `migrations_url` to be set in config file with a relation path, eg. `file:./cmd/migration/scripts`
func Migrate(pgConf config.Postgres, log infraLogger.InfraLogger) (err error) {
	databaseURL := fmt.Sprintf(`postgresql://%s:%s@%s:%d/%s?sslmode=%s&application_name=%s`,
		pgConf.Username,
		pgConf.Password,
		pgConf.Host,
		pgConf.Port,
		pgConf.Name,
		pgConf.Ssl,
		pgConf.AppName,
	)

	if pgConf.MigrationsURL == "" {
		log.Panicf("migration_url cannot be empty, set it with a relation path, eg. `file:./cmd/migration/scripts`\n")
	}

	m, err := migrate.New(pgConf.MigrationsURL, databaseURL)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			err = srcErr
			return
		}

		if dbErr != nil {
			err = dbErr
			return
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		log.Error(err.Error())
		return err
	}

	return nil
}
