"""
Generate a workload file.
"""
import random
import string
import sys

import numpy as np

def usage():
    usage = """
    [taskmaster Workload Generator]
    Usage:
        python generator.py {filename} {n} {mean} {variance} {function_file}
    Params:
        filename: name of workload file to generate
        n: number of serverless function calls to invoke
        mean: the mean time for a sleep, gaussian distribution
        variance: the variance of the sleep, gaussian distribution
        function_file: the name of the function
    Note:
        All functions in function_file must be ones that are supported by Openwhisk. We do not check for this so care
    """
    print(usage)
    exit(0)
    
    
def generate_random_parameter(param_type):
    if param_type == "string":
        N = random.randint(1, 100)
        return ''.join(random.choices(string.ascii_uppercase + string.digits, k=N))
    elif param_type == "int":
        return random.randint(1, 256)
    else:
        print(param_type)
        raise ValueError("Unsupported param type!")
    

if __name__ == "__main__":
    if len(sys.argv) != 6:
        usage()
    filename = sys.argv[1]
    n = int(sys.argv[2])
    mean = float(sys.argv[3])
    var = float(sys.argv[4])
    function_filename = sys.argv[5]
    
    # Populate functions map
    functions = {}  # map containing function names as well as parameters
    functionNames = []  # the keys() of functions is not subscriptable which we need for rand.choice
    with open(function_filename, "r") as inf:
        for line in inf.readlines():
            line = line.strip()
            # Format is [openwhisk action] [function filename (ignored)] [params (possibly empty)]
            line_split = line.split(" ")
            fnName = line_split[0]
            params = {}
            if len(line_split) > 2:
                for param in line_split[2].split(","):
                    param_name, param_type = param.split("=")
                    params[param_name] = param_type
            functions[fnName] = params
            functionNames.append(fnName)

    with open(filename, "w") as outf:
        for _ in range(n):
            delta = random.gauss(mu=mean, sigma=var)
            fn = random.choice(functionNames)
            s = f"{delta},{fn}"
            params = []
            for param, param_type in functions[fn].items():
                param_value = generate_random_parameter(param_type)
                s += f",{param}:{param_value}"
            outf.write(s+"\n")