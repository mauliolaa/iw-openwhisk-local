"""Read in a workload file and send requests in specified order"""
import sys
import requests


def usage():
    print("[USAGE]: python invoker.py {workload_file} {faas_url}")
    exit()

if __name__ == "__main__":
    if len(sys.argv) != 3:
        usage()
    filename = sys.argv[1]
    faas_url = sys.argv[2]
    
    # Should check if faas service is running
    
    with open(filename, "r") as inf:
        for line in inf.readlines():
            # Should split into timestamp, function name, param
            # Then use requests to make a GET request
            print(line)