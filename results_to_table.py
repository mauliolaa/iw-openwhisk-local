"""Generate the latex formatted table. Some manual work still needs to be done to merge the strategies into one cell"""

from tabulate import tabulate
import json

# Strat, Period, Cold, Warm, Pre, Recreated, Java, Python, Javascript, Ruby, Php
header = ["Strat", "Period", "Cold", "Warm", "Prewarmed", "Recreated", "Java", "py", "Js", "rb", "php"]
table = []

strats = [
    "lru",
    "mfe",
    "mru",
    "pq",
    "naive",
]

with open("overall_results.json", "r") as inf:
    results = json.load(inf)
    
print(results)
for strat in results.keys():
    for period in results[strat].keys():
        r = results[strat][period]
        table.append([strat, period, r["cold"], r["warm"], r["prewarmed"], r["recreated"], r["jar"], r["py"], r["js"], r["rb"], r["php"]])

latex_s = tabulate(table, header, tablefmt="latex")
with open("latex_table", "w") as outf:
    outf.write(latex_s)