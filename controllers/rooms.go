package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nankeen/vwes-backend/game"
	"github.com/rs/xid"
	"log"
	"net/http"
)

// RoomController handles CRUD operations for rooms
type RoomController struct {
	games      map[string]*game.Game
	wsupgrader *websocket.Upgrader
}

func NewRoomController() RoomController {
	return RoomController{
		games: make(map[string]*game.Game),
		wsupgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// GetGameByID retrives a game session by it's room ID
func (rc *RoomController) GetGameByID(c *gin.Context) {
	id := c.Param("room")

	game := rc.games[id]
	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": "Unable to find that room",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "Room is online",
		"playerCount": game.PlayerCount(),
		"joinable":    game.PlayerCount() < 2,
	})
}

func (rc *RoomController) JoinRoom(c *gin.Context) {
	var msg game.Message
	id := c.Param("room")

	game := rc.games[id]
	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": "Can't find that room",
		})
		return
	}

	conn, err := rc.wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	playerID, err := game.AddPlayer(conn)
	if err != nil {
		conn.WriteJSON(gin.H{
			"status": "room is full",
		})
		conn.Close()
		return
	}

	// Send hello handshake with player ID
	conn.WriteJSON(gin.H{
		"status": "connected",
		"player": *playerID,
	})

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		msg.Player = *playerID
		game.Swing(msg)
	}
}

func (rc *RoomController) CreateRoom(c *gin.Context) {
	gc, err := rc.wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	// Generate room id
	uid := xid.New()
	if rc.games[uid.String()] != nil {
		log.Println("UID clash")
		gc.Close()
		return
	}

	// Create a game
	g := game.NewGame(gc)
	rc.games[uid.String()] = &g
	// Return room id

	for {
		_, _, err := gc.ReadMessage()
		if err != nil {
			log.Println("Game %+v disconnected", uid)
			g.End()
			return
		}
	}
}
