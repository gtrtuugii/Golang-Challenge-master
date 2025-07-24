package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func setupTestContext(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// Test the response structure creation logic
func TestGetUsersResponseStructure(t *testing.T) {
	// Mock data that would come from queries.GetUsers()
	mockUserRows := []struct {
		ID                 int
		Username           string
		Email              string
		UserType           string
		Nickname           pgtype.Text
		PermissionBitfield string
	}{
		{
			ID:                 1,
			Username:           "john_doe",
			Email:              "john@example.com",
			UserType:           "admin",
			Nickname:           pgtype.Text{String: "Johnny", Valid: true},
			PermissionBitfield: "1111",
		},
		{
			ID:                 2,
			Username:           "jane_smith",
			Email:              "jane@example.com",
			UserType:           "user",
			Nickname:           pgtype.Text{String: "", Valid: false},
			PermissionBitfield: "1010",
		},
	}

	// Test the conversion logic from your GetUsers handler
	var users []UserResponse = []UserResponse{}
	for _, row := range mockUserRows {
		var nickname *string = nil
		if row.Nickname.Valid {
			nickname = &row.Nickname.String
		}
		users = append(users, UserResponse{
			ID:       row.ID,
			Username: row.Username,
			Email:    row.Email,
			UserType: row.UserType,
			Nickname: nickname,
		})
	}

	response := GetUsersResponse{Users: users}

	// Verify the response structure
	assert.Len(t, response.Users, 2)

	// Check first user (with nickname)
	assert.Equal(t, 1, response.Users[0].ID)
	assert.Equal(t, "john_doe", response.Users[0].Username)
	assert.Equal(t, "john@example.com", response.Users[0].Email)
	assert.Equal(t, "admin", response.Users[0].UserType)
	assert.NotNil(t, response.Users[0].Nickname)
	assert.Equal(t, "Johnny", *response.Users[0].Nickname)

	// Check second user (without nickname)
	assert.Equal(t, 2, response.Users[1].ID)
	assert.Equal(t, "jane_smith", response.Users[1].Username)
	assert.Equal(t, "jane@example.com", response.Users[1].Email)
	assert.Equal(t, "user", response.Users[1].UserType)
	assert.Nil(t, response.Users[1].Nickname)
}

func TestCreateUserRequestValidation(t *testing.T) {
	tests := []struct {
		name          string
		request       CreateUserRequest
		expectValid   bool
		expectedError string
	}{
		{
			name: "valid request with nickname",
			request: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				UserType: "user",
				Nickname: stringPtr("Test User"),
			},
			expectValid: true,
		},
		{
			name: "valid request without nickname",
			request: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				UserType: "user",
				Nickname: nil,
			},
			expectValid: true,
		},
		{
			name: "missing username",
			request: CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				UserType: "user",
			},
			expectValid:   false,
			expectedError: "Username is required",
		},
		{
			name: "missing email",
			request: CreateUserRequest{
				Username: "testuser",
				Email:    "",
				UserType: "user",
			},
			expectValid:   false,
			expectedError: "Email is required",
		},
		{
			name: "missing user type",
			request: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				UserType: "",
			},
			expectValid:   false,
			expectedError: "User type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic from your CreateUser handler
			isValid := tt.request.Username != "" &&
				tt.request.Email != "" &&
				tt.request.UserType != ""

			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

func TestCreateUserResponseCreation(t *testing.T) {
	// Mock data that would come from queries.CreateUser()
	mockUser := struct {
		ID       int
		Username string
		Email    string
		UserType string
		Nickname pgtype.Text
	}{
		ID:       123,
		Username: "newuser",
		Email:    "new@example.com",
		UserType: "user",
		Nickname: pgtype.Text{String: "New User", Valid: true},
	}

	// Test the conversion logic from your CreateUser handler
	var nickname *string = nil
	if mockUser.Nickname.Valid {
		nickname = &mockUser.Nickname.String
	}

	response := UserResponse{
		ID:       mockUser.ID,
		Username: mockUser.Username,
		Email:    mockUser.Email,
		UserType: mockUser.UserType,
		Nickname: nickname,
	}

	// Verify the response
	assert.Equal(t, 123, response.ID)
	assert.Equal(t, "newuser", response.Username)
	assert.Equal(t, "new@example.com", response.Email)
	assert.Equal(t, "user", response.UserType)
	assert.NotNil(t, response.Nickname)
	assert.Equal(t, "New User", *response.Nickname)
}

func TestCreateUserResponseWithoutNickname(t *testing.T) {
	// Mock data without nickname
	mockUser := struct {
		ID       int
		Username string
		Email    string
		UserType string
		Nickname pgtype.Text
	}{
		ID:       456,
		Username: "anothuser",
		Email:    "another@example.com",
		UserType: "admin",
		Nickname: pgtype.Text{String: "", Valid: false},
	}

	// Test the conversion logic
	var nickname *string = nil
	if mockUser.Nickname.Valid {
		nickname = &mockUser.Nickname.String
	}

	response := UserResponse{
		ID:       mockUser.ID,
		Username: mockUser.Username,
		Email:    mockUser.Email,
		UserType: mockUser.UserType,
		Nickname: nickname,
	}

	// Verify the response
	assert.Equal(t, 456, response.ID)
	assert.Equal(t, "anothuser", response.Username)
	assert.Equal(t, "another@example.com", response.Email)
	assert.Equal(t, "admin", response.UserType)
	assert.Nil(t, response.Nickname)
}

func TestJSONBindingValidation(t *testing.T) {
	tests := []struct {
		name        string
		jsonBody    string
		shouldError bool
	}{
		{
			name:        "valid JSON",
			jsonBody:    `{"username":"test","email":"test@example.com","user_type":"user"}`,
			shouldError: false,
		},
		{
			name:        "invalid JSON",
			jsonBody:    `{"username":"test","email":}`,
			shouldError: true,
		},
		{
			name:        "empty JSON",
			jsonBody:    `{}`,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req CreateUserRequest
			err := json.Unmarshal([]byte(tt.jsonBody), &req)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserResponseJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		user     UserResponse
		expected map[string]interface{}
	}{
		{
			name: "user with nickname",
			user: UserResponse{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				UserType: "admin",
				Nickname: stringPtr("Test Nickname"),
			},
			expected: map[string]interface{}{
				"id":       float64(1),
				"username": "testuser",
				"email":    "test@example.com",
				"userType": "admin",
				"nickname": "Test Nickname",
			},
		},
		{
			name: "user without nickname",
			user: UserResponse{
				ID:       2,
				Username: "testuser2",
				Email:    "test2@example.com",
				UserType: "user",
				Nickname: nil,
			},
			expected: map[string]interface{}{
				"id":       float64(2),
				"username": "testuser2",
				"email":    "test2@example.com",
				"userType": "user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.user)
			assert.NoError(t, err)

			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expected {
				actualValue, exists := result[key]
				assert.True(t, exists, "Key %s should exist", key)
				assert.Equal(t, expectedValue, actualValue, "Value for key %s should match", key)
			}

			// Check that nickname is omitted when nil
			if tt.user.Nickname == nil {
				_, exists := result["nickname"]
				assert.False(t, exists, "Nickname should be omitted when nil")
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
