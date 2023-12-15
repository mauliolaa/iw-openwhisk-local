"""This script obtains the average latencies for each language for each periodicity and strategy"""
import json

periods = [5, 10, 15]
strats = [
    "lru",
    "mfe",
    "mru",
    "pq",
    "naive",
]
# {action: {period: | {language
# Results table
# strategies periodicity cold warm prewarmed recreated py java rb php js
results = {
    strat: {period: {} for period in periods} for strat in strats
}

print(results)

action_to_language_mapping = {}

with open("taskmaster/functions_test", "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        action = line.split(" ")[0]
        language = line.split(" ")[1].split(".")[1]
        action_to_language_mapping[action] = language

for strat in strats:
    for period in periods:
        if strat == "naive":
            filename = f"{strat}/results.json"
        else:
            filename = f"{strat}_{period}/results.json"
        with open(filename, "r") as f:
            sub_data = json.load(f)
            for action, language in action_to_language_mapping.items():
                timings = sub_data[action]["elapsedTimes"]
                if language not in results[strat][period]:  # First time
                    results[strat][period][language] = timings
                else:  # All files accounted for
                    results[strat][period][language].extend(timings)
                    elapsed_timings = results[strat][period][language]
                    results[strat][period][language] = sum(elapsed_timings) / len(elapsed_timings)

print(results)
with open("latencies.json", "w") as outf:
    json.dump(results, outf)