package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var addr = ":12345"

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn = make(map[string]*websocket.Conn)
)

func Echo(c echo.Context) error {
	userID := c.Param("id")

	if len(userID) == 0 {
		return c.String(http.StatusBadRequest, "Param error")
	}

	connect, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "upgrade: "+err.Error())
	}
	defer connect.Close()

	log.Printf("Connected: %s", userID)
	conn[userID] = connect

	for {
		mt, message, err := connect.ReadMessage()
		if err != nil {
			log.Println("read:", err)

			break
		}

		log.Printf("[%d] recv: %s", mt, message)
	}

	return c.String(http.StatusBadRequest, "Error")
}

func Alert(c echo.Context) error {
	id := c.FormValue("id")
	msg := c.FormValue("message")
	userName := c.FormValue("username")
	amount := c.FormValue("amount")
	coin := c.FormValue("coin")
	duration := c.FormValue("duration")
	fadeout := c.FormValue("fadeout")
	fadein := c.FormValue("fadein")
	img := c.FormValue("img")

	sendData, _ := json.Marshal(map[string]interface{}{
		"id":       id,
		"message":  msg,
		"username": userName,
		"amount":   amount,
		"coin":     coin,
		"duration": duration,
		"fadein":   fadein,
		"fadeout":  fadeout,
		"img":      img,
	})

	err := conn[id].WriteMessage(websocket.TextMessage, sendData)
	if err != nil {
		return c.String(http.StatusBadRequest, "write: "+err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

func main() {
	log.Println("Running on " + addr)
	e := echo.New()

	e.GET("/echo/:id", Echo)
	e.POST("/alert", Alert)

	e.Logger.Fatal(e.Start(addr))
}
