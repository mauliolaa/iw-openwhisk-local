"""This script reads in the workload/function file and plots out the pie chart for the action/language distribution"""

from collections import defaultdict
import matplotlib.pyplot as plt

total_duration = 0.0
action_to_language = {}
action_counts = defaultdict(int)
language_counts = defaultdict(int)

with open("functions_test", "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        components = line.split(" ")
        action = components[0]
        language = components[1].split(".")[1]
        action_to_language[action] = language
print(f"action to language: ", action_to_language)
with open("test_workload", "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        components = line.split(",")
        total_duration += float(components[0])
        action_counts[components[1]] += 1
        language_counts[action_to_language[components[1]]] += 1

print(f"Total duration: {total_duration}s")

a_labels = []
a_counts = []
for label, count in action_counts.items():
    a_labels.append(label)
    a_counts.append(count)

plt.pie(a_counts, labels=a_labels, autopct='%1.1f%%', startangle=90)
plt.axis("equal")
plt.savefig(f"Action counts.png")
plt.clf()

l_labels = []
l_counts = []
print(language_counts)
for label, count in language_counts.items():
    l_labels.append(label)
    l_counts.append(count)

plt.pie(l_counts, labels=l_labels, autopct='%1.1f%%', startangle=90)
plt.axis("equal")
plt.savefig(f"Language counts.png")