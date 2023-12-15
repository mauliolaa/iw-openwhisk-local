"""Obtain the overall results file by combining average latencies, container states"""
import json

periods = [5, 10, 15]

strats = [
    "lru",
    # "mfe",
    "mru",
    "pq",
    "naive"
]

languages = [
    "jar",
    "py",
    "js",
    "rb",
    "php",
]

results = {
    strat: {period: {"cold": 0, "warm": 0, "prewarmed": 0, "recreated": 0} for period in periods} for strat in strats
}


with open("latencies.json", "r") as lf:
    latencies = json.load(lf)
    for l_periods_strategy in latencies.keys():
        for period in latencies[l_periods_strategy].keys():
            for language in languages:
                results[l_periods_strategy][int(period)][language] = latencies[l_periods_strategy][period][language]


for strat in strats:
    for period in periods:
        if strat == "naive":
            prefix = "naive"
        else:
            prefix = f"{strat}_{period}"
            prefix = f"{strat}_{period}"
        with open(f"./{prefix}/results.json", "r") as f:
            json_data = json.load(f)
        for fun_name in json_data.keys():
            if fun_name == "languages":
                continue
            print(json_data[fun_name])
            results[strat][period]["cold"] += json_data[fun_name]["coldContainerCount"]
            results[strat][period]["warm"] += json_data[fun_name]["warmedContainerCount"]
            results[strat][period]["prewarmed"] += json_data[fun_name]["prewarmedContainerCount"]
            results[strat][period]["recreated"] += json_data[fun_name]["recreatedContainerCount"]

with open("overall_results.json", "w") as outf:
    outf.write(json.dumps(results, indent=4))
