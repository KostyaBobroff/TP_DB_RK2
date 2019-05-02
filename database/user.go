package database

import (
	// "bytes"
	// "github.com/jackc/pgx"	
	"github.com/bozaro/tech-db-forum/generated/models"
	// "fmt"
)

const (
	createUserSQL = `
		INSERT
		INTO users ("nickname", "fullname", "email", "about")
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING
	`
	getUserByNicknameOrEmailSQL = `
		SELECT "nickname", "fullname", "email", "about"
		FROM users
		WHERE "email" = $1 OR "nickname" = $2
	`
	getUserSQL = `
		SELECT "nickname", "fullname", "email", "about"
		FROM users
		WHERE "nickname" = $1
	`
	updateUserSQL = `
		UPDATE users
		SET fullname = coalesce(nullif($2, ''), fullname),
			about    = coalesce(nullif($3, ''), about),
			email    = coalesce(nullif($4, ''), email)
		WHERE "nickname" = $1
		RETURNING fullname, about, email, nickname
	`
)

// /user/{nickname}/create Создание нового пользователя
func CreateUserDB(u *models.User) (*models.Users, error)  {
	rows, err := DB.pool.Exec(
		createUserSQL,
		&u.Nickname,
		&u.Fullname,
		&u.Email,
		&u.About,
	)
	if err != nil {
		return nil, err
	}

	// if it returns 0 - user existed, else user was created
	if rows.RowsAffected() == 0 {
		users := models.Users{}
		queryRows, err := DB.pool.Query(getUserByNicknameOrEmailSQL, &u.Email, &u.Nickname)
		defer queryRows.Close()

		if err != nil {
			return nil, err
		}

		for queryRows.Next() {
			user := models.User{}
			queryRows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			users = append(users, &user)
		}
		return &users, UserIsExist
	}

	return nil, nil
}

// /user/{nickname}/profile Получение информации о пользователе
func GetUserDB(nickname string) (*models.User, error) {
	user := models.User{}

	err := DB.pool.QueryRow(getUserSQL, nickname).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil {
		return nil, UserNotFound
	}

	return &user, nil
}

// /user/{nickname}/profile Изменение данных о пользователе
func UpdateUserDB(user *models.User) error {
	err := DB.pool.QueryRow(updateUserSQL,
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	).Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
		&user.Nickname,
	)

	if err != nil {
		if ErrorCode(err) == pgxOK {
			return UserUpdateConflict
		}
		return UserNotFound
	}

	return nil
}