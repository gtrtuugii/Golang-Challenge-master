package handlers

type CreateMessageRequest struct {
	UserID  int    `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type GetMessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}
