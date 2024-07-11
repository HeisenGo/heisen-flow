package presenter

import (
	"server/internal/board"
	userboardrole "server/internal/user_board_role"
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

func UserBoardToBoard(userBoard *UserBoard, userID uuid.UUID) (*board.Board, *userboardrole.UserBoardRole) {
	b := &board.Board{
		Name: userBoard.Name,
		Type: userBoard.Type,
	}
	ubr := &userboardrole.UserBoardRole{
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

func InviteUserToBoardToUserBoardRole(inviteUserToBoard *InviteUserToBoard) *userboardrole.UserBoardRole {
	return &userboardrole.UserBoardRole{
		Role:    inviteUserToBoard.Role,
		BoardID: inviteUserToBoard.BoardID,
	}
}
