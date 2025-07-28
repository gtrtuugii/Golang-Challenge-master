package queries

import "github.com/jackc/pgx/v5/pgtype"

type GetUsersQueryRow struct {
	ID                 int         `db:"id"`
	Username           string      `db:"username"`
	Email              string      `db:"email"`
	UserType           string      `db:"user_type"`
	Nickname           pgtype.Text `db:"nickname"`
	PermissionBitfield string      `db:"permission_bitfield"`
	MessageCount       int32       `db:"message_count"`
}

type CreateUserParams struct {
	Username     string `db:"username"`
	Email        string `db:"email"`
	UserType     string `db:"user_type"`
	Nickname     *string
	MessageCount int32 `db:"message_count"`
}

type UpdateUserParams struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	UserType *string `json:"user_type,omitempty"`
	Nickname *string `json:"nickname"` // Nullable on purpose
}
