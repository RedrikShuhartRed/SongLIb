package db

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/RedrikShuhartRed/EfMobSongLib/config"
)

type Database struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDatabase(logger *zap.Logger, cfg *config.Config) (*Database, error) {

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s", cfg.LibraryUser, cfg.LibraryPassword, cfg.LibraryHost, cfg.LibraryDBPort, cfg.LibrarySSLMode)
	fmt.Println(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return nil, err
	}
	query := fmt.Sprintf("CREATE DATABASE %s", cfg.LibraryName)
	existsErr := fmt.Sprintf("pq: database \"%s\" already exists", cfg.LibraryName)
	_, err = db.Exec(query)
	if err != nil && err.Error() != existsErr {
		logger.Error("error creat database", zap.Error(err))
		return nil, err
	}
	defer db.Close()

	connStr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", cfg.LibraryUser, cfg.LibraryPassword, cfg.LibraryHost, cfg.LibraryDBPort, cfg.LibraryName, cfg.LibrarySSLMode)
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		logger.Error("errror conntcting to database", zap.Error(err))
		return nil, err
	}
	return &Database{db: db, logger: logger}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}

func (d *Database) RunMigrate() error {

	sqlDB := d.db.DB
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://./db/migrations", "postgres", driver)
	if err != nil {
		d.logger.Error("error create migrations:", zap.Error(err))
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		d.logger.Error("error migrating database", zap.Error(err))
		return err
	}
	return nil
}
