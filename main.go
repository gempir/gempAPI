package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/op/go-logging"
	"gopkg.in/redis.v3"
	"net/http"
	"os"
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

	e := echo.New()
	e.Get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Get("/user/:username/last", getLastMessage)
	e.Get("/channel/:channel/user/:username", getDatedChannelLogs)
	e.Get("/channel/:channel/user/:username/:year/:month", getDatedChannelLogs)
	e.Get("/channel/:channel/user/:username/random", getRandomquote)
	e.Get("/user/:username", getUser)
	e.Get("/twitch/followage/channel/:channel/user/:username", getFollowage)

	log.Info("starting webserver on 1323")
	e.Run(standard.New(webserverPort))
}
