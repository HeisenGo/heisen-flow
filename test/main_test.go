package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"

	http_server "server/api/http"
	"server/config"

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

	// Clear database tables
	ClearDatabase(app.RawDBConnection())

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func ClearDatabase(db *gorm.DB) error {
	// Disable foreign key checks
	db.Exec("SET session_replication_role = 'replica';")

	// Clear all tables
	err := db.Exec("TRUNCATE TABLE user_board_roles, users, boards, columns, tasks RESTART IDENTITY CASCADE;").Error

	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = 'origin';")

	return err
}

func readConfig() config.Config {
	cfg, err := config.ReadStandard(configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func CreateUser(user MockUser) UserCreationResult {
	url := fmt.Sprintf("%s%s", ServerURL, Register)

	reqBody, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	res := new(Response)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	return UserCreationResult{
		StatusCode: resp.StatusCode,
		Message:    res.Message,
	}
}

func LoginAndGetToken(t *testing.T, user MockUserLogin) (string, error) {
	reqBody, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal reqBody: %v", err)
	}

	url := fmt.Sprintf("%s%s", ServerURL, Login)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status code 200 OK, got %s", resp.Status)
	}

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	if res.Data == nil {
		return "", fmt.Errorf("response data is nil")
	}

	tokenData, ok := res.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to convert response data to map[string]interface{}")
	}

	authToken, ok := tokenData["auth_token"].(string)
	if !ok {
		return "", fmt.Errorf("auth_token not found or not a string")
	}

	return authToken, nil
}
