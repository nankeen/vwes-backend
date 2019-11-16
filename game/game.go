package game

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
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
			if err := g.gameConn.WriteJSON(msg); err != nil {
				log.Println("Error sending player action to game", err)
			}
		}
	}
}

func (g *Game) PlayerCount() int {
	c := 0
	for _, pc := range g.playerConns {
		if pc != nil {
			c++
		}
	}
	return c
}

func (g *Game) AddPlayer(p *websocket.Conn) (int, error) {

	for i, conn := range g.playerConns {
		if conn == nil {
			g.playerConns[i] = p
			return i, nil
		}
	}

	return -1, errors.New("Room is full")
}

func (g *Game) Swing(msg Message) {
	g.broadcast <- msg
}

func (g *Game) RemovePlayer(id int) {
	g.playerConns[id].Close()
	g.playerConns[id] = nil
}

func (g *Game) End() {
	close(g.broadcast)
}
