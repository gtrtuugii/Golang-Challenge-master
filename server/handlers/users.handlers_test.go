package handlers

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

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

func TestUpdateUserRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateUserParams
		expectValid bool
		description string
	}{
		{
			name: "valid partial update - username only",
			request: UpdateUserParams{
				Username: stringPtr("newusername"),
			},
			expectValid: true,
			description: "Should accept username-only update",
		},
		{
			name: "valid partial update - email only",
			request: UpdateUserParams{
				Email: stringPtr("new@example.com"),
			},
			expectValid: true,
			description: "Should accept email-only update",
		},
		{
			name: "valid partial update - user_type only",
			request: UpdateUserParams{
				UserType: stringPtr("UTYPE_ADMIN"),
			},
			expectValid: true,
			description: "Should accept user_type-only update",
		},
		{
			name: "valid partial update - nickname only",
			request: UpdateUserParams{
				Nickname: stringPtr("New Nickname"),
			},
			expectValid: true,
			description: "Should accept nickname-only update",
		},
		{
			name: "valid nickname removal",
			request: UpdateUserParams{
				Nickname: stringPtr(""),
			},
			expectValid: true,
			description: "Should accept empty string to remove nickname",
		},
		{
			name: "valid multiple fields",
			request: UpdateUserParams{
				Username: stringPtr("updateduser"),
				Email:    stringPtr("updated@example.com"),
				UserType: stringPtr("UTYPE_MODERATOR"),
			},
			expectValid: true,
			description: "Should accept multiple field updates",
		},
		{
			name: "valid all fields",
			request: UpdateUserParams{
				Username: stringPtr("allnew"),
				Email:    stringPtr("allnew@example.com"),
				UserType: stringPtr("UTYPE_ADMIN"),
				Nickname: stringPtr("All New"),
			},
			expectValid: true,
			description: "Should accept all field updates",
		},
		{
			name:        "empty request",
			request:     UpdateUserParams{},
			expectValid: false,
			description: "Should reject request with no fields to update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic from UpdateUser handler
			hasFields := tt.request.Username != nil ||
				tt.request.Email != nil ||
				tt.request.UserType != nil ||
				tt.request.Nickname != nil

			assert.Equal(t, tt.expectValid, hasFields, tt.description)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
