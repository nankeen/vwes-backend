package game

// Message represents a message sent by player
type Message struct {
	Player int    `json:"playerId"`
	Action string `json:"action"`
}

type RoomInfo struct {
	RoomID           string `json:"roomId"`
	PlayersConnected int    `json:"playersConnected"`
}
