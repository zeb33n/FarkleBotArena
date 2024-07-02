from collections.abc import Callable
from pipe_client import PipeClient
import random
import time
import os
import subprocess
import json
from dataclasses import asdict, dataclass
from farkle_scorer import calculate_score

# TODO
# save write proper pybot exebot functions for calling the bots wiht jsons and capturing the outpur
# write gamestate to json function
# calculate scores properly and remove dice from pool
# make output pretty

LINES = 3
BOT_DIR_LOC = f'{__file__.rsplit("/", 1)[0]}/bots'


@dataclass
class GameState:
    bots: dict[str, int]
    num_dice: int
    round_score: int
    roll: list[int]
    turn: str

    def to_bot(self) -> bytes:
        bot_info = ("bots", "num_dice", "round_score")
        out_json = {k: v for k, v in asdict(self).items() if k in bot_info}
        return bytes(json.dumps(out_json).encode("utf-8"))

    def to_tui(self) -> str:
        return json.dumps(asdict(self))


def make_bot(name: str, extension: str) -> Callable[[bytes], bool]:
    match extension:
        case "py":
            process = ["python", f"{BOT_DIR_LOC}/{name}.py"]
        case "exe":
            process = [f"{BOT_DIR_LOC}/{name}.exe"]
        case _:
            raise ValueError(f"bot {name} has unsupported extension .{extension}")

    def bot(json: bytes) -> bool:
        try:
            out = subprocess.run(
                process + [json],
                capture_output=True,
                check=True,
            )
            return bool(int(out.stdout))

        except subprocess.CalledProcessError as e:
            print(e.stderr.decode())
            return False

    return bot


class App:
    def __init__(self):
        input("press enter to start game\n")
        self.bots = self.load_bots()
        self.pipe_client = PipeClient()
        self.game_state = GameState({name: 0 for name in self.bots}, 6, 0, [], "")

    def load_bots(self) -> dict[str, Callable[[bytes], bool]]:
        bot_info = tuple(botfile.rsplit(".", 1) for botfile in os.listdir(BOT_DIR_LOC))
        return {name: make_bot(name, extension) for name, extension in bot_info}

    def roll_dice(self):
        return [random.randint(1, 6) for _ in range(self.game_state.num_dice)]

    def main(self):
        while max(self.game_state.bots.values()) < 10000:
            for bot_name, run_bot in self.bots.items():
                self.game_state.round_score = 0
                while run_bot(self.game_state.to_bot()):
                    self.game_state.turn = bot_name
                    self.game_state.roll = self.roll_dice()
                    score, self.game_state.num_dice = calculate_score(
                        self.game_state.roll
                    )
                    self.game_state.round_score += score
                    self.pipe_client.pipe(self.game_state.to_tui())
                    if score == 0:
                        self.game_state.num_dice = 6
                        break
                    if self.game_state.num_dice == 0:
                        self.game_state.num_dice = 6
                    time.sleep(1)
                else:
                    self.game_state.bots[bot_name] += self.game_state.round_score
        self.pipe_client.cleanup()


if __name__ == "__main__":
    App().main()
