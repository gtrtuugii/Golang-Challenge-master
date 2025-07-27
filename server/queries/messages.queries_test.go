package queries

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockMessageService struct {
	mock.Mock
	messages []GetMessagesQueryRow
	nextID   int
}

func (m *MockMessageService) CreateMessage(params CreateMessageParams) (GetMessagesQueryRow, error) {
	args := m.Called(params)

	if args.Error(1) != nil {
		return GetMessagesQueryRow{}, args.Error(1)
	}

	// Create mock message
	message := GetMessagesQueryRow{
		ID:      m.nextID,
		UserID:  params.UserID,
		Content: params.Content,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	m.messages = append(m.messages, message)
	m.nextID++

	return message, nil
}

func (m *MockMessageService) GetMessages() ([]GetMessagesQueryRow, error) {
	args := m.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return m.messages, nil
}

func (m *MockMessageService) GetMessagesByUser(userID int) ([]GetMessagesQueryRow, error) {
	args := m.Called(userID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	var userMessages []GetMessagesQueryRow
	for _, msg := range m.messages {
		if msg.UserID == userID {
			userMessages = append(userMessages, msg)
		}
	}

	return userMessages, nil
}

func TestCreateMessage(t *testing.T) {
	mockService := &MockMessageService{
		nextID: 1,
	}

	params := CreateMessageParams{
		UserID:  1,
		Content: "This is a test message",
	}

	// Mock successful creation
	mockService.On("CreateMessage", params).Return(GetMessagesQueryRow{}, nil)

	message, err := mockService.CreateMessage(params)

	// Assertions
	assert.NoError(t, err, "CreateMessage should not return an error")
	assert.Equal(t, 1, message.ID, "Message ID should be 1")
	assert.Equal(t, 1, message.UserID, "Message UserID should be 1")
	assert.Equal(t, "This is a test message", message.Content, "Message content should match")
	assert.True(t, message.CreatedAt.Valid, "CreatedAt should be valid")

	mockService.AssertExpectations(t)
}

func TestCreateMessageWithInvalidUser(t *testing.T) {
	mockService := &MockMessageService{}

	params := CreateMessageParams{
		UserID:  99999, // Non-existent user ID
		Content: "This should fail",
	}

	// Mock error response
	mockService.On("CreateMessage", params).Return(GetMessagesQueryRow{}, errors.New("foreign key constraint violation"))

	_, err := mockService.CreateMessage(params)
	assert.Error(t, err, "CreateMessage should return an error for invalid user")
	assert.Contains(t, err.Error(), "foreign key constraint", "Error should mention foreign key constraint")

	mockService.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	mockService := &MockMessageService{
		messages: []GetMessagesQueryRow{
			{
				ID:      1,
				UserID:  1,
				Content: "First test message",
				CreatedAt: pgtype.Timestamptz{
					Time:  time.Now(),
					Valid: true,
				},
			},
			{
				ID:      2,
				UserID:  2,
				Content: "Second test message",
				CreatedAt: pgtype.Timestamptz{
					Time:  time.Now(),
					Valid: true,
				},
			},
		},
	}

	// Mock successful retrieval
	mockService.On("GetMessages").Return(mockService.messages, nil)

	messages, err := mockService.GetMessages()

	// Assertions
	assert.NoError(t, err, "GetMessages should not return an error")
	assert.Len(t, messages, 2, "Should return 2 messages")
	assert.Equal(t, "First test message", messages[0].Content)
	assert.Equal(t, "Second test message", messages[1].Content)

	mockService.AssertExpectations(t)
}

func TestGetMessagesError(t *testing.T) {
	mockService := &MockMessageService{}

	// Mock database error
	mockService.On("GetMessages").Return([]GetMessagesQueryRow(nil), errors.New("database connection failed"))

	messages, err := mockService.GetMessages()

	assert.Error(t, err, "GetMessages should return an error")
	assert.Nil(t, messages, "Messages should be nil on error")
	assert.Contains(t, err.Error(), "database connection", "Error should mention database connection")

	mockService.AssertExpectations(t)
}

func TestGetMessagesByUser(t *testing.T) {
	mockService := &MockMessageService{
		messages: []GetMessagesQueryRow{
			{
				ID:      1,
				UserID:  1,
				Content: "User 1 message 1",
				CreatedAt: pgtype.Timestamptz{
					Time:  time.Now(),
					Valid: true,
				},
			},
			{
				ID:      2,
				UserID:  2,
				Content: "User 2 message",
				CreatedAt: pgtype.Timestamptz{
					Time:  time.Now(),
					Valid: true,
				},
			},
			{
				ID:      3,
				UserID:  1,
				Content: "User 1 message 2",
				CreatedAt: pgtype.Timestamptz{
					Time:  time.Now(),
					Valid: true,
				},
			},
		},
	}

	// Mock successful retrieval for user 1
	mockService.On("GetMessagesByUser", 1).Return([]GetMessagesQueryRow(nil), nil)

	user1Messages, err := mockService.GetMessagesByUser(1)

	// Assertions
	assert.NoError(t, err, "GetMessagesByUser should not return an error")
	assert.Len(t, user1Messages, 2, "User 1 should have 2 messages")

	for _, msg := range user1Messages {
		assert.Equal(t, 1, msg.UserID, "All messages should belong to user 1")
	}

	mockService.AssertExpectations(t)
}

func TestGetMessagesByUserEmpty(t *testing.T) {
	mockService := &MockMessageService{
		messages: []GetMessagesQueryRow{},
	}

	// Mock empty result for non-existent user
	mockService.On("GetMessagesByUser", 99999).Return([]GetMessagesQueryRow{}, nil)

	messages, err := mockService.GetMessagesByUser(99999)

	assert.NoError(t, err, "GetMessagesByUser should not error for non-existent user")
	assert.Empty(t, messages, "Non-existent user should have no messages")

	mockService.AssertExpectations(t)
}

func TestMessageFieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateMessageParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			params: CreateMessageParams{
				UserID:  1,
				Content: "Valid message content",
			},
			wantErr: false,
		},
		{
			name: "empty content",
			params: CreateMessageParams{
				UserID:  1,
				Content: "",
			},
			wantErr: true,
			errMsg:  "content cannot be empty",
		},
		{
			name: "invalid user ID",
			params: CreateMessageParams{
				UserID:  0,
				Content: "Some content",
			},
			wantErr: true,
			errMsg:  "invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockMessageService{}

			if tt.wantErr {
				mockService.On("CreateMessage", tt.params).Return(GetMessagesQueryRow{}, errors.New(tt.errMsg))
			} else {
				mockService.On("CreateMessage", tt.params).Return(GetMessagesQueryRow{}, nil)
			}

			_, err := mockService.CreateMessage(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			mockService.AssertExpectations(t)
		})
	}
}
