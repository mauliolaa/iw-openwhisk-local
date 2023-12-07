"""
This python script scrapes the openwhisk log file, keeping count of how many cold, recreated and warm containers there are.
"""
from dataclasses import dataclass, field
import sys
import dateutil.parser
import datetime
from datetime import datetime as dt


@dataclass
class ActionMetric:
    actionName: str
    language: str
    elapsedTimes: list[float] = field(default_factory=list)
    prewarmedContainerCount: int = 0
    warmedContainerCount: int = 0
    coldContainerCount: int = 0
    recreatedContainerCount: int = 0
    
    
@dataclass
class TrackingTID:
    actionName: str
    tid: str
    startingTime: datetime.datetime


if len(sys.argv) != 3:
    usage = """[experiment logs]
    Usage:
        python get_experiment_metrics.py [logfile] [functions_test]
    """
    print(usage)
    exit(1)
    
metrics_of_interest = {
    "languages": [],
    # Elapsed time taken for each function
}

functions_workload_filename = sys.argv[2]
    
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

with open(log_file_name, "r") as f:
    # OpenWhisk spawns a few container pools that separately keep track of their cold, prewarmed and warm state so 
    # it's best to just increment per encounter instead of use their reported values.
    current_time = None
    current_action = None
    tracking_tids = {}
    lines = f.readlines()
    i = 0
    while i < len(lines):
        line = lines[i]
        # Encountered prewarmed, warm, cold container
        if "containerStart containerState" in line:
            # Format is 
            # [2023-12-06T20:05:39.324Z] [34m[INFO][0;39m [[1m#tid_sDmMc7WcvvA0fbhaGF6ubefjKu0Ef4mL[0m] [[36mContainerPool[0m] containerStart containerState: recreated container: None activations: 1 of max 1 action: hello_world_py namespace: guest activationId: 60399bef79e54791b99bef79e5b79169 [marker:invoker_containerStart.recreated_counter:57]
            # Get the tid
            tid = line.split(" ")[2].replace("\x1b", "")[1:-1]
            if tid not in tracking_tids.keys():
                i += 1
                continue
            current_action = tracking_tids[tid].actionName
            if current_action not in metrics_of_interest.keys():
                i += 1
                continue
            print(line)
            if "containerState: prewarmed container" in line:
                metrics_of_interest[current_action].prewarmedContainerCount += 1
            elif "containerState: cold container" in line:
                metrics_of_interest[current_action].coldContainerCount += 1
            elif "containerState: warmed container" in line:
                metrics_of_interest[current_action].warmedContainerCount += 1
            elif "containerState: recreated container" in line:
                metrics_of_interest[current_action].recreatedContainerCount += 1
            else:
                raise ValueError("Something went wrong???")
        # New function invocation, keep track of the tid
        if "POST" in line:
            # Follows format: 
            # [2023-12-06T20:04:40.980Z] [34m[INFO][0;39m [#tid_2Jdj2cFsvQilyMAXt0L12eW7htROfU4u] POST /api/v1/namespaces/_/actions/hello_world_go blocking=true&result=true
            components = line.split(" ")
            # Get the time component and remove the brackets
            current_time = dateutil.parser.isoparse(components[0][1:-1])
            # Get the action name
            action_name = components[-2].split("/")[-1]
            # Do not keep track if not in list of action names
            if action_name not in metrics_of_interest.keys():
                i += 1
                continue
            # Now we want to get the unique tid of the invoked function and keep track of it    
            # Based off eyeball analysis, it's always the line immediately after and does not have any relevant metric
            i += 1
            line = lines[i]
            # Follows format:
            # [2023-12-06T20:04:40.981Z] [34m[INFO][0;39m [[1m#tid_2Jdj2cFsvQilyMAXt0L12eW7htROfU4u[0m] [[36mIdentity[0m] [GET] serving from cache: CacheKey(23bc46b1-71f6-4ed5-8c54-816aa4f8c502) [marker:database_cacheHit_counter:1]
            # tid of interest is the 3rd item in the split()
            tid = line.split(" ")[2].replace("\x1b", "")[1:-1]
            # We use the time the HTTP request was received
            tracking_tids[tid] = TrackingTID(actionName=action_name, tid=tid, startingTime=current_time)
        # End of a tid
        if "completion ack" in line:
            # Line is of the format
            # [2023-12-06T20:04:44.734Z] [34m[INFO][0;39m [[1m#tid_BA18kojTAysLPE7nH1MmICemWu1Rn8HC[0m] [[36mLeanBalancer[0m] received completion ack for '20a3294eefb649f9a3294eefb6f9f967', system error=false
            components = line.split(" ")
            tid = components[2].replace("\x1b", "")[1:-1]
            if tid not in tracking_tids.keys():
                i += 1
                continue
            end_time = dateutil.parser.isoparse(components[0][1:-1])
            elapsed_timedelta = end_time - tracking_tids[tid].startingTime
            action_name = tracking_tids[tid].actionName
            metrics_of_interest[action_name].elapsedTimes.append(elapsed_timedelta.total_seconds())
            del tracking_tids[tid]
        i += 1        
                
print(f"{metrics_of_interest.items()}")
with open(f"results.txt", "w") as f:
    f.write("Metrics of interest\n")
    for k, v in metrics_of_interest.items():
        f.write(f"{k=} {v=}\n")
    f.write("Lines of interest\n")
    for line in lines_of_interest:
        f.write(f"{line}\n")