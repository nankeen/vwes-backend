package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nankeen/vwes-backend/controllers"
)

// SetupRouter creates the gin router and binds handlers
func SetupRouter() *gin.Engine {
	r := gin.Default()

	roomController := controllers.NewRoomController()
	r.GET("/api/rooms/:room/", roomController.GetGameByID)
	r.GET("/ws/new/", roomController.CreateRoom)
	r.GET("/ws/rooms/:room/", roomController.JoinRoom)
	return r
}
