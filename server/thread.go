package server

import (
	"db_forum/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)


func CreateThread(c echo.Context) error{
	forum, err := SelectForumBySlug(c.Param("slug"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	newThread := models.Thread{}

	defer c.Request().Body.Close()
	err = json.NewDecoder(c.Request().Body).Decode(&newThread)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	newThread.Forum = forum.Slug

	user, err := SelectUserByNickname(newThread.Author)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	newThread.Author = user.Nickname

	thread, err := InsertThread(newThread)
	if err != nil {
		thread, err = SelectThreadBySlug(newThread.Slug.String)
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusConflict, thread)
	}
	return c.JSON(http.StatusCreated, thread)
}

func GetThreadDetails(c echo.Context) error{
	thread := models.Thread{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		thread, err = SelectThreadBySlug(c.Param("id"))
	} else {
		thread, err = SelectThreadById(id)
	}
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find thread\n",
			},
		)
	}
	return c.JSON(http.StatusOK, thread)
}

func UpdateThread(c echo.Context) error{
	threadUpdate := models.ThreadUpdate{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&threadUpdate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thread := models.Thread{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		thread, err = UpdateDBThreadBySlug(c.Param("id"), threadUpdate)
	} else {
		thread, err = UpdateDBThreadById(id, threadUpdate)
	}

	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find thread\n",
			},
		)
	}
	return c.JSON(http.StatusOK, thread)
}

func GetThreadPosts(c echo.Context) error{
	var thread models.Thread
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		thread, err = SelectThreadBySlug(c.Param("id"))
	} else {
		thread, err = SelectThreadById(id)
	}
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Thread not found\n",
			},
		)
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	since, err := strconv.Atoi(c.QueryParam("since"))
	if err != nil {
		since = 0
	}
	var sort string
	switch c.QueryParam("sort") {
	case "tree":
		sort = "tree"
	case "parent_tree":
		sort = "parent_tree"
	default:
		sort = "flat"
	}

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	var posts []models.Post
	posts, err = GetThreadPostsById(thread.Id, limit, since, sort, desc)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Thread not found\n",
			},
		)
	}
	return c.JSON(http.StatusOK, posts)

}

func VoteThread(c echo.Context) error{
	vote := models.Vote{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&vote)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if vote.Voice != 1 && vote.Voice != -1 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			models.UserIDError{
				Message: "Can't find thread\n",
			},
			)
	}

	thread := models.Thread{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		thread, err = SelectThreadBySlug(c.Param("id"))
	} else {
		thread.Id = id
	}
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find thread\n",
			},
		)
	}

	err = InsertVote(thread.Id, vote)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find thread or user\n",
			},
		)
	}

	thread, err = SelectThreadById(thread.Id)

	return c.JSON(http.StatusOK, thread)
}


func GetForumThreads(c echo.Context) error{
	slug := c.Param("slug")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	sinceParam := c.QueryParam("since")
	var since time.Time
	if sinceParam != "" {
		since, err = time.Parse("2006-01-02T15:04:05Z", sinceParam)
		if err != nil {
			sinceParam = ""
		}
	}

	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	threads, err := GetThreadsForumBySlug(slug, limit, since, desc, sinceParam != "")
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find user\n",
			},
		)
	}
	return c.JSON(http.StatusOK, threads)
}