package main

import (
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

const TaskmasterOutputFile = "taskmaster_activation_ids.txt"

func DumpData(w http.ResponseWriter, r *http.Request) {
	log.Print("Dumping data.")
	f, err := os.Create(TaskmasterOutputFile)
	if err != nil {
		log.Panic("Unable to create file with err: ", err)
	}
	defer f.Close()

	for _, activationId := range activationList {
		_, err = f.WriteString(activationId)
		if err != nil {
			log.Panic("Error writing to log: ", err)
		}
	}
}

var Config taskmasterConfig = taskmasterConfig{}
var strategy predictor.Predictor

// Calls the Openwhisk interface
// assumes that Openwhisk has been set up and that wsk cli utility exists on the system
func CallFn(fnName string, parameters map[string]string, logData bool) {
	// make a call to the faascli with the requested fnName
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
	if logData {
		activationMU.Lock()
		components := strings.Split(string(output), " ")
		activationId := components[len(components)-1]
		activationList = append(activationList, activationId)
		activationMU.Unlock()
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
	switch runtime.GOOS {
	case "windows":
		cmdPrompt = "cmd"
	default:
		cmdPrompt = "bash"
	}
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
	case "mfe":
		strategy = predictor.NewMFE()
	default:
		log.Fatalf("Strategy not specified")
	}
	activationList = make([]string, 0)
	// launch the periodic scheduling algorithm only if the PollingPeriodicity is greater than 0
	if Config.PollingPeriodicity > 0 {
		go schedule()
	}
}

func usage() {
	if len(os.Args) != 3 {
		usage := `[Usage]: [general_config] [predictor_config]
		[general_config]: a yaml file that contains the following parameters
			pollingPeriodicity: a float
			strategy: Choose from 'lru', 'pq', 'ml', 'mfe'
		[predictor_config]: a yaml file that corresponds to the predictor.
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
