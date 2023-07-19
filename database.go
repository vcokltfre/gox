package gox

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormConfig = gorm.Config{
	Logger: logger.Default.LogMode(logger.Silent),
}

func (g *Gox) UseSQLite(addr string) error {
	if g.DB != nil {
		return fmt.Errorf("database is already configured")
	}

	db, err := gorm.Open(sqlite.Open(addr), &GormConfig)
	if err != nil {
		return err
	}

	g.goxLogger.Info("Using SQLite database.")

	g.DB = db

	return nil
}

func (g *Gox) UsePostgreSQL(addr string) error {
	if g.DB != nil {
		return fmt.Errorf("database is already configured")
	}

	db, err := gorm.Open(postgres.Open(addr), &GormConfig)
	if err != nil {
		return err
	}

	g.goxLogger.Info("Using PostgreSQL database.")

	g.DB = db

	return nil
}

func (g *Gox) UseDatabase(db *gorm.DB) error {
	if g.DB != nil {
		return fmt.Errorf("database is already configured")
	}

	g.goxLogger.Info("Using external database.")

	g.DB = db

	return nil
}

func (g *Gox) AddModel(model any) {
	g.models = append(g.models, model)
}

func (g *Gox) Migrate() error {
	if g.DB == nil {
		return fmt.Errorf("database is not configured")
	}

	g.goxLogger.Info("Migrating database...")

	err := g.DB.AutoMigrate(g.models...)
	if err != nil {
		return err
	}

	g.goxLogger.Info("Database migration complete.")

	return nil
}
