package queries

import "github.com/jackc/pgx/v5/pgtype"

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
