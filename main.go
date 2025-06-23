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

// input-output.json
type InputRequest struct {
	From   string  `json:"From"`
	To     string  `json:"To"`
	Weight float64 `json:"Weight"`
}

type PossiblePrice struct {
	RateID string  `json:"RateID"`
	Price  float64 `json:"Price"`
}

func main() {
	fmt.Println("Hello, Flip!")
}
