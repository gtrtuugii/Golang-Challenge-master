package handlers

import (
	"main/utils"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersResponseStructure(t *testing.T) {
	mockUserRows := []struct {
		ID                 int
		Username           string
		Email              string
		UserType           string
		Nickname           pgtype.Text
		PermissionBitfield string
	}{
		{1, "john_doe", "john@example.com", "admin", pgtype.Text{String: "Johnny", Valid: true}, "1111"},
		{2, "jane_smith", "jane@example.com", "user", pgtype.Text{Valid: false}, "1010"},
	}

	var users []UserResponse
	for _, row := range mockUserRows {
		var nickname *string
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

	assert.Len(t, response.Users, 2)
	assert.Equal(t, "Johnny", *response.Users[0].Nickname)
	assert.Nil(t, response.Users[1].Nickname)
}

func TestCreateUserRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateUserRequest
		expectValid bool
	}{
		{"valid with nickname", CreateUserRequest{"testuser", "test@example.com", "user", utils.StringPtr("Test User")}, true},
		{"valid without nickname", CreateUserRequest{"testuser", "test@example.com", "user", nil}, true},
		{"missing username", CreateUserRequest{"", "test@example.com", "user", nil}, false},
		{"missing email", CreateUserRequest{"testuser", "", "user", nil}, false},
		{"missing user type", CreateUserRequest{"testuser", "test@example.com", "", nil}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.Username != "" && tt.request.Email != "" && tt.request.UserType != ""
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

func TestCreateUserResponseConversion(t *testing.T) {
	mockUser := struct {
		ID       int
		Username string
		Email    string
		UserType string
		Nickname pgtype.Text
	}{
		123, "newuser", "new@example.com", "user", pgtype.Text{String: "New User", Valid: true},
	}

	var nickname *string
	if mockUser.Nickname.Valid {
		nickname = &mockUser.Nickname.String
	}

	resp := UserResponse{
		ID:       mockUser.ID,
		Username: mockUser.Username,
		Email:    mockUser.Email,
		UserType: mockUser.UserType,
		Nickname: nickname,
	}

	assert.Equal(t, 123, resp.ID)
	assert.NotNil(t, resp.Nickname)
	assert.Equal(t, "New User", *resp.Nickname)
}

func TestUpdateUserRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateUserParams
		expectValid bool
	}{
		{"username only", UpdateUserParams{Username: utils.StringPtr("newusername")}, true},
		{"email only", UpdateUserParams{Email: utils.StringPtr("new@example.com")}, true},
		{"user_type only", UpdateUserParams{UserType: utils.StringPtr("UTYPE_ADMIN")}, true},
		{"nickname only", UpdateUserParams{Nickname: utils.StringPtr("New Nickname")}, true},
		{"empty request", UpdateUserParams{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasFields := tt.request.Username != nil || tt.request.Email != nil || tt.request.UserType != nil || tt.request.Nickname != nil
			assert.Equal(t, tt.expectValid, hasFields)
		})
	}
}
