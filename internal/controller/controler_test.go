package controller

import (
	"bigJson/config"
	"bigJson/internal/model"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var mockData = `[
	{
		"postcode": "10120",
		"recipe": "Potato",
		"delivery": "10AM-3PM"
	},
	{
		"postcode": "10120",
		"recipe": "Veggie",
		"delivery": "10AM-3PM"
	},
	{
		"postcode": "20230",
		"recipe": "Mushroom",
		"delivery": "10AM-3PM"
	}
]`

func TestCLIProcess_Process(t *testing.T) {
	// Define the expected output
	expected := &model.Output{
		UniqueRecipeCount: 3,
		CountPerRecipe: []model.CountPerRecipe{
			{Recipe: "Mushroom", Count: 1},
			{Recipe: "Potato", Count: 1},
			{Recipe: "Veggie", Count: 1},
		},
		BusiestPostcode: model.PostcodeDeliveryCount{
			Postcode:      "10120",
			DeliveryCount: 2,
		},
		CountPerPostcodeAndTime: model.PostcodeDeliveryCount{
			Postcode:      "10120",
			From:          "10AM",
			To:            "3PM",
			DeliveryCount: 2,
		},
		MatchByName: []string{"Potato", "Veggie"},
	}

	config := &config.Application{
		FilePath:        "",
		Postcode:        "10120",
		StartTimeString: "10AM",
		EndTimeString:   "3PM",
		KeywordsString:  []string{"Potato", "Veggie"},
		StartTime:       config.ParseHour("10AM"),
		EndTime:         config.ParseHour("3PM"),
	}

	// Create a new jsoniter.Iterator with mockData
	iter := jsoniter.ParseString(jsoniter.ConfigFastest, mockData)

	// Create a new CLIProcess with mock config and iterator
	cli := NewCliProcess(config, iter)

	// Call the Process method
	output := cli.Process()

	// Compare the output to the expected output
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %+v, got %+v", expected, output)
	}
}

func TestParseDeliveryTime(t *testing.T) {
	tests := []struct {
		name      string
		delivery  string
		startHour int
		endHour   int
	}{
		{"Morning to afternoon", "10AM-3PM", 10, 15},
		{"Afternoon to evening", "3PM-8PM", 15, 20},
		{"Around midnight", "11PM-1AM", 23, 1},
		{"Noon to evening", "12PM-6PM", 12, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := parseDeliveryTime(tt.delivery)
			if start != tt.startHour || end != tt.endHour {
				t.Errorf("parseDeliveryTime() = %v, %v; expected %v, %v", start, end, tt.startHour, tt.endHour)
			}
		})
	}
}

func TestContainsKeyword(t *testing.T) {
	tests := []struct {
		name     string
		recipe   string
		keywords []string
		want     bool
	}{
		{
			name:     "keyword is present",
			recipe:   "Delicious Mushroom Soup",
			keywords: []string{"Mushroom", "Onion"},
			want:     true,
		},
		{
			name:     "keyword is absent",
			recipe:   "Delicious Mushroom Soup",
			keywords: []string{"Chicken", "Onion"},
			want:     false,
		},
		{
			name:     "case insensitive",
			recipe:   "Delicious mushroom Soup",
			keywords: []string{"Mushroom", "Onion"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsKeyword(tt.recipe, tt.keywords); got != tt.want {
				t.Errorf("containsKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}
