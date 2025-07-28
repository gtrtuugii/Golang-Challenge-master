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
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(params interface{}) (MockUser, error) {
	args := m.Called(params)
	return args.Get(0).(MockUser), args.Error(1)
}

func (m *MockUserService) GetUsers() ([]MockUser, error) {
	args := m.Called()
	return args.Get(0).([]MockUser), args.Error(1)
}

type MockUser struct {
	ID       int
	Username string
	Email    string
	UserType string
	Nickname pgtype.Text
}

func setupUserTestRouter(mockService *MockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", CreateUserMock(mockService))
	router.GET("/users", GetUsersMock(mockService))
	return router
}

func CreateUserMock(mockService *MockUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if req.Username == "" || req.Email == "" || req.UserType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		user, err := mockService.CreateUser(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var nickname *string
		if user.Nickname.Valid {
			nickname = &user.Nickname.String
		}

		c.JSON(http.StatusCreated, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			UserType: user.UserType,
			Nickname: nickname,
		})
	}
}

func GetUsersMock(mockService *MockUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := mockService.GetUsers()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var responseUsers []UserResponse
		for _, user := range users {
			var nickname *string
			if user.Nickname.Valid {
				nickname = &user.Nickname.String
			}

			responseUsers = append(responseUsers, UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				UserType: user.UserType,
				Nickname: nickname,
			})
		}

		c.JSON(http.StatusOK, GetUsersResponse{Users: responseUsers})
	}
}

func TestCreateUser(t *testing.T) {
	mockService := &MockUserService{}
	router := setupUserTestRouter(mockService)

	nickname := "Test User"
	payload := CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		UserType: "user",
		Nickname: &nickname,
	}

	mockUser := MockUser{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		UserType: "user",
		Nickname: pgtype.Text{String: "Test User", Valid: true},
	}

	mockService.On("CreateUser", payload).Return(mockUser, nil)

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")

	mockService.AssertExpectations(t)
}

func TestCreateUserValidation(t *testing.T) {
	mockService := &MockUserService{}
	router := setupUserTestRouter(mockService)

	payload := CreateUserRequest{
		Username: "",
		Email:    "test@example.com",
		UserType: "user",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}

func TestGetUsers(t *testing.T) {
	mockService := &MockUserService{}
	router := setupUserTestRouter(mockService)

	mockUsers := []MockUser{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
			UserType: "user",
			Nickname: pgtype.Text{String: "User One", Valid: true},
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
			UserType: "admin",
			Nickname: pgtype.Text{Valid: false},
		},
	}

	mockService.On("GetUsers").Return(mockUsers, nil)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user1")
	assert.Contains(t, w.Body.String(), "user2")

	mockService.AssertExpectations(t)
}
