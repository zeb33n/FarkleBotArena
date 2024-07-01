// package main

// import (
// 	"fmt"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// type Player struct {
// 	name  string `json:bot_name`
// 	score int    `json:round_score`
// }

// type GameState struct {
// 	players    []Player `json:bots`
// 	Numdice    int      `json:num_dice`
// 	RoundScore int      `json:round_score`
// 	Roll       []int    `json:roll`
// 	Turn       string   `json:turn`
// }

// type BoardModel struct {
// 	game  GameState
// 	board string
// }

// func InitialBoardModel() BoardModel {
// 	// This currently populates with a default valued screen whilst we
// 	// await a connection from the client
// 	defaultGameState := GameState{
// 		players: []Player{
// 			{name: "player 1", score: 0},
// 			{name: "player 2", score: 0},
// 			{name: "player 3", score: 0},
// 			{name: "player 4", score: 0},
// 		},
// 		Numdice:    6,
// 		RoundScore: 00000,
// 		Roll:       []int{},
// 		Turn:       "waiting for connections",
// 	}

// 	return BoardModel{
// 		game: defaultGameState,
// 	}

// }

// func (m BoardModel) BuildBoard() string {

// 	board := (`

// 	.=-=-=-=-=-=
// 	|
// 	|      %s
// 	|
// 	|
// 	|

// 	`)

// 	return fmt.Sprintf(board, m.game.players[0].name)

// }

// func (m BoardModel) Init() tea.Cmd {
// 	return nil
// }

// func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	return nil, nil
// }

//	func (m BoardModel) View() string {
//		return m.BuildBoard()
//	}
package main
