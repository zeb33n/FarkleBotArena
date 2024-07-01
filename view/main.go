package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Player struct {
	Name  string `json:bot_name`
	Score int    `json:round_score`
}

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
			{Name: "Player", Score: 0},
			{Name: "reallylongplayernamewhatthe", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Numdice:    6,
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

	lines := []string{}

	var baseHeader string = `/---------------------------------------------\`

	var leftName string
	var rightName string

	if len(State.Players[0].Name) > 12 {
		fmt.Printf("Max username length is 11. %s is too long", leftName)
		leftName = State.Players[0].Name[:12]
	} else {
		leftName = State.Players[0].Name
		for i := 0; i < (11 - len(leftName)); i++ {
			leftName += " "
		}
	}

	if len(State.Players[1].Name) > 12 {
		fmt.Printf("Max username length is 11. %s is too long", rightName)
		rightName = State.Players[1].Name[:12]
	} else {
		rightName = State.Players[1].Name
		for i := 0; i < (12 - len(rightName)); i++ {
			rightName += " "
		}
	}

	headerString := leftName + baseHeader + rightName

	titleLine := `.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.`

	lines = append(lines, headerString, titleLine)

	baseString := ""

	for _, line := range lines {
		baseString += (line + "\n")

	}

	return baseString

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
