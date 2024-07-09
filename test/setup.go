package test

import (
	"log"
	"os"
	"testing"
	"gorm.io/gorm"
	"server/config"
	"server/pkg/adapters/storage"
	"server/service"
)

var (
	TestDB     *gorm.DB
	AppConfig  config.Config
	AppService *service.AppContainer
)

func TestMain(m *testing.M) {
	// Load configuration
	cfg, err := config.ReadStandard("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}
	AppConfig = cfg

	// Set up SQLite in-memory database for testing
	TestDB, err = storage.NewSQLiteGormConnection()
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	storage.Migrate(TestDB)

	// Initialize application service container
	AppService, err = service.NewAppContainer(AppConfig)
	if err != nil {
		log.Fatalf("Failed to initialize AppContainer: %v", err)
	}

	AppService.RawDBConnection().Exec("PRAGMA foreign_keys = ON;")

	// Run the tests
	code := m.Run()

	os.Exit(code)
}
