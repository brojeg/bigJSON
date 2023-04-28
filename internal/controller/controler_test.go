package controller

import (
	"bigJson/internal/model"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestProcess(t *testing.T) {
	testData := `
	[
		{
			"recipe": "Veggie Potato Salad",
			"postcode": "10120",
			"delivery": "Tuesday 10AM - 1PM"
		},
		{
			"recipe": "Mushroom Risotto",
			"postcode": "10120",
			"delivery": "Tuesday 11AM - 2PM"
		}
	]`

	iter := jsoniter.ParseString(jsoniter.ConfigFastest, testData)
	keywords := []string{"Potato", "Veggie", "Mushroom"}

	output := Process(iter, "10120", "10AM", "3PM", keywords)

	expectedOutput := model.Output{
		UniqueRecipeCount: 2,
		CountPerRecipe: []model.CountPerRecipe{
			{Recipe: "Mushroom Risotto", Count: 1},
			{Recipe: "Veggie Potato Salad", Count: 1},
		},
		BusiestPostcode: model.PostcodeDeliveryCount{
			Postcode:      "10120",
			DeliveryCount: 2,
		},
		CountPerPostcodeAndTime: model.PostcodeDeliveryCount{
			Postcode:      "10120",
			From:          "11AM",
			To:            "3PM",
			DeliveryCount: 1,
		},
		MatchByName: []string{"Mushroom Risotto", "Veggie Potato Salad"},
	}

	if output.UniqueRecipeCount != expectedOutput.UniqueRecipeCount {
		t.Errorf("Expected unique recipe count to be %d, got %d", expectedOutput.UniqueRecipeCount, output.UniqueRecipeCount)
	}

	if len(output.CountPerRecipe) != len(expectedOutput.CountPerRecipe) {
		t.Errorf("Expected count per recipe length to be %d, got %d", len(expectedOutput.CountPerRecipe), len(output.CountPerRecipe))
	}

	defer os.Remove("output.json")

}
