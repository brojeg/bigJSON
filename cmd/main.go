package main

import (
	"bigJson/internal/controller"
	"flag"
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func main() {

	var keywordsString, postcode, timeWindow string
	flag.StringVar(&keywordsString, "keywords", "Potato,Veggie,Mushroom", "Comma-separated list of keywords to search for in recipe names")
	flag.StringVar(&postcode, "postcode", "10120", "Postcode to search for")
	flag.StringVar(&timeWindow, "time", "10AM-3PM", "Time window to search for in format 'start-end' (e.g., '10AM-3PM')")
	flag.Parse()

	keywords := strings.Split(keywordsString, ",")
	timeParts := strings.Split(timeWindow, "-")
	if len(timeParts) != 2 {
		fmt.Println("Invalid time window format. Expected 'start-end', got", timeWindow)
		return
	}
	startTime, endTime := timeParts[0], timeParts[1]

	file, err := os.Open("bigJSON.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bufferSize := 32768 // or another size based on your needs
	iter := jsoniter.Parse(jsoniter.ConfigFastest, file, bufferSize)
	controller.Process(iter, postcode, startTime, endTime, keywords)
}
