package main

import (
	"fmt"
	"io"
	"os"
)

// type Player struct {
// 	name  string `json:bot_name`
// 	score int    `json:round_score`
// }

// type GameState struct {
// 	players    []Player `json:bots`
// 	Numedice   int      `json:num_dice`
// 	RoundScore int      `json:round_score`
// 	Roll       []int    `json:roll`
// 	Turn       string   `json:turn`
// }

func main() {
	fifo := "../tmp/tmpFifo"

	fr, err := os.OpenFile(fifo, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fr.Close()

	data, err := io.ReadAll(fr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}
