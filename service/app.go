package service

import (
	"context"
	"log"
	"server/config"
	"server/internal/board"
	"server/internal/comment"
	"server/internal/user"
	userboardrole "server/internal/user_board_role"
	"server/pkg/adapters/storage"
	"server/pkg/valuecontext"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg            config.Config
	dbConn         *gorm.DB
	authService    *AuthService
	boardService   *BoardService
	commentService *CommentService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()

	app.setAuthService()
	app.setBoardService()
	app.setCommentService()

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
		log.Fatal("Cerate extention failed: ", err)
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
	)
}

func (a *AppContainer) setBoardService() {
	if a.boardService != nil {
		return
	}
	a.boardService = NewBoardService(user.NewOps(storage.NewUserRepo(a.dbConn)), board.NewOps(storage.NewBoardRepo(a.dbConn)), userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)))
}

func (a *AppContainer) setCommentService() {
	if a.commentService != nil {
		return
	}
	a.commentService = NewCommentService(
		comment.NewOps(storage.NewCommentRepo(a.dbConn)),
		userboardrole.NewOps(storage.NewUserBoardRepo(a.dbConn)),
	)
}

func (a *AppContainer) CommentServiceFromCtx(ctx context.Context) *CommentService {
	tx, ok := valuecontext.TryGetTxFromContext(ctx)
	if !ok {
		return a.commentService
	}

	gc, ok := tx.Tx().(*gorm.DB)
	if !ok {
		return a.commentService
	}

	return NewCommentService(
		comment.NewOps(storage.NewCommentRepo(gc)),
		userboardrole.NewOps(storage.NewUserBoardRepo(gc)),
	)
}
