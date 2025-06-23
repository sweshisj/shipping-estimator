# Freight Cost Calculator

This repository contains a Go program designed to calculate freight shipping costs based on defined geographical zones and rates. It processes a stream of events to establish the pricing rules and then provides cost estimates for various shipments.

---

## Problem Overview

A warehouse needs to estimate shipping costs for freight to different locations. The costs are determined by a set of rules (rates) that depend on:

  -The shipment's weight.
  -The origin and destination postal codes, which map to predefined zones.
  -The program's goal is to calculate and list all possible prices for a given shipment, considering all applicable rates.

---

## Requirements Fulfilled

This solution addresses the following core requirements:

  -Input Handling: Accepts shipment details (weight, from postcode, to postcode).
  -Rate Calculation: Calculates all possible prices for the given input based on defined rates.
  -Zone Mapping: Correctly maps postcodes to their respective zones.
  -Event Processing: Processes a stream of ZoneDefined and RateDefined events to build an in-memory representation of zones and rates.
  -Output Generation: Prints all possible prices to the console and writes a structured JSON output file (output.json) detailing the input request and its corresponding applicable rates.
  -Standard Go Tooling: Developed using standard Golang features and practices.
  -Test Data Usage: Utilizes provided testdata for events and inputs.

---

## Getting Started

### Prerequisites

- Go


### Clone the repo

```bash
git clone https://github.com/sweshisj/shipping-estimator
cd shipping-estimator
```

### Running the App

#### Ensure go.mod is initialized

```bash
go mod init shipping_estimator
```

#### Run the application

```bash
go run main.go
```
---

## Input and Output

  -Input: The program reads shipment requests from testdata/input.json. This file is an array of JSON objects, each specifying From postcode, To postcode, and Weight.
  
  -Output: The program will print a summary of processed requests and their applicable rates to the console.
A file named output.json will be generated in the root of the project directory. This file will contain a structured JSON array, where each element represents an input request along with an array of all PossiblePrice objects found for that request.

---

## Time Spent

I spent approximately 2.5 hours working on this problem. This time was distributed as follows:

0.5 hours: Understanding the problem, reading documentation, and setting up the Go module.
1 hour: Defining data structures (structs), implementing JSON unmarshaling for events and inputs, and basic event processing.
0.5 hours: Implementing the core pricing logic (calculatePrices function) and zone lookup.
0.5 hours: Modifying the main function for reading from input.json, writing to output.json, and refining console output.

---

## Contact

For questions or support, please contact [sweshisj@gmail.com](mailto:sweshisj@gmail.com).
