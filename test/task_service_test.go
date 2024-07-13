package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreation(t *testing.T) {
	url := fmt.Sprintf("%s%s", ServerURL, TaskPost)

	// Create mock user
	user := MockUser{
		FirstName: "task",
		LastName:  "task",
		Email:     "task@gmail.com",
		Password:  "12@Amir###90",
	}

	userResult, userData, err := CreateUserWithResp(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if userResult.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create user. Status code: %d, Response message: %s", userResult.StatusCode, userResult.Message)
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

	// Parse assignee user ID to uuid.UUID
	assigneeUUID, err := uuid.Parse(userData.UserID)
	if err != nil {
		t.Fatalf("Failed to parse assignee user ID: %v", err)
	}

	boardUUID, err := uuid.Parse(boardData.BoardID)
	if err != nil {
		t.Fatalf("Failed to parse board ID: %v", err)
	}

	// Define mock task scenarios
	mockScenarios := []struct {
		name               string
		task               MockTask
		expectedStatusCode int
	}{
		{
			name: "ValidTaskCreation",
			task: MockTask{
				Title:          "Complete Documentation",
				Description:    "Finish writing the project documentation.",
				AssigneeUserID: assigneeUUID,
				StoryPoint:     5,
				BoardID:        boardUUID,
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "MissingTitle",
			task: MockTask{
				Description:    "Finish writing the project documentation.",
				AssigneeUserID: assigneeUUID,
				StoryPoint:     5,
				BoardID:        boardUUID,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "MissingAssignee",
			task: MockTask{
				Title:       "Complete Documentation",
				Description: "Finish writing the project documentation.",
				StoryPoint:  5,
				BoardID:     boardUUID,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "UnauthorizedRequest",
			task: MockTask{
				Title:          "Complete Documentation",
				Description:    "Finish writing the project documentation.",
				AssigneeUserID: assigneeUUID,
				StoryPoint:     5,
				BoardID:        boardUUID,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "InvalidBoardID",
			task: MockTask{
				Title:          "Complete Documentation",
				Description:    "Finish writing the project documentation.",
				AssigneeUserID: assigneeUUID,
				StoryPoint:     5,
				BoardID:        uuid.Nil, // Invalid UUID
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, scenario := range mockScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Marshal task to JSON
			taskJSON, err := json.Marshal(scenario.task)
			if err != nil {
				t.Fatalf("Failed to marshal task to JSON: %v", err)
			}

			// Create HTTP request to create a task
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskJSON))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			taskResp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer taskResp.Body.Close()

			// Read response body
			body, err := io.ReadAll(taskResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Verify if task creation was successful
			assert.Equal(t, scenario.expectedStatusCode, taskResp.StatusCode, "Expected status code")
			if taskResp.StatusCode == http.StatusCreated {
				// Optionally, you can validate the response body further if needed
				var res Response
				err = json.Unmarshal(body, &res)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				// Add more assertions if needed for successful creation scenario
			}
		})
	}
}
