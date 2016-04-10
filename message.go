package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Quote basic response
type Quote struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Duration  string `json:"duration"`
}

func getLastMessage(c echo.Context) error {
	channel := c.Param("channel")
	username := c.Param("username")
	channel = "#" + channel

	rows, err := db.Query("SELECT channel, timestamp, username, message  FROM chatlogs WHERE channel = ? AND username = ? ORDER BY timestamp DESC LIMIT 1", channel, username)
	checkErr(err)

	quote := new(Quote)

	for rows.Next() {
		var channel string
		var timestamp string
		var username string
		var message string
		err = rows.Scan(&channel, &timestamp, &username, &message)
		checkErr(err)
		quote.Channel = channel
		quote.Timestamp = timestamp
		timeparsed, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		quote.Duration = formatDiff(diff(timeparsed, time.Now()))
		quote.Username = username
		quote.Message = message
	}

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}

func getRandomquote(c echo.Context) error {

	channel := c.Param("channel")
	username := c.Param("username")
	channel = "#" + channel

	rows, err := db.Query(`
        SELECT channel, timestamp, username, message
        FROM chatlogs AS r1 JOIN
           (SELECT CEIL(RAND() *
                (SELECT MAX(id) FROM chatlogs)) AS id)
                AS r2
        WHERE r1.id >= r2.id
        AND channel = ?
        AND username = ?
        ORDER BY r1.id ASC
        LIMIT 1`, channel, username)
	checkErr(err)

	quote := new(Quote)

	for rows.Next() {
		var channel string
		var timestamp string
		var username string
		var message string
		err = rows.Scan(&channel, &timestamp, &username, &message)
		checkErr(err)
		quote.Channel = channel
		quote.Timestamp = timestamp
		quote.Username = username
		quote.Message = message
	}

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}
