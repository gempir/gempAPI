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

// Msg struct to define a simple message
type Msg struct {
	Channel       string `json:"channel"`
	Username      string `json:"username"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
	UnixTimestamp string `json:"unix_timestamp"`
	Duration      string `json:"duration"`
}

func getCurrentChanneLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year    := strconv.Itoa(time.Now().Year())
	month   := time.Now().Month().String()
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	redirectURL := fmt.Sprintf("/channel/%s/user/%s/%s/%s", channel, username, year, month)
	return c.Redirect(303, redirectURL)
}

func getDatedChannelLogs(c echo.Context) error {
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))

	if year == "" || month == "" {
		year = strconv.Itoa(time.Now().Year())
		month = time.Now().Month().String()
	}

	file := fmt.Sprintf(logsfile+"%s/%s/%s/%s.txt", channel, year, month, username)
	log.Debug(file)

	return c.File(file)
}

func getLastMessage(c echo.Context) error {
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))
	results, err := rclient.HGet("user:lastmessage", username).Result()
	if err != nil {
		log.Error(err)
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	split := strings.Split(results, "[|]")

	msg := new(Msg)
	msg.Channel = split[1]
	msg.Timestamp = split[0]
	msg.Username = split[2]
	msg.Message = split[3]

	timeObj, err := time.Parse(DateTime, msg.Timestamp)

	checkErr(err)

	msg.Duration = formatDiff(diff(timeObj, time.Now()))

	msg.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)

	return c.JSON(http.StatusOK, msg)

}

func getRandomquote(c echo.Context) error {
	username := c.Param("username")
	username = strings.ToLower(strings.TrimSpace(username))
	channel := strings.ToLower(c.Param("channel"))
	channel = strings.TrimSpace(channel)

	var userlogs []string
	var lines []string

	years, _ := ioutil.ReadDir(logsfile + channel)
	for _, yearDir := range years {
		year := yearDir.Name()
		months, _ := ioutil.ReadDir(logsfile + channel + "/" + year + "/")
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
	if len(userlogs) < 1 {
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
	if len(lines) < 1 {
		errJSON := new(ErrorJSON)
		errJSON.Error = "error finding logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	ranNum := rand.Intn(len(lines))
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
