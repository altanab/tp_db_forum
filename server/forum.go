package server

import (
	"db_forum/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateForum(c echo.Context) error{
	newForum := models.Forum{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&newForum)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := SelectUserByNickname(newForum.User)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	newForum.User = user.Nickname
	forum, err := InsertForum(newForum)
	if err != nil {
		forum, err := SelectForumBySlug(newForum.Slug)
		if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		return c.JSON(http.StatusConflict, forum)
	}
	return c.JSON(http.StatusCreated, forum)
}

func GetForumDetails(c echo.Context) error {
	slug := c.Param("slug")
	forum, err := SelectForumBySlug(slug)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find forum\n",
			},
			)
	}
	return c.JSON(http.StatusOK, forum)
}


