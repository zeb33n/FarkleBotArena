package main

import "testing"

func TestMakeBoard(t *testing.T) {

	t.Run("Test board", func(t *testing.T) {
		print(BuildBoard(GameState{
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
		}))
	})
}
