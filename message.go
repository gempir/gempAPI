package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// Quote basic response
type Quote struct {
	Channel       string `json:"channel"`
	Username      string `json:"username"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
	UnixTimestamp string `json:"unix_timestamp"`
	Duration      string `json:"duration"`
}

func getLastMessage(c echo.Context) error {
	channel := c.Param("channel")
	username := c.Param("username")
	channel = "#" + channel

	rows, err := db.Query("SELECT channel, timestamp, username, message  FROM gempLog WHERE channel = ? AND username = ? ORDER BY timestamp DESC LIMIT 1", channel, username)
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
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		quote.Duration = formatDiff(diff(timeObj, time.Now()))
		quote.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)
	}

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}

func getLastGlobalMessage(c echo.Context) error {
	username := c.Param("username")

	rows, err := db.Query("SELECT channel, timestamp, username, message  FROM gempLog WHERE username = ? ORDER BY timestamp DESC LIMIT 1", username)
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
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		quote.Duration = formatDiff(diff(timeObj, time.Now()))
		quote.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)
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
        FROM gempLog AS r1 JOIN
           (SELECT CEIL(RAND() *
                (SELECT MAX(id) FROM gempLog)) AS id)
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
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		quote.Duration = formatDiff(diff(timeObj, time.Now()))
		quote.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)
	}

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}

func getGlobalRandomquote(c echo.Context) error {
	username := c.Param("username")

	rows, err := db.Query(`
        SELECT channel, timestamp, username, message
        FROM gempLog AS r1 JOIN
           (SELECT CEIL(RAND() *
                (SELECT MAX(id) FROM gempLog)) AS id)
                AS r2
        WHERE r1.id >= r2.id
        AND username = ?
        ORDER BY r1.id ASC
        LIMIT 1`, username)
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
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		quote.Duration = formatDiff(diff(timeObj, time.Now()))
		quote.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)
	}

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}

// Logs basic logs response
type Logs struct {
	Messages []Msg `json:"messages"`
}

// Msg struct to define a simple message
type Msg struct {
	Channel       string `json:"channel"`
	Username      string `json:"username"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
	UnixTimestamp string `json:"unix_timestamp"`
	Duration      string `json:"duration"`
}

func getLogs(c echo.Context) error {

	channel := c.Param("channel")
	username := c.Param("username")
	limit := c.Param("limit")
	channel = "#" + channel
	limitInt, err := strconv.Atoi(limit)
	checkErr(err)

	if limitInt > 250 {
		limit = "250"
	}
	rows, err := db.Query(`
        SELECT channel, timestamp, username, message
        FROM gempLog
		WHERE channel = ?
		AND username = ?
		ORDER BY timestamp DESC
		LIMIT ?`, channel, username, limit)
	checkErr(err)

	logs := new(Logs)

	for rows.Next() {
		var channel string
		var timestamp string
		var username string
		var message string
		err = rows.Scan(&channel, &timestamp, &username, &message)
		checkErr(err)
		msg := new(Msg)
		msg.Channel = channel
		msg.Timestamp = timestamp
		msg.Username = username
		msg.Message = message
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		msg.Duration = formatDiff(diff(timeObj, time.Now()))
		msg.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)

		logs.Messages = append(logs.Messages, *msg)

	}

	defer rows.Close()
	return c.JSON(http.StatusOK, logs)
}

func getGlobalLogs(c echo.Context) error {
	username := c.Param("username")
	limit := c.Param("limit")
	limitInt, err := strconv.Atoi(limit)
	checkErr(err)

	if limitInt > 250 || limitInt < 1 || err != nil {
		limit = "250"
	}

	rows, err := db.Query(`
        SELECT channel, timestamp, username, message
        FROM gempLog
		WHERE username = ?
		ORDER BY timestamp DESC
		LIMIT ?`, username, limit)
	checkErr(err)

	logs := new(Logs)

	for rows.Next() {
		var channel string
		var timestamp string
		var username string
		var message string
		err = rows.Scan(&channel, &timestamp, &username, &message)
		checkErr(err)
		msg := new(Msg)
		msg.Channel = channel
		msg.Timestamp = timestamp
		msg.Username = username
		msg.Message = message
		timeObj, err := time.Parse(DateTime, timestamp)
		checkErr(err)
		msg.Duration = formatDiff(diff(timeObj, time.Now()))
		msg.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)

		logs.Messages = append(logs.Messages, *msg)

	}

	defer rows.Close()
	return c.JSON(http.StatusOK, logs)
}
