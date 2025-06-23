package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

type InputOutputPair struct {
	Input  InputRequest    `json:"Input"`
	Output []PossiblePrice `json:"Output"`
}

// internal data structures
type Zone struct {
	Name      string
	Postcodes map[string]struct{} // Using a map for O(1) postcode lookup
}

type Rate struct {
	ID        string
	MaxWeight float64
	Cost      float64
	FromZone  string
	ToZone    string
}

type ApplicationState struct {
	Zones map[string]Zone
	Rates []Rate
}

func loadEvents(filename string) (*ApplicationState, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading events file: %w", err)
	}

	var events []Event
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling events JSON: %w", err)
	}

	appState := &ApplicationState{
		Zones: make(map[string]Zone),
		Rates: []Rate{},
	}

	for _, event := range events {
		switch event.EventType {
		case "ZoneDefined":
			var zoneData ZoneDefinedData
			if err := json.Unmarshal(event.Data, &zoneData); err != nil {
				return nil, fmt.Errorf("error unmarshaling ZoneDefinedData: %w", err)
			}
			postcodeMap := make(map[string]struct{})
			for _, postcode := range zoneData.Postcodes {
				postcodeMap[postcode] = struct{}{}
			}
			appState.Zones[zoneData.Name] = Zone{
				Name:      zoneData.Name,
				Postcodes: postcodeMap,
			}
		case "RateDefined":
			var rateData RateDefinedData
			if err := json.Unmarshal(event.Data, &rateData); err != nil {
				return nil, fmt.Errorf("error unmarshaling RateDefinedData: %w", err)
			}
			appState.Rates = append(appState.Rates, Rate(rateData))
		default:
			fmt.Printf("Warning: Unknown event type encountered: %s\n", event.EventType)
		}
	}
	return appState, nil
}

func main() {
	_, err := loadEvents("testdata/events.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load events: %v\n", err)
		os.Exit(1)
	}

	// Load input-output.json
	inputOutputBytes, err := os.ReadFile("testdata/input-output.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input-output file: %v\n", err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(strings.NewReader(string(inputOutputBytes)))

	var inputOutputs []InputOutputPair
	for {
		var pair InputOutputPair
		err := decoder.Decode(&pair)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding input-output pair: %v\n", err)
			os.Exit(1)
		}
		inputOutputs = append(inputOutputs, pair)
		fmt.Printf("Input: %+v\n", pair.Input)
		fmt.Printf("Output: %+v\n", pair.Output)
	}
}
