package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// This represents what will recieve from the game logic
type Player struct {
	name  string
	score int
}

type GameState struct {
	players      []Player
	dice         []int
	CurrentScore []int
}

type model struct {
	GameState *GameState
}

// This is a bubbletea method to define the initial state of the application. Atm we're using a mock model
func initialModel(state *GameState) *model {
	return &model{
		GameState: state,
	}
}

// Perform initial I/O
func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case model:
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl + c", "q":
			return m, tea.Quit
		}

	}

	return m, nil

}

func (m model) View() string {
	s := "Farkle Arena \n"
	for i, player := range m.GameState.players {
		s += fmt.Sprintf("Name: %s, Total Score: %v, Round Score: %v \n", player.name, player.score, m.GameState.CurrentScore[i])
		s += strconv.Itoa(m.GameState.dice[i])
	}

	return s
}

func main() {

	state := flag.String("stateBytes", "", "byte string for current game state")

	var gs GameState

	if err := json.Unmarshal([]byte(*state), &gs); err != nil {
		//slog
		panic(fmt.Sprintf("couldn't load state %s", err))
	}

	p := tea.NewProgram(initialModel(&gs))
	if _, err := p.Run(); err != nil {
		fmt.Printf("err %s", err)
		os.Exit(1)
	}

	p.Kill()

}
