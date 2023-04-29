package main

import (
	"bigJson/config"
	"bigJson/internal/controller"
	"fmt"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func main() {

	configApp := config.Init()

	file, err := os.Open(configApp.FilePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bufferSize := 100 * 1024 * 1024
	iter := jsoniter.Parse(jsoniter.ConfigFastest, file, bufferSize)

	app := controller.NewCliProcess(configApp, iter)
	result := app.Process()
	fmt.Println(result)

}
