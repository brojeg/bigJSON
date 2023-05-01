package controller

import (
	"bigJson/config"
	"bigJson/internal/model"
	"bigJson/pkg"
	"log"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type CLIProcess struct {
	config *config.Application
	iter   *jsoniter.Iterator
}

func NewCliProcess(config *config.Application, iter *jsoniter.Iterator) *CLIProcess {
	var app CLIProcess
	app.iter = iter
	app.config = config
	return &app
}

func (cli CLIProcess) Process() (*model.Output, error) {
	recipeCounts := make(map[string]int)
	postcodeCounts := make(map[string]int)
	var maxPostcode string
	maxCount := 0
	deliveries10120 := 0

	recipesWithKeywords := make(map[string]bool)

	for cli.iter.ReadArray() {
		var record model.Recipe
		cli.iter.ReadVal(&record)
		if cli.iter.Error != nil {
			log.Println("Error reading JSON:", cli.iter.Error)
			return nil, cli.iter.Error
		}

		// Increment the count for this recipe name
		recipeCounts[record.Recipe]++

		// Increment the count for this postcode
		postcodeCounts[record.Postcode]++

		// If this postcode has more deliveries than the current max, update the max
		if count := postcodeCounts[record.Postcode]; count > maxCount {
			maxCount = count
			maxPostcode = record.Postcode
		}

		// Count the number of deliveries to postcode
		if record.Postcode == cli.config.Postcode {
			startHour, endHour := parseDeliveryTime(record.Delivery)
			if endHour < startHour { // delivery spans over midnight
				if (startHour <= cli.config.StartTime() && cli.config.StartTime() < 24) || // start before midnight
					(0 <= cli.config.EndTime() && cli.config.EndTime() <= endHour) { // end after midnight
					deliveries10120++
				}
			} else if startHour <= cli.config.StartTime() && endHour >= cli.config.EndTime() {
				deliveries10120++
			}
		}

		// Add recipe names containing any of the specified keywords to the set
		if containsKeyword(record.Recipe, cli.config.KeywordsString) {
			recipesWithKeywords[record.Recipe] = true
		}
	}

	// Get the keys of the map (the unique recipe names) and sort them
	recipes := make([]string, 0, len(recipeCounts))
	for recipe := range recipeCounts {
		recipes = append(recipes, recipe)
	}
	sort.Strings(recipes)

	// Get the keys of the map (the unique recipe names containing keywords) and sort them
	recipesKeywords := make([]string, 0, len(recipesWithKeywords))
	for recipe := range recipesWithKeywords {
		recipesKeywords = append(recipesKeywords, recipe)
	}
	sort.Strings(recipesKeywords)

	output := model.Output{
		UniqueRecipeCount: len(recipeCounts),
		BusiestPostcode: model.PostcodeDeliveryCount{
			Postcode:      maxPostcode,
			DeliveryCount: maxCount,
		},
		CountPerPostcodeAndTime: model.PostcodeDeliveryCount{
			Postcode:      cli.config.Postcode,
			From:          cli.config.StartTimeString,
			To:            cli.config.EndTimeString,
			DeliveryCount: deliveries10120,
		},
		MatchByName: recipesKeywords,
	}

	for _, recipe := range recipes {
		output.CountPerRecipe = append(output.CountPerRecipe, model.CountPerRecipe{
			Recipe: recipe,
			Count:  recipeCounts[recipe],
		})
	}

	return &output, nil
}

func parseDeliveryTime(delivery string) (startHour, endHour int) {
	parts := strings.Split(delivery, "-")
	startHour = pkg.ParseHour(parts[0])
	endHour = pkg.ParseHour(parts[1])
	return startHour, endHour
}

func containsKeyword(recipe string, keywords []string) bool {
	lowerCaseRecipe := strings.ToLower(recipe)
	for _, keyword := range keywords {
		lowerCaseKeyword := strings.ToLower(keyword)
		if strings.Contains(lowerCaseRecipe, lowerCaseKeyword) {
			return true
		}
	}
	return false
}
