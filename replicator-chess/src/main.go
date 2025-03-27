package main

import (
	"context"
	"flag"
	"fmt"

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
	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	al, _ := NewArchiveList(*player)
	for _, archiveStr := range al.Archives {
		archive, err := NewArchive(archiveStr)
		if err != nil {
			err = fmt.Errorf("NewArchive(): %w", err)
			cmn.LoggedError(err)
		} else {
			log := fmt.Sprintf("Created Archive: %s", archiveStr)
			cmn.LogOut.Info(log)
		}

		record, _ := archive.eloTsdb(*player)

		err = writeAPI.WriteRecord(context.Background(), record)
		if err != nil {
			err = fmt.Errorf("writeAPI.WriteRecord(): %w", err)
			cmn.LoggedError(err)
		} else {
			log := fmt.Sprintf("Query written: %s", archiveStr)
			cmn.LogOut.Info(log)
		}
	}

	// Ensures background processes finishes
	client.Close()
}
