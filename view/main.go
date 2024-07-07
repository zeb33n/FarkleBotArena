package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

/*

TODO:
[] Change playercolour based on turn
[] tcp client
[] figure how we want the board to be represented to best allow access and change in the view
   during the loop of the game
[]
*/

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

type tcpClient struct {
	conn net.Conn
}

func NewTCPClient() (*tcpClient, error) {

	conn, err := net.Dial("tcp", "localhost:8990")
	if err != nil {
		return &tcpClient{}, err
	}
	return &tcpClient{conn: conn}, nil

}

func (c tcpClient) Read() ([]byte, error) {
	buffer := make([]byte, 256)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return []byte{}, err
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

	return cleanedBuff, nil

}

type BoardModel struct {
	game   GameState
	screen string
	// TcpModel tcpClient
}

func InitialBoardModel() BoardModel {
	// This currently populates with a default valued screen whilst we
	// await a connection from the client
	defaultGameState := GameState{
		Players: []Player{
			{Name: "something", Score: 0},
			{Name: "22", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Numdice:    4,
		RoundScore: 00000,
		Roll:       []int{1, 2, 3, 4, 5, 6},
		Turn:       "waiting for connections",
	}

	return BoardModel{
		game:   defaultGameState,
		screen: BuildBoard(defaultGameState),
	}

}

func BuildBoard(State GameState) string {

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, State.Players[0].Name, State.Players[1].Name))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, State.Players[0].Score, State.Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(State.Roll))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`|%11d\---------------------------------------------/%-11d|`, State.Players[2].Score, State.Players[2].Score))
	sb.WriteString("\n")
	// being lazy, needs padding, maybe need a padding func
	sb.WriteString(`.=-=-=-=-=-=\              Press R or P               /=-=-=-=-=-=.` + "\n")
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, State.Players[2].Name, State.Players[2].Name))

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

	client, _ := NewTCPClient()

	// bm := InitialBoardModel()
	// fmt.Println(bm.screen)

	welcomeMessage := "waiting for game to start"

	for {
		var gs GameState
		response, _ := client.Read()
		if string(response) == "" {
			fmt.Printf("empty")
			continue
		}
		if string(response) == welcomeMessage {
			fmt.Printf("%s\n", (string(response)))
		} else {

			if err := json.Unmarshal(response, &gs); err != nil {
				fmt.Printf("failed decoding %s", err)
			}

			fmt.Print(BuildBoard(gs))

		}

	}

}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}
