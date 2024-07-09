package presenter

import (
	"server/internal/column"

	"github.com/google/uuid"
)

type Column struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	BoardID uuid.UUID `json:"board_id"`
	Order   uint      `json:"order"`
}

func ColumnToResponse(c *column.Column) *Column {
	if c == nil {
		return nil
	}
	return &Column{
		ID:      c.ID,
		Name:    c.Name,
		BoardID: c.BoardID,
		Order:   c.Order,
	}
}

func RequestToColumn(name string, boardID uuid.UUID, order uint) *column.Column {
	return &column.Column{
		ID:      uuid.New(),
		Name:    name,
		BoardID: boardID,
		Order:   order,
	}
}
