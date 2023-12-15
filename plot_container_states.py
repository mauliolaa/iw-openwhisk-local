import sys
import json
import numpy as np
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
    "naive"
]

periods = [5, 10, 15]

total_cold_container_stats = {
    5: {},
    10: {},
    15: {},
}

total_warm_container_stats = {
    5: {},
    10: {},
    15: {},
}

total_prewarmed_container_stats = {
    5: {},
    10: {},
    15: {},
}

total_recreated_container_stats = {
    5: {},
    10: {},
    15: {},
}

for strat in strats:
    for period in periods:
        if strat == "naive":
            prefix = "naive"
        else:
            prefix = f"{strat}_{period}"
        with open(f"./{prefix}/results.json", "r") as f:
            json_data = json.load(f)
        for fun_name in json_data.keys():
            if fun_name == "languages":
                continue
            if strat in total_cold_container_stats[period]:
                total_cold_container_stats[period][strat] += json_data[fun_name]["coldContainerCount"]
            else:
                total_cold_container_stats[period][strat] = 0

            if strat in total_warm_container_stats[period]:
                total_warm_container_stats[period][strat] += json_data[fun_name]["warmedContainerCount"]
            else:
                total_warm_container_stats[period][strat] = 0

            if strat in total_prewarmed_container_stats[period]:
                total_prewarmed_container_stats[period][strat] += json_data[fun_name]["prewarmedContainerCount"]
            else:
                total_prewarmed_container_stats[period][strat] = 0
                
            if strat in total_recreated_container_stats[period]:
                total_recreated_container_stats[period][strat] += json_data[fun_name]["recreatedContainerCount"]
            else:
                total_recreated_container_stats[period][strat] = 0

pprint(total_cold_container_stats)

    
for period in periods:
    period = int(period)
    fig_title = f"Period {period}"
    fig, axs = plt.subplots(1, 1)
    bottom = np.zeros(len(strats))
    print(total_cold_container_stats.keys())
    print(total_cold_container_stats[period].values())
    axs.bar(total_cold_container_stats[period].keys(), total_cold_container_stats[period].values(), label="Cold", color="blue")
    axs.bar(total_warm_container_stats[period].keys(), total_warm_container_stats[period].values(), bottom=np.array(list(total_cold_container_stats[period].values())), label="Warmed", color="red")
    axs.bar(total_prewarmed_container_stats[period].keys(), total_prewarmed_container_stats[period].values(), bottom=np.array(list(total_cold_container_stats[period].values())) + np.array(list(total_warm_container_stats[period].values())),label="Prewarmed", color="salmon")
    axs.bar(total_recreated_container_stats[period].keys(), total_recreated_container_stats[period].values(), label="Recreated", bottom=np.array(list(total_cold_container_stats[period].values())) + np.array(list(total_warm_container_stats[period].values())) + np.array(list(total_prewarmed_container_stats[period].values())), color="dodgerblue")
    plt.title(f"Total number of containers for period {period}")
    plt.xlabel("Strategy")
    plt.ylabel("Number of container starts")
    plt.legend()
    plt.savefig(fig_title + ".png", bbox_inches="tight")
    plt.clf()
