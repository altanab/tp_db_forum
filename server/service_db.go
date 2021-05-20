package server

import (
	"context"
	"db_forum/models"
)

func ClearDB() error {
	_, err := models.DBConn.Exec(
		context.Background(),
		"TRUNCATE users, forums, threads, posts, votes, forum_users;",
		)
	return err
}

func GetDBInfo() (models.Status, error) {
	var status models.Status
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users;",
		).Scan(&status.User)
	if err != nil {
		return models.Status{}, err
	}
	err = models.DBConn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM forums;",
	).Scan(&status.Forum)
	if err != nil {
		return models.Status{}, err
	}
	err = models.DBConn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM threads;",
	).Scan(&status.Thread)
	if err != nil {
		return models.Status{}, err
	}
	err = models.DBConn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM posts;",
	).Scan(&status.Post)
	if err != nil {
		return models.Status{}, err
	}
	return status, nil
}