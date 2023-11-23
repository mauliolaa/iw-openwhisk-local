"""Read in a workload file and send requests in specified order"""
import sys
import requests
import time

from workload_format import csv_params

# This is a map associating function numbers with function names,
# formatted as a URL path. 
# For example,
# {0: '/myfunction-1', 1: '/otherfunction', 2: '/finalfunction'} 
#
# This could be defined in code, or read in from a file. For now, I am 
# hardcoding an example function map. inovke with faas_url='https://google.com'
function_map = {
    0: "/travel",
    1: "/imghp",
    2: "/maps"
}

def usage():
    print("[USAGE]: python invoker.py {workload_file} {faas_url}")
    exit()

if __name__ == "__main__":
    if len(sys.argv) != 3:
        usage()
    filename = sys.argv[1]
    # Faas gateway URL
    faas_url = sys.argv[2]
    
    # Should check if faas service is running
    r = requests.get(faas_url)
    if r.status_code != 200:
        raise RuntimeError("Faas gateway returned non-200 status code.\nURL: {}\nCode: {}".format(r.url,r.status_code))
    
    faas_requests = []
    # Program execution starts at timestamp 0.
    prev_time = 0
    with open(filename, "r") as inf:
        lines = inf.read().splitlines()
        for line in lines:
            # Should split into timestamp, function name, param
            # Then use requests to make a GET request
            # print(line)
            csvs = line.split(sep=',')
            p = csv_params

            timestamp = float(csvs[p["timestamp"]])
            function_number = int(csvs[p["function_number"]])
            
            print("Timestamp: {}".format(timestamp))
            print("Function number: {}".format(function_number))

            # I don't think this sleep code is totally correct, maybe it would be easier to
            # define workloads in terms of time deltas rather than timestamps?
            delta = timestamp - prev_time
            print("Sleep delta: {}".format(delta))
            time.sleep(max(0, delta))

            r = requests.get(faas_url + function_map[int(function_number)])
            if r.status_code != 200:
                raise RuntimeError("Function returned non-200 status code.\nURL: {}\nCode: {}".format(r.url,r.status_code))
            
            # Time between request send and response in miliseconds.
            elasped_time = round(r.elapsed.total_seconds() * 1000)
            print("Elapsed time for request {}:\n{} ms\n".format(r.url, elasped_time))
            faas_requests.append(elasped_time)
            

            prev_time = timestamp + delta



    # Now, graph elapsed time data
    print(faas_requests)
            