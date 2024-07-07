package service

import (
	"log"
	"server/config"
	"server/internal/user"
	"server/pkg/adapters/storage"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg         config.Config
	dbConn      *gorm.DB
	authService *AuthService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	err := storage.Migrate(app.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	app.setAuthService()

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

func (a *AppContainer) AuthService() *AuthService {
	return a.authService
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = NewAuthService(user.NewOps(storage.NewUserRepo(a.dbConn)), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}
