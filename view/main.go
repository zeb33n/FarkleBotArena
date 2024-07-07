package main

import (
	"encoding/json"
	"fmt"
	rand "math/rand"
	"net"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/*

TODO:
[] Change playercolour based on turn
[] tcp client
[] figure how we want the board to be represented to best allow access and change in the view
   during the loop of the game
[]

<<<<<<< Updated upstream
*/
=======
func main() {
	fifo := "../tmp/pipe"
>>>>>>> Stashed changes

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
	roll int
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
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, State.Players[0].Name, State.Players[1].Name))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, State.Players[0].Score, State.Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(State.Roll))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`|%11d\---------------------------------------------/%-11d|`, State.Players[2].Score, State.Players[3].Score))
	sb.WriteString("\n")
	// being lazy, needs padding, maybe need a padding func
	sb.WriteString(`.=-=-=-=-=-=\              Press R or P               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, State.Players[2].Name, State.Players[3].Name))

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

<<<<<<< Updated upstream
	fPosition := fmt.Sprintf("%*s|=====|", leftPadding, " ")

	rPadding := (boardWidth - len(fPosition)) - 2 // -2 accounts for the two empty strings we add ??

	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceHeader, rPadding, " "))))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceBodyStr, rPadding, " "))))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("|%s|", (fmt.Sprintf("%*s%s%*s", leftPadding, " ", diceHeader, rPadding, " "))))

	return sb.String()
}

type tcpError struct {
	err error
}

func (t tcpError) Error() string {
	return t.err.Error()
}

type stateResponse struct{ res GameState }

func connTCP() tea.Msg {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return tcpError{err: err}
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	var gs GameState

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return tcpError{err: err}
		}

		if err := json.Unmarshal(buffer[:n], &gs); err != nil {
			return nil
		}
		return stateResponse{res: gs}
	}

}

func (m BoardModel) Init() tea.Cmd {
	return connTCP

	// all i/o (tcp connections etc) need to be returned by the init function
	// all cmds are called by the tea runtime when needed.
	// cmds run in routines andd sent to the update function for handling.

	// ares are functions that dont take args and return a msg of type any
	// if your commands need args you can include them in a closure

	return nil
}

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// The update function wants to know how to deal with the commands we
	// made and returned from init based on the message we expect back from the
	// commands

	// so we can do something like instead of having default players we could have
	// "awaiting conn" and then the update changes the model (board string) and
	// that will trigger the view to update on tcp connection

	// also need to add in some kind of button to roll the dice

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "r":
			m.game.Roll[0] = rand.Intn(6)
			return m, nil
		case "p":
			// send message saying we've passed
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

=======
	for {
		data, err := io.ReadAll(fr)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(data))
	}
>>>>>>> Stashed changes
}
