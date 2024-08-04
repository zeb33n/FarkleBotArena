import os


class PipeClient:
    def __init__(self, pipename: str = "pipe") -> None:
        self._pipename = pipename
        if not os.path.exists(self._pipename):
            os.mkfifo(self._pipename)

    def pipe(self, msg: str):
        with open(self._pipename, "w") as p:
            print(msg)
            p.write(f"{msg};")

    def await_pipe(self):
        with open(self._pipename) as p:
            while True:
                data = p.read()
                if data:
                    break

    def cleanup(self):
        os.remove(self._pipename)
