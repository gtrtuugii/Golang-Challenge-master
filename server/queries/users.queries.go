package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// GetUsers retrieves all users from the database with their associated permissions and message counts.
// Returns:
//   - []GetUsersQueryRow: Slice of user records with permissions and message counts.
//   - error: Database error if query fails.
func GetUsers() ([]GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

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

	var users []GetUsersQueryRow
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

	return users, nil
}

// CreateUser inserts a new user into the database and returns the created record.
// Params:
//   - params: User details (username, email, type, optional nickname).
//
// Returns:
//   - GetUsersQueryRow: The newly created user with permissions.
//   - error: Database error if insertion fails.
func CreateUser(params CreateUserParams) (GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	var nickname pgtype.Text
	if params.Nickname != nil && *params.Nickname != "" {
		nickname = pgtype.Text{String: *params.Nickname, Valid: true}
	} else {
		nickname = pgtype.Text{Valid: false}
	}

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

	err = conn.QueryRow(context.TODO(), `
		SELECT permission_bitfield::text 
		FROM public.user_types 
		WHERE type_key = $1
	`, user.UserType).Scan(&user.PermissionBitfield)
	if err != nil {
		return GetUsersQueryRow{}, err
	}

	user.MessageCount = 0
	return user, nil
}

// UpdateUser modifies an existing user's fields (username, email, type, or nickname).
// Params:
//   - userID: ID of the user to update.
//   - params: Fields to update (nil fields are ignored).
//
// Returns:
//   - GetUsersQueryRow: Updated user record with permissions.
//   - error: Database error or "no fields to update" if params are empty.
func UpdateUser(userID int, params UpdateUserParams) (GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

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
	if params.Nickname != nil {
		var nickname pgtype.Text
		if *params.Nickname != "" {
			nickname = pgtype.Text{String: *params.Nickname, Valid: true}
		} else {
			nickname = pgtype.Text{Valid: false}
		}
		setParts = append(setParts, fmt.Sprintf("nickname = $%d", argCount))
		args = append(args, nickname)
		argCount++
	}

	if len(setParts) == 0 {
		return GetUsersQueryRow{}, fmt.Errorf("no fields to update")
	}

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

	err = conn.QueryRow(context.TODO(), `
		SELECT permission_bitfield::text 
		FROM public.user_types 
		WHERE type_key = $1
	`, user.UserType).Scan(&user.PermissionBitfield)
	if err != nil {
		return GetUsersQueryRow{}, fmt.Errorf("failed to get permission bitfield: %w", err)
	}

	user.MessageCount = 0 // Optional: can query actual count
	return user, nil
}
