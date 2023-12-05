metrics_of_interest = {
    "invoker_containerStart.cold_counter:": 0,
    "invoker_containerStart.recreated_counter:": 0,
    "invoker_containerStart.warmed_counter:": 0,
}

with open("openwhisk/openwhisk_out", "r") as f:
    lines = []
    # FIXME: Very unsure how Openwhisk does metric reporting but we are just gonna stick to what it says for now
    # Probably because the containers take awhile to report the metrics back, we cannot simply just return [-1] and call it a day
    # as the log's cold_counter is not a strictly increasing sequence
    for line in f.readlines():
        for metric in metrics_of_interest.keys():
            if metric in line:
                metrics_of_interest[metric] = max(metrics_of_interest[metric], int(line.split(metric)[-1].split("]")[0]))
print(f"{metrics_of_interest.items()}")