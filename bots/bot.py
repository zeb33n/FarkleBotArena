import sys
import json

game_state_json = json.loads(sys.argv[1])

if game_state_json["round_score"] < 500:
    print(1)
else:
    print(0)
