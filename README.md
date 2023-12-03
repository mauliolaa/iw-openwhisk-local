# Instructions

First, you want to install the following stuff

1. Openwhisk
   1. Requires java, gradle and docker
2. wsk cli

You want to follow the openwhisk instructions for configuring the wsk cli to work with Openwhisk.

## How to run stuff

1. Run openwhisk in standalone mode

In a separate terminal process:

```bash
cd openwhisk
./gradlew :core:standalone:build
```

This will take a while to run on your first build but subsequent builds will be faster.

2. Run taskmaster

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

4. Run the simulator.py

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
