import datetime
import sys
import matplotlib.pyplot as plt

if len(sys.argv) != 5:
    usage = '''[Usage]: python3 visualize_events.py [a/w] [events_filename] [functions_filename] [output_filename]
    activations/workload: whether its activations or workload, expect a str of the form "a" or "w"
    events_filename: name of activations of workload file
    functions_filename: name of functions suite
    output_filename: name of figure to output
    
    This script scrapes the event occurrences and plots them on a matplotlib.pyplot eventplot figure. This is 
    helpful for visualizing how the prediction strategy and the simulated workload go hand in hand and for 
    doing a quick sanity check if its bugged or something.
    '''
    print(usage)
    exit(1)

functions = {}
mode = sys.argv[1]
if mode not in ["w", "a"]:
    print("Mode must either be w or a")
    exit(1)
events_filename = sys.argv[2]
functions_filename = sys.argv[3]
with open(functions_filename, "r") as inf:
    for line in inf.readlines():
        line = line.strip()
        fn_name = line.split(" ")[0]
        functions[fn_name] = []

def scrape_workload(functions, filename):
    current_dt = datetime.timedelta(seconds=0)
    # The expected format is
    # 5.294115286254328,hello_js,name:3UJXK20PB7MJJ97ACK4QAUA3N5SMULZTJ6WBS,place:QXZ27O76JTVN760A4IQ1X2FZQ7QDXEODI5LUY0DNE37LAN6CAFD
    with open(filename, "r") as inf:
        for line in inf.readlines():
            line = line.strip()
            # The timedelta is already given for free
            components = line.split(",")
            current_dt = datetime.timedelta(seconds=float(components[0])) + current_dt
            fn_name = components[1]
            functions[fn_name].append(current_dt)
    return functions
    

def scrape_predictions(functions, filename):
    start_time = None
    # The expected format is
    # 2023-12-12 06:56:40.041918 -0500 EST m=+10.002772418 hello_js
    with open(filename, "r") as inf:
        for line in inf.readlines():
            line = line.strip()
            # We will plot by milliseconds time delta from the start
            components = line.split(" ")
            if len(components) != 6:  # A nil prediction
                continue
            dt_str = components[0] + " " + "".join(components[1].split(".")[0])
            fn_name = components[-1]
            print(dt_str, fn_name)
            dt_obj = datetime.datetime.strptime(dt_str, "%Y-%m-%d %H:%M:%S")
            if not start_time:
                start_time = dt_obj
            dt = dt_obj - start_time
            functions[fn_name].append(dt)
    return functions
    
if mode == "w":
    functions = scrape_workload(functions, events_filename)
else:
    functions = scrape_predictions(functions, events_filename)
    
    
fig_title = sys.argv[4] + " " + ("workload" if mode == "w" else "ping")
fn_names = list(functions.keys())
X = [[dt.total_seconds() for dt in functions[k]] for k in fn_names]
labels = [fn_name for fn_name in fn_names]
fig, axs = plt.subplots(1, 1)
axs.eventplot(X, 
              orientation="horizontal",
              lineoffsets=1,
              colors=[f"C{i}" for i in range(len(X))],
              )
axs.legend(labels, bbox_to_anchor=(0., 1.0, 1., .10), loc=3,ncol=3, mode="expand", borderaxespad=0.)
plt.gcf().subplots_adjust(bottom=0.15)
plt.title(fig_title)
plt.savefig(fig_title + ".png", bbox_inches="tight")