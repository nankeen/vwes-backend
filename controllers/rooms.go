package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/matoous/go-nanoid"
	"github.com/nankeen/vwes-backend/game"
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
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
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
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	playerID, err := game.AddPlayer(conn)
	if err != nil {
		log.Println("Room %v is full", id)
		conn.WriteJSON(gin.H{
			"status": "room is full",
		})
		conn.Close()
		return
	}

	// Send hello handshake with player ID
	err = conn.WriteJSON(gin.H{
		"status": "connected",
		"player": playerID,
	})

	if err != nil {
		log.Println("Can't send status to client", err)
	}

	log.Printf("Player %v connected to room %v\n", playerID, id)

	for {
		err := conn.ReadJSON(&msg)
		log.Printf("Got message from player: %+v", msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("Player disconnected: %v\n", playerID)
				game.RemovePlayer(playerID)
				break
			}
			log.Println("Error parsing player message", err)
			continue
		}
		msg.Player = playerID
		game.Swing(msg)
	}
}

func (rc *RoomController) CreateRoom(c *gin.Context) {
	gc, err := rc.wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	// Generate room id
	uid, err := gonanoid.Generate("0123456789", 4)
	if err != nil {
		log.Println(err)
		gc.Close()
		return
	}

	if rc.games[uid] != nil {
		log.Println("UID clash")
		gc.Close()
		return
	}

	// Create a game
	g := game.NewGame(gc)
	rc.games[uid] = &g
	// Return room id

	gc.WriteJSON(game.RoomInfo{
		RoomID:           uid,
		PlayersConnected: g.PlayerCount(),
	})

	for {
		_, _, err := gc.ReadMessage()
		if err != nil {
			log.Println("Game %+v disconnected", uid)
			g.End()
			return
		}
	}
}
