package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/josephchapman/datasources/cmn"
)

type ArchiveList struct {
	Archives []string `json:"archives"`
}

func (al ArchiveList) url(player string) (url string) {
	url = fmt.Sprintf("https://api.chess.com/pub/player/%s/games/archives", player)

	return url
}

func (al ArchiveList) current() (current string) {
	// Get the last element of the Archives slice
	current = al.Archives[len(al.Archives)-1]

	// TODO rework this to confirm that the archive is the latest chronologically, not just the last element in the slice
	return current
}

func (al ArchiveList) printToConsole() (err error) {
	data, err := json.MarshalIndent(al, "", "  ")
	if err != nil {
		err = fmt.Errorf("json.MarshalIndent(): %w", err)
		return cmn.LoggedError(err)
	}
	fmt.Println(string(data))
	return nil
}

func NewArchiveList(player string) (al ArchiveList, err error) {
	url := al.url(player)
	apiData, err := cmn.QueryAPI(url)
	if err != nil {
		err = fmt.Errorf("cmn.QueryAPI(): %w", err)
		return ArchiveList{}, cmn.LoggedError(err)
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(apiData)
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return ArchiveList{}, cmn.LoggedError(err)
	}

	// Unmarshal the JSON data into the ArchiveList struct
	err = json.Unmarshal(jsonData, &al)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return ArchiveList{}, cmn.LoggedError(err)
	}

	return al, cmn.LoggedError(err)
}

type Color struct {
	Rating   int    `json:"rating"`
	Result   string `json:"result"`
	Id       string `json:"@id"`
	Username string `json:"username"`
	Uuid     string `json:"uuid"`
}

type Game struct {
	Url          string `json:"url"`
	Pgn          string `json:"pgn"`
	TimeControl  string `json:"time_control"`
	EndTime      int    `json:"end_time"`
	Rated        bool   `json:"rated"`
	Tcn          string `json:"tcn"`
	Uuid         string `json:"uuid"`
	InitialSetup string `json:"initial_setup"`
	Fen          string `json:"fen"`
	StartTime    int    `json:"start_time"`
	TimeClass    string `json:"time_class"`
	Rules        string `json:"rules"`
	White        Color  `json:"white"`
	Black        Color  `json:"black"`
	Eco          string `json:"eco"`
}

type Archive struct {
	Url    string `json:"url"`
	Player string `json:"player"`
	Year   string `json:"year"`
	Month  string `json:"month"`
	Games  []Game `json:"games"`
}

func (a Archive) printToConsole() (err error) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		err = fmt.Errorf("json.MarshalIndent(): %w", err)
		return cmn.LoggedError(err)
	}
	fmt.Println(string(data))
	return nil
}

func (a Archive) eloTsdb() (record string, err error) {
	line := ""
	for _, game := range a.Games {
		if strings.EqualFold(game.White.Username, a.Player) {
			line = fmt.Sprintf("rating,player=%s,time_class=%s elo=%d %d000000000\n", a.Player, game.TimeClass, game.White.Rating, game.EndTime)
		} else if strings.EqualFold(game.Black.Username, a.Player) {
			line = fmt.Sprintf("rating,player=%s,time_class=%s elo=%d %d000000000\n", a.Player, game.TimeClass, game.Black.Rating, game.EndTime)
		}
		record += line
	}
	fmt.Println(record)
	return record, nil
}

func NewArchive(url string) (a Archive, err error) {
	a.Url = url

	player, year, month, err := playerYearMonth(url)
	if err != nil {
		err = fmt.Errorf("playerYearMonth(): %w", err)
		return Archive{}, cmn.LoggedError(err)
	}

	a.Player = player
	a.Year = year
	a.Month = month

	apiData, err := cmn.QueryAPI(url)
	if err != nil {
		err = fmt.Errorf("cmn.QueryAPI(): %w", err)
		return Archive{}, cmn.LoggedError(err)
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(apiData)
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return Archive{}, cmn.LoggedError(err)
	}

	// Unmarshal the JSON data into the ArchiveList struct
	err = json.Unmarshal(jsonData, &a)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return Archive{}, cmn.LoggedError(err)
	}

	return a, cmn.LoggedError(err)
}

type Record map[string]Player
type Player map[string]Year
type Year []string

func (r Record) printToConsole() {
	encodedData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}

	fmt.Println(string(encodedData))
}

func playerYearMonth(archiveUrl string) (player string, yearStr string, monthStr string, err error) {
	// expected format of archiveDataUrl string:
	//   https://api.chess.com/pub/player/PLAYERNAME/games/2025/02
	//                                    ^^^^^^^^^^       ^^^^ ^^
	//                                    player           year month
	//   ^^^^^   ^^^^^^^^^^^^^ ^^^ ^^^^^^ ^^^^^^^^^^ ^^^^^ ^^^^ ^^
	//   0       2             3   4      5          6     7    8
	parts := strings.Split(archiveUrl, "/")

	player = parts[5]
	yearStr = parts[7]
	monthStr = parts[8]

	// Check if yearStr is a four-digit string
	if len(yearStr) != 4 || !isDigitsOnly(yearStr) {
		err = fmt.Errorf("invalid year format: %s", yearStr)
		return "", "", "", cmn.LoggedError(err)
	}

	// Check if monthStr is a two-digit string
	if len(monthStr) != 2 || !isDigitsOnly(monthStr) {
		err = fmt.Errorf("invalid month format: %s", monthStr)
		return "", "", "", cmn.LoggedError(err)
	}

	return player, yearStr, monthStr, nil
}

// Helper function to check if a string contains only digits
func isDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
