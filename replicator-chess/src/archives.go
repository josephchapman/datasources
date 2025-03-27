package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/josephchapman/datasources/cmn"
)

type ArchiveList struct {
	Archives []string `json:"archives"`
}

func (al ArchiveList) url(player string) (url string) {
	url = fmt.Sprintf("https://api.chess.com/pub/player/%s/games/archives", player)

	return url
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

type Player struct {
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
	White        Player `json:"white"`
	Black        Player `json:"black"`
	Eco          string `json:"eco"`
}

type Archive struct {
	Games []Game `json:"games"`
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

func (a Archive) eloTsdb(player string) (record string, err error) {
	line := ""
	// for each game
	for _, game := range a.Games {
		if strings.EqualFold(game.White.Username, player) {
			line = fmt.Sprintf("rating,player=%s,time_class=%s elo=%d %d000000000\n", game.White.Username, game.TimeClass, game.White.Rating, game.EndTime)
		} else if strings.EqualFold(game.Black.Username, player) {
			line = fmt.Sprintf("rating,player=%s,time_class=%s elo=%d %d000000000\n", game.Black.Username, game.TimeClass, game.Black.Rating, game.EndTime)
		}
		record += line
	}
	// check white/black = player, and use it
	// create string
	// append string to record string
	fmt.Println(record)
	return record, nil
}

func NewArchive(url string) (a Archive, err error) {
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
