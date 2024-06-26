from collections.abc import Callable
import random
import time
import os
import subprocess
import json
from dataclasses import asdict, dataclass

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

    def to_json(self) -> bytes:
        return bytes(json.dumps(asdict(self)).encode("utf-8"))


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
        self.bots = self.load_bots()
        self.game_state = GameState({name: 0 for name in self.bots}, 6, 0)
        self.initialise_screen()
        self.update_scores()

    def load_bots(self) -> dict[str, Callable[[bytes], bool]]:
        bot_info = tuple(botfile.rsplit(".", 1) for botfile in os.listdir(BOT_DIR_LOC))
        return {name: make_bot(name, extension) for name, extension in bot_info}

    def update_scores(self):
        magic = "\033[F"
        update = "\n".join(
            [f"{name}: {score}" for name, score in self.game_state.bots.items()]
        )
        magic = magic * (update.count("\n") + (LINES + 1))
        print(f"{magic}{update}{'\n' * LINES}")

    def initialise_screen(self):
        print("SCORES")
        for _ in self.game_state.bots:
            print()
        for _ in range(LINES):
            print()

    def roll_dice(self):
        return [random.randint(1, 6) for _ in range(self.game_state.num_dice)]

    def main(self):
        while True:
            for bot_name, run_bot in self.bots.items():
                self.game_state.round_score = 0
                while run_bot(self.game_state.to_json()):
                    self.game_state.round_score += 100
                    print("\033[F" * 4)
                    print((" " * 50 + "\n") * 2)
                    print("\033[F" * 4)
                    print(f"{bot_name}'s turn", flush=True)
                    print(f"current score: {self.game_state.round_score}")
                    print(f"rolled {self.roll_dice()}")
                    time.sleep(1)
                else:
                    self.game_state.bots[bot_name] += self.game_state.round_score
                    self.update_scores()


if __name__ == "__main__":
    App().main()
