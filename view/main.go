package main

import (
	"bytes"
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	Game "github.com/lregs/FarkleBotArena/Game"
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

type BaseModel struct {
	GameClient
	Display string
}

func InitialBaseModel() *BaseModel {
	return &BaseModel{
		GameClient: Game.NewClient(),
		Display:    "Press C To Connect To a New Game!",
	}
}

// method on pointer because we're reading from the channels I THINK?!s
// will return a msg when data is recieved in the channel and trigger the bt update functions
// which will loop back into this
// func (m *BoardModel) monitorChannels() tea.Cmd {
// 	return func() tea.Msg {
// 		select {
// 		case data := <-m.tcpDataChan:
// 			return tcpResponse(data)
// 		case err := <-m.tcpErrChan:
// 			return tcpReadError(err.Error())
// 		}
// 	}

// }

// func (m *BoardModel) readCmd() tea.Cmd {
// 	return func() tea.Msg {
// 		m.log.Print("reading")
// 		buffer := make([]byte, 256)
// 		n, err := m.tcp.Read(buffer)
// 		if err != nil {
// 			return tcpReadError(err.Error())
// 		}

// 		cleanedBuff := []byte{}

// 		for _, b := range buffer[:n] {
// 			if !(b == 0) {
// 				cleanedBuff = append(cleanedBuff, b)
// 			} else {
// 				break
// 			}
// 		}

// 		m.log.Print((string(tcpResponse(cleanedBuff))))

// 		return tcpResponse(cleanedBuff)

// 	}
// }

//

func main() {

	log := NewLogger("log.txt")

	log.Println("new log")

	// includes placeholder values and initialised tcp connection on the mode
	m := InitialBaseModel()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Printf("something bad happened %s", err)
	}

}

// init starts reading immediately, server and game(?) need to be started atm for it to work as it will
// return a nil pointer panic if there is no connection for it to read from
func (m BaseModel) Init() tea.Cmd {
	m.startReading()
	return m.monitorChannels()
}

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "1":
			return m, m.sendResponse()

		}

	// case startReading:
	// 	return m, m.readCmd()

	case tcpResponse:
		var gs GameState
		r := bytes.NewReader(msg)
		if err := json.NewDecoder(r).Decode(&gs); err != nil {
			m.log.Printf("decoding failed %s", err)
		}
		m.screen = BuildBoard(gs)
		return m, m.monitorChannels()

	}
	return m, nil
}

func (m BaseModel) View() string {
	// m.log.Print(m.screen)
	return m.Display
}

// default model to be displayed by bt
// func InitialBoardModel(log *log.Logger) (BoardModel, error) {

// 	conn, err := net.Dial("tcp", "localhost:4123")
// 	if err != nil {
// 		return BoardModel{screen: "failed to connect"}, err
// 	}

// 	log.Printf("connected success %v", conn)

// 	defaultGameState := c.GameState{
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
// 		game:        defaultGameState,
// 		screen:      BuildBoard(defaultGameState),
// 		tcp:         conn,
// 		log:         log,
// 		tcpDataChan: make(chan []byte),
// 		tcpErrChan:  make(chan error),
// 	}, nil

// }
