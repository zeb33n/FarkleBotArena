from collections.abc import Callable
import random
import time
import os
import subprocess
from dataclasses import dataclass

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


def make_pybot(name: str) -> Callable[[dict[str, str]], bool]:
    def pybot(json: dict[str, str]) -> bool:
        out = subprocess.run(
            ["python", f"{BOT_DIR_LOC}/{name}.py"],
            capture_output=True,
        )
        return int(out.stdout)

    return pybot


def make_exebot(name: str) -> Callable[[dict[str, str]], bool]:
    def exebot(json: dict[str, str]) -> bool:
        out = subprocess.run(
            [f"{BOT_DIR_LOC}/{name}.exe"],
            capture_output=True,
        )
        return int(out.stdout)

    return exebot


class App:
    def __init__(self):
        self.bots = self.load_bots()
        self.game_state = GameState({name: 0 for name in self.bots}, 6, 0)
        self.initialise_screen()
        self.update_scores()

    def load_bots(self) -> dict[str, Callable[[dict[str, str]], bool]]:
        bot_info = tuple(botfile.rsplit(".", 1) for botfile in os.listdir(BOT_DIR_LOC))
        out = {}
        for bot_name, extension in bot_info:
            match extension:
                case "py":
                    out[bot_name] = make_pybot(bot_name)
                case "exe":
                    out[bot_name] = make_exebot(bot_name)
                case _:
                    raise ValueError(
                        f"bot {bot_name} has unsupported extension {extension}"
                    )
        return out

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
                while run_bot("json TODO"):
                    time.sleep(1)
                    print("\033[F" * 4)
                    print((" " * 50 + "\n") * 2)
                    print("\033[F" * 4)
                    print(f"{bot_name}'s turn", flush=True)
                    print(f"current score: {self.game_state.round_score}")
                    print(f"rolled {self.roll_dice()}")
                    self.game_state.round_score += 100
                else:
                    self.game_state.bots[bot_name] += self.game_state.round_score
                    self.update_scores()


if __name__ == "__main__":
    App().main()
