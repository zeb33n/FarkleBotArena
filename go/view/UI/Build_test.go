package ui

import (
	"log"
	"strings"
	"testing"

	c "github.com/lregs/FarkleBotArena/common"
)

func TestBuildUI_Appearance(t *testing.T) {
	// Create a logger for the UI (can be nil for testing)
	logger := log.New(nil, "", log.LstdFlags)

	// Create a default game state for testing
	defaultState := c.GameState{
		Players: []c.Player{
			{Name: "player 1", Score: 0},
			{Name: "player 2", Score: 0},
			{Name: "player 3", Score: 0},
			{Name: "player 4", Score: 0},
		},
		Roll: []int{1, 2, 3, 4, 5},
		Turn: "player 1",
	}

	// Initialize the UI
	ui := NewUI(logger)

	// Test BuildUI with the default game state
	result := ui.BuildUI(defaultState)

	// Define the expected output for the given game state
	expected := `
player 1/---------------------------------------------\player 2                   147
.=-=-=-=-=-=\              FARKLE BOT ARENA               /=-=-=-=-=-=.
|         0/---------------------------------------------\         0|
|                              |=====||=====||=====||=====||=====|                              |
|                              |  1  ||  2  ||  3  ||  4  ||  5  |                              |
|                              |=====||=====||=====||=====||=====|                              |
|         0         0|
.=-=-=-=-=-=\              R OR P               /=-=-=-=-=-=.
player 3/---------------------------------------------\player 4
`
	// Clean up the formatting
	expected = cleanOutput(expected)
	result = cleanOutput(result)

	// Compare the output to the expected output
	if result != expected {
		t.Errorf("Expected UI output:\n%s\nBut got:\n%s", expected, result)
	}
}

func TestBuildDice_Appearance(t *testing.T) {
	// Test buildDice with a sample roll
	dice := []int{1, 2, 3, 4, 5}
	result := buildDice(dice)

	// Define the expected output for the given dice roll
	expected := `
|                              |=====||=====||=====||=====||=====|                              |
|                              |  1  ||  2  ||  3  ||  4  ||  5  |                              |
|                              |=====||=====||=====||=====||=====|                              |
`
	// Clean up the formatting
	expected = cleanOutput(expected)
	result = cleanOutput(result)

	// Compare the output to the expected output
	if result != expected {
		t.Errorf("Expected dice output:\n%s\nBut got:\n%s", expected, result)
	}
}

// cleanOutput helps to remove unnecessary leading/trailing whitespace and standardize newlines.
func cleanOutput(s string) string {
	return strings.TrimSpace(s)
}
