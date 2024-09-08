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

type GameClient interface {
	Connect(addr string) error
	Read() (<-chan []byte, <-chan error) // returns channel holding state of the game or error
	Respond([]byte) error                // writes a response to the server - 1 to play atm

}

type tcpResponse []byte

type tcpReadError string

// Base model will hold the UI and Game client. It will be monitored within bubbletea loop
// and call ui and game methods based on user input

type BaseModel struct {
	client     GameClient // should be using the GameClient interfacE?!
	UI         *UI.UI     // this naming is horrible :)
	Display    string
	GameData   <-chan []byte
	GameErrors <-chan error
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

		err := m.client.Respond([]byte{'1'})
		if err != nil {
			return FailedRespondingToServer{}
		}
		return SuccessfulResponse{}
	}

}

func (m *BaseModel) monitorChannels() tea.Cmd {
	return func() tea.Msg {
		select {
		case data := <-m.GameData:
			return tcpResponse(data)
		case err := <-m.GameErrors:
			return tcpReadError(err.Error())
		}
	}

}

func (m *BaseModel) Init() tea.Cmd {
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
		m.UI.CurrState = UI.SuccessfulConnection
		m.GameData, m.GameErrors = m.client.Read()
		// do we want to be returning and looking up channels each update, or is there a better way to do this?
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
