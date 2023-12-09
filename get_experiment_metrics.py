"""
This python script scrapes the openwhisk log file, keeping count of how many cold, recreated and warm containers there are.
"""
from dataclasses import dataclass, field
import sys
import dateutil.parser
from datetime import datetime as dt
import requests
import subprocess
import json

@dataclass
class ActionMetric:
    actionName: str
    language: str
    elapsedTimes: list[float] = field(default_factory=list)
    prewarmedContainerCount: int = 0
    warmedContainerCount: int = 0
    coldContainerCount: int = 0
    recreatedContainerCount: int = 0


if len(sys.argv) != 4:
    usage = """[experiment logs]
    Usage:
        python get_experiment_metrics.py [logfile] [functions_test] [faas_url]
    """
    print(usage)
    exit(1)
    
metrics_of_interest = {
    "languages": [],
    # Elapsed time taken for each function
}
TASKMASTER_ACTIVATION_LIST = "taskmaster/taskmaster_activation_ids.txt"

functions_workload_filename = sys.argv[2]
faas_url = sys.argv[3]
    
# Parse functions workload
with open(functions_workload_filename, "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        components = line.split()
        action_name = components[0]
        function_name = components[1]
        language = function_name.split(".")[-1]
        if language not in metrics_of_interest["languages"]:
            metrics_of_interest["languages"].append(language)
        metrics_of_interest[action_name] = ActionMetric(actionName=action_name, language=language)

log_file_name = sys.argv[1]
lines_of_interest = []


def parse_iso_time(line):
    # Line is of the format
    # [2023-12-06T20:04:44.734Z] [34m[INFO][0;39m [[1m#tid_BA18kojTAysLPE7nH1MmICemWu1Rn8HC[0m] [[36mLeanBalancer[0m] received completion ack for '20a3294eefb649f9a3294eefb6f9f967', system error=false
    time = dateutil.parser.isoparse(components[0][1:-1])
    return time

# First get activation ids of interest
tracking_activation_ids = {}
# NOTE: Assumes that dump data has not yet been called
query_url = faas_url + "/dumpData"
requests.get(query_url)
with open(TASKMASTER_ACTIVATION_LIST, "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        tracking_activation_ids[line] = []
        
# NOTE: We are only scraping the log file to determine the state of the container since wsk activation get {id} does not tell us that
with open(log_file_name, "r") as f:
    # We just want to obtain all lines associated with each activation id
    # We can't do it sequentially because there is no sequential guarantee in the log file
    for line in f.readlines():
        for tracking_id in tracking_activation_ids.keys():
            if tracking_id in line:
                tracking_activation_ids[tracking_id].append(line)

for tracking_id in tracking_activation_ids.keys():
    print(f"Handling {tracking_id}")
    command = f"wsk activation get {tracking_id}"
    result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if result.returncode != 0:
        print(result.stderr)
        exit(1)
    json_result = result.stdout
    # This will be of the form 
    # ok: got activation {activation} \n
    # json result
    # We strip the first line by splitting then joining back into a string and parsing into a dict. No better way I think
    json_result = json.loads("".join(result.stdout.split("\n")[1:]))
    print(json_result)
    action_name = json_result["name"]
    # In miliseconds, convert to seconds
    duration = json_result["duration"] * 0.001  # TODO: Find unit time measurement
    metrics_of_interest[action_name].elapsedTimes.append(duration)
    # Obtain container state
    for line in tracking_activation_ids[tracking_id]:
        if "containerState: prewarmed container" in line:
            metrics_of_interest[action_name].prewarmedContainerCount += 1
        elif "containerState: cold container" in line:
            metrics_of_interest[action_name].coldContainerCount += 1
        elif "containerState: warmed container" in line:
            metrics_of_interest[action_name].warmedContainerCount += 1
        elif "containerState: recreated container" in line:
            metrics_of_interest[action_name].recreatedContainerCount += 1
                
print(f"{metrics_of_interest.items()}")
with open(f"results.txt", "w") as f:
    f.write("Metrics of interest\n")
    for k, v in metrics_of_interest.items():
        f.write(f"{k=} {v=}\n")
    f.write("Lines of interest\n")
    for line in lines_of_interest:
        f.write(f"{line}\n")