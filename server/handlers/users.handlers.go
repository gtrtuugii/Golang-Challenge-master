package handlers

import (
	"main/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserQueries interface {
	GetUsers() ([]queries.GetUsersQueryRow, error)
	CreateUser(params queries.CreateUserParams) (queries.GetUsersQueryRow, error)
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
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	UserType string  `json:"userType"`
	Nickname *string `json:"nickname,omitempty"`
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
			ID:       row.ID,
			Username: row.Username,
			Email:    row.Email,
			UserType: row.UserType,
			Nickname: nickname,
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
		Username: req.Username,
		Email:    req.Email,
		UserType: req.UserType,
		Nickname: req.Nickname,
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
