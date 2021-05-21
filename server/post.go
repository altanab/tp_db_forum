package server

import (
	"db_forum/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func GetPostDetails(c echo.Context) error{
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	post, err := SelectPostById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.UserIDError{
			Message: "Can't find post\n",
		})
	}
	postDetails := map[string]interface{}{}
	postDetails["post"] = post

	related := strings.Split(c.QueryParam("related"), ",")
	for _, param := range related{
		switch param {
		case "user":
			user, err := SelectUserByNickname(post.Author)
			if err != nil {
				return c.JSON(http.StatusNotFound, models.UserIDError{
					Message: "Can't find user\n",
				})
			}
			postDetails["author"] = user
		case "thread":
			thread, err := SelectThreadById(post.Thread)
			if err != nil {
				return c.JSON(http.StatusNotFound, models.UserIDError{
					Message: "Can't find thread\n",
				})
			}
			postDetails["thread"] = thread
		case "forum":
			forum, err := SelectForumBySlug(post.Forum)
			if err != nil {
				return c.JSON(http.StatusNotFound, models.UserIDError{
					Message: "Can't find forum\n",
				})
			}
			postDetails["forum"] = forum
		}
	}

	return c.JSON(http.StatusOK, postDetails)
}

func UpdatePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	postUpdate := models.PostUpdate{
		Id : id,
	}

	defer c.Request().Body.Close()
	err = json.NewDecoder(c.Request().Body).Decode(&postUpdate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	post, err := UpdatePostById(postUpdate)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.UserIDError{
			Message: "Can't find post\n",
		})
	}
	return c.JSON(http.StatusOK, post)
}

func CreatePost(c echo.Context) error{
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

	var posts = make([]models.Post, 0, 0)

	defer c.Request().Body.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&posts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(posts) == 0 {
		return c.JSON(http.StatusCreated, posts)
	}

	posts, err = InsertPosts(posts, thread)
	if err != nil {
		switch err.(type) {
		case models.UserIDError:
			return c.JSON(http.StatusNotFound, err)
		default:
			return c.JSON(
				http.StatusConflict,
				models.UserIDError{
					Message: "Can't find parent\n",
				},
			)
		}
	}
	return c.JSON(http.StatusCreated, posts)
}
