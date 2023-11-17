package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ImageRequest is a struct that represents an image request
type ImageRequest struct {
	ImageName string `json:"imageName"`
}

var predictor Predictor

func ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var imageRequest ImageRequest
	if err = json.Unmarshal(body, &imageRequest); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	// add to timeseries database
}

func PredictImage(w http.ResponseWriter, r *http.Request) {
	// what should be sent here?
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// send prediction back
	image := predictor.Predict()
	_, err = w.Write([]byte(image))
	if err != nil {
		http.Error(w, "Error writing to response writer", http.StatusInternalServerError)
		return
	}
}

// initialize reads in the config file and initializes the scheduler accordingly
func initialize() {
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		log.Panicf("Error reading configuration file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// format of configuration file should be image specifications
		// some notion of the image's resource requirements should be taken into account
		// first line should be Predictor choice
		// lines afterwards should be resources available?
		// lines afterwards should be images?
		line := scanner.Text()
		_ = line
	}

	if err := scanner.Err(); err != nil {
		log.Panicf("Error scanning file: %v", err)
	}
}

func usage() {
	if len(os.Args) != 2 {
		log.Panic("[Usage]: [config]")
	}
}

func main() {
	usage()
	initialize()
	port := 1024

	http.HandleFunc("/receive", ReceiveEvent)
	http.HandleFunc("/predict", PredictImage)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
