package presenter

import (
	"server/internal/column"
	"server/pkg/adapters/storage/entities"

	"github.com/google/uuid"
)

type CreateColumnsRequest struct {
	BoardID uuid.UUID          `json:"board_id"`
	Columns []CreateColumnItem `json:"columns"`
}

type CreateColumnItem struct {
	Name string `json:"name"`
}

type CreateColumnsResponse struct {
	Data    []ColumnResponseItem `json:"data"`
	Message string               `json:"message"`
}

type ColumnResponseItem struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Order uint      `json:"order"`
}

type GetColumnsResponse struct {
	Columns []ColumnResponseItem `json:"columns"`
	Message string               `json:"message"`
}

func CreateColumnsRequestToEntities(req CreateColumnsRequest, maxOrder uint) []entities.Column {
	columns := make([]entities.Column, len(req.Columns))
	for i, col := range req.Columns {
		columns[i] = entities.Column{
			ID:      uuid.New(),
			Name:    col.Name,
			BoardID: req.BoardID,
			Order:   maxOrder + uint(i) + 1,
		}
	}
	return columns
}

func EntitiesToCreateColumnsResponse(columns []entities.Column) CreateColumnsResponse {
	respItems := make([]ColumnResponseItem, len(columns))
	for i, col := range columns {
		respItems[i] = ColumnResponseItem{
			ID:    col.ID,
			Name:  col.Name,
			Order: col.Order,
		}
	}
	return CreateColumnsResponse{
		Data:    respItems,
		Message: "Columns successfully created.",
	}
}

func EntityToColumnResponse(c column.Column) ColumnResponseItem {
	return ColumnResponseItem{
		ID:    c.ID,
		Name:  c.Name,
		Order: c.Order,
	}
}

func EntitiesToGetColumnsResponse(columns []column.Column) GetColumnsResponse {
	columnResponses := make([]ColumnResponseItem, len(columns))
	for i, col := range columns {
		columnResponses[i] = EntityToColumnResponse(col)
	}
	return GetColumnsResponse{
		Columns: columnResponses,
		Message: "Columns fetched successfully",
	}
}
