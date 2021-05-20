package server

import (
	"context"
	"db_forum/models"
	"github.com/jackc/pgconn"
)

func InsertForum(forum models.Forum) (models.Forum, error) {
	err := models.DBConn.QueryRow(
		context.Background(),
		"INSERT INTO forums (title, username, slug) VALUES ($1, $2, $3) RETURNING *;",
		forum.Title,
		forum.User,
		forum.Slug,
		).Scan(
			&forum.Title,
			&forum.User,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads,
			)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "forums_pkey" {
				return models.Forum{}, models.AlreadyExists{
					Message: "Forum already exists",
				}
			} else if pgerr.ConstraintName == "forums_username_fkey" {
				return models.Forum{}, models.UserIDError{
					Message: "Can't find user\n",
				}
			}
		}
		return models.Forum{}, err
	}

	return forum, nil
}

func SelectForumBySlug(slug string) (models.Forum, error) {
	forum := models.Forum{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM forums WHERE LOWER(slug)=LOWER($1) LIMIT 1;",
		slug,
	).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
		)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

