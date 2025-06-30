package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type GeoPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"` // [lat, long]
}

type DriverLocation struct {
	DriverID string   `json:"driverId"`
	Location GeoPoint `json:"location"`
}

type Drivers struct {
	Drivers []DriverLocation `json:"drivers"`
}

const (
	csvFilePath = "./Coordinates.csv"
	batchSize   = 10000
	url         = "http://127.0.0.1:10001/v1/drivers"

	// authenticated = true
	jwtKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjp0cnVlfQ.fZZU3vlwkRPOXmoCdR2RjvCRiY5h7LLjRBWgYBpRJjM"

	// authenticated = false
	// jwtKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoZW50aWNhdGVkIjpmYWxzZX0.ThlNUHzVvYc8Zo2Y8sTqb2DtvqgLp70IlI2LFvXvMkA"
)

func main() {
	f, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("err openin csv file: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("err reading all records: %v", err)
	}

	// Skip header [Latitude Longitude]
	records = records[1:]

	var batch Drivers
	for _, row := range records {
		lat, _ := strconv.ParseFloat(row[0], 64)
		long, _ := strconv.ParseFloat(row[1], 64)

		driver := DriverLocation{
			DriverID: "driver-" + uuid.NewString(),
			Location: GeoPoint{
				Type:        "Point",
				Coordinates: [2]float64{lat, long},
			},
		}

		batch.Drivers = append(batch.Drivers, driver)

		if len(batch.Drivers) >= batchSize {
			sendBatch(batch)
			batch.Drivers = batch.Drivers[:0]
		}

	}

	if len(batch.Drivers) > 0 {
		sendBatch(batch)
	}
	log.Println("Import completed.")
}

func sendBatch(batch Drivers) {
	dat, err := json.Marshal(batch)
	if err != nil {
		log.Printf("err marshalling batch: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(dat))
	if err != nil {
		log.Printf("err creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// TODO: add Authorization header here
	req.Header.Set("Authorization", "Bearer "+jwtKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("api error: %v", resp.Status)
	} else {
		log.Printf("Sent batch of %d drivers", len(batch.Drivers))
	}

}
