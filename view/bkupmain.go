package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// // This represents what will recieve from the game logic
// type Player struct {
// 	name  string
// 	score int
// }

// type GameState struct {
// 	players      []Player
// 	dice         []int
// 	CurrentScore int
// }

// type model struct {
// 	GameState *GameState
// }

// // This is a bubbletea method to define the initial state of the application. Atm we're using a mock model
// func initialModel(state *GameState) *model {
// 	return &model{
// 		GameState: state,
// 	}
// }

// // Perform initial I/O
// func (m model) Init() tea.Cmd {

// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case model:
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl + c", "q":
// 			return m, tea.Quit
// 		}

// 	}

// 	return m, nil

// }

// func (m model) View() string {
// 	s := "Farkle Arena \n"
// 	for i, player := range m.GameState.players {
// 		s += fmt.Sprintf("Name: %s, Total Score: %v, Round Score: %v \n", player.name, player.score, m.GameState.CurrentScore)
// 		s += strconv.Itoa(m.GameState.dice[i])
// 	}

// 	return s
// }

// func main() {

// 	pipeName := "game_pipe"
// 	file, err := os.OpenFile(pipeName, os.O_RDONLY, 0)
// 	if err != nil {
// 		// could manage more gracefully with maybe a "waiting for game connection timer"
// 		panic(err)
// 	}

// 	fmt.Println("pipe open")

// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	for {
// 		line, err := r.ReadBytes('\n')
// 		if err != nil {
// 			// actually once we have finished reading from the buffer we don't
// 			// neccassarily want to exit because the game will be calculated
// 			// quicker than we care to render it

// 			// so we want to read all the game data into the buffer, but display
// 			// at our own pace and exit when we decide pls
// 			os.Exit(0)
// 		}

// 		var welcome string
// 		if err := json.Unmarshal(line, &welcome); err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(welcome)

// 		// var gs GameState

// 		// if err := json.Unmarshal(line, &gs); err != nil {
// 		// 	//slog
// 		// 	panic(fmt.Sprintf("couldn't load state %s", err))
// 		// }
// 		// fmt.Println(gs)

// 		// p := tea.NewProgram(initialModel(&gs))
// 		// if _, err := p.Run(); err != nil {
// 		// 	fmt.Printf("err %s", err)
// 		// 	os.Exit(1)
// 		// }

// 		// p.Kill()

// 	}

// }
