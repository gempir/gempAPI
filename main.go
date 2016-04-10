package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var db, err = sql.Open("mysql", mysql)

func main() {
	log.SetOutput(os.Stdout)
	e := echo.New()
	e.Get("/v1/channel/:channel/user/:username/messages/random", getRandomquote)
	e.Get("/v1/channel/:channel/user/:username/messages/last", getLastMessage)
	e.Get("/v1/twitch/channel/:channel/user/:username/followage", getFollowage)
	e.Run(standard.New(":1323"))
}

func httpRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return contents, nil
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Quote basic response
type Quote struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
	Message   string `json:"message"`
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
		quote.Username = username
		quote.Message = message
	}

	log.Println(quote)

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

	log.Println(quote)

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}
