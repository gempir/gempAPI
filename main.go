package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/gempbot")

func main() {
	connectDB()
	e := echo.New()
	e.Get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Get("/v1/:channel/:username", getRandomquote)

	defer e.Run(standard.New(":1323"))
}

func connectDB() {
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Quote basic response
type Quote struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func getRandomquote(c echo.Context) error {

	channel := c.Param("channel")
	username := c.Param("username")
	channel = "#" + channel

	rows, err := db.Query(`
	SELECT username, message
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
		var username string
		var message string
		err = rows.Scan(&username, &message)
		checkErr(err)
		quote.Username = username
		quote.Message = message
	}

	log.Println(quote)

	defer rows.Close()
	return c.JSON(http.StatusOK, quote)
}
