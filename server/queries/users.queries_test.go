package queries

import (
	"errors"
	"main/utils"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
	users []GetUsersQueryRow
}

func TestGetUsersQueryRow_Structure(t *testing.T) {
	// Test creating and validating GetUsersQueryRow structure
	row := GetUsersQueryRow{
		ID:                 1,
		Username:           "testuser",
		Email:              "test@example.com",
		UserType:           "user",
		Nickname:           pgtype.Text{String: "Test User", Valid: true},
		PermissionBitfield: "1010",
	}

	assert.Equal(t, 1, row.ID)
	assert.Equal(t, "testuser", row.Username)
	assert.Equal(t, "test@example.com", row.Email)
	assert.Equal(t, "user", row.UserType)
	assert.True(t, row.Nickname.Valid)
	assert.Equal(t, "Test User", row.Nickname.String)
	assert.Equal(t, "1010", row.PermissionBitfield)
}

func TestGetUsersQueryRow_WithInvalidNickname(t *testing.T) {
	// Test with invalid/null nickname
	row := GetUsersQueryRow{
		ID:                 2,
		Username:           "testuser2",
		Email:              "test2@example.com",
		UserType:           "admin",
		Nickname:           pgtype.Text{String: "", Valid: false},
		PermissionBitfield: "1111",
	}

	assert.Equal(t, 2, row.ID)
	assert.Equal(t, "testuser2", row.Username)
	assert.Equal(t, "test2@example.com", row.Email)
	assert.Equal(t, "admin", row.UserType)
	assert.False(t, row.Nickname.Valid)
	assert.Equal(t, "", row.Nickname.String)
	assert.Equal(t, "1111", row.PermissionBitfield)
}
func TestCreateUserParams_WithNickname(t *testing.T) {
	nickname := "Test Nickname"
	params := CreateUserParams{
		Username: "newuser",
		Email:    "new@example.com",
		UserType: "user",
		Nickname: &nickname,
	}

	assert.Equal(t, "newuser", params.Username)
	assert.Equal(t, "new@example.com", params.Email)
	assert.Equal(t, "user", params.UserType)
	assert.NotNil(t, params.Nickname)
	assert.Equal(t, "Test Nickname", *params.Nickname)
}

func TestCreateUserParams_WithoutNickname(t *testing.T) {
	params := CreateUserParams{
		Username: "newuser2",
		Email:    "new2@example.com",
		UserType: "admin",
		Nickname: nil,
	}

	assert.Equal(t, "newuser2", params.Username)
	assert.Equal(t, "new2@example.com", params.Email)
	assert.Equal(t, "admin", params.UserType)
	assert.Nil(t, params.Nickname)
}

