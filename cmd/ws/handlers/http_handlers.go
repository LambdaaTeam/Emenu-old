package handlers

import (
	"net/http"

	"github.com/LambdaaTeam/Emenu/cmd/ws/services"
	"github.com/LambdaaTeam/Emenu/cmd/ws/shared"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*Client]bool)

func Notify(c *gin.Context) {
	var packet shared.Packet

	err := c.BindJSON(&packet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	res := services.HandlePacket(ctx, &packet)

	if res.Type == shared.Error {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	BroadcastMessage(res)

	c.JSON(http.StatusOK, res)
}

func UpgradeConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	go HandleWebSocketConnection(ctx, conn)

}
