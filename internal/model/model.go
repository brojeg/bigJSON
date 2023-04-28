package model

type Recipe struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

type RecordData struct {
	Recipe      string
	Delivery    string
	LocalCounts map[string]int
}

type Output struct {
	UniqueRecipeCount       int                   `json:"unique_recipe_count"`
	CountPerRecipe          []CountPerRecipe      `json:"count_per_recipe"`
	BusiestPostcode         PostcodeDeliveryCount `json:"busiest_postcode"`
	CountPerPostcodeAndTime PostcodeDeliveryCount `json:"count_per_postcode_and_time"`
	MatchByName             []string              `json:"match_by_name"`
}

type CountPerRecipe struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type PostcodeDeliveryCount struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
	From          string `json:"from,omitempty"`
	To            string `json:"to,omitempty"`
}
