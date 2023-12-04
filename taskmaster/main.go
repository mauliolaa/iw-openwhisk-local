package main

import (
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"taskmaster/predictor"
	"time"

	"gopkg.in/yaml.v2"
)

// taskmasterConfig is the configuration file parsed from the yaml config file
type taskmasterConfig struct {
	PollingPeriodicity int64  `yaml:"pollingPeriodicity"`
	Strategy           string `yaml:"strategy"`
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

var Config taskmasterConfig = taskmasterConfig{}
var strategy predictor.Predictor

// Calls the Openwhisk interface
// assumes that Openwhisk has been set up and that wsk cli utility exists on the system
func CallFn(fnName string, parameters map[string]string, logData bool) {
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
	if logData {
		logFn(fnName, elapsed)
	}
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
	// Acknowledge immediately
	w.WriteHeader(http.StatusOK)
	go CallFn(fnName, params, true)
	// Update the predictor
	updateStrategy(fnName, params)
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
	// Parse yaml file into Config struct
	config_filepath := os.Args[1]
	predictor_filepath := os.Args[2]
	yamlFile, err := os.ReadFile(config_filepath)
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Error parsing yaml file: %v", err)
	}
	fmt.Printf("%+v", Config)
	fmt.Printf("Config periodicity %d\n", Config.PollingPeriodicity)
	// initialize the strategy
	switch Config.Strategy {
	case "lru":
		strategy = predictor.NewLRU(predictor_filepath)
	default:
		log.Fatalf("Strategy not specified")
	}
	// launch the periodic scheduling algorithm
	go schedule()
	// initialize logging
	fnTimings = make(map[string][]time.Duration)
}

func usage() {
	if len(os.Args) != 3 {
		log.Panic("[Usage]: [general_config] [predictor_config]")
	}
}

func updateStrategy(fnName string, fnParams map[string]string) {
	log.Printf("Updating strategy with %s\n", fnName)
	info := make(map[string]any)
	var fnRequest predictor.FnRequest
	fnRequest = predictor.FnRequest{
		FnName:       fnName,
		FnParameters: fnParams,
	}
	// Add more information as we go along...
	info["fnRequest"] = fnRequest
	strategy.Update(info)
}

func predict() {
	log.Println("Predict called!")
	response := strategy.Predict()
	if response == predictor.NilPrediction {
		log.Println("No function to be called!")
		return
	}
	log.Printf("Pinging %s\n", response)
	CallFn(response.FnName, response.FnParameters, false)
	// TODO: Add metrics here to observe usefulness of pinging
}

func schedule() {
	// use the ticker from the time package
	ticker := time.NewTicker(time.Duration(Config.PollingPeriodicity) * time.Second)

	for {
		select {
		case <-ticker.C:
			predict()
		}
	}
}

func main() {
	usage()
	initialize()
	port := 1024

	// default handlers
	http.HandleFunc("/receive", ReceiveEvent)
	http.HandleFunc("/predict", PredictImage)
	http.HandleFunc("/dumpData", DumpData)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
