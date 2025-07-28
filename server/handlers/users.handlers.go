package handlers

import (
	"main/queries"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUsers handles GET /users requests.
// Response:
//   - 200: JSON list of all users.
//   - 400: Error if database query fails.
func GetUsers(c *gin.Context) {
	userRows, err := queries.GetUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve users: " + err.Error()})
		return
	}

	users := make([]UserResponse, 0, len(userRows))
	for _, row := range userRows {
		var nickname *string
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

	c.JSON(http.StatusOK, GetUsersResponse{Users: users})
}

// CreateUser handles POST /users requests to create a new user.
// Validates required fields (username, email, user_type).
// Response:
//   - 201: JSON of the created user.
//   - 400: Error if validation or database insertion fails.
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if req.UserType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User type is required"})
		return
	}

	params := queries.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		UserType:     req.UserType,
		Nickname:     req.Nickname,
		MessageCount: 0,
	}

	user, err := queries.CreateUser(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	var nickname *string
	if user.Nickname.Valid {
		nickname = &user.Nickname.String
	}

	resp := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
		Nickname: nickname,
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateUser handles PATCH /users/:user_id requests.
// Validates:
//   - user_id as integer.
//   - user_type (if provided) must be a valid role.
//
// Response:
//   - 200: JSON of the updated user.
//   - 400: Error if input validation fails.
//   - 404: Error if user is not found.
func UpdateUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if req.UserType != nil {
		valid := slices.Contains([]string{"UTYPE_USER", "UTYPE_ADMIN", "UTYPE_MODERATOR"}, *req.UserType)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type. Must be one of: UTYPE_USER, UTYPE_ADMIN, UTYPE_MODERATOR"})
			return
		}
	}

	updateParams := queries.UpdateUserParams{
		Username: req.Username,
		Email:    req.Email,
		UserType: req.UserType,
		Nickname: req.Nickname,
	}

	user, err := queries.UpdateUser(userID, updateParams)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	var nickname *string
	if user.Nickname.Valid {
		nickname = &user.Nickname.String
	}

	resp := UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		UserType:     user.UserType,
		Nickname:     nickname,
		MessageCount: user.MessageCount,
	}

	c.JSON(http.StatusOK, resp)
}
