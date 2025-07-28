package queries

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) CreateMessage(params CreateMessageParams) (GetMessagesQueryRow, error) {
	args := m.Called(params)
	return args.Get(0).(GetMessagesQueryRow), args.Error(1)
}

func (m *MockMessageService) GetMessages() ([]GetMessagesQueryRow, error) {
	args := m.Called()
	return args.Get(0).([]GetMessagesQueryRow), args.Error(1)
}

func TestCreateMessage(t *testing.T) {
	mockService := &MockMessageService{}

	params := CreateMessageParams{
		UserID:  1,
		Content: "Test message",
	}

	expectedMessage := GetMessagesQueryRow{
		ID:      1,
		UserID:  1,
		Content: "Test message",
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	mockService.On("CreateMessage", params).Return(expectedMessage, nil)

	message, err := mockService.CreateMessage(params)

	assert.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, 1, message.UserID)
	assert.Equal(t, "Test message", message.Content)
	assert.True(t, message.CreatedAt.Valid)

	mockService.AssertExpectations(t)
}

func TestCreateMessageError(t *testing.T) {
	mockService := &MockMessageService{}

	params := CreateMessageParams{
		UserID:  999,
		Content: "Test",
	}

	mockService.On("CreateMessage", params).Return(GetMessagesQueryRow{}, errors.New("user not found"))

	_, err := mockService.CreateMessage(params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	mockService.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	mockService := &MockMessageService{}

	expectedMessages := []GetMessagesQueryRow{
		{
			ID:      1,
			UserID:  1,
			Content: "First message",
			CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		},
		{
			ID:      2,
			UserID:  2,
			Content: "Second message",
			CreatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		},
	}

	mockService.On("GetMessages").Return(expectedMessages, nil)

	messages, err := mockService.GetMessages()

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "First message", messages[0].Content)
	assert.Equal(t, "Second message", messages[1].Content)

	mockService.AssertExpectations(t)
}
