package controller

import (
	"fmt"
	"net/http"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var conns = map[string]*websocket.Conn{}

func GetAllMSg(c *gin.Context) {
	var body struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)

	msg1 := []models.Message{}
	// msg2 := []models.Message{}
	msg2 := []models.Message{}
	config.DB.Where("sender_id = ?", body.From).Where("recipient_id = ?", body.To).Find(&msg1)
	config.DB.Where("sender_id = ?", body.To).Where("recipient_id= ?", body.From).Find(&msg2)
	// allmsg := append(msg1, msg2...)
	c.JSON(200, gin.H{
		"data1": msg1,
		"data2": msg2,
	})
	return
	// config.DB.Find(&msg)
}
func GetAllShopMsg(c *gin.Context) {
	var body struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)

	msg1 := []models.Message{}
	// msg2 := []models.Message{}
	msg2 := []models.Message{}
	config.DB.Where("sender_id = ?", body.From).Where("recipient_id = ?", body.To).Where("type = ?", "shop").Find(&msg1)
	config.DB.Where("sender_id = ?", body.To).Where("recipient_id= ?", body.From).Where("type = ?", "shop").Find(&msg2)
	// allmsg := append(msg1, msg2...)
	c.JSON(200, gin.H{
		"data1": msg1,
		"data2": msg2,
	})
	return
	// config.DB.Find(&msg)
}
func SendingMessage(c *gin.Context) {
	var ReqFrom string
	h := http.Header{}
	fmt.Println("connected")
	// fmt.Println(h)
	for _, sub := range websocket.Subprotocols(c.Request) {
		h.Set("Sec-Websocket-Protocol", sub)
		// fmt.Println(sub)
		ReqFrom = sub
		// fmt.Println(sub)
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, h)
	if err != nil {
		fmt.Println(err)
	}
	conns[ReqFrom] = ws
	for {
		var req models.Message
		err = ws.ReadJSON(&req)
		if req.SenderID != "" {
			config.DB.Create(&req)
		}

		conns[req.SenderID] = ws
		if con, ok := conns[req.RecipientID]; ok {
			err = con.WriteJSON(&req)
			if err != nil {
				fmt.Println(err)
			}
		}
		if con, ok := conns[req.SenderID]; ok {
			err = con.WriteJSON(&req)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
func GetAllMsg(c *gin.Context){
	messages := []models.Message{}
	config.DB.Find(&messages)
	c.JSON(200,&messages)
}
