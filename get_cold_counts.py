with open("openwhisk/openwhisk_out", "r") as f:
    lines = []
    for line in f.readlines():
        if "cold_counter" in line:
            lines.append(line)
    cold_count = 0
    # FIXME: Very unsure how Openwhisk does metric reporting but we are just gonna stick to what it says for now
    # Probably because the containers take awhile to report the metrics back, we cannot simply just return [-1] and call it a day
    # as the log's cold_counter is not a strictly increasing sequence
    for line in lines:
        ss = line.split("cold_counter:")[-1].split("]")[0]
        cold_count = max(cold_count, int(ss))
print(f"Cold count is {cold_count}")