package main

import (
	"encoding/json"
	"fmt"
)

// events.json
type Event struct {
	EventType string          `json:"Event"`
	Data      json.RawMessage `json:"Data"` // Use RawMessage to unmarshal later based on EventType
}

type ZoneDefinedData struct {
	Name      string   `json:"Name"`
	Postcodes []string `json:"Postcodes"`
}

type RateDefinedData struct {
	ID        string  `json:"ID"`
	MaxWeight float64 `json:"MaxWeight"`
	Cost      float64 `json:"Cost"`
	FromZone  string  `json:"FromZone"`
	ToZone    string  `json:"ToZone"`
}

func main() {
	fmt.Println("Hello, Flip!")
}
