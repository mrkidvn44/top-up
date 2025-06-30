package db

import (
	"os"
	"top-up-api/config"
	"top-up-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Database *gorm.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	query, err := os.ReadFile("sql/init.sql")
	if err != nil {
		return nil, err
	}
	db.Exec(string(query))

	models := model.GetModels()
	if cfg.Env == "dev" {
		db.Migrator().DropTable(models...)
	}

	db.AutoMigrate(models...)
	if cfg.Env == "dev" {
		files, err := os.ReadDir("sql/data")
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			query, err = os.ReadFile("sql/data/" + file.Name())
			if err != nil {
				return nil, err
			}
			db.Exec(string(query))
		}

	}
	return &DB{Database: db}, nil
}