func TestCreateUserParams_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateUserParams
		isValid bool
	}{
		{
			name: "valid params with nickname",
			params: CreateUserParams{
				Username: "validuser",
				Email:    "valid@example.com",
				UserType: "user",
				Nickname: utils.StringPtr("Valid User"),
			},
			isValid: true,
		},
		{
			name: "valid params without nickname",
			params: CreateUserParams{
				Username: "validuser2",
				Email:    "valid2@example.com",
				UserType: "admin",
				Nickname: nil,
			},
			isValid: true,
		},
		{
			name: "invalid - empty username",
			params: CreateUserParams{
				Username: "",
				Email:    "valid@example.com",
				UserType: "user",
				Nickname: nil,
			},
			isValid: false,
		},
		{
			name: "invalid - empty email",
			params: CreateUserParams{
				Username: "validuser",
				Email:    "",
				UserType: "user",
				Nickname: nil,
			},
			isValid: false,
		},
		{
			name: "invalid - empty user type",
			params: CreateUserParams{
				Username: "validuser",
				Email:    "valid@example.com",
				UserType: "",
				Nickname: nil,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic (mimics what should be done before DB call)
			isValid := tt.params.Username != "" &&
				tt.params.Email != "" &&
				tt.params.UserType != ""

			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestNicknameHandling(t *testing.T) {
	tests := []struct {
		name             string
		inputNickname    *string
		expectedPgType   pgtype.Text
		expectedResponse *string
	}{
		{
			name:          "with nickname",
			inputNickname: utils.StringPtr("Test User"),
			expectedPgType: pgtype.Text{
				String: "Test User",
				Valid:  true,
			},
			expectedResponse: utils.StringPtr("Test User"),
		},
		{
			name:          "without nickname",
			inputNickname: nil,
			expectedPgType: pgtype.Text{
				String: "",
				Valid:  false,
			},
			expectedResponse: nil,
		},
		{
			name:          "empty nickname",
			inputNickname: utils.StringPtr(""),
			expectedPgType: pgtype.Text{
				String: "",
				Valid:  true,
			},
			expectedResponse: utils.StringPtr(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test input to pgtype.Text conversion (mimics CreateUser logic)
			var pgNickname pgtype.Text
			if tt.inputNickname != nil {
				pgNickname = pgtype.Text{String: *tt.inputNickname, Valid: true}
			} else {
				pgNickname = pgtype.Text{Valid: false}
			}

			assert.Equal(t, tt.expectedPgType.Valid, pgNickname.Valid)
			if pgNickname.Valid {
				assert.Equal(t, tt.expectedPgType.String, pgNickname.String)
			}

			// Test pgtype.Text to response conversion (mimics GetUsers logic)
			var responseNickname *string = nil
			if pgNickname.Valid {
				responseNickname = &pgNickname.String
			}

			if tt.expectedResponse == nil {
				assert.Nil(t, responseNickname)
			} else {
				assert.NotNil(t, responseNickname)
				assert.Equal(t, *tt.expectedResponse, *responseNickname)
			}
		})
	}
}

func TestUserTypesAndPermissions(t *testing.T) {
	// Test different user types and their expected structure
	userTypes := []struct {
		userType         string
		expectedBitfield string
	}{
		{"admin", "1111"},
		{"user", "1010"},
		{"moderator", "1100"},
		{"guest", "0001"},
	}

	for _, ut := range userTypes {
		t.Run("user_type_"+ut.userType, func(t *testing.T) {
			row := GetUsersQueryRow{
				ID:                 1,
				Username:           "test_" + ut.userType,
				Email:              ut.userType + "@example.com",
				UserType:           ut.userType,
				Nickname:           pgtype.Text{String: "Test " + ut.userType, Valid: true},
				PermissionBitfield: ut.expectedBitfield,
			}

			assert.Equal(t, ut.userType, row.UserType)
			assert.Equal(t, ut.expectedBitfield, row.PermissionBitfield)
		})
	}
}

func (m *MockUserService) UpdateUser(userID int, params UpdateUserParams) (GetUsersQueryRow, error) {
	args := m.Called(userID, params)

	if args.Error(1) != nil {
		return GetUsersQueryRow{}, args.Error(1)
	}

	// Find user to update
	for i, user := range m.users {
		if user.ID == int(userID) {
			// Update fields if provided
			if params.Username != nil {
				user.Username = *params.Username
			}
			if params.Email != nil {
				user.Email = *params.Email
			}
			if params.UserType != nil {
				user.UserType = *params.UserType
			}
			if params.Nickname != nil {
				if *params.Nickname == "" {
					user.Nickname = pgtype.Text{Valid: false}
				} else {
					user.Nickname = pgtype.Text{String: *params.Nickname, Valid: true}
				}
			}

			// Mock permission bitfield based on user type
			switch user.UserType {
			case "UTYPE_ADMIN":
				user.PermissionBitfield = "10000000"
			case "UTYPE_MODERATOR":
				user.PermissionBitfield = "01000000"
			default:
				user.PermissionBitfield = "00000000"
			}

			user.MessageCount = 0
			m.users[i] = user
			return user, nil
		}
	}

	return GetUsersQueryRow{}, errors.New("user not found")
}
