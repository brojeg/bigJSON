package controller

import (
	"bigJson/internal/model"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Process(iter *jsoniter.Iterator, postcode, startTime, endTime string, keywords []string) *model.Output {
	recipeCounts := make(map[string]int)
	postcodeCounts := make(map[string]int)
	var maxPostcode string
	maxCount := 0
	deliveries10120 := 0

	recipesWithKeywords := make(map[string]bool)

	for iter.ReadArray() {
		var record model.Recipe
		iter.ReadVal(&record)
		if iter.Error != nil {
			fmt.Println("Error reading JSON:", iter.Error)
			return nil
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

		// Count the number of deliveries to postcode 10120 between 10AM and 3PM
		if record.Postcode == postcode {
			startHour, endHour := parseDeliveryTime(record.Delivery)
			if startHour >= parseHour(startTime) && endHour <= parseHour(endTime) {
				deliveries10120++
			}
		}

		// Add recipe names containing any of the specified keywords to the set
		if containsKeyword(record.Recipe, keywords) {
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

	// Print the recipe names containing any of the specified keywords, in alphabetical order

	output := model.Output{
		UniqueRecipeCount: len(recipeCounts),
		BusiestPostcode: model.PostcodeDeliveryCount{
			Postcode:      maxPostcode,
			DeliveryCount: maxCount,
		},
		CountPerPostcodeAndTime: model.PostcodeDeliveryCount{
			Postcode:      "10120",
			From:          "11AM",
			To:            "3PM",
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

	fileout, _ := os.Create("output.json")
	defer fileout.Close()

	jsoniter.NewEncoder(fileout).Encode(output)

	return &output
}

func parseDeliveryTime(delivery string) (startHour, endHour int) {
	parts := strings.Split(delivery, " ")
	startHour, _ = strconv.Atoi(parts[1][:len(parts[1])-2])
	endHour, _ = strconv.Atoi(parts[3][:len(parts[3])-2])
	return startHour, endHour
}
func containsKeyword(recipe string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(recipe, keyword) {
			return true
		}
	}
	return false
}

func parseHour(timeString string) int {
	hour, _ := strconv.Atoi(timeString[:len(timeString)-2])
	if strings.HasSuffix(timeString, "PM") && hour != 12 {
		hour += 12
	}
	return hour
}
