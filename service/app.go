package service

import (
	"context"
	"log"
	"server/config"
	"server/internal/board"
	"server/internal/column"
	"server/internal/task"
	"server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage"
	"server/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg           config.Config
	dbConn        *gorm.DB
	authService   *AuthService
	boardService  *BoardService
	taskService   *TaskService
	columnService *ColumnService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()

	app.setAuthService()
	app.setBoardService()
	app.setTaskService()

	app.setColumnService()

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

	err = storage.AddExtension(a.dbConn)
	if err != nil {
		log.Fatal("Create extension failed: ", err)
	}

	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
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
		column.NewOps(storage.NewColumnRepo(gc)),
	)
}

func (a *AppContainer) ColumnService() *ColumnService {
	return a.columnService
}

func (a *AppContainer) ColumnServiceFromCtx(ctx context.Context) *ColumnService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.columnService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.columnService
	}

	return NewColumnService(
		column.NewOps(storage.NewColumnRepo(gc)),
		userboardrole.NewOps(storage.NewUserBoardRepo(gc)),
		board.NewOps(storage.NewBoardRepo(gc)),
	)
}

func (a *AppContainer) setBoardService() {
	if a.boardService != nil {
		return
	}
	a.boardService = NewBoardService(user.NewOps(storage.NewUserRepo(a.dbConn)), board.NewOps(storage.NewBoardRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)), column.NewOps(storage.NewColumnRepo(a.dbConn)))
}

func (a *AppContainer) setColumnService() {
	if a.columnService != nil { // Fixed condition
		return
	}
	a.columnService = NewColumnService(column.NewOps(storage.NewColumnRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)),
		board.NewOps(storage.NewBoardRepo(a.dbConn)))
}

func (a *AppContainer) TaskService() *TaskService {
	return a.taskService
}

func (a *AppContainer) TaskServiceFromCtx(ctx context.Context) *TaskService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.taskService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.taskService
	}

	return NewTaskService(
		user.NewOps(storage.NewUserRepo(gc)),
		board.NewOps(storage.NewBoardRepo(gc)),
		userboardrole.NewOps(storage.NewUserBoardRepo(gc)),
		task.NewOps(storage.NewTaskRepo(gc)),
		column.NewOps(storage.NewColumnRepo(gc)),
	)
}

func (a *AppContainer) setTaskService() {
	if a.taskService != nil {
		return
	}
	a.taskService = NewTaskService(user.NewOps(storage.NewUserRepo(a.dbConn)), board.NewOps(storage.NewBoardRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)), task.NewOps(storage.NewTaskRepo(a.dbConn)), column.NewOps(storage.NewColumnRepo(a.dbConn)))
}
