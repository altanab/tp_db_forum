package server

import (
	"db_forum/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateUser(c echo.Context) error{
	nickname := c.Param("nickname")
	if nickname == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty nickname")
	}
	newUser := models.User{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser.Nickname = nickname
	user, err := InsertUser(newUser)
	if err != nil {
		users, err := SelectUsers(nickname, newUser.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusConflict, users)
	}
	return c.JSON(http.StatusCreated, user)
}

func GetUserProfile(c echo.Context) error{
	nickname := c.Param("nickname")
	user, err := SelectUserByNickname(nickname)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				Message: "Can't find user\n",
			},
		)
	}
	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error{
	nickname := c.Param("nickname")
	userUpdate := models.User{}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&userUpdate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userUpdate.Nickname = nickname

	user, err := UpdateDBUser(userUpdate)
	if err != nil {
		switch err.(type) {
		case models.AlreadyExists:
			return c.JSON(http.StatusConflict, err)
		default:
			return c.JSON(
				http.StatusNotFound,
				models.UserIDError{
					Message: "Can't find user\n",
				},
			)
		}
	}
	return c.JSON(http.StatusOK, user)
}


func GetForumUsers(c echo.Context) error{
	slug := c.Param("slug")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}
	since := c.QueryParam("since")
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	users, err := GetUsersForumBySlug(slug, limit, since, desc)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.UserIDError{
				"Can't find user\n",
			},
			)
	}
	return c.JSON(http.StatusOK, users)
}