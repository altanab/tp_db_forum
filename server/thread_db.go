package server

import (
	"db_forum/models"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)


func InsertThread(thread models.Thread) (models.Thread, error) {
	err := models.DBConn.QueryRow(
		context.Background(),
		"INSERT INTO threads (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING " +
			"id, title, author, forum, message, slug, created;",
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "threads_sluq_key" {
				return models.Thread{}, models.AlreadyExists{
					Message: "Thread already exists",
				}
			} else if pgerr.ConstraintName == "threads_author_fkey" ||  pgerr.ConstraintName == "threads_forum_fkey"{
				return models.Thread{}, models.UserIDError{
					Message: "Can't find user\n",
				}
			}
			return models.Thread{}, err
		}
	}
	return thread, nil
}

func SelectThreadBySlug(slug string) (models.Thread, error){
	thread := models.Thread{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM threads WHERE LOWER(slug)=LOWER($1) LIMIT 1;",
		slug,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func SelectThreadById(id int) (models.Thread, error) {
	thread := models.Thread{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM threads WHERE id=$1 LIMIT 1;",
		id,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func UpdateDBThreadById(id int, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	thread := models.Thread{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"UPDATE threads SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) WHERE id=$3 RETURNING *",
		threadUpdate.Title,
		threadUpdate.Message,
		id,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func UpdateDBThreadBySlug(slug string, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	thread := models.Thread{}
	err := models.DBConn.QueryRow(
		context.Background(),
		"UPDATE threads SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) WHERE slug=$3 RETURNING *",
		threadUpdate.Title,
		threadUpdate.Message,
		slug,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func GetThreadPostsById(id, limit, since int, sort string, desc bool) ([]models.Post, error) {
	posts := make([]models.Post, 0, 0)
	var rows pgx.Rows
	var err error
	switch {
	case desc && since != 0 :
		switch sort {
		case "flat":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 AND id < $2 ORDER BY id DESC LIMIT $3;",
				id,
				since,
				limit,
			)
		case "tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 AND " +
					"m_path < (SELECT m_path FROM posts WHERE id = $2) " +
					"ORDER BY m_path DESC LIMIT $3;",
				id,
				since,
				limit,
			)
		case "parent_tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE m_path[1] IN " +
					"(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL AND " +
					"m_path[1] < (SELECT m_path[1] FROM posts WHERE id = $2) " +
					"ORDER BY id DESC LIMIT $3) " +
					"ORDER BY m_path[1] DESC, m_path;",
				id,
				since,
				limit,
			)

		}
	case desc :
		switch sort {
		case "flat":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 ORDER BY id DESC LIMIT $2;",
				id,
				limit,
			)
		case "tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 " +
					"ORDER BY m_path DESC LIMIT $2;",
				id,
				limit,
			)
		case "parent_tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE m_path[1] IN " +
					"(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL " +
					"ORDER BY id DESC LIMIT $2) " +
					"ORDER BY m_path[1] DESC, m_path;",
				id,
				limit,
			)

		}
	case since != 0 :
		switch sort {
		case "flat":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 AND id > $2 ORDER BY id LIMIT $3;",
				id,
				since,
				limit,
			)
		case "tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 AND " +
					"m_path > (SELECT m_path FROM posts WHERE id = $2) " +
					"ORDER BY m_path LIMIT $3;",
				id,
				since,
				limit,
			)
		case "parent_tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE m_path[1] IN " +
					"(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL AND " +
					"m_path[1] > (SELECT m_path[1] FROM posts WHERE id = $2) " +
					"ORDER BY id LIMIT $3) " +
					"ORDER BY m_path;",
				id,
				since,
				limit,
			)

		}
	default:
		switch sort {
		case "flat":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 ORDER BY id LIMIT $2;",
				id,
				limit,
			)
		case "tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE thread=$1 " +
					"ORDER BY m_path LIMIT $2;",
				id,
				limit,
			)
		case "parent_tree":
			rows, err = models.DBConn.Query(
				context.Background(),
				"SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts WHERE m_path[1] IN " +
					"(SELECT id FROM posts WHERE thread=$1 AND parent IS NULL " +
					"ORDER BY id LIMIT $2) " +
					"ORDER BY m_path;",
				id,
				limit,
			)
		}
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()
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
		posts = append(posts, post)
	}
	return posts, nil
}

func InsertVote(id int, newVote models.Vote) error {
	vote := models.Vote{}
	var exists bool
	err := models.DBConn.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT * FROM votes WHERE nickname=$1 AND thread=$2 LIMIT 1);",
		newVote.Nickname,
		id,
	).Scan(
		&exists,
		)
	if err != nil {
		return err
	}

	if exists {
		err = models.DBConn.QueryRow(
			context.Background(),
			"UPDATE votes SET voice=$1 WHERE nickname=$2 and thread=$3 RETURNING *;",
			newVote.Voice,
			newVote.Nickname,
			id,
		).Scan(
			&vote.Nickname,
			&vote.Voice,
			&id,
			)

	} else {
		err = models.DBConn.QueryRow(
			context.Background(),
			"INSERT INTO votes(nickname, voice, thread)  VALUES($1, $2, $3) RETURNING *;",
			newVote.Nickname,
			newVote.Voice,
			id,
		).Scan(
			&vote.Nickname,
			&vote.Voice,
			&id,
		)
	}

	return err
}

func GetThreadsForumBySlug(slug string, limit int, since time.Time, desc bool, sinceParam bool) ([]models.Thread, error) {
	_, err := SelectForumBySlug(slug)
	if err != nil {
		return []models.Thread{}, err
	}

	threads := make([]models.Thread, 0, 0)
	var rows pgx.Rows

	if desc && sinceParam {
		rows, err = models.DBConn.Query(
			context.Background(),
			"SELECT * FROM threads " +
				"WHERE forum=$1 AND created <= $2 " +
				"ORDER BY created DESC " +
				"LIMIT $3;",
			slug,
			since,
			limit,
		)

	} else if sinceParam {
		rows, err = models.DBConn.Query(
			context.Background(),
			"SELECT * FROM threads " +
				"WHERE forum=$1 AND created >= $2 " +
				"ORDER BY created " +
				"LIMIT $3;",
			slug,
			since,
			limit,
		)
	} else if desc {
		rows, err = models.DBConn.Query(
			context.Background(),
			"SELECT * FROM threads " +
				"WHERE forum=$1 " +
				"ORDER BY created DESC " +
				"LIMIT $2;",
			slug,
			limit,
		)

	} else {
		rows, err = models.DBConn.Query(
			context.Background(),
			"SELECT * FROM threads " +
				"WHERE forum=$1 " +
				"ORDER BY created " +
				"LIMIT $2;",
			slug,
			limit,
		)
	}
	if err != nil {
		return threads, err
	}

	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		err = rows.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			return threads, err
		}
		threads = append(threads, thread)
	}
	return threads, nil

}