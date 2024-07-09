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

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// If I end up using Byte[] I should decode the strings into byte slices to save one copy

type GameState struct {
	Players    []Player `json:"bots"`
	Numdice    int      `json:"num_dice"`
	RoundScore int      `json:"round_score"`
	Roll       []int    `json:"roll"`
	Turn       string   `json:"turn"`
}

// im guesing well move tcp out of this at some point and leave this to handle ui
type BoardModel struct {
	log     *log.Logger
	game    GameState
	screen  string
	tcp     net.Conn
	tcpData chan []byte
	tcpErr  chan error
}

type startReading struct{}

type tcpResponse []byte

type tcpReadError string

func (m *BoardModel) readCmd() tea.Cmd {
	return func() tea.Msg {
		m.log.Print("reading")
		buffer := make([]byte, 256)
		n, err := m.tcp.Read(buffer)
		if err != nil {
			return tcpReadError(err.Error())
		}

		// var gs GameState

		// if err := json.Unmarshal(buffer[:n], &gs); err != nil {
		// 	return GameState{}, err
		// }

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
		game:    defaultGameState,
		screen:  BuildBoard(defaultGameState),
		tcp:     conn,
		log:     log,
		tcpData: make(chan []byte),
		tcpErr:  make(chan error),
	}, nil

}

func BuildBoard(State GameState) string {

	Players := make([]Player, len(State.Players))
	copy(Players, State.Players)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[0].Name, Players[1].Name))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, Players[0].Score, Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(State.Roll))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`|%11d\---------------------------------------------/%-11d|`, Players[2].Score, Players[2].Score))
	sb.WriteString("\n")
	// being lazy, needs padding, maybe need a padding func
	sb.WriteString(`.=-=-=-=-=-=\              Press R or P               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[2].Name, Players[2].Name))

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

func (m BoardModel) Init() tea.Cmd {
	return m.readCmd()
}
func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case startReading:
		return m, m.readCmd()

	case tcpResponse:
		var gs GameState
		r := bytes.NewReader(msg)
		if err := json.NewDecoder(r).Decode(&gs); err != nil {
			m.log.Printf("decoding failed %s", err)
		}
		m.screen = BuildBoard(gs)
		return m, m.readCmd()

	}
	return m, nil
}

func (m BoardModel) View() string {
	return m.screen
}
