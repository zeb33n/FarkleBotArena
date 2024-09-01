package main

import (
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
				// connect
			}
		}

		// case startReading:
		// 	return m, m.readCmd()

		// case tcpResponse:
		// 	var gs GameData
		// 	r := bytes.NewReader(msg)
		// 	if err := json.NewDecoder(r).Decode(&gs); err != nil {
		// 		m.log.Printf("decoding failed %s", err)
		// 	}
		// 	m.screen = BuildBoard(gs)
		// 	return m, m.monitorChannels()

		// }
		return m, nil
	}

	return m, nil
}

func (m *BaseModel) View() string {
	// m.log.Print(m.screen)
	return m.UI.Render()
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

// init starts reading immediately, server and game(?) need to be started atm for it to work as it will
// return a nil pointer panic if there is no connection for it to read from

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
