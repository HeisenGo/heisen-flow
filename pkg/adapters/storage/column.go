package storage

import (
	"context"
	"server/internal/column"
	"server/pkg/adapters/storage/entities"
	"server/pkg/adapters/storage/mappers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type columnRepo struct {
	db *gorm.DB
}

func NewColumnRepo(db *gorm.DB) column.Repo {
	return &columnRepo{
		db: db,
	}
}

func (r *columnRepo) SetDoneAsDefault(ctx context.Context, column *column.Column) error {
	columnEntity := mappers.ColumnDomainToEntity(*column)
	if err := r.db.WithContext(ctx).Save(&columnEntity).Error; err != nil {
		return err
	}
	column.ID = columnEntity.ID
	return nil
}

func (r *columnRepo) Create(ctx context.Context, col *column.Column) (*column.Column, error) {
	columnEntity := mappers.ColumnDomainToEntity(*col)
	if err := r.db.WithContext(ctx).Save(&columnEntity).Error; err != nil {
		return nil, err
	}
	*col = mappers.ColumnEntityToDomain(columnEntity)
	return col, nil
}

func (r *columnRepo) GetByID(ctx context.Context, id uuid.UUID) (*column.Column, error) {
	var colEntity entities.Column
	if err := r.db.WithContext(ctx).First(&colEntity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	col := mappers.ColumnEntityToDomain(colEntity)
	return &col, nil
}

func (r *columnRepo) GetMaxOrderForBoard(ctx context.Context, boardID uuid.UUID) (uint, error) {
	var maxOrder uint
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(\"order_num\"), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}
func (r *columnRepo) GetMinOrderColumn(ctx context.Context, boardID uuid.UUID) (*column.Column, error) {
	// Query to find the column with the minimum order
	var minOrderColumn entities.Column
	if err := r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Order("order_num ASC").
		First(&minOrderColumn).Error; err != nil {
		return nil, err
	}
	domainColumn := mappers.ColumnEntityToDomain(minOrderColumn)
	return &domainColumn, nil
}
func (r *columnRepo) CreateBatch(ctx context.Context, cols []column.Column) ([]column.Column, error) {
	columnEntities := mappers.ColumnDomainsToEntities(cols)
	if err := r.db.WithContext(ctx).Create(&columnEntities).Error; err != nil {
		return nil, err
	}
	return mappers.BatchColumnEntitiesToDomain(columnEntities), nil
}

func (r *columnRepo) Delete(ctx context.Context, columnID uuid.UUID) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("column_id = ?", columnID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return column.ErrColumnNotEmpty
	}

	result := r.db.WithContext(ctx).Where("id = ?", columnID).Delete(&entities.Column{})
	if result.RowsAffected == 0 {
		return column.ErrColumnNotFound
	}
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) GetByBoardID(ctx context.Context, boardID uuid.UUID) ([]column.Column, error) {
	var colEntities []entities.Column
	if err := r.db.WithContext(ctx).Where("board_id = ?", boardID).Find(&colEntities).Error; err != nil {
		return nil, err
	}
	return mappers.BatchColumnEntitiesToDomain(colEntities), nil
}

func (r *columnRepo) UpdateColumns(ctx context.Context, columns []column.Column) ([]column.Column, error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, col := range columns {
		colEntity := mappers.ColumnDomainToEntity(col)
		if err := tx.Model(&entities.Column{}).Where("id = ?", col.ID).Updates(&colEntity).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *columnRepo) ReorderColumns(ctx context.Context, boardID uuid.UUID, newOrder map[uuid.UUID]uint) error {
	var columns []entities.Column
	if err := r.db.WithContext(ctx).Where("board_id = ?", boardID).Find(&columns).Error; err != nil {
		return column.ErrFailedToFetchColumns
	}
	if len(newOrder) != len(columns) {
		return column.ErrLengthMismatch
	}
	for columnID := range newOrder {
		found := false
		for _, col := range columns {
			if col.ID == columnID {
				found = true
				break
			}
		}
		if !found {
			return column.ErrInvalidColumnID
		}
	}

	var maxOrder uint
	for _, col := range columns {
		if col.OrderNum > maxOrder {
			maxOrder = col.OrderNum
		}
	}
	tempOrder := maxOrder + 1

	for _, col := range columns {
		if err := r.db.WithContext(ctx).Model(&col).Update("order_num", tempOrder).Error; err != nil {
			return column.ErrFailedToUpdateColumn
		}
		tempOrder++
	}

	for _, col := range columns {
		newOrderNum, exists := newOrder[col.ID]
		if !exists {
			continue
		}
		if err := r.db.WithContext(ctx).Model(&col).Update("order_num", newOrderNum).Error; err != nil {
			return column.ErrFailedToUpdateColumn
		}
	}

	return nil
}

func (r *columnRepo) GetColumns(ctx context.Context, boardID uuid.UUID) ([]column.Column, error) {
	var columns []entities.Column
	err := r.db.WithContext(ctx).Where("board_id = ?", boardID).
		Order("order_num ASC").
		Find(&columns).Error
	if err != nil {
		return nil, column.ErrFailedToFetchColumns
	}
	cols := mappers.BatchColumnEntitiesToDomain(columns)
	return cols, nil
}
