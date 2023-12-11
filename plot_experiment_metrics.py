"""This script reads in a result.json file and plots some very nice graphs"""

import matplotlib.pyplot as plt
import json
import sys
import os

if len(sys.argv) != 3:
    print("Usage: python plot_experiment_metrics.py [results.json] [folder_name]")
    exit(1)
    

filename = sys.argv[1]
output_folder = sys.argv[2]
os.makedirs(output_folder, exist_ok=True)
with open(filename, "r") as inf:
    results = json.load(inf)
    for funcName in results.keys():
        if funcName not in ["languages", "num_fns_completed"]:
            funcResults = results[funcName]
            plt.hist(funcResults["elapsedTimes"], bins=10, color='blue', alpha=0.7)
            plot_name = f"Histogram of {funcName}-{funcResults['language']}"
            plt.title(plot_name)
            plt.xlabel(f"Time (milliseconds)")
            plt.ylabel("Frequency")
            text_message = f"Cold: {funcResults['coldContainerCount']}"
            plt.text(0.05, 0.95, text_message, transform=plt.gca().transAxes, bbox=dict(facecolor='white', edgecolor='black', boxstyle='round,pad=0.5'))

            plt.grid(True)
            plt.savefig(os.path.join(output_folder, plot_name))
            plt.clf()