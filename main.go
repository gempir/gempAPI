package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"gopkg.in/redis.v3"
)

var (
	rclient *redis.Client
	log     = logging.MustGetLogger("gempAPI")
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

	rclient = redis.NewClient(&redis.Options{
		Addr:     redisaddress,
		Password: redispass, // no password set
		DB:       0,         // use default DB
	})

	channel := NewChannel()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/user/:username/last", getLastMessage)
	e.GET("/channel/:channel/user/:username", getCurrentChanneLogs)
	e.GET("/channel/:channel/user/:username/:year/:month", getDatedChannelLogs)
	e.GET("/channel/:channel/user/:username/random", getRandomquote)
	e.GET("/user/:username", getUser)
	e.GET("/twitch/followage/channel/:channel/user/:username", getFollowage)
	e.GET("/channel/:channel/commands", channel.getCommands)
	e.GET("/channel/:channel/levels", channel.getLevels)

	log.Info("starting webserver on 1323")
	e.Logger.Fatal(e.Start(webserverPort))
}
