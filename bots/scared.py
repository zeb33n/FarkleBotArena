import json
import sys

gs = json.loads(sys.argv[1])

if max(gs["bots"].values()) > 2:
    print(1)

if gs["num_dice"] < 4:
    print(0)
else:
    print(1)
