package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/josephchapman/datasources/cmn"
)

type runElement struct {
	Place int `json:"place"`
	Run   run `json:"run"`
}

type run struct {
	Id        string              `json:"id"`
	Weblink   string              `json:"weblink"`
	Game      string              `json:"game"`
	Category  string              `json:"category"`
	Videos    videos              `json:"videos"`
	Comment   string              `json:"comment"`
	Status    map[string]string   `json:"status"`
	Players   []map[string]string `json:"players"`
	Date      string              `json:"date"`
	Submitted string              `json:"submitted"`
	Times     times               `json:"times"`
	System    system              `json:"system"`
	Values    map[string]string   `json:"values"`
}

type videos struct {
	Links []link `json:"links"`
}

type link struct {
	Uri string `json:"uri"`
}

type times struct {
	Primary    string  `json:"primary"`
	Primary_t  float64 `json:"primary_t"`
	Realtime   string  `json:"realtime"`
	Realtime_t float64 `json:"realtime_t"`
}

type system struct {
	Platform string `json:"platform"`
	Emulated bool   `json:"emulated"`
	Region   string `json:"region"`
}

type leaderboard struct {
	Weblink  string            `json:"weblink"`
	Game     string            `json:"game"`
	Category string            `json:"category"`
	Values   map[string]string `json:"values"`
	Runs     []runElement      `json:"runs"`
}

type performed struct {
	Date    string `json:"date"`
	DaysAgo int    `json:"days_ago"`
}

type currentRecord struct {
	Game      string    `json:"game"`
	Category  string    `json:"category"`
	Player    string    `json:"player"`
	Time      string    `json:"time"`
	Performed performed `json:"performed"`
}

func (cr currentRecord) log() (err error) {
	cmn.LogOut.Info("Current record retrieved", slog.Any("current_record", cr))
	return err
}

func (l leaderboard) endpoint() (url string, err error) {
	baseURL := fmt.Sprintf("https://www.speedrun.com/api/v1/leaderboards/%s/category/%s?&top=1", l.Game, l.Category)

	if len(l.Values) > 0 {
		for key, value := range l.Values {
			baseURL += fmt.Sprintf("&var-%s=%s", key, value)
		}
	}

	return baseURL, nil
}

// update API Data
func (l *leaderboard) updateAPI() (err error) {
	// Get the endpoint from the location
	url, err := l.endpoint()
	if err != nil {
		err = fmt.Errorf("w.Location.endpoint(): %w", err)
		return cmn.LoggedError(err)
	}

	// Query the endpoint to receive updated data
	data, err := cmn.QueryAPI(url)
	if err != nil {
		err = fmt.Errorf("cmn.QueryAPI(): %w", err)
		return cmn.LoggedError(err)
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(data["data"])
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return cmn.LoggedError(err)
	}

	// Unmarshal the JSON data into the leaderboard struct
	err = json.Unmarshal(jsonData, &l)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return cmn.LoggedError(err)
	}

	return nil
}

func (l leaderboard) NewCurrentRecord() (cr currentRecord, err error) {
	g := game{
		Id: l.Game,
	}
	g.updateAPI()

	c := category{
		Id: l.Category,
	}
	c.updateAPI()

	u := user{
		Id: l.Runs[0].Run.Players[0]["id"],
	}
	u.updateAPI()

	// Convert the primary time to a human-readable format
	duration := time.Duration(l.Runs[0].Run.Times.Primary_t * float64(time.Second))
	humanReadableTime := fmt.Sprintf("%02dh %02dm %02ds", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)

	// Calculate the time since submission
	date, err := time.Parse("2006-01-02", l.Runs[0].Run.Date)
	if err != nil {
		err = fmt.Errorf("time.Parse: %w", err)
		return cr, cmn.LoggedError(err)
	}
	timeSince := time.Since(date)
	days := int(timeSince.Hours() / 24)
	timeSinceStr := days

	cr = currentRecord{
		Game:     g.Names["international"],
		Category: c.Name,
		Player:   u.Names["international"],
		Time:     humanReadableTime,
		Performed: performed{
			Date:    l.Runs[0].Run.Date,
			DaysAgo: timeSinceStr,
		},
	}

	return cr, nil
}
