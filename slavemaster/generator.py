"""
Generate a workload file.
"""
import sys

import numpy as np

def usage():
    print("python generator.py {filename} {n} {mean} {variance} {num_functions}")
    print("ideally, mean and variance should be such that sampled durations are not negative")
    exit(0)
    

if __name__ == "__main__":
    if len(sys.argv) != 6:
        usage()
    filename = sys.argv[1]
    n = int(sys.argv[2])
    mean = float(sys.argv[3])
    var = float(sys.argv[4])
    num_funcs = int(sys.argv[5])

    # Note, what units are these timestamps in?
    timestamps = np.random.normal(mean, var, n)
    timestamps = np.cumsum(timestamps)
    timestamps = np.round(timestamps)
    functions = np.random.randint(low=0, high=num_funcs, size=n)
    with open(filename, "w") as outf:
        for timestamp, func in zip(timestamps, functions):
            # Note: I think that our workload files should be in CSV format, since CSV
            # is easy to read and a bit less of a hassle if we want to start adding 
            # more parameters to functions, etcetera.

            # If we are going to store function number, then we will also need to have a
            # map associating faas functions with thier number. For example,
            # 0 = '/myfunction-1', 1 = '/otherfunction', 2 = '/finalfunction' 

            # Format will be something like this:
            # timestamp, function_number
            outf.write(f"{str(timestamp)},{func}\n")