package main

import (
	"bufio"
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
	"strconv"
)

// ImageRequest is a struct that represents an image request
type ImageRequest struct {
	ImageName string `json:"imageName"`
}

// Logging of experimental results
var fnTimings map[string][]time.Duration

func logFn(fnName string, elapsedTime time.Duration) {
	_, contains := fnTimings[fnName]
	if !contains {
		fnTimings[fnName] = make([]time.Duration, 0)
	}
	elapsedTime.Hours()
	fnTimings[fnName] = append(fnTimings[fnName], elapsedTime)
}

func DumpData(w http.ResponseWriter, r *http.Request) {
	dumpData("sampleData.txt")
	w.WriteHeader(http.StatusOK)
}

func dumpData(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error creating file: ", err)
	}
	defer f.Close()

	for fnName := range fnTimings {
		resultString := fnName + ": ["
		for _, timing := range fnTimings[fnName] {
			resultString = resultString + " " + strconv.FormatInt(timing.Microseconds(), 10)
		}
		resultString = resultString + "]\n"
		_, err := f.WriteString(resultString)
		if err != nil {
			log.Fatal("Error writing string to file: ", err)
		}
	}
}


// var predictor Predictor


// Calls the Openwhisk interface
// assumes that Openwhisk has been set up and that wsk cli utility exists on the system
func CallFn(fnName string, parameters map[string]string) {
	// make a call to the faascli with the requested fnName
	cmd := fmt.Sprintf("wsk action invoke %s --result ", fnName)
	fmt.Println(parameters)
	for param, value := range parameters {
		cmd = cmd + fmt.Sprintf("--param %s %s ", param, value)
	}
	
	start := time.Now()
	fmt.Println("Executing command ", cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("could not run command!", err)
	}
	fmt.Println("Output: ", string(out))
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Elapsed time was %s\n", elapsed)
	logFn(fnName, elapsed)
}

// Gateway function that receives a fnRequest from the workload
func ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Receive event called!")
	body, err := io.ReadAll(r.Body)
	_ = body
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	if !r.URL.Query().Has("fnName") {
		http.Error(w, "No fnName specified", http.StatusInternalServerError)
		return
	}
	fnName := r.URL.Query().Get("fnName")

	// Handle serverless function parameters
	params := make(map[string]string)
	for param, value := range r.URL.Query() {
		if param == "fnName" {
			continue
		}
		params[param] = value[0]
	}
	fmt.Println(fnName)
	CallFn(fnName, params)
	_ = fnName

	// var imageRequest ImageRequest
	// if err = json.Unmarshal(body, &imageRequest); err != nil {
	// 	http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
	// 	return
	// }

	// add to timeseries database
	w.WriteHeader(http.StatusOK)
}

func PredictImage(w http.ResponseWriter, r *http.Request) {
	// what should be sent here?
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// send prediction back
	// image := predictor.Predict()
	image := "Jotham"
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

	// initialize logging
	fnTimings = make(map[string][]time.Duration)
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
	http.HandleFunc("/dumpData", DumpData)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
