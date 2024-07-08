package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Assume this is your existing handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user mockUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registered successfully"))
}

func TestRegisterUser(t *testing.T) {
	// Initialize the test server with your existing handler
	server := httptest.NewServer(http.HandlerFunc(registerHandler))
	defer server.Close()

	newMockUser := mockUser{
		FirstName: "john",
		LastName:  "johnny",
		Email:     "john@gmail.com",
		Password:  "1233455678",
	}

	requestBody, err := json.Marshal(newMockUser)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify the response body
	expectedBody := "Registered successfully"
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != expectedBody {
		t.Fatalf("Expected response body %q, got %q", expectedBody, buf.String())
	}
}
