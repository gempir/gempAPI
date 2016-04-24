package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"bufio"
	"os"
	"github.com/labstack/echo"
	"compress/gzip"
	"strings"
	"math/rand"
	"io/ioutil"
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

func getLastGlobalLogs(c echo.Context) error {
	limit, err    := strconv.Atoi(c.Param("limit"))
	if err != nil {
		limit = 1
	} else if limit > 500 {
		limit = 500
	}
	username := c.Param("username")
	username  = strings.ToLower(username)
	month    := time.Now().Month()
	year     := time.Now().Year()


	var lines []string

	file := fmt.Sprintf(logsfile + "%d/%s/%s.txt", year, month, username)
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
		errJSON.Error = "error reading logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	txt := ""

	for i := len(lines)-1; i >= 0; i--  {
		line := lines[i]
		if limit == 0 {
			break
		}
		split := strings.Split(line, "[|]")
 		txt += "[" + split[0] + " UTC] "
		txt += "[" + split[1] + "] "
		txt += split[2] + ": "
		txt += split[3] + "\r\n"
		limit--
	}

	return c.String(http.StatusOK, txt)
}

func getDatedGlobalLogs(c echo.Context) error {
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := strings.ToLower(c.Param("username"))

	file := fmt.Sprintf(logsfile + "%s/%s/%s.txt", year, month, username)
	var lines []string

	log.Debug(file)
	f, err := os.Open(file)
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
		errJSON.Error = "error reading logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	txt := ""

	for i := len(lines)-1; i >= 0; i--  {
		line := lines[i]
		split := strings.Split(line, "[|]")
		txt += "[" + split[0] + " UTC] "
		txt += "[" + split[1] + "] "
		txt += split[2] + ": "
		txt += split[3] + "\r\n"
	}

	return c.String(http.StatusOK, txt)
}

func getDatedChannelLogs(c echo.Context) error {
	channel := "#" + strings.ToLower(c.Param("channel"))
	year := c.Param("year")
	month := strings.Title(c.Param("month"))
	username := strings.ToLower(c.Param("username"))

	file := fmt.Sprintf(logsfile + "%s/%s/%s.txt", year, month, username)
	var lines []string

	log.Debug(file)
	f, err := os.Open(file)
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
		errJSON.Error = "error reading logs"
		return c.JSON(http.StatusNotFound, errJSON)
	}

	txt := ""

	for i := len(lines)-1; i >= 0; i--  {
		line := lines[i]
		split := strings.Split(line, "[|]")
		if split[1] != channel {
			continue
		}
		txt += "[" + split[0] + " UTC] "
		txt += split[2] + ": "
		txt += split[3] + "\r\n"
	}

	return c.String(http.StatusOK, txt)
}


func getLastChannelLogs(c echo.Context) error {
	limit, err    := strconv.Atoi(c.Param("limit"))
	if err != nil {
		limit = 1
	} else if limit > 500 {
		limit = 500
	}
	channel  := c.Param("channel")
	channel   = "#" + channel
	channel  = strings.ToLower(channel)
	username := c.Param("username")
	username  = strings.ToLower(username)
	month    := time.Now().Month()
	year     := time.Now().Year()


	var lines []string

	file := fmt.Sprintf(logsfile + "%d/%s/%s.txt", year, month, username)
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

	for i := len(lines)-1; i >= 0; i--  {
		line := lines[i]
		if limit == 0 {
			break
		}
		split := strings.Split(line, "[|]")
		if split[1] != channel {
			continue
		}
		txt += "[" + split[0] + " UTC] "
		txt += "[" + split[1] + "] "
		txt += split[2] + ": "
		txt += split[3] + "\r\n"
		limit--
	}

	return c.String(http.StatusOK, txt)
}

func getRandomquote(c echo.Context) error {
	username := c.Param("username")
	username  = strings.ToLower(username)

	var userlogs []string
	var lines    []string

	years, _ := ioutil.ReadDir(logsfile)
    for _, yearDir := range years {
		year := yearDir.Name()
        months, _ := ioutil.ReadDir(logsfile + year + "/")
		for _, monthDir := range months {
			month := monthDir.Name()
			path := fmt.Sprintf("%s%s/%s/%s.txt", logsfile, year, month, username)
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

	logs := new(Logs)

	ranNum := rand.Intn(len(lines)-1)

	split := strings.Split(lines[ranNum], "[|]")

	msg := new(Msg)
	msg.Channel = split[1]
	msg.Timestamp = split[0]
	msg.Username = split[2]
	msg.Message = split[3]
	timeObj, err := time.Parse(DateTime, msg.Timestamp)
	checkErr(err)
	msg.Duration = formatDiff(diff(timeObj, time.Now()))
	msg.UnixTimestamp = strconv.FormatInt(timeObj.Unix(), 10)
	logs.Messages = append(logs.Messages, *msg)

	return c.JSON(http.StatusOK, logs)
}


func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
