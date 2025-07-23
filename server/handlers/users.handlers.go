package handlers

import (
	"main/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
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
