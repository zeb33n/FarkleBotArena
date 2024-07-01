// package main

// import tea "github.com/charmbracelet/bubbletea"

// // wrapped so we could add something like a welcome screen easily
// // even if it looks uneccassary atm
// type BaseModel struct {
// 	board BoardModel
// }

// func InitialModel() *BaseModel {
// 	return &BaseModel{
// 		board: InitalBoardModel(),
// 	}
// }

// func (m BaseModel) Init() tea.Cmd {
// 	return nil
// }

// func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		}
// 	}

// 	return nil, nil
// }

//	func (m BaseModel) View() string {
//		return m.board.BuildBoard()
//	}
package main
