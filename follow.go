package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"io/ioutil"

	"github.com/labstack/echo"
)

const (
	// UTCtimestamp for time.Parse
	UTCtimestamp = "2006-01-02T15:04:05-07:00"
	// DateTime for time.Parse out of the db
	DateTime = "2006-01-2 15:04:05"
)

// TwitchFollowage json for the followage response from twitch
type TwitchFollowage struct {
	CreatedAt string `json:"created_at"`
}

// Followage basic json for the followage reponse
type Followage struct {
	CreatedAt     string `json:"created_at"`
	Timestamp     string `json:"timestamp"`
	Duration      string `json:"duration"`
	UnixTimestamp string `json:"unix_timestamp"`
}

func getFollowage(c echo.Context) error {
	channel := c.Param("channel")
	username := c.Param("username")

	response, _ := httpRequest("https://api.twitch.tv/kraken/users/" + username + "/follows/channels/" + channel)
	var twitchFollowage TwitchFollowage
	json.Unmarshal(response, &twitchFollowage)

	followage := new(Followage)
	followage.CreatedAt = twitchFollowage.CreatedAt

	timestamp, err := time.Parse(UTCtimestamp, followage.CreatedAt)
	if err != nil {
		fail := new(ErrorJSON)
		fail.Error = "not following"
		return c.JSON(404, fail)
	}
	followage.Timestamp = timestamp.Format("2006-01-02 15:04:05")
	followage.UnixTimestamp = strconv.FormatInt(timestamp.Unix(), 10)

	followage.Duration = formatDiff(diff(timestamp, time.Now()))

	return c.JSON(http.StatusOK, followage)
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

func formatDiff(years, months, days, hours, mins, secs int) string {
	since := ""
	if years > 0 {
		switch years {
		case 1:
			since += fmt.Sprintf("%d year ", years)
			break
		default:
			since += fmt.Sprintf("%d years ", years)
			break
		}
	}
	if months > 0 {
		switch months {
		case 1:
			since += fmt.Sprintf("%d month ", months)
			break
		default:
			since += fmt.Sprintf("%d months ", months)
			break
		}
	}
	if days > 0 {
		switch days {
		case 1:
			since += fmt.Sprintf("%d day ", days)
			break
		default:
			since += fmt.Sprintf("%d days ", days)
			break
		}
	}
	if hours > 0 {
		switch hours {
		case 1:
			since += fmt.Sprintf("%d hour ", hours)
			break
		default:
			since += fmt.Sprintf("%d hours ", hours)
			break
		}
	}
	if mins > 0 && days == 0 && months == 0 && years == 0 {
		switch mins {
		case 1:
			since += fmt.Sprintf("%d min ", mins)
			break
		default:
			since += fmt.Sprintf("%d mins ", mins)
			break
		}
	}
	if secs > 0 && days == 0 && months == 0 && years == 0 && hours == 0 {
		switch secs {
		case 1:
			since += fmt.Sprintf("%d sec ", secs)
			break
		default:
			since += fmt.Sprintf("%d secs ", secs)
			break
		}
	}
	return strings.TrimSpace(since)
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}
