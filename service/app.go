package service

import (
	"log"
	"server/config"
	"server/pkg/adapters/storage"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg    config.Config
	dbConn *gorm.DB
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage.Migrate(app.dbConn)

	return app, nil
}

func (a *AppContainer) RawDBConnection() *gorm.DB {
	return a.dbConn
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db
}
