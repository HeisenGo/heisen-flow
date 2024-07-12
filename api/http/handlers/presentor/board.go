package presenter

import (
	"server/internal/board"
	"server/internal/column"
	"server/internal/task"
	"server/internal/user"
	userBoardRole "server/internal/user_board_role"
	"server/pkg/fp"
	"time"

	"github.com/google/uuid"
)

type UserBoard struct {
	ID        uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type BoardUserResp struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}
type BoardColumnResp struct {
	ID    uuid.UUID       `json:"id"`
	Name  string          `json:"name"`
	Order uint            `json:"order"`
	Tasks []BoardTaskResp `json:"tasks"`
}
type BoardTaskResp struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
	StoryPoint uint      `json:"story_at"`
}
type FullBoardResp struct {
	ID        uuid.UUID         `json:"board_id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	CreatedAt time.Time         `json:"created_at"`
	Users     []BoardUserResp   `json:"users"`
	Columns   []BoardColumnResp `json:"columns"`
}

func userToBoardUserResp(u user.User) BoardUserResp {
	return BoardUserResp{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role.String(),
	}
}
func TaskToBoardTaskResp(t task.Task) BoardTaskResp {
	return BoardTaskResp{
		ID:         t.ID,
		Title:      t.Title,
		StartAt:    t.StartAt,
		EndAt:      t.EndAt,
		StoryPoint: t.StoryPoint,
	}
}

func columnToBoardColumnResp(c column.Column) BoardColumnResp {
	tasksResp := BatchTaskToBoardTaskResp(c.Tasks)
	return BoardColumnResp{
		ID:    c.ID,
		Name:  c.Name,
		Order: c.OrderNum,
		Tasks: tasksResp,
	}
}

func BatchTaskToBoardTaskResp(t []task.Task) []BoardTaskResp {
	return fp.Map(t, TaskToBoardTaskResp)
}
func BatchUserToBoardUserResp(u []user.User) []BoardUserResp {
	return fp.Map(u, userToBoardUserResp)
}
func BatchColumnToBoardUserResp(c []column.Column) []BoardColumnResp {
	return fp.Map(c, columnToBoardColumnResp)
}

func BoardToFullBoardResp(b board.Board) FullBoardResp {
	usersResp := BatchUserToBoardUserResp(b.Users)
	columnsResp := BatchColumnToBoardUserResp(b.Columns)
	return FullBoardResp{
		ID:        b.ID,
		Name:      b.Name,
		Type:      b.Type,
		CreatedAt: b.CreatedAt,
		Users:     usersResp,
		Columns:   columnsResp,
	}
}

func boardToUserBoard(b board.Board) UserBoard {
	return UserBoard{
		ID:        b.ID,
		Name:      b.Name,
		Type:      b.Type,
		CreatedAt: b.CreatedAt,
	}
}

func BatchBoardsToUserBoard(boards []board.Board) []UserBoard {
	return fp.Map(boards, boardToUserBoard)
}

func UserBoardToBoard(userBoard *UserBoard, userID uuid.UUID) (*board.Board, *userBoardRole.UserBoardRole) {
	b := &board.Board{
		Name: userBoard.Name,
		Type: userBoard.Type,
	}
	ubr := &userBoardRole.UserBoardRole{
		UserID: userID,
	}
	return b, ubr
}

type InviteUserToBoard struct {
	ID      uuid.UUID `json:"user_board_role_id"`
	Email   string    `json:"email"`
	BoardID uuid.UUID `json:"board_id"`
	Role    string    `json:"role"`
}

func InviteUserToBoardToUserBoardRole(inviteUserToBoard *InviteUserToBoard) *userBoardRole.UserBoardRole {
	return &userBoardRole.UserBoardRole{
		Role:    inviteUserToBoard.Role,
		BoardID: inviteUserToBoard.BoardID,
	}
}
