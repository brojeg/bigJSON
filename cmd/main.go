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

	configApp, err := config.Init()
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open(configApp.FilePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bufferSize := 100 * 1024 * 1024
	iter := jsoniter.Parse(jsoniter.ConfigFastest, file, bufferSize)

	app := controller.NewCliProcess(configApp, iter)
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
