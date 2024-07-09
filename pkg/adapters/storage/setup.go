package storage

import (
	"fmt"
	"server/config"
	"server/pkg/adapters/storage/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormConnection(dbConfig config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.DBName, dbConfig.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AddExtension(db *gorm.DB) error {
	return db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

func Migrate(db *gorm.DB) error {
	migrator := db.Migrator()

	err := migrator.AutoMigrate(&entities.User{},
		&entities.Board{}, &entities.UserBoardRole{},
		&entities.Task{}, &entities.TaskDependency{})
	if err != nil {
		return err
	}
	err = migrator.AutoMigrate(&entities.Board{}, &entities.UserBoardRole{})
	if err != nil {
		return err
	}
	err = migrator.AutoMigrate(&entities.Column{})
	if err != nil {
		return err
	}
	//err = migrator.AutoMigrate(&entities.UserBoardRole{})
	// err = migrator.AutoMigrate(&entities.Board{}, &entities.UserBoardRole{})
	// if err != nil {
	// 	return err
	// }
	// err = migrator.AutoMigrate(&entities.Task{}, &entities.TaskDependency{})

	// if err != nil {
	// 	return err
	// }
	return nil
}
