package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// GetMessagesQueryRow represents a row in the messages query result
// It includes fields for the message ID, user ID, content, and creation timestamp.
type GetMessagesQueryRow struct {
	ID        int                `db:"id"`
	UserID    int                `db:"user_id"`
	Content   string             `db:"content"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
}

type CreateMessageParams struct {
	UserID  int    `db:"user_id"`
	Content string `db:"content"`
}

// GetMessages retrieves all messages from the database.
// It returns a slice of GetMessagesQueryRow and an error if any occurs.
func GetMessages() ([]GetMessagesQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	rows, err := conn.Query(context.TODO(), `
		SELECT id, user_id, content, created_at 
		FROM public.messages
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []GetMessagesQueryRow = []GetMessagesQueryRow{}
	for rows.Next() {
		var message GetMessagesQueryRow
		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func CreateMessage(params CreateMessageParams) (GetMessagesQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	var message GetMessagesQueryRow
	err := conn.QueryRow(context.TODO(), `
		INSERT INTO public.messages (user_id, content) 
		VALUES ($1, $2) 
		RETURNING id, user_id, content, created_at
	`, params.UserID, params.Content).Scan(
		&message.ID,
		&message.UserID,
		&message.Content,
		&message.CreatedAt,
	)
	if err != nil {
		return GetMessagesQueryRow{}, err
	}

	return message, nil
}

func GetMessagesByUser(userID int32) ([]GetMessagesQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	rows, err := conn.Query(context.TODO(), `
		SELECT id, user_id, content, created_at
		FROM public.messages
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []GetMessagesQueryRow = []GetMessagesQueryRow{}
	for rows.Next() {
		var message GetMessagesQueryRow
		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
