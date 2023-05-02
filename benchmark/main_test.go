package main

import (
	"bigJson/config"
	"bigJson/internal/controller"
	"fmt"
	"log"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func BenchmarkJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := config.Application{FilePath: "../bigJSON.json", Postcode: "10213", KeywordsString: []string{"Tex-Mex", "Tilapia"}, StartTimeString: "8AM", EndTimeString: "7PM"}
		file, err := os.Open(conf.FilePath)
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		bufferSize := 100 * 1024 * 1024
		iter := jsoniter.Parse(jsoniter.ConfigFastest, file, bufferSize)

		app := controller.NewCliProcess(&conf, iter)
		result, err := app.Process()
		if err != nil {
			log.Println(err)
		}
		e, err := jsoniter.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(e))
	}
}
