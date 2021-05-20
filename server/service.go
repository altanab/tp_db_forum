package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Clear(c echo.Context) error{
	err := ClearDB()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "DataBase successfully cleared")
}

func GetInfo(c echo.Context) error{
	dbInfo, err := GetDBInfo()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dbInfo)

}
