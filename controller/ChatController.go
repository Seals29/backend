package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

var upgrader = websocket.Upgrader{}

func WebSocket(c *gin.Context) {

}
func SendMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message1": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message2": conn,
	})
}
