package service

import (
	"context"
	"log"
	"server/config"
	"server/internal/board"
	"server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage"
	"server/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg          config.Config
	dbConn       *gorm.DB
	authService  *AuthService
	boardService *BoardService
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
	app.setBoardService()

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

func (a *AppContainer) BoardService() *BoardService {
	return a.boardService
}

func (a *AppContainer) BoardServiceFromCtx(ctx context.Context) *BoardService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.boardService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.boardService
	}

	return NewBoardService(
		user.NewOps(storage.NewUserRepo(gc)),
		board.NewOps(storage.NewBoardRepo(gc)),
		userboardrole.NewOps(storage.NewUserBoardRepo(gc)),
	)
}

func (a *AppContainer) setBoardService() {
	if a.boardService != nil {
		return
	}
	a.boardService = NewBoardService(user.NewOps(storage.NewUserRepo(a.dbConn)), board.NewOps(storage.NewBoardRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)))
}
