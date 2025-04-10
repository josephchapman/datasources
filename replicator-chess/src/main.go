package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/josephchapman/datasources/cmn"
)

const applicationName = "replicator-chess"

func main() {
	cmn.SetApplicationName(applicationName)

	player := flag.String("player", "", "Pull ratings for this player into a TSDB")
	flag.Parse()

	// create writeAPI for influx
	bucket := "my-bucket"
	org := "my-org"
	token := "my-super-secret-auth-token"
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	// Create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(url, token)

	// A record of which archives are already in the database
	var record = make(Record)

	queryAPI := client.QueryAPI(org)
	// Define query
	fluxQuery := fmt.Sprintf(`from(bucket: "%s")
	      |> range(start: 0)
        |> filter(fn: (r) => r["_measurement"] == "pull")
        |> keep(columns: ["player", "archive_year", "archive_month"])
        |> distinct()
        |> yield(name: "pulls")`, bucket)
	// Get QueryTableResult
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Type assertions, because the values are returned as interface{}
			player := result.Record().ValueByKey("player").(string)
			archiveYear := result.Record().ValueByKey("archive_year").(string)
			archiveMonth := result.Record().ValueByKey("archive_month").(string)

			if _, ok := record[player]; !ok {
				record[player] = make(Player)
			}
			if _, ok := record[player][archiveYear]; !ok {
				record[player][archiveYear] = make(Year, 0)
			}
			record[player][archiveYear] = append(record[player][archiveYear], archiveMonth)
		}
		// Check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	record.printToConsole()

	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	al, _ := NewArchiveList(*player)
	for _, archiveStr := range al.Archives {
		player, year, month, err := playerYearMonth(archiveStr)
		if err != nil {
			err = fmt.Errorf("playerYearMonth(): %w", err)
			cmn.LoggedError(err)
		} else {
			log := fmt.Sprintf("Archive found %s-%s-%s", player, year, month)
			cmn.LogOut.Info(log)
		}

		// Check if record[player][year][month] exists
		if record[player] != nil && record[player][year] != nil && contains(record[player][year], month) {
			log := fmt.Sprintf("Skipping %s-%s-%s, already in database", player, year, month)
			cmn.LogOut.Info(log)
		} else {
			log := fmt.Sprintf("Adding %s-%s-%s to database", player, year, month)
			cmn.LogOut.Info(log)

			// Create a new archive object
			archive, err := NewArchive(archiveStr)
			if err != nil {
				err = fmt.Errorf("NewArchive(): %w", err)
				cmn.LoggedError(err)
			} else {
				log := fmt.Sprintf("Archive created %s-%s-%s", player, year, month)
				cmn.LogOut.Info(log)
			}

			// Write the TSDB record
			entry, _ := archive.eloTsdb()
			err = writeAPI.WriteRecord(context.Background(), entry)
			if err != nil {
				err = fmt.Errorf("writeAPI.WriteRecord(): %w", err)
				cmn.LoggedError(err)
			} else {
				log := fmt.Sprintf("Query written: %s", archiveStr)
				cmn.LogOut.Info(log)
			}

			// if the archive is not currently being modified (it's complete), record the pull
			if archiveStr != al.current() {
				p := influxdb2.NewPoint("pull",
					map[string]string{"player": player, "archive_year": year, "archive_month": month},
					map[string]interface{}{"pulled": 1},
					time.Now())

				// Write point immediately
				err = writeAPI.WritePoint(context.Background(), p)
				if err != nil {
					err = fmt.Errorf("writeAPI.WriteRecord(): %w", err)
					cmn.LoggedError(err)
				} else {
					log := fmt.Sprintf("Archive complete. Pull recorded %s-%s-%s", player, year, month)
					cmn.LogOut.Info(log)
				}
			} else {
				log := fmt.Sprintf("Archive incomplete. Pull not recorded %s-%s-%s", player, year, month)
				cmn.LogOut.Info(log)
			}

		}
	}

	// Ensures background processes finishes
	client.Close()
}
