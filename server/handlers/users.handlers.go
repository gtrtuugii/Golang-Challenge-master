package handlers

import (
	"main/queries"
	"net/http"
	"strconv"

	"slices"

	"github.com/gin-gonic/gin"
)

type UserQueries interface {
	GetUsers() ([]queries.GetUsersQueryRow, error)
	CreateUser(params queries.CreateUserParams) (queries.GetUsersQueryRow, error)
	UpdateUser(params queries.UpdateUserParams) (queries.GetUsersQueryRow, error)
}

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required"`
	UserType string  `json:"user_type" binding:"required"`
	Nickname *string `json:"nickname,omitempty"`
}

type UserResponse struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	UserType     string  `json:"userType"`
	Nickname     *string `json:"nickname,omitempty"`
	MessageCount int32   `json:"message_count"`
}

type UpdateUserParams struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	UserType *string `json:"user_type,omitempty"`
	Nickname *string `json:"nickname"` // Note: no omitempty to allow explicit null
}

func GetUsers(c *gin.Context) {

	userRows, err := queries.GetUsers()

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	var users []UserResponse = []UserResponse{}
	for _, row := range userRows {

		var nickname *string = nil
		if row.Nickname.Valid {
			nickname = &row.Nickname.String
		}

		users = append(users, UserResponse{
			ID:           row.ID,
			Username:     row.Username,
			Email:        row.Email,
			UserType:     row.UserType,
			Nickname:     nickname,
			MessageCount: row.MessageCount,
		})
	}

	c.JSON(http.StatusOK, GetUsersResponse{
		Users: users,
	})
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username is required",
		})
		return
	}
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}
	if req.UserType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User type is required",
		})
		return
	}

	params := queries.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		UserType:     req.UserType,
		Nickname:     req.Nickname,
		MessageCount: 0, // New user starts with 0 messages
	}

	user, err := queries.CreateUser(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	// Convert nickname from pgtype.Text to *string
	var nickname *string = nil
	if user.Nickname.Valid {
		nickname = &user.Nickname.String
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
		Nickname: nickname,
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateUser(c *gin.Context) {
	userIDStr := c.Param("user_id") // Use c.Param, not c.Params.Get
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate user_type if provided
	if req.UserType != nil {
		validTypes := []string{"UTYPE_USER", "UTYPE_ADMIN", "UTYPE_MODERATOR"}
		valid := slices.Contains(validTypes, *req.UserType)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user type. Must be one of: UTYPE_USER, UTYPE_ADMIN, UTYPE_MODERATOR",
			})
			return
		}
	}

	// Map handler's UpdateUserParams to queries.UpdateUserParams
	updateParams := queries.UpdateUserParams{
		Username: req.Username,
		Email:    req.Email,
		UserType: req.UserType,
		Nickname: req.Nickname,
	}

	user, err := queries.UpdateUser(int32(userID), updateParams)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user: " + err.Error(),
		})
		return
	}

	// Convert nickname from pgtype.Text to *string
	var nickname *string = nil
	if user.Nickname.Valid {
		nickname = &user.Nickname.String
	}

	response := UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		UserType:     user.UserType,
		Nickname:     nickname,
		MessageCount: user.MessageCount,
	}

	c.JSON(http.StatusOK, response)
}
