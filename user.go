package main

import (
    "net/http"
	"github.com/labstack/echo"
)

type User struct {
    Username      string  `json:"username"`
	Lines         float64 `json:"lines"`
}

func getUser(c echo.Context) error {
    username := c.Param("username")
    user := new(User)
    user.Username = username
    result := rclient.ZScore("user:lines", username)
    user.Lines, _ = result.Result()

    return c.JSON(http.StatusOK, user)
}
