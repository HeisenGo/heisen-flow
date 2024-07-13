package test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// Create mock user
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

	// Define mock board scenarios
	mockScenarios := []struct {
		name               string
		board              MockBoard
		expectedStatusCode int
	}{
		{
			name: "ValidBoardCreation",
			board: MockBoard{
				Name: "Test Board",
				Type: "private",
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "MissingName",
			board: MockBoard{
				Type: "private",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "UnauthorizedRequest",
			board: MockBoard{
				Name: "Test Board",
				Type: "private",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, scenario := range mockScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Marshal board to JSON
			boardJSON, err := json.Marshal(scenario.board)
			if err != nil {
				t.Fatalf("Failed to marshal board to JSON: %v", err)
			}

			// Create HTTP POST request to create a board
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(boardJSON))
			if err != nil {
				t.Fatalf("Failed to create HTTP request: %v", err)
			}

			// Set authorization token in the request header
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			// Assertions to validate response status code
			assert.Equal(t, scenario.expectedStatusCode, resp.StatusCode, "Expected status code")

			// Print additional information for clarity
			fmt.Printf("Scenario: %s\n", scenario.name)
			fmt.Printf("Request URL: %s\n", url)
			fmt.Printf("Request Body: %s\n", boardJSON)
			fmt.Printf("Response Status: %d\n", resp.StatusCode)
			fmt.Println("---")
		})
	}
}
