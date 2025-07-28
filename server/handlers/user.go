package handlers

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
	Nickname *string `json:"nickname"`
}
