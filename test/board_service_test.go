package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardCreation(t *testing.T) {
	url := fmt.Sprintf("%s%s", ServerURL, BoardPost)
	// Mock user credentials
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

	// Obtain authentication token
	token, err := LoginAndGetToken(t, MockUserLogin{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Construct request body for board creation
	newMockBoard := MockBoard{
		Name: "Test Board",
		Type: "private",
	}

	reqBody, err := json.Marshal(newMockBoard)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create HTTP POST request to create a board
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set authorization token in the request header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Perform request using Fiber app handler
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Unmarshal response body into Response struct
	res := new(Response)
	err = json.Unmarshal(body, &res)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Assertions to validate response
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
	assert.Equal(t, "Board created successfully", res.Message, "Expected message 'Board created successfully'")
}
