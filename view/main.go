package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	Game "github.com/lregs/FarkleBotArena/Game"
	UI "github.com/lregs/FarkleBotArena/UI"
	c "github.com/lregs/FarkleBotArena/common"
)

// pls can we make a customise tab where we can change the dice colour or something would be awesome

type GameClient interface {
	Connect(addr string) error // establish a connection type within implementation
	Read()                     // gets the first state from the game
	Respond([]byte) error      // writes a response to the server - 1 to play atm
}

type startReading struct{}

type tcpResponse []byte

type tcpReadError string

type userinput string

// Base model will hold the UI and Game client. It will be monitored within bubbletea loop
// and call ui and game methods based on user input

type BaseModel struct {
	client  *Game.Client
	UI      *UI.UI // this naming is horrible :)
	Display string
}

func InitialBaseModel(log *log.Logger) *BaseModel {
	return &BaseModel{
		client: Game.NewClient(),
		UI:     UI.NewUI(log), // shouldnt be creating a new logger here
	}
}

type ConnectionSuccess struct{}

// currently this will get stuck on connection failed and like not listen for a new connection or
// something - not sure if its to do with ui state conditional within the update function
type ConnectionFailed struct{ err string }

func (m *BaseModel) AttemptConnection(addr string) tea.Msg {
	err := m.client.Connect(addr)
	if err != nil {
		return ConnectionFailed{err: err.Error()}
	}
	return ConnectionSuccess{}
}

type FailedRespondingToServer struct{}
type SuccessfulResponse struct{}

func (m *BaseModel) sendResponse() tea.Cmd {
	return func() tea.Msg {

		err := m.client.Respond([]byte{1})
		if err != nil {
			return FailedRespondingToServer{}
		}
		return SuccessfulResponse{}
	}

}

func (m *BaseModel) monitorChannels() tea.Cmd {
	return func() tea.Msg {
		select {
		case data := <-m.client.DataCh:
			return tcpResponse(data)
		case err := <-m.client.ErrCh:
			return tcpReadError(err.Error())
		}
	}

}

func (m *BaseModel) Init() tea.Cmd {
	// returning nil because nothing is needed at the begining
	return nil
}

func (m *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "c":
			if m.UI.CurrState == UI.WelcomeState {
				// Pass the returned message to the Bubble Tea framework
				return m, func() tea.Msg {
					return m.AttemptConnection("localhost:4123")
				}
			}
		case "1":
			if m.UI.CurrState == UI.SuccessfulConnection {
				return m, m.sendResponse()
			}
		}
		return m, nil
	case ConnectionFailed:
		m.UI.CurrState = UI.FailedConnection
		return m, nil
	case ConnectionSuccess:
		// we want to render the board and start waiting for the game to start basically
		m.UI.CurrState = UI.SuccessfulConnection
		m.client.Read()
		return m, m.monitorChannels()

	case FailedRespondingToServer:

	case tcpResponse:
		var gs c.GameData
		r := bytes.NewReader(msg)
		if err := json.NewDecoder(r).Decode(&gs); err != nil {
			fmt.Print("do something")
		}
		m.UI.Data = gs
		m.UI.CurrState = UI.GameLive
		return m, m.monitorChannels()

	}

	return m, nil
}

func (m *BaseModel) View() string {
	// m.log.Print(m.screen)
	return m.UI.Render()
}

func main() {

	log := c.NewLogger("log.txt")

	log.Println("new log")

	// includes placeholder values and initialised tcp connection on the mode
	m := InitialBaseModel(log)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Printf("something bad happened %s", err)
	}

}

// default model to be displayed by bt
// func InitialBoardModel(log *log.Logger) (BoardModel, error) {

// 	conn, err := net.Dial("tcp", "localhost:4123")
// 	if err != nil {
// 		return BoardModel{screen: "failed to connect"}, err
// 	}

// 	log.Printf("connected success %v", conn)

// 	defaultGameData := c.GameData{
// 		Players: []Player{
// 			{Name: "player 1", Score: 0},
// 			{Name: "player 2", Score: 0},
// 			{Name: "player 3", Score: 0},
// 			{Name: "player 4", Score: 0},
// 		},
// 		Numdice:    6,
// 		RoundScore: 00000,
// 		Roll:       []int{1, 2, 3, 4, 5, 6},
// 		Turn:       "waiting for connections",
// 	}

// 	return BoardModel{
// 		game:        defaultGameData,
// 		screen:      BuildBoard(defaultGameData),
// 		tcp:         conn,
// 		log:         log,
// 		tcpDataChan: make(chan []byte),
// 		tcpErrChan:  make(chan error),
// 	}, nil

// }
