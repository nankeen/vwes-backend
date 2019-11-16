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
	r.GET("/ws/:room/", roomController.JoinRoom)
	r.POST("/api/rooms/", roomController.CreateRoom)
	return r
}
