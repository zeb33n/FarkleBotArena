package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Player struct {
	Name  string `json:bot_name`
	Score int    `json:round_score`
}

// If I end up using Byte[] I should decode the strings into byte slices to save one copy

type GameState struct {
	Players    []Player `json:bots`
	Numdice    int      `json:num_dice`
	RoundScore int      `json:round_score`
	Roll       []int    `json:roll`
	Turn       string   `json:turn`
}

type BoardModel struct {
	game GameState
}

func InitialBoardModel() BoardModel {
	// This currently populates with a default valued screen whilst we
	// await a connection from the client
	defaultGameState := GameState{
		Players: []Player{
			{Name: "hello", Score: 0},
			{Name: "22", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Numdice:    1,
		RoundScore: 00000,
		Roll:       []int{},
		Turn:       "waiting for connections",
	}

	return BoardModel{
		game: defaultGameState,
	}

}

func BuildBoard(State GameState) string {

	// need to find a way of adding the string in insert mode.
	// the best way for this might be add this string into 2d list of strings
	// which each element in the list representing an element on the ui
	// this would introduce name length caps though?!
	// either way then we could add rows and columns for the ui and each line stored
	// but it does feel like that makes it more complex
	// but mybe opens up more interactivity if needed in the future

	// or we need to give in and create some kind of string generator because %s will displace the string by len

	// board := (`

	//  Player One /---------------------------------------------\ Player Two
	// .=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.
	// |   10000   /---------------------------------------------\   10000   |
	// |                                                                     |

	//             |=====| |=====| |=====| |=====| |=====| |=====|
	//             |  1  | |  2  | |  3  | |  4  | |  5  | |  6  |
	//             |=====| |=====| |=====| |=====| |=====| |=====|

	// `)

	// names can be 11 characters long. Names under 11 characters will be printed from left or right respectively. len(name) - 11 amount of whitespaces will be added on to the end (or beginning)
	// of each name to build the first string

	// The 'game section' (from the score down will remain the same mostly (depending on dice numbers) A reduction in the number of dice will not cause the arena to become smaller)
	// but the dice need to be centreed

	boardWidth := 71
	boardWidth += 0

	formattedNames := make([]string, len(State.Players))

	for i, player := range State.Players {
		if len(player.Name) > 12 {
			formattedNames[i] = player.Name[:12]
		} else {
			formattedNames[i] = fmt.Sprintf("%12s", State.Players[0].Name)
		}
	}

	var sb strings.Builder
	// I dont understand why -12 isn't flipping the indent for the string>?!??!?!
	sb.WriteString(fmt.Sprintf((`%12s/---------------------------------------------\%-12s`), formattedNames[0], formattedNames[1]) + "\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, State.Players[0].Score, State.Players[1].Score) + "\n")

	// diceHead := (strings.TrimRight(strings.Repeat(fmt.Sprintf("%32s|=====| ", " "), State.Numdice), " "))

	// sb.WriteString(diceHead)

	switch State.Numdice {
	case 1:
		sb.WriteString(fmt.Sprintf("|%s|", (strings.TrimRight(strings.Repeat(fmt.Sprintf("%32s|=====|%-32s", " ", " "), State.Numdice), " ")))) // 7

	}

	return sb.String()

}

func (m BoardModel) Init() tea.Cmd {
	return nil
}

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m BoardModel) View() string {
	return BuildBoard(m.game)
}

func main() {
	// fifo := "../pipe"

	// fr, err := os.OpenFile(fifo, os.O_RDONLY, os.ModeNamedPipe)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// defer fr.Close()

	// data, err := io.ReadAll(fr)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(data))

	m := InitialBoardModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

}
