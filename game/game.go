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
	g := Game{
		gameConn:  gc,
		broadcast: bc,
	}
	go g.Start()
	return g
}

func (g *Game) Start() {
	for {
		select {
		case msg := <-g.broadcast:
			g.gameConn.WriteJSON(msg)
		}
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

func (g *Game) End() {
}
