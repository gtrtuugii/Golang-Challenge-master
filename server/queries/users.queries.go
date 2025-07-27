package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

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
	Username     string  `db:"username"`
	Email        string  `db:"email"`
	UserType     string  `db:"user_type"`
	Nickname     *string // Using pointer for optional field
	MessageCount int32   `db:"message_count"`
}

// UpdateUserParams struct for PATCH operations
type UpdateUserParams struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	UserType *string `json:"user_type,omitempty"`
	Nickname *string `json:"nickname"` // Note: no omitempty to allow explicit null
}

func GetUsers() ([]GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	// Duplicates prevented by adding DISTINCT on the query below:
	rows, err := conn.Query(context.TODO(), `
		SELECT DISTINCT
			u.id, 
			u.username, 
			u.email, 
			u.user_type, 
			u.nickname,
			ut.permission_bitfield::text,
			COUNT(m.id) as message_count
		FROM public.users u
		LEFT JOIN public.user_types ut ON ut.type_key = u.user_type
		LEFT JOIN public.messages m ON m.user_id = u.id
		GROUP BY u.id, u.username, u.email, u.user_type, u.nickname, ut.permission_bitfield
		ORDER BY u.id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []GetUsersQueryRow = []GetUsersQueryRow{}
	for rows.Next() {
		var user GetUsersQueryRow
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.UserType,
			&user.Nickname,
			&user.PermissionBitfield,
			&user.MessageCount,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func CreateUser(params CreateUserParams) (GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	var nickname pgtype.Text
	// TODO: Nickname validation
	// if params.Nickname != nil && *params.Nickname != "" {
	// 	nickname = pgtype.Text{String: *params.Nickname, Valid: true}
	// } else {
	// 	nickname = pgtype.Text{Valid: false}
	// }

	var user GetUsersQueryRow
	err := conn.QueryRow(context.TODO(), `
		INSERT INTO public.users (username, email, user_type, nickname) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, username, email, user_type, nickname, message_count
	`, params.Username, params.Email, params.UserType, nickname).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.UserType,
		&user.Nickname,
		&user.MessageCount,
	)
	if err != nil {
		return GetUsersQueryRow{}, err
	}

	// New user has no messages
	user.MessageCount = 0

	// Get permission bitfield for the response
	err = conn.QueryRow(context.TODO(), `
		SELECT permission_bitfield::text 
		FROM user_types 
		WHERE type_key = $1
	`, user.UserType).Scan(&user.PermissionBitfield)
	if err != nil {
		return GetUsersQueryRow{}, err
	}

	return user, nil
}

func UpdateUser(userID int32, params UpdateUserParams) (GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if params.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argCount))
		args = append(args, *params.Username)
		argCount++
	}

	if params.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *params.Email)
		argCount++
	}

	if params.UserType != nil {
		setParts = append(setParts, fmt.Sprintf("user_type = $%d", argCount))
		args = append(args, *params.UserType)
		argCount++
	}

	// Handle nickname - can be set to null explicitly
	if params.Nickname != nil {
		var nickname pgtype.Text
		if *params.Nickname != "" {
			nickname = pgtype.Text{String: *params.Nickname, Valid: true}
		} else {
			nickname = pgtype.Text{Valid: false} // NULL
		}
		setParts = append(setParts, fmt.Sprintf("nickname = $%d", argCount))
		args = append(args, nickname)
		argCount++
	}

	if len(setParts) == 0 {
		return GetUsersQueryRow{}, fmt.Errorf("no fields to update")
	}

	// Add userID as the last parameter
	args = append(args, userID)

	query := fmt.Sprintf(`
		UPDATE public.users 
		SET %s 
		WHERE id = $%d 
		RETURNING id, username, email, user_type, nickname
	`, strings.Join(setParts, ", "), argCount)

	var user GetUsersQueryRow
	err := conn.QueryRow(context.TODO(), query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.UserType,
		&user.Nickname,
	)
	if err != nil {
		return GetUsersQueryRow{}, err
	}

	// Get permission bitfield
	err = conn.QueryRow(context.TODO(), `
		SELECT permission_bitfield::text 
		FROM user_types 
		WHERE type_key = $1
	`, user.UserType).Scan(&user.PermissionBitfield)
	if err != nil {
		return GetUsersQueryRow{}, fmt.Errorf("failed to get permission bitfield: %w", err)
	}

	// You might want to get actual message count from database
	user.MessageCount = 0 // Or query from messages table

	return user, nil
}
