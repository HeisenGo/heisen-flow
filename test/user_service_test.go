package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	url := fmt.Sprintf("%s%s", ServerURL, Register)

	t.Run("successful register", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "status code should be 201")
	})

	t.Run("email not provided", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("password not provided", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("invalid email format", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@invalid",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("weak password", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "123",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("email already registered", func(t *testing.T) {
		// Assuming the email is already registered in the system
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode, "status code should be 409")
	})

	t.Run("missing first name", func(t *testing.T) {
		newMockUser := MockUser{
			LastName: "johnny",
			Email:    "john@gmail.com",
			Password: "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("missing last name", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("short first name", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "j",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("short last name", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "j",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("password without special character", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12Amir90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("password without uppercase letter", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("password without number", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "@Amir###",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("empty request body", func(t *testing.T) {
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		invalidJSON := []byte(`{"first_name": "john"}`)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(invalidJSON))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be 400")
	})

	t.Run("unsupported method (GET)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatalf("Failed to create GET request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "status code should be 405")
	})

	t.Run("server error", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@Amir###90",
		}

		reqBody, err := json.Marshal(newMockUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		// Change ServerURL to a non-existent address to force a server error
		url := fmt.Sprintf("%s%s", "http://localhost:9999", Register)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "status code should be 500")
	})
}
