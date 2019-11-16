package game

// Message represents a message sent by player
type Message struct {
	Player int    `json:"-"`
	Action string `json:"action"`
}
