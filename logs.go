package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func getDatedChannelLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := strings.ToLower(c.Param("username"))

	file := fmt.Sprintf(logsfile+"%s/%s/%s/%s.txt", channel, year, month, username)
	log.Debug(file)

	return c.File(file)
}

func getLastChannelLogs(c echo.Context) error {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		limit = 1
	} else if limit > 500 {
		limit = 500
	}
	channel := c.Param("channel")
	channel = strings.ToLower(channel)
	username := c.Param("username")
	username = strings.ToLower(username)
	month := time.Now().Month()
	year := time.Now().Year()

	var lines []string

	file := fmt.Sprintf(logsfile+"%s/%d/%s/%s.txt", channel, year, month, username)
	log.Debug(file)
	f, err := os.Open(file)
	if err != nil {
		log.Error(err)
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

	}

	if err := scanner.Err(); err != nil {
		log.Error(scanner.Err())
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	txt := ""

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if limit == 0 {
			break
		}
		txt += line + "\r\n"
		limit--
	}

	return c.String(http.StatusOK, txt)
}

func getRandomquote(c echo.Context) error {
	username := c.Param("username")
	username = strings.ToLower(username)
	channel := strings.ToLower(c.Param("channel"))

	var userlogs []string
	var lines []string

	years, _ := ioutil.ReadDir(logsfile + channel)
	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(logsfile + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s%s/%s/%s/%s.txt", logsfile, channel, year, month, username)
			if _, err := os.Stat(path); err == nil {
				userlogs = append(userlogs, path)
			} else if _, err := os.Stat(path + ".gz"); err == nil {
				userlogs = append(userlogs, path)
			}
		}
	}

	if len(userlogs) == 0 {
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	file := userlogs[rand.Intn(len(userlogs))]
	log.Debug(file, len(userlogs))

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Error(err)
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}
	scanner := bufio.NewScanner(f)

	if strings.HasSuffix(file, ".gz") {
		gz, err := gzip.NewReader(f)
		scanner = bufio.NewScanner(gz)
		if err != nil {
			log.Error(err)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

	}

	if err := scanner.Err(); err != nil {
		log.Error(scanner.Err())
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}
	ranNum := rand.Intn(len(lines) - 1)
	line := lines[ranNum]
	log.Debug(line)
	lineSplit := strings.SplitN(line, "]", 2)
	return c.String(http.StatusOK, lineSplit[1])
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
