package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateColumns(t *testing.T) {
	// Create mock board
	user := MockUser{
		FirstName: "board",
		LastName:  "board",
		Email:     "board@gmail.com",
		Password:  "12@Amir###90",
	}

	result := CreateUser(user)
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create user. Status code: %d, Response message: %s", result.StatusCode, result.Message)
	}

	// Login to obtain authentication token
	token, err := LoginAndGetToken(t, MockUserLogin{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Create mock board
	newMockBoard := MockBoard{
		Name: "Task Board",
		Type: "public",
	}
	boardResp, boardData, err := CreateBoard(token, newMockBoard)
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	if boardResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create board. Status code: %d", boardResp.StatusCode)
	}

	// Extract board ID from boardData
	boardID, err := uuid.Parse(boardData.BoardID)
	if err != nil {
		t.Fatalf("Failed to parse board ID: %v", err)
	}

	// Define columns payload
	columns := []map[string]string{
		{"name": "todo"},
		{"name": "in progress"},
		{"name": "in review"},
	}
	payload := map[string]interface{}{
		"board_id": boardID,
		"columns":  columns,
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload to JSON: %v", err)
	}

	// URL for creating columns
	url := fmt.Sprintf("%s%s", ServerURL, "/api/v1/columns")

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")

}

// MockBoard and CreateBoard function should be defined as per your application implementation
// Adjust the Column struct and assertions as needed based on your API's response structure
