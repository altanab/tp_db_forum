package server

import (
	"db_forum/models"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
	"context"
)

func InsertPosts(posts []models.Post, thread models.Thread) ([]models.Post, error)  {
	sql := "INSERT INTO posts(parent, author, message, forum, thread) VALUES "

	var values []interface{}

	for i, post := range posts {
		sql += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d), ",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			)
		values = append(values, post.Parent, post.Author, post.Message, thread.Forum, thread.Id)
	}
	sql = strings.TrimSuffix(sql, ", ")
	sql += " RETURNING id, parent, author, message, is_edited, forum, thread, created"

	rows, err := models.DBConn.Query(
		context.Background(),
		sql,
		values...,
		)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "posts_author_fkey" {
				return nil, models.UserIDError{
					Message: "user doesn't exists",
				}
			}
		}
		return nil, err
	}

	defer rows.Close()
	insertedPosts := make([]models.Post, 0, 0)
	for rows.Next() {
		post := models.Post{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			)
		if err != nil {
			return nil, err
		}
		insertedPosts = append(insertedPosts, post)
	}

	if err = rows.Err(); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "posts_author_fkey" {
				return nil, models.UserIDError{
					Message: "user doesn't exists",
				}
			}
		}
		return nil, err
	}

	return insertedPosts, nil

}

func SelectPostById(id int) (models.Post, error) {
	var post models.Post
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE id=$1;",
		id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func UpdatePostById(postUpdate models.PostUpdate) (models.Post, error)  {
	var post models.Post
	err := models.DBConn.QueryRow(
		context.Background(),
		"UPDATE posts SET message=COALESCE(NULLIF($1, ''), message), is_edited=CASE WHEN $1='' OR message=$1 THEN FALSE ELSE TRUE END WHERE id=$2 RETURNING id, parent, author, message, is_edited, forum, thread, created",
		postUpdate.Message,
		postUpdate.Id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}