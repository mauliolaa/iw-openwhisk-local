import datetime
import sys
import json
import matplotlib.pyplot as plt
from pprint import pprint

if len(sys.argv) != 1:
    usage = '''[Usage]: python3 plot_results.py 
    '''
    print(usage)
    exit(1)
strats = [
    "lru",
    "mru",
    "mfe",
    "pq",
    "rs",
    "naive"
]
periods = ["5", "10", "15"]

total_cold_container_stats = {
}

total_warm_container_stats = {
}

total_prewarmed_container_stats = {
}

total_recreated_container_stats = {
}

for strat in strats:
    for period in periods:
        prefix = "{}_{}".format(strat, period)
        with open("./{}/results.json".format(prefix), "r") as f:
            json_data = json.load(f)

        for fun_name in json_data.keys():
            if fun_name == "languages":
                continue
            if strat in total_cold_container_stats:
                total_cold_container_stats[strat] += json_data[fun_name]["coldContainerCount"]
            else:
                total_cold_container_stats[strat] = 0

            if strat in total_warm_container_stats:
                total_warm_container_stats[strat] += json_data[fun_name]["warmedContainerCount"]
            else:
                total_warm_container_stats[strat] = 0

            if strat in total_prewarmed_container_stats:
                total_prewarmed_container_stats[strat] += json_data[fun_name]["prewarmedContainerCount"]
            else:
                total_prewarmed_container_stats[strat] = 0
                
            if strat in total_recreated_container_stats:
                total_recreated_container_stats[strat] += json_data[fun_name]["recreatedContainerCount"]
            else:
                total_recreated_container_stats[strat] = 0
        
        if strat == "naive":
            break


pprint(total_cold_container_stats)

    
fig_title = "hello"
fig, axs = plt.subplots(1, 1)
axs.bar(total_cold_container_stats.keys(), total_cold_container_stats.values())
# plt.gcf().subplots_adjust(bottom=0.15)
plt.title("Total number of cold container starts by strategy")
plt.xlabel("Strategy")
plt.ylabel("Number of container starts")
plt.savefig(fig_title + ".png", bbox_inches="tight")
