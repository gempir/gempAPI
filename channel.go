package main

import (
    "github.com/labstack/echo"
	"net/http"
    "encoding/json"
    "strconv"
    "strings"
)

type Channel struct {

}

func NewChannel() *Channel {
	return &Channel {

	}
}

type Levels struct {
    Levels map[string]int `json:"levels"`
}

type Commands struct {
    Commands  map[string]Command `json:"commands"`
}

type Oddshots struct {
    Oddshots  []Oddshot `json:"oddshots"`
}

type Oddshot struct {
    Link string `json:"link"`
    Unixtimestamp int `json:"unixtimestamp"`
}

type Command struct {
    Name     string `json:"name"`
    Message  string `json:"message"`
    Cd       int `json:"cd"`
    Function string `json:"function"`
    Response bool `json:"response"`
    Level    int `json:"level"`
    Description string `json:"description"`
}


func (channel *Channel) getCommands(c echo.Context) error {
    current := c.Param("channel")
    current = "#" + current

    results, err := rclient.HGetAllMap(current + ":commands").Result()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusNotFound, "not found")
	}
    coms := new(Commands)
    coms.Commands = make(map[string]Command)
    for name, command := range results {
        var com Command
        err := json.Unmarshal([]byte(command), &com)
    	if err != nil {
            log.Notice(command)
    		log.Error(err)
    	}
        coms.Commands[name] = com
    }

    return c.JSON(http.StatusOK, coms)
}

func (channel *Channel) getLevels(c echo.Context) error {
    current := c.Param("channel")
    current = "#" + current

    results, err := rclient.HGetAllMap(current + ":levels").Result()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusNotFound, "not found")
	}
    lvls := new(Levels)
    lvls.Levels = make(map[string]int)
    for name, level := range results {
        lvls.Levels[name], err = strconv.Atoi(level)
        if err != nil {
            log.Error(err)
        }
    }

    return c.JSON(http.StatusOK, lvls)
}

func (channel *Channel) getOddshots(c echo.Context) error {
    current := c.Param("channel")
    current = "#" + current

    results, err := rclient.HGetAllMap(current + ":oddshots").Result()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusNotFound, "not found")
	}
    shots := new(Oddshots)
    for shot, timestamp := range results {
        timestamp = strings.TrimSpace(timestamp)
        oddshot := new(Oddshot)
        oddshot.Link = shot
        oddshot.Unixtimestamp, err = strconv.Atoi(timestamp)
        if err != nil {
            log.Error(err)
            continue
        }
        shots.Oddshots = append(shots.Oddshots, *oddshot)
    }
    return c.JSON(http.StatusOK, shots)
}
