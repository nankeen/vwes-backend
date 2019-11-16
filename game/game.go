package game

import (
	"errors"
	"github.com/gorilla/websocket"
)

// Game handles a game's communication
type Game struct {
	gameConn    *websocket.Conn
	playerConns [2]*websocket.Conn
	broadcast   chan Message
}

func NewGame(gc *websocket.Conn) Game {
	bc := make(chan Message)
	return Game{
		gameConn:  gc,
		broadcast: bc,
	}
}

func (g *Game) PlayerCount() int {
	return len(g.playerConns)
}

func (g *Game) AddPlayer(p *websocket.Conn) (*int, error) {

	for i, conn := range g.playerConns {
		if conn == nil {
			g.playerConns[i] = p
			return &i, nil
		}
	}

	return nil, errors.New("Room is full")
}

func (g *Game) Swing(msg Message) {
	g.broadcast <- msg
}
