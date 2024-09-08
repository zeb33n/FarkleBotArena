package common

type GameData struct {
	Players    []Player `json:"bots"`
	Numdice    int      `json:"num_dice"`
	RoundScore int      `json:"round_score"`
	Roll       []int    `json:"roll"`
	Turn       string   `json:"turn"`
}
