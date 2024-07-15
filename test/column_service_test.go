package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
	// Create mock user
	user := MockUser{
		FirstName: "column",
		LastName:  "column",
		Email:     "column@gmail.com",
		Password:  "12@Amir###90",
	}

	result, _, err := CreateUserWithResp(user)
	if err != nil || result.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create user. Status code: %d, Response message: %s, Error: %v", result.StatusCode, result.Message, err)
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
		Name: "Column Board",
		Type: "public",
	}
	boardResp, boardData, err := CreateBoard(token, newMockBoard)
	if err != nil || boardResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create board. Status code: %d, Error: %v", boardResp.StatusCode, err)
	}

	// Extract board ID from boardData
	boardID, err := uuid.Parse(boardData.BoardID)
	if err != nil {
		t.Fatalf("Failed to parse board ID: %v", err)
	}

	// Helper function to create columns
	createColumns := func(payload map[string]interface{}, expectedStatusCode int, token string) {
		// Marshal payload to JSON
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload to JSON: %v", err)
		}

		// URL for creating columns
		url := fmt.Sprintf("%s%s", ServerURL, ColumnPost)

		// Create HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Log response body for debugging
		t.Logf("Response Body: %s", body)

		assert.Equal(t, expectedStatusCode, resp.StatusCode, fmt.Sprintf("Expected status code %d", expectedStatusCode))
	}

	// Test cases
	tests := []struct {
		name               string
		payload            map[string]interface{}
		expectedStatusCode int
	}{
		{
			name: "InvalidBoardID",
			payload: map[string]interface{}{
				"board_id": "invalid-board-id",
				"columns": []map[string]string{
					{"name": "todo"},
					{"name": "in progress"},
					{"name": "in review"},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "DuplicateNames",
			payload: map[string]interface{}{
				"board_id": boardID,
				"columns": []map[string]string{
					{"name": "todo"},
					{"name": "todo"},
					{"name": "in review"},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "InvalidPayload",
			payload: map[string]interface{}{
				"board_id": "invalid-id",
				"columns":  "invalid-column",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "ExcessivelyLongNames",
			payload: map[string]interface{}{
				"board_id": boardID,
				"columns": []map[string]string{
					{"name": "thisisaverylongcolumnnamethatshouldfailbecauseitiswaytoolong"},
					{"name": "anotherexcessivelylongcolumnname"},
					{"name": "yetagainaverylongcolumnnamethatshouldnotbeallowed"},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "WithoutAuthorization" {
				createColumns(tt.payload, tt.expectedStatusCode, "")
			} else {
				createColumns(tt.payload, tt.expectedStatusCode, token)
			}
		})
	}
}

func TestCreateColumns(t *testing.T) {
	// Create mock user
	user := MockUser{
		FirstName: "column2",
		LastName:  "column2",
		Email:     "column2@gmail.com",
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
		Name: "Column2 Board",
		Type: "public",
	}
	boardResp, boardData, err := CreateBoard(token, newMockBoard)
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	if boardResp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create board. Status code: %d", boardResp.StatusCode)
	}

	// Ensure boardData is not nil and has expected fields
	if boardData.BoardID == "" {
		t.Fatalf("Board ID is empty")
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
	url := fmt.Sprintf("%s%s", ServerURL, "/columns")

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Log response body for debugging
	t.Logf("Response Body: %s", body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
}
