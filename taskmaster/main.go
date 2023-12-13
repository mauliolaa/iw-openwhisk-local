package main

import (
	"bufio"
	"strconv"

	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"taskmaster/predictor"
	"time"

	"gopkg.in/yaml.v2"
)

// taskmasterConfig is the configuration file parsed from the yaml config file
type taskmasterConfig struct {
	PollingPeriodicity int64  `yaml:"pollingPeriodicity"`
	Strategy           string `yaml:"strategy"`
}

// dependent on OS
var cmdPrompt string

// Lock for writing to activationID list
var activationMU sync.Mutex
var activationList []string

type PredictionLogRecord struct {
	predictedFn string
	timestamp   time.Time
}

var predictionRecords []PredictionLogRecord
var simulationStarted = false

// Filled in by config
var experimentName string

const TaskmasterActivationFile = "taskmaster_activation_ids.txt"
const TaskmasterPingFile = "taskmaster_ping.txt"

func DumpData(w http.ResponseWriter, r *http.Request) {
	log.Print("Dumping activation ids data.")
	activationFile, err := os.Create(experimentName + "_" + TaskmasterActivationFile)
	if err != nil {
		log.Panic("Unable to create file with err: ", err)
	}
	defer activationFile.Close()

	for _, activationId := range activationList {
		_, err = activationFile.WriteString(activationId)
		if err != nil {
			log.Panic("Error writing to log: ", err)
		}
	}

	pingFile, err := os.Create(experimentName + "_" + TaskmasterPingFile)
	if err != nil {
		log.Panic("Unable to create file with err: ", err)
	}
	defer pingFile.Close()

	for _, pingRecord := range predictionRecords {
		s := pingRecord.timestamp.String() + " " + pingRecord.predictedFn
		_, err = pingFile.WriteString(s + "\n") // \n is needed for the timestamp
		if err != nil {
			log.Panic("Error writing to log: ", err)
		}
	}

}

var Config taskmasterConfig = taskmasterConfig{}
var strategy predictor.Predictor

// mapping of actions to important information, current just language
var actionMapping map[string]ActionInfo

type ActionInfo struct {
	language string
}

// Calls the Openwhisk interface
// assumes that Openwhisk has been set up and that wsk cli utility exists on the system
func CallFn(fnName string, parameters map[string]string, logData bool) {
	// make a call to the wsk cli with the requested fnName
	cmd := fmt.Sprintf("wsk action invoke %s ", fnName)
	for param, value := range parameters {
		cmd = cmd + fmt.Sprintf("--param %s %s ", param, value)
	}

	log.Println("Executing command ", cmd)
	output, err := exec.Command(cmdPrompt, "-c", cmd).Output()

	if err != nil {
		log.Println("could not run command!", err)
	}
	// We want to store the activation id if it's not a ping!
	// The format of the output is
	// ok: invoked /_/simple_math with id 1d6b703e1e5d4d39ab703e1e5dcd3984
	// We also want to initialize the predictionLogRecords now since the simulation has started for realz
	if logData {
		activationMU.Lock()
		components := strings.Split(string(output), " ")
		activationId := components[len(components)-1]
		activationList = append(activationList, activationId)
		activationMU.Unlock()
	}
	if logData && !simulationStarted {
		log.Println("Initialized prediction log records")
		simulationStarted = true
	}
}

// Gateway function that receives a fnRequest from the workload
func ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive event called!")
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
	log.Println(fnName)
	// Acknowledge immediately
	w.WriteHeader(http.StatusOK)
	go CallFn(fnName, params, true)
	// Update the predictor
	updateStrategy(fnName, params)
}

// initialize reads in the config file and initializes the scheduler accordingly
func initialize() {
	switch runtime.GOOS {
	case "windows":
		cmdPrompt = "powershell"
	default:
		cmdPrompt = "bash"
	}
	// Parse yaml file into Config struct
	configFilepath := os.Args[1]
	predictorFilepath := os.Args[2]
	yamlFile, err := os.ReadFile(configFilepath)
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Error parsing yaml file: %v", err)
	}
	fmt.Printf("%+v", Config)
	fmt.Printf("Config periodicity %d\n", Config.PollingPeriodicity)

	// read in function file and populate action map
	functionsFilepath := os.Args[3]
	functionsFile, err := os.Open(functionsFilepath)
	if err != nil {
		log.Fatalf("Error reading functions file: %v", err)
	}
	defer functionsFile.Close()

	functionsScanner := bufio.NewScanner(functionsFile)
	functionsScanner.Split(bufio.ScanLines)
	actionMapping = make(map[string]ActionInfo)
	for functionsScanner.Scan() {
		line := functionsScanner.Text()
		components := strings.Split(line, " ")
		actionName := components[0]
		language := strings.Split(components[1], ".")[1]
		actionMapping[actionName] = ActionInfo{language: language}
	}

	// initialize the strategy
	switch Config.Strategy {
	case "naive":
		strategy = nil
	case "lru":
		strategy = predictor.NewLRU(predictorFilepath)
	case "mfe":
		strategy = predictor.NewMFE()
	case "pq":
		strategy = predictor.NewPriorityQueue()
	case "rs":
		strategy = predictor.NewRS()
	case "mru":
		strategy = predictor.NewMRU(predictorFilepath)
	default:
		log.Fatalf("Strategy not specified")
	}
	activationList = make([]string, 0)
	predictionRecords = make([]PredictionLogRecord, 0)
	experimentName = Config.Strategy + "_" + strconv.FormatInt(Config.PollingPeriodicity, 10)
	// launch the periodic scheduling algorithm only if the PollingPeriodicity is greater than 0
	if Config.PollingPeriodicity > 0 {
		go schedule()
	}
}

func usage() {
	if len(os.Args) != 4 {
		usage := `[Usage]: [general_config] [predictor_config]
		[general_config]: a yaml file that contains the following parameters
			pollingPeriodicity: a float
			strategy: Choose from 'lru', 'pq', 'rs', 'mfe', 'mru'
		[predictor_config]: a yaml file that corresponds to the predictor.
		[functions_file]: file consisting of functions to be called
		`
		log.Panic(usage)
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
	info["language"] = actionMapping[fnName].language
	strategy.Update(info)
}

func predict() {
	// Do not waste prediction on nothing
	if !simulationStarted {
		log.Println("Simulation not started!")
		return
	}
	log.Println("Predict called!")
	response := strategy.Predict()
	predictionRecord := PredictionLogRecord{
		predictedFn: response.FnName,
		timestamp:   time.Now(),
	}
	predictionRecords = append(predictionRecords, predictionRecord)
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
	http.HandleFunc("/dumpData", DumpData)
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
