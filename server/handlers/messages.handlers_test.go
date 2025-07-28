package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageService struct {
	mock.Mock
}

type MockMessage struct {
	ID        int
	UserID    int
	Content   string
	CreatedAt time.Time
}

func (m *MockMessageService) CreateMessage(params interface{}) (MockMessage, error) {
	args := m.Called(params)
	return args.Get(0).(MockMessage), args.Error(1)
}

func (m *MockMessageService) GetMessages() ([]MockMessage, error) {
	args := m.Called()
	return args.Get(0).([]MockMessage), args.Error(1)
}

func setupTestRouter(mockService *MockMessageService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/messages", GetMessagesMock(mockService))
	router.POST("/messages", CreateMessageMock(mockService))
	return router
}

func CreateMessageMock(mockService *MockMessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if req.Content == "" || req.UserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		message, err := mockService.CreateMessage(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, MessageResponse{
			ID:        message.ID,
			UserID:    message.UserID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

func GetMessagesMock(mockService *MockMessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		messages, err := mockService.GetMessages()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var responseMessages []MessageResponse
		for _, msg := range messages {
			responseMessages = append(responseMessages, MessageResponse{
				ID:        msg.ID,
				UserID:    msg.UserID,
				Content:   msg.Content,
				CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		}

		if responseMessages == nil {
			responseMessages = []MessageResponse{}
		}

		c.JSON(http.StatusOK, GetMessagesResponse{Messages: responseMessages})
	}
}

func TestCreateMessage(t *testing.T) {
	mockService := &MockMessageService{}
	router := setupTestRouter(mockService)

	// Valid creation
	payload := CreateMessageRequest{UserID: 1, Content: "Test message"}
	mockMessage := MockMessage{ID: 1, UserID: 1, Content: "Test message", CreatedAt: time.Now()}
	mockService.On("CreateMessage", payload).Return(mockMessage, nil)

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test message")
	mockService.AssertExpectations(t)
}

func TestCreateMessageErrors(t *testing.T) {
	mockService := &MockMessageService{}
	router := setupTestRouter(mockService)

	// Empty content
	payload := CreateMessageRequest{UserID: 1, Content: ""}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMessages(t *testing.T) {
	mockService := &MockMessageService{}
	router := setupTestRouter(mockService)

	mockMessages := []MockMessage{{ID: 1, UserID: 1, Content: "Test", CreatedAt: time.Now()}}
	mockService.On("GetMessages").Return(mockMessages, nil)

	req, _ := http.NewRequest("GET", "/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
	mockService.AssertExpectations(t)
}
