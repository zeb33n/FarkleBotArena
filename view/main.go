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
			{Name: "something", Score: 0},
			{Name: "22", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Numdice:    4,
		RoundScore: 00000,
		Roll:       []int{1, 2, 3, 4, 5, 6},
		Turn:       "waiting for connections",
	}

	return BoardModel{
		game: defaultGameState,
	}

}

func BuildBoard(State GameState) string {

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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, "player one", "somestr"))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, State.Players[0].Score, State.Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(State.Roll))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`|%11d\---------------------------------------------/%-11d|`, State.Players[2].Score, State.Players[3].Score))
	sb.WriteString("\n")
	// will change to current score but cba with the spacing this second
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, "plyr3", "somestr"))

	return sb.String()

}

func buildDice(dice []int) string {
	boardWidth := 71

	var sb strings.Builder

	diceHeader := strings.Repeat("|=====|", len(dice))
	diceBodyStr := ""
	for _, die := range dice {
		diceBodyStr += fmt.Sprintf("|  %d  |", die)
	}

	//middle dice placement = middle - 3?
	leftPadding := ((boardWidth / 2) - (len(diceHeader) / 2))
	// rightPadding := (boardWidth - (leftPadding - 4))

	fPosition := fmt.Sprintf("%*s|=====|", leftPadding, " ")

	rPadding := (boardWidth - len(fPosition)) - 2 // -2 accounts for the two empty strings we add ??

	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceHeader, rPadding, " "))))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceBodyStr, rPadding, " "))))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceHeader, rPadding, " "))))

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

	m := InitialBoardModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

}
