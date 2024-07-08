package test

import (
	"flag"
	"log"
	"os"
	"testing"

	"gorm.io/gorm"

	http_server "server/api/http"
	"server/config"
	"server/internal/user"
	"server/service"
)

var (
	TestDB     *gorm.DB
	AppConfig  config.Config
	AppService *service.AppContainer
)

var configPath = flag.String("config", "", "configuration path")

func TestMain(m *testing.M) {
	// Load configuration
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http_server.Run(cfg.Server, app)
	}()

	// Run tests
	code := m.Run()

	os.Exit(code)

}

func ClearDatabaseTables(db *gorm.DB) error {
	tables := []interface{}{
		&user.User{},
	}

	for _, model := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
			return err
		}
	}

	return nil
}

func readConfig() config.Config {
	flag.Parse()

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	if len(*configPath) == 0 {
		log.Fatal("configuration file not found")
	}

	cfg, err := config.ReadStandard(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
