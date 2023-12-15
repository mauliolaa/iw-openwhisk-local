# Instructions

All set up was done on an M1 Mac machine. No other devices have been tested.

First, you want to install the following stuff

1. Openwhisk
   1. Requires java, gradle and docker
2. wsk cli

You want to follow the openwhisk instructions for configuring the wsk cli to work with Openwhisk.

## Project Structure

Results are excluded.

```
├── README.md  # You are reading this
├── bugreport.md  # Antiquated, only had one minor bug that we decided to ignore by doing a different approach
├── generate_overall_results.py  # Used to generate a json file of the overall results used for the latex table
├── get_experiment_metrics.py  # Used to scrape the experimental results from the openwhisk log file
├── plot_container_states.py  # Used to obtain the stacked bar plots of the container states
├── plot_experiment_metrics.py  # Used to obtain the histogram of the function timings 
├── obtain_latencies.py  # Obtain the latencies for each programming language
├── report  # The folder for our report
├── results_to_table.py  # Used to produce the latex string for the table of results. Further manual work needed to merge cells
├── taskmaster
│   ├── confs            # The various configuration files used for the experiments
│   ├── functions        # The suite of functions for use by the simulation
│   ├── functions_test   # The file describing the action names (for openwhisk), their corresponding file and parameters
│   ├── generator.py     # Generates a workload for the simulation given a functions file and their respective serverless functions
│   ├── main.go          # Taskmaster
│   ├── plot_workload.py # Plots the action and language counts for a given workload
│   ├── predictor        # The strategies used for the simulation alongside their test cases. Names are self explanatory
│   │   ├── lru.go
│   │   ├── lru_test.go
│   │   ├── mfe.go
│   │   ├── mfe_test.go
│   │   ├── mru.go
│   │   ├── mru_test.go
│   │   ├── pqueue.go
│   │   ├── pqueue_test.go
│   │   ├── predictor.go
│   │   ├── rs.go
│   │   └── rs_test.go
│   ├── run_exps.ps1     # Untested powershell script for automating it. Use at your own risk
│   ├── simulate.py      # Used to simulate the workload on taskmaster
│   ├── test_workload    # The actual workload we used. Feel free to use this to reproduce our results
│   ├── test_workload_short  # A shorter workload for making sure it works before committing 1 hours and 20 minutes
├── visualize_events.py  # Generates the event timeline for either the workload or the "pings"
├── workload_events workload.png   # The event timeline for our test workload
```

## How to run stuff

1. Run openwhisk in standalone mode

In a separate terminal process:

```bash
cd openwhisk
./gradlew :core:standalone:build # This creates a runnable jar
java -jar bin/openwhisk-standalone.jar > openwhisk_out
```

Note: may need to reserve extra jvm memory for openwhisk with -Xmx option. Example:
```
java -Xmx4096m -jar bin/openwhisk-standalone.jar > openwhisk_out
```

This will take a while to run on your first build but subsequent builds will be faster. Please follow all instructions on wskcli using the default namespace option for MacOS.


1. Run taskmaster

You will need to have the configuration files ready. You may refer to the sample config [here](taskmaster/sampleconfig.yaml) as well as the sample lru predictor strategy [here](taskmaster/lru_config.yaml). Then in a new terminal process

```bash
cd taskmaster
go run main.go sampleconfig.yaml lru_config.yaml
```

This spawns the taskmaster http server and it should wait on port localhost 1024

3. Prepare functions and workload file

You should have some sample functions ready or simply use our functions at taskmaster/functions.
Prepare the functions test file, refer to [this](taskmaster/functions_test).
To generate a random workload file, refer to the [generator](taskmaster/generator.py)

4. Run the commands in experiments section, using the correct workload and functions file

5. Plot results and look

In a new terminal process, run

```bash
cd taskmaster
python simulate.py test_workload functions_test http://127.0.0.1:1024
```

replacing test_workload and functions_test with your own ones if preferrable.

When the workload is finished, run the following command

```bash
curl -X get "localhost:1024/dumpData"
```

This causes taskmaster to dump the experimental metrics into a log file for further analysis.

## File formats

### Format of the function workload file

[name of function in openwhisk] [filename of function] [params comma delimited]

## Metrics of interest

Container Start

openwhisk.counter.invoker_containerStart.cold_counter (counter) - Count of number of cold starts.
openwhisk.counter.invoker_containerStart.recreated_counter (counter) - Count of number of times container is recreated.
openwhisk.counter.invoker_containerStart.warm_counter (counter) - Count of number of times a warm container is used.
openwhisk.counter.invoker_containerStart.prewarmed_counter (counter) - Count of number of times a prewarmed container is used.

## Misc Sources

https://mikhail.io/serverless/coldstarts/aws/
https://mikhail.io/serverless/coldstarts/aws/languages/


## Getting results

Getting experimental results is non-trivial as Openwhisk CLI does not have a good API for displaying results.
It either shows the time/parameters and not the cold/warm/prewarmed status or it only shows the cold/warm/prewarmed status and not the parameters!
So we have to obtain our experiment metrics in a rather roundabout way

Taskmaster will invoke all functions without receiving results and receiving the activation id.
With the activation id, we can keep track of what are pings and what are legitimate non-activations.
Then we can tabulate the warm/cold counters as well as elapsed timing.

## Experiments

Parameters: 
- Polling Periodicity
  - 5s
  - 10s
  - 15s
