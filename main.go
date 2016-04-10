package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/op/go-logging"
)

var (
	db, err = sql.Open("mysql", mysql)
	log     = logging.MustGetLogger("example")
	format  = logging.MustStringFormatter(
		`%{color}[%{time:2006-01-02 15:04:05}] [%{level:.4s}] %{color:reset}%{message}`,
	)
)

// ErrorJSON simple json for default error response
type ErrorJSON struct {
	Error string `json:"Error"`
}

func main() {
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

	e := echo.New()
	e.Get("/v1/channel/:channel/user/:username/messages/random", getRandomquote)
	e.Get("/v1/channel/:channel/user/:username/messages/last", getLastMessage)
	e.Get("/v1/user/:username/messages/last", getLastGlobalMessage)
	e.Get("/v1/twitch/followage/channel/:channel/user/:username", getFollowage)

	log.Info("starting webserver on 1323")
	e.Run(standard.New(":1323"))
}

func httpRequest(url string) ([]byte, error) {
	log.Debugf("httpRequest %s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return contents, nil
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
