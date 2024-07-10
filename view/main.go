package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// pls can we make a customise tab where we can change the dice colour or something would be awesome

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type GameState struct {
	Players    []Player `json:"bots"`
	Numdice    int      `json:"num_dice"`
	RoundScore int      `json:"round_score"`
	Roll       []int    `json:"roll"`
	Turn       string   `json:"turn"`
}

// im guesing well move tcp out of this at some point and leave this to handle ui
type BoardModel struct {
	log         *log.Logger
	game        GameState
	screen      string
	tcp         net.Conn
	tcpDataChan chan []byte
	tcpErrChan  chan error
}

type startReading struct{}

type tcpResponse []byte

type tcpReadError string

// reads from the tcp connection inside a go routine and sends the data it recieves to channels on the model struct
func (m *BoardModel) startReading() {
	go func() {
		buffer := make([]byte, 256)
		for {
			n, err := m.tcp.Read(buffer)
			if err != nil {
				m.tcpErrChan <- err
				// do we want to return if we have an error reading?
			}

			cleanedBuff := []byte{}
			for _, b := range buffer[:n] {
				if !(b == 0) {
					cleanedBuff = append(cleanedBuff, b)
				} else {
					break
				}
			}
			m.tcpDataChan <- cleanedBuff

		}
	}()

}

// method on pointer because we're reading from the channels I THINK?!s
// will return a msg when data is recieved in the channel and trigger the bt update functions
// which will loop back into this
func (m *BoardModel) monitorChannels() tea.Cmd {
	return func() tea.Msg {
		select {
		case data := <-m.tcpDataChan:
			return tcpResponse(data)
		case err := <-m.tcpErrChan:
			return tcpReadError(err.Error())
		}
	}

}

func (m *BoardModel) readCmd() tea.Cmd {
	return func() tea.Msg {
		m.log.Print("reading")
		buffer := make([]byte, 256)
		n, err := m.tcp.Read(buffer)
		if err != nil {
			return tcpReadError(err.Error())
		}

		cleanedBuff := []byte{}

		for _, b := range buffer[:n] {
			if !(b == 0) {
				cleanedBuff = append(cleanedBuff, b)
			} else {
				break
			}
		}

		m.log.Print((string(tcpResponse(cleanedBuff))))

		return tcpResponse(cleanedBuff)

	}
}

// default model to be displayed by bt
func InitialBoardModel(log *log.Logger) (BoardModel, error) {

	conn, err := net.Dial("tcp", "localhost:8990")
	if err != nil {
		return BoardModel{screen: "failed to connect"}, err
	}

	log.Printf("connected success %v", conn)

	defaultGameState := GameState{
		Players: []Player{
			{Name: "player 1", Score: 0},
			{Name: "player 2", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Numdice:    6,
		RoundScore: 00000,
		Roll:       []int{1, 2, 3, 4, 5, 6},
		Turn:       "waiting for connections",
	}

	return BoardModel{
		game:        defaultGameState,
		screen:      BuildBoard(defaultGameState),
		tcp:         conn,
		log:         log,
		tcpDataChan: make(chan []byte),
		tcpErrChan:  make(chan error),
	}, nil

}

func BuildBoard(State GameState) string {

	// we dont need to this anymore and are just creating copies for joke every time to board is made :D
	Players := make([]Player, len(State.Players))
	copy(Players, State.Players)

	// hardcoded minimum of two players atm
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[0].Name, Players[1].Name))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, Players[0].Score, Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(State.Roll))
	sb.WriteString("\n")

	// im br0ke
	// I think I just need to re think how this board is being made entirely cause it sucks rn and its v rigid
	// and the spacings are awkward, dice padding isnt working and its kind of a mess

	// also how can we just edit the little bits of strings that change is that better I dunno if a string builder
	// is the most ideal or we go back to our idea of representing the board in an array with each
	switch len(State.Players) {
	case 3:
		sb.WriteString(fmt.Sprintf(`|%11d`, Players[2].Score))
		sb.WriteString(fmt.Sprintf(`%11d|`, 0))
		sb.WriteString("\n")
		sb.WriteString(`.=-=-=-=-=-=\              Press R or P               /=-=-=-=-=-=.` + "\n")
		sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[2].Name, "No Player"))

	case 4:
		sb.WriteString(fmt.Sprintf(`|%11d`, Players[2].Score))
		sb.WriteString(fmt.Sprintf(`%11d|`, Players[3].Score))
		sb.WriteString("\n")
		sb.WriteString(`.=-=-=-=-=-=\              Press R or P               /=-=-=-=-=-=.` + "\n")
		sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[2].Name, Players[3].Name))
	}

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

func main() {

	log := NewLogger("log.txt")

	log.Println("new log")

	// includes placeholder values and initialised tcp connection on the mode
	m, err := InitialBoardModel(log)
	log.Printf("new model %v", m)
	if err != nil {
		log.Printf("tcp failed %s", err)
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Printf("something bad happened %s", err)
	}

}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}

// init starts reading immediately, server and game(?) need to be started atm for it to work as it will
// return a nil pointer panic if there is no connection for it to read from
func (m BoardModel) Init() tea.Cmd {
	m.startReading()
	return m.monitorChannels()
}

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
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

func (m BoardModel) View() string {
	return m.screen
}
