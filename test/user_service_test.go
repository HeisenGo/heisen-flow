package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type mockUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func TestRegisterUser(t *testing.T) {
	// Initialize the test server with your existing handler
	newMockUser := mockUser{
		FirstName: "john",
		LastName:  "johnny",
		Email:     "john@gmail.com",
		Password:  "12@Amir###90",
	}

	requestBody, err := json.Marshal(newMockUser)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/api/v1/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify the response body
	expectedBody := "user successfully registered."
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != expectedBody {
		t.Fatalf("Expected response body %q, got %q", expectedBody, buf.String())
	}
}