- Predictors
  - LRU
  - MRU
  - PQueue
  - Baseline
  - MFE
- Workload
  - 5000 seconds = 83 mins


Number of experiments to run

### Experiments to run

#### Naive

Status: Done

```bash
cd taskmaster
go run main.go naiveconfig.yaml lru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 naive_0_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json naive
# No point
# python visualize_events.py a lru_10_taskmaster_ping.txt functions_test lru_10_ping
```

#### LRU


##### LRU 5

```bash
cd taskmaster
go run main.go confs/lru_5.yaml confs/lru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/lru_5_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json lru_5
python visualize_events.py a taskmaster/lru_5_taskmaster_ping.txt taskmaster/functions_test lru_5/lru_5_ping
```

##### LRU 10


```bash
cd taskmaster
go run main.go lru_10_config.yaml lru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/lru_10_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json lru_10
python visualize_events.py a taskmaster/lru_10_taskmaster_ping.txt taskmaster/functions_test lru_10/lru_10_ping
```

##### LRU 15

```bash
cd taskmaster
go run main.go lru_15_config.yaml lru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/lru_15_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json lru_15
python visualize_events.py a taskmaster/lru_15_taskmaster_ping.txt taskmaster/functions_test lru_15/lru_15_ping
```

#### MFE

##### MFE 5

```bash
cd taskmaster
go run main.go mfe_5_config.yaml mfe_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mfe_5_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mfe_5
python visualize_events.py a taskmaster/mfe_5_taskmaster_ping.txt taskmaster/functions_test mfe_5/mfe_5_ping
```

##### MFE 10


```bash
cd taskmaster
go run main.go mfe_10_config.yaml lru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mfe_10_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mfe_10
python visualize_events.py a taskmaster/mfe_10_taskmaster_ping.txt taskmaster/functions_test mfe_10/mfe_10_ping
``` 

##### MFE 15

```bash
cd taskmaster
go run main.go mfe_15_config.yaml mfe_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mfe_15_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mfe_15
python visualize_events.py a taskmaster/mfe_15_taskmaster_ping.txt taskmaster/functions_test mfe_15/mfe_15_ping
```

#### MRU

##### MRU 5

```bash
cd taskmaster
go run main.go mru_5_config.yaml mru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mru_5_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mru_5
python visualize_events.py a taskmaster/mru_5_taskmaster_ping.txt taskmaster/functions_test mru_5/mru_5_ping
```

##### MRU 10


```bash
cd taskmaster
go run main.go mru_10_config.yaml mru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mru_10_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mru_10
python visualize_events.py a taskmaster/mru_10_taskmaster_ping.txt taskmaster/functions_test mru_10/mru_10_ping
``` 

##### MRU 15

```bash
cd taskmaster
go run main.go mru_15_config.yaml mru_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/mru_15_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json mru_15
python visualize_events.py a taskmaster/mru_15_taskmaster_ping.txt taskmaster/functions_test mru_15/mru_15_ping
```

#### PQ

##### PQ 5

The pq_config.yaml file does not matter

```bash
cd taskmaster
go run main.go confs/pq_5.yaml pq_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/pq_5_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json pq_5
python visualize_events.py a taskmaster/pq_5_taskmaster_ping.txt taskmaster/functions_test pq_5/pq_5_ping
```

##### PQ 10

The pq_config.yaml file does not matter

```bash
cd taskmaster
go run main.go confs/pq_10.yaml pq_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/pq_10_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json pq_10
python visualize_events.py a taskmaster/pq_10_taskmaster_ping.txt taskmaster/functions_test pq_10/pq_10_ping
```

##### PQ 15

The pq_config.yaml file does not matter

```bash
cd taskmaster
go run main.go confs/pq_15.yaml pq_config.yaml functions_test
python simulate.py test_workload functions_test http://127.0.0.1:1024
cd ..
python get_experiment_metrics.py openwhisk/outty.log taskmaster/functions_test http://127.0.0.1:1024 taskmaster/pq_15_taskmaster_activation_ids.txt
python plot_experiment_metrics.py results.json pq_15
python visualize_events.py a taskmaster/pq_15_taskmaster_ping.txt taskmaster/functions_test pq_15/pq_15_ping
```


## Notes

### Compiling Java

Compiling Java actions requre you to compile with --release 8 flag as the Openwhisk Java runtime only supports Java 8. Not doing so will result in an application error.

```bash
cd functions
javac --release 8 -cp Gson\ 2.10.1.jar Hello.java
jar cvf Hello.jar Hello.class
wsk action create HelloJava Hello.jar --main Hello
```

### Go cannot find binary

TODO: Cannot get go openwhisk actions to invoke successfully. No idea why and googling doesn't help. Logs also do not show up. Invoking with debug gives no useful information. No issue on GitHub either. Temporarily ignore all go functions for now. 

## Random musings

Some trends to take note of:

1. Does the polling periodicity actually make a difference?
   1. What is the trend as you go up with relation to container states and average latency of function response
2. Write down that a polling periodicity of 1s causes Openwhisk to be overwhelmed, unfortunately no time to run full suite of experiments
3. Which are the most effective (and harmful) strategies?
4. Does an eyeball test on the time axis event plot suggest that the strategies are working as intended?
   1. Mostly yes so far
   2. PQ seems most promising
   3. Must look at MFE
      1. Since distribution is uniform, maybe will seem similar
5. How would the results differ if the workload distribution was not Gaussian?
   1. Potential future explorations
