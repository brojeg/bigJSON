package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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
	StartTime,
	EndTime int
}

func Init() *Application {
	var config Application

	var keywordsString, postcode, timeWindow, path string
	flag.StringVar(&keywordsString, "keywords", "Potato,Veggie,Mushroom", "Comma-separated list of keywords to search for in recipe names")
	flag.StringVar(&postcode, "postcode", "10120", "Postcode to search for")
	flag.StringVar(&timeWindow, "time", "10AM-3PM", "Time window to search for in format 'start-end' (e.g., '10AM-3PM')")
	flag.StringVar(&path, "file", "deliveries.json", "Path to file location")
	flag.Parse()

	err := validatePostcode(postcode)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	keywords, err := validateKeywords(keywordsString)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	err = validateTimeWindow(timeWindow)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	err = validateFilePath(path)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	config.KeywordsString = keywords
	config.Postcode = postcode
	timeParts := strings.Split(timeWindow, "-")
	if len(timeParts) != 2 {
		log.Println("Invalid time window format. Expected 'start-end', got", timeWindow)
		return nil
	}
	config.StartTimeString = timeParts[0]
	config.EndTimeString = timeParts[1]

	config.StartTime = ParseHour(config.StartTimeString)
	config.EndTime = ParseHour(config.EndTimeString)
	config.FilePath = path
	return &config

}

func validateFilePath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	} else if err != nil {
		return fmt.Errorf("error accessing file: %s", err)
	}
	return nil
}

func ParseHour(timeString string) int {
	hour, _ := strconv.Atoi(timeString[:len(timeString)-2])
	if strings.HasSuffix(timeString, "PM") && hour != 12 {
		hour += 12
	}
	if strings.HasSuffix(timeString, "AM") && hour == 12 {
		hour = 0
	}
	return hour
}
func validateKeywords(keywords string) ([]string, error) {
	// Split the keywords string into a slice
	keywordsSlice := strings.Split(keywords, ",")

	// Create a map to check for duplicates
	keywordsMap := make(map[string]bool)

	// Create a slice to store the validated keywords
	validatedKeywords := []string{}

	for _, keyword := range keywordsSlice {
		// Trim spaces and convert to lowercase for case-insensitive comparison
		trimmedKeyword := strings.TrimSpace(strings.ToLower(keyword))

		if trimmedKeyword == "" {
			// Skip empty keywords
			continue
		}

		// Check for duplicates
		if _, exists := keywordsMap[trimmedKeyword]; !exists {
			keywordsMap[trimmedKeyword] = true
			validatedKeywords = append(validatedKeywords, trimmedKeyword)
		}
	}

	if len(validatedKeywords) == 0 {
		return nil, fmt.Errorf("no valid keywords provided")
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
	times := strings.Split(timeWindow, "-")
	if len(times) != 2 {
		return fmt.Errorf("time window must be in format 'start-end'")
	}

	// Trim spaces
	startTimeStr := strings.TrimSpace(times[0])
	endTimeStr := strings.TrimSpace(times[1])

	// Check that both times contain either "AM" or "PM"
	if !(strings.Contains(startTimeStr, "AM") || strings.Contains(startTimeStr, "PM")) ||
		!(strings.Contains(endTimeStr, "AM") || strings.Contains(endTimeStr, "PM")) {
		return errors.New("both start and end times must specify either 'AM' or 'PM'")
	}

	// Validate start time
	if _, err := time.Parse("3PM", startTimeStr); err != nil {
		return fmt.Errorf("invalid start time: %v", err)
	}

	// Validate end time
	if _, err := time.Parse("3PM", endTimeStr); err != nil {
		return fmt.Errorf("invalid end time: %v", err)
	}

	startTime, err := time.Parse("3PM", startTimeStr)
	if err != nil {
		return fmt.Errorf("invalid start time: %v", err)
	}
	endTime, err := time.Parse("3PM", endTimeStr)
	if err != nil {
		return fmt.Errorf("invalid end time: %v", err)
	}

	// Check that the start time is earlier than the end time
	if startTime.After(endTime) {
		return errors.New("start time must be earlier than end time")
	}

	return nil
}
