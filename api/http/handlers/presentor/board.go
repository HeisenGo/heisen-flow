package presenter

import (
	"server/internal/board"
	userboardrole "server/internal/user_board_role"

	"github.com/google/uuid"
)

type UserBoard struct {
	ID        uuid.UUID `json:"board_id"`
	CreatedAt Timestamp
	Name      string `json:"name"`
	Type      string `json:"type"`
}

func BoardToUserBoard(b board.Board) UserBoard {
	return UserBoard{
		ID:        b.ID,
		CreatedAt: Timestamp(b.CreatedAt),
		Name:      b.Name,
		Type:      b.Type,
	}
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
