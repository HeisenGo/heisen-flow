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

		url := fmt.Sprintf("%s%s", ServerURL, Register)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// parse the response
		res := new(Response)
		err = json.Unmarshal(body, &res)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "status code should be 201")
		assert.Equal(t, res.Message, "user successfully registered", "message should be 'user successfully registered'")
	})

	t.Run("", func(t *testing.T) {})
}

// func TestBoard(t *testing.T) {
// 	t.Run("successful board creation", func(t *testing.T) {
// 		newMockUser := MockUserLogin{
// 			Email:    "john@gmail.com",
// 			Password: "12@Amir###90",
// 		}
// 		token := LoginAndGetToken(t, newMockUser)

// 		reqBody, err := json.Marshal(newMockUser)
// 		if err != nil {
// 			t.Fatalf("Failed to marshal request: %v", err)
// 		}

// 		url := fmt.Sprintf("%s%s", ServerURL, BoardPost)
// 		req, err := http.NewRequest("Post", url, bytes.NewReader(reqBody))
// 		req.Header.Set("Authorization", "Bearer "+token)
// 		var response io.Writer
// 		// send request
// 		req.Write(response)
// 		defer req.Response.Body.Close()

// 		body, err := io.ReadAll(req.Response.Body)
// 		if err != nil {
// 			t.Fatalf("Failed to read response body: %v", err)
// 		}

// 		var resp Response
// 		err = json.Unmarshal(body, &resp)
// 		if err != nil {
// 			t.Fatalf("Failed to read response body: %v", err)
// 		}

// 		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

// 	})
// }
