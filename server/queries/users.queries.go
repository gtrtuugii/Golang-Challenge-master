package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetUsersQueryRow struct {
	ID                 int         `db:"id"`
	Username           string      `db:"username"`
	Email              string      `db:"email"`
	UserType           string      `db:"user_type"`
	Nickname           pgtype.Text `db:"nickname"`
	PermissionBitfield string      `db:"permission_bitfield"`
}

type CreateUserParams struct {
	Username string  `db:"username"`
	Email    string  `db:"email"`
	UserType string  `db:"user_type"`
	Nickname *string // Using pointer for optional field
}

func GetUsers() ([]GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	// Duplicates prevented by adding DISTINCT on the query below:
	rows, err := conn.Query(context.TODO(), `
		SELECT DISTINCT
			public.users.id, 
			public.users.username, 
			public.users.email,
			utype.type_key as "user_type",
			public.users.nickname,
			utype.permission_bitfield::text as "permission_bitfield"
		from 
			public.users 
		left join 
			user_types utype
		on
			utype.type_key = public.users.user_type
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
		RETURNING id, username, email, user_type, nickname
	`, params.Username, params.Email, params.UserType, nickname).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.UserType,
		&user.Nickname,
	)
	if err != nil {
		return GetUsersQueryRow{}, err
	}

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
