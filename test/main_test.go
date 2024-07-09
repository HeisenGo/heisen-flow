package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	http_server "server/api/http"
	"server/config"
	"server/internal/user"
	"server/service"
)

var (
	TestDB *gorm.DB
)

const (
	ServerURL  = "http://0.0.0.0:8080/api/v1"
	Register   = "/register"
	Login      = "/login"
	BoardPost  = "/boards"
	configPath = "test_config.yaml"
)

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

	// wait for server to start
	time.Sleep(2 * time.Second)

	// clear the database
	ClearDatabaseTables(app.RawDBConnection())

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
	cfg, err := config.ReadStandard(configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func LoginAndGetToken(t *testing.T, user MockUserLogin) string {
	reqBody, err := json.Marshal(user)
	if err != nil {
		t.Fatal("failed to marshal reqBody")
	}

	url := fmt.Sprintf("%s%s", ServerURL, Login)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.Status, "Expected status code to be 200 OK")

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		t.Fatalf("Failed to decode token response: %v", err)
	}

	type Token struct {
		AuthToken    string `json:"auth_token"`
		RefreshToken string `json:"refresh_token"`
	}

	token := res.Data.(Token).AuthToken

	return token
}
