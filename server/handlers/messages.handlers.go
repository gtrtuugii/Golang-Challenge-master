package handlers

import (
	"net/http"
	"strconv"

	"main/queries"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	messageRows, err := queries.GetMessages()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve messages: " + err.Error(),
		})
		return
	}

	var messages []MessageResponse = []MessageResponse{}
	for _, row := range messageRows {
		messages = append(messages, MessageResponse{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, GetMessagesResponse{
		Messages: messages,
	})
}

func GetMessagesByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	messageRows, err := queries.GetMessagesByUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve messages: " + err.Error(),
		})
		return
	}

	var messages []MessageResponse = []MessageResponse{}
	for _, row := range messageRows {
		messages = append(messages, MessageResponse{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, GetMessagesResponse{
		Messages: messages,
	})
}

func CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Note: Since you're using binding:"required" tags, the validation
	// will be handled by Gin's ShouldBindJSON, so you can remove the manual checks
	// or keep them for additional validation

	params := queries.CreateMessageParams{
		UserID:  req.UserID,
		Content: req.Content,
	}

	message, err := queries.CreateMessage(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create message: " + err.Error(),
		})
		return
	}

	response := MessageResponse{
		ID:        message.ID,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}
