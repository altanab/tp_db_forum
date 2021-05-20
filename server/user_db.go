package server

import (
	"context"
	"db_forum/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

//TODO: scan null
func InsertUser(user models.User) (models.User, error) {
	err := models.DBConn.QueryRow(
		context.Background(),
		"INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING *;",
		user.Nickname,
		user.Fullname,
		user.Email,
		user.About,
		).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.Email,
			&user.About,
			)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func SelectUsers(nickname, email string) ([]models.User, error) {
	users := make([]models.User, 0, 0)
	rows, err := models.DBConn.Query(
		context.Background(),
		"SELECT * FROM users WHERE LOWER(nickname)=LOWER($1) OR LOWER(email)=LOWER($2) LIMIT 2;",
		nickname,
		email,
	)

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.Email,
			&user.About,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func SelectUserByNickname(nickname string) (models.User, error) {
	user := models.User{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE LOWER(nickname)=LOWER($1) LIMIT 1;",
		nickname,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
		)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UpdateDBUser(userUpdate models.User) (models.User, error) {
	user := models.User{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"UPDATE users SET fullname=$1, about=$2, email=$3 WHERE LOWER(nickname)=LOWER($4) RETURNING *",
		userUpdate.Fullname,
		userUpdate.About,
		userUpdate.Email,
		userUpdate.Nickname,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
	)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_email_key" {
				return models.User{}, models.AlreadyExists{
					Message: "Email already exists",
				}
			}
		}
		return models.User{}, err
	}
	return user, nil
}


func GetUsersForumBySlug(slug string, limit int, since string, desc bool) ([]models.User, error) {
	_, err := SelectForumBySlug(slug)
	if err != nil {
		return []models.User{}, err
	}

	users := make([]models.User, 0, 0)
 	var rows pgx.Rows

	if desc {
		if since != "" {
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users "+
					"AS f JOIN users AS u ON LOWER(f.forum_user)=LOWER(u.nickname) "+
					"WHERE LOWER(f.forum)=LOWER($1) AND LOWER(u.nickname) < LOWER($2) "+
					"ORDER BY u.nickname DESC "+
					"LIMIT $3;",
				slug,
				since,
				limit,
			)
		} else {
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users "+
					"AS f JOIN users AS u ON LOWER(f.forum_user)=LOWER(u.nickname) "+
					"WHERE LOWER(f.forum)=LOWER($1) "+
					"ORDER BY u.nickname DESC "+
					"LIMIT $2;",
				slug,
				limit,
			)

		}
	} else {
		rows, err = models.DBConn.Query(
			context.Background(),
			"SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users " +
				"AS f JOIN users AS u ON LOWER(f.forum_user)=LOWER(u.nickname) " +
				"WHERE LOWER(f.forum)=LOWER($1) AND LOWER(u.nickname) > LOWER($2) " +
				"ORDER BY u.nickname " +
				"LIMIT $3",
			slug,
			since,
			limit,
		)
	}
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.Email,
			&user.About,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}