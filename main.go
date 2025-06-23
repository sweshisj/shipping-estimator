package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// events.json
type Event struct {
	EventType string          `json:"Event"`
	Data      json.RawMessage `json:"Data"`
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

type CalculatedResult struct {
	Input  InputRequest    `json:"Input"`
	Output []PossiblePrice `json:"Output"`
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

func calculatePrices(req InputRequest, appState *ApplicationState) []PossiblePrice {
	var applicablePrices []PossiblePrice

	// 1. Find the zones for the given postcodes
	var fromZoneName, toZoneName string
	for _, zone := range appState.Zones {
		if _, ok := zone.Postcodes[req.From]; ok {
			fromZoneName = zone.Name
		}
		if _, ok := zone.Postcodes[req.To]; ok {
			toZoneName = zone.Name
		}
	}

	// If either postcode doesn't belong to a defined zone, no rates can apply
	if fromZoneName == "" || toZoneName == "" {
		return applicablePrices
	}

	// 2. Filter rates
	for _, rate := range appState.Rates {
		// Check weight
		if req.Weight > rate.MaxWeight {
			continue
		}

		// Check fromZone and toZone
		if rate.FromZone == fromZoneName && rate.ToZone == toZoneName {
			applicablePrices = append(applicablePrices, PossiblePrice{
				RateID: rate.ID,
				Price:  rate.Cost,
			})
		}
	}

	return applicablePrices
}

func main() {
	appState, err := loadEvents("testdata/events.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load events: %v\n", err)
		os.Exit(1)
	}

	inputRequestsBytes, err := os.ReadFile("testdata/input.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input.json: %v\n", err)
		os.Exit(1)
	}

	var inputRequests []InputRequest
	err = json.Unmarshal(inputRequestsBytes, &inputRequests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling input.json: %v\n", err)
		os.Exit(1)
	}

	var allCalculatedResults []CalculatedResult

	for i, req := range inputRequests {
		fmt.Printf("Processing Request %d: From %s, To %s, Weight %.2fkg\n",
			i+1, req.From, req.To, req.Weight)

		possiblePrices := calculatePrices(req, appState)

		allCalculatedResults = append(allCalculatedResults, CalculatedResult{
			Input:  req,
			Output: possiblePrices,
		})

		if len(possiblePrices) > 0 {
			fmt.Println("  Possible Rates:")
			for _, price := range possiblePrices {
				fmt.Printf("    - RateID: %s, Price: %.2f\n", price.RateID, price.Price)
			}
		} else {
			fmt.Println("  No applicable rates found.")
		}
		fmt.Println()
	}

	// Write all calculated results to output.json
	outputJSON, err := json.MarshalIndent(allCalculatedResults, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling output to JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile("output.json", outputJSON, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Calculations complete. Results written to output.json")
}
