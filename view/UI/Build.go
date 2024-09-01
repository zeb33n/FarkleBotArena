package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	c "github.com/lregs/FarkleBotArena/common"
)

// I dont know if I want to track state at base model or ui level but atm its the state of the ui
// so I keep here I think
type state int

const (
	WelcomeState = iota
	Playing
	FailedConnection
	SuccessfulConnection
	FailedResponse
	GameLive
)

// This struct will hold all the methods for rendering the different pages within the UI.
type UI struct {
	CurrState state
	Data      c.GameData
	log       *log.Logger
}

func NewUI(log *log.Logger) *UI {
	return &UI{
		log: log,
	}
}

func (u *UI) Update(newState state, data c.GameData) {
	u.CurrState = newState
	u.Data = data
}

func (u *UI) Render() string {
	switch u.CurrState {
	case WelcomeState:
		return u.renderWelcomeState()
	case FailedConnection:
		return u.renderFailedConnection()
	case SuccessfulConnection:
		return u.renderSuccessfulConnection()
	case FailedResponse:
		return u.renderFailedResponse()
	default:
		return "you have failed to assign the state properly mate :)"
	}

}

func (u *UI) renderFailedResponse() string {
	return "responding to server failed"
}

func (u *UI) renderGameLive() string {
	return u.BuildUI()
}

func (u *UI) renderFailedConnection() string {
	return "Failed to connect, please try again by pressing c"
}

func (u *UI) renderSuccessfulConnection() string {
	return "Great Success, Press 1 to Start!"
}

func (u *UI) renderWelcomeState() string {
	// Todo - can make checkboxes and enter with "available" servers to join?
	return "Welcome, Press C To Connect"
}

func (u *UI) BuildUI() string {

	// width, height, err := term.GetSize(0)
	// if err != nil {
	// 	return " "
	// }

	PlayerMessage := "R OR P"
	if u.Data.Turn == "waiting for connections" {
		PlayerMessage = "Press 1 To Start Game"
	}
	// we dont need to this anymore and are just creating copies for joke every time to board is made :D
	Players := make([]c.Player, len(u.Data.Players))
	copy(Players, u.Data.Players)

	// hardcoded minimum of two players atm
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[0].Name, Players[1].Name))
	sb.WriteString(strconv.Itoa(len(sb.String())))
	sb.WriteString("\n")
	sb.WriteString(`.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.`)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`|%11d/---------------------------------------------\%-11d|`, Players[0].Score, Players[1].Score))
	sb.WriteString("\n")
	sb.WriteString(buildDice(u.Data.Roll))
	sb.WriteString("\n")

	// im br0ke
	// I think I just need to re think how this board is being made entirely cause it sucks rn and its v rigid
	// and the spacings are awkward, dice padding isnt working and its kind of a mess

	// also how can we just edit the little bits of strings that change is that better I dunno if a string builder
	// is the most ideal or we go back to our idea of representing the board in an array with each
	switch len(u.Data.Players) {
	case 3:
		sb.WriteString(fmt.Sprintf(`|%11d`, Players[2].Score))
		sb.WriteString(fmt.Sprintf(`%11d|`, 0))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`.=-=-=-=-=-=\              %s               /=-=-=-=-=-=.`, PlayerMessage) + "\n")
		sb.WriteString(fmt.Sprintf(`%12s/---------------------------------------------\%-12s`, Players[2].Name, "No Player"))

	case 4:
		sb.WriteString(fmt.Sprintf(`|%11d`, Players[2].Score))
		sb.WriteString(fmt.Sprintf(`%11d|`, Players[3].Score))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(`.=-=-=-=-=-=\              %s               /=-=-=-=-=-=.`, PlayerMessage) + "\n")
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
