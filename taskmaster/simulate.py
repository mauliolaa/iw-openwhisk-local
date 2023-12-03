"""Read in a workload file and send requests in specified order"""
import sys
import requests
import asyncio
import aiohttp  # Used to make asynchronous http requests because we don't actually want to wait for the http request to be done
import time
import subprocess

from workload_format import csv_params

def usage():
    usage = """
    [taskmaster Simulate]
    Usage:
        python invoker.py {workload_file} {function_file} {taskmaster_url}
    Params:
        workload_file: name of the workload file that the simulator uses
        function_file: name of function file that contains functions that are invoked on Openwhisk
        taskmaster_url: the url with which taskmaster is hosted on
    """
    print(usage)
    exit()
    
async def main():
    if len(sys.argv) != 4:
        usage()
    workload_filename = sys.argv[1]
    function_filename = sys.argv[2]
    faas_url = sys.argv[3]
    
    # Remove all existing functions in openwhisk first
    command = "wsk action list"
    result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if result.returncode != 0:
        print(result.stderr)
        exit(1)
    
    for line in result.stdout.splitlines()[1:]:
        action = line.split()[0].split("/")[-1]
        command = f"wsk action delete {action}"
        result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        if result.returncode != 0:
            print(result.stderr)
            exit(1)
        
    # Install all functions into openwhisk
    with open(function_filename, "r") as inf:
        lines = inf.read().splitlines()
        # Format of the file should be [fnName in whisk] [filename] [params (ignored)]
        for line in lines:
            fnName, filename = line.split(" ")[:2]
            command = f"wsk action create {fnName} {filename}"
            result = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
            if result.returncode != 0:
                print(f"Error: {result.stderr}")
                exit(1)

    start_time = time.time()
    async with aiohttp.ClientSession() as session:
        with open(workload_filename, "r") as inf:
            lines = inf.read().splitlines()
            for i, line in enumerate(lines):
                # Format is [time_delta],[action],[params...]
                # Should split into timestamp, function name, param
                print(f"Executing line {i}: {line}")
                p = line.split(sep=',')
                delta = float(p[0])
                functionName = p[1]
                
                params = {}
                for faas_param in p[2:]:
                    param_name, param_value = faas_param.split(":")
                    params[param_name] = param_value
                print(f"Sleeping for {delta} seconds")
                time.sleep(max(0, delta))
                query_url = faas_url + "/receive?fnName=" + functionName
                
                for param_name, param_value in params.items():
                    query_url += "&" + param_name + "=" + param_value
                    
                print(f"Query url is {query_url}") 
                async with session.get(query_url) as r:
                    _ = await r.json(content_type=None)  # disable decoding of json cos we don't really care
                    # if stuff.status_code != 200:
                    #     raise RuntimeError(f"Function returned non-200 status code.\nURL: {r.url}\nCode: {r.status_code}")

    # Now, graph elapsed time data
    elapsed_time = time.time() - start_time
    print(f"Simulation elapsed time: %s", elapsed_time)

if __name__ == "__main__":
    asyncio.run(main())