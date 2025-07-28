package queries

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(params CreateUserParams) (GetUsersQueryRow, error) {
	args := m.Called(params)
	return args.Get(0).(GetUsersQueryRow), args.Error(1)
}

func (m *MockUserService) GetUsers() ([]GetUsersQueryRow, error) {
	args := m.Called()
	return args.Get(0).([]GetUsersQueryRow), args.Error(1)
}

func (m *MockUserService) UpdateUser(userID int, params UpdateUserParams) (GetUsersQueryRow, error) {
	args := m.Called(userID, params)
	return args.Get(0).(GetUsersQueryRow), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockService := &MockUserService{}

	nickname := "Test User"
	params := CreateUserParams{
		Username: "testuser",
		Email:    "test@example.com",
		UserType: "user",
		Nickname: &nickname,
	}

	expectedUser := GetUsersQueryRow{
		ID:                 1,
		Username:           "testuser",
		Email:              "test@example.com",
		UserType:           "user",
		Nickname:           pgtype.Text{String: "Test User", Valid: true},
		PermissionBitfield: "1010",
	}

	mockService.On("CreateUser", params).Return(expectedUser, nil)

	user, err := mockService.CreateUser(params)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.True(t, user.Nickname.Valid)

	mockService.AssertExpectations(t)
}

func TestCreateUserError(t *testing.T) {
	mockService := &MockUserService{}

	params := CreateUserParams{
		Username: "duplicate",
		Email:    "test@example.com",
		UserType: "user",
		Nickname: nil,
	}

	mockService.On("CreateUser", params).Return(GetUsersQueryRow{}, errors.New("username already exists"))

	_, err := mockService.CreateUser(params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	mockService.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockService := &MockUserService{}

	newEmail := "updated@example.com"
	params := UpdateUserParams{
		Email: &newEmail,
	}

	expectedUser := GetUsersQueryRow{
		ID:       1,
		Username: "testuser",
		Email:    "updated@example.com",
		UserType: "user",
	}

	mockService.On("UpdateUser", 1, params).Return(expectedUser, nil)

	user, err := mockService.UpdateUser(1, params)

	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", user.Email)

	mockService.AssertExpectations(t)
}
