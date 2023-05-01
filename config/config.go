package config

import (
	"bigJson/pkg"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Application struct {
	FilePath,
	Postcode,
	StartTimeString,
	EndTimeString string
	KeywordsString []string
}

func (a *Application) StartTime() int {
	return pkg.ParseHour(a.StartTimeString)
}

func (a *Application) EndTime() int {
	return pkg.ParseHour(a.EndTimeString)
}
func Init() (*Application, error) {
	var config Application

	var keywordsString, postcode, timeWindow, path string
	flag.StringVar(&keywordsString, "keywords", "Potato,Veggie,Mushroom", "Comma-separated list of keywords to search for in recipe names")
	flag.StringVar(&postcode, "postcode", "10120", "Postcode to search for")
	flag.StringVar(&timeWindow, "time", "10AM-3PM", "Time window to search for in format 'start-end' (e.g., '10AM-3PM')")
	flag.StringVar(&path, "file", "deliveries.json", "Path to file location")
	flag.Parse()

	err := validatePostcode(postcode)
	if err != nil {
		return nil, err
	}
	keywords, err := validateKeywords(keywordsString)
	if err != nil {
		return nil, err
	}
	err = validateTimeWindow(timeWindow)
	if err != nil {
		return nil, err
	}
	err = validateFilePath(path)
	if err != nil {
		return nil, err
	}
	config.KeywordsString = keywords
	config.Postcode = postcode
	timeParts := strings.Split(timeWindow, "-")
	if len(timeParts) != 2 {
		return nil, err
	}
	config.StartTimeString = timeParts[0]
	config.EndTimeString = timeParts[1]

	config.FilePath = path
	return &config, nil

}

func validateFilePath(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("file path is empty")
	}

	if filepath.Ext(path) != ".json" {
		return fmt.Errorf("file is not a json file: %s", path)
	}

	return nil
}

func validateKeywords(keywords string) ([]string, error) {
	// Split the keywords string into a slice
	keywordsSlice := strings.Split(keywords, ",")

	// Create a map to check for duplicates and store validated keywords
	keywordsMap := make(map[string]bool)

	for _, keyword := range keywordsSlice {
		// Trim spaces and convert to lowercase for case-insensitive comparison
		trimmedKeyword := strings.TrimSpace(strings.ToLower(keyword))

		if trimmedKeyword == "" {
			// Skip empty keywords
			continue
		}

		// Add to map if not already present (map automatically avoids duplicates)
		keywordsMap[trimmedKeyword] = true
	}

	// Return an error if no valid keywords
	if len(keywordsMap) == 0 {
		return nil, fmt.Errorf("no valid keywords provided")
	}

	// Generate the validatedKeywords slice from the keys of the map
	validatedKeywords := make([]string, 0, len(keywordsMap))
	for k := range keywordsMap {
		validatedKeywords = append(validatedKeywords, k)
	}

	return validatedKeywords, nil
}

func validatePostcode(postcode string) error {
	if len(postcode) != 5 {
		return fmt.Errorf("postcode must be 5 digits long")
	}

	for _, char := range postcode {
		if !unicode.IsDigit(char) {
			return fmt.Errorf("postcode must only contain digits")
		}
	}

	return nil
}

func validateTimeWindow(timeWindow string) error {
	// Split the timeWindow string into start and end times
	timeWindow = strings.ReplaceAll(timeWindow, " ", "")
	times := strings.Split(timeWindow, "-")
	if len(times) != 2 {
		return fmt.Errorf("time window must be in format 'start-end'")
	}

	// Trim spaces
	startTimeStr := times[0]
	endTimeStr := times[1]

	// Parse and validate start time
	_, err := time.Parse("3PM", startTimeStr)
	if err != nil {
		return fmt.Errorf("invalid start time: %v", err)
	}

	// Parse and validate end time
	_, err = time.Parse("3PM", endTimeStr)
	if err != nil {
		return fmt.Errorf("invalid end time: %v", err)
	}

	return nil
}
