package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "status code should be 201")
		assert.Equal(t, "user successfully registered", res.Message, "message should be 'user successfully registered'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "Key: 'UserRegisterReq.Email' Error:Field validation for 'Email' failed on the 'required' tag", "error message should contain 'Key: 'UserRegisterReq.Email' Error:Field validation for 'Email' failed on the 'required' tag'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "Key: 'UserRegisterReq.Password' Error:Field validation for 'Password' failed on the 'required' tag", "error message should contain 'Key: 'UserRegisterReq.Password' Error:Field validation for 'Password' failed on the 'required' tag'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid email format")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid password format", "error message should contain 'invalid password format'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusConflict, resp.StatusCode, "status code should be %s", fiber.StatusConflict)
		assert.Contains(t, res.Error, "email already exists", "error message should contain 'email already exists'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}
		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "Key: 'UserRegisterReq.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag", "error message should contain 'Key: 'UserRegisterReq.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "Key: 'UserRegisterReq.LastName' Error:Field validation for 'LastName' failed on the 'required' tag", "error message should contain 'Key: 'UserRegisterReq.LastName' Error:Field validation for 'LastName' failed on the 'required' tag'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid password format", "error message should contain 'invalid password format'")
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid password format", "error message should contain 'invalid password format'")
	})

	t.Run("password without lowercase letter", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "12@AMIR###90",
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid password format", "error message should contain 'invalid password format'")
	})

	t.Run("password without number", func(t *testing.T) {
		newMockUser := MockUser{
			FirstName: "john",
			LastName:  "johnny",
			Email:     "john@gmail.com",
			Password:  "@Amir###xyz",
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "status code should be %s", fiber.StatusBadRequest)
		assert.Contains(t, res.Error, "invalid password format", "error message should contain 'invalid password format'")
	})
}
