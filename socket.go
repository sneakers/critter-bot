package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

func createConnection() *websocket.Conn {
	socketURL := "wss://boxcritters.herokuapp.com/socket.io/?EIO=4&transport=websocket"
	u, _ := url.Parse(socketURL)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	return c
}

func doLogin(c *websocket.Conn, sessionTicket interface{}) {
	sessionID := fmt.Sprintf("%v", sessionTicket)

	c.WriteMessage(websocket.TextMessage, []byte(`40`))
	c.ReadMessage()

	c.WriteMessage(websocket.TextMessage, []byte(`42["login","`+sessionID+`"]`))
	fmt.Println("Logging in...")
	c.ReadMessage()

	c.WriteMessage(websocket.TextMessage, []byte(`42["joinRoom", "port"]`))
	fmt.Println("Joining room...")
	c.ReadMessage()
	fmt.Println("Joined room")
}

func messageListener(c *websocket.Conn) {
	for {
		_, message, _ := c.ReadMessage()
		msg := string(message)

		if msg == "2" {
			c.WriteMessage(websocket.TextMessage, []byte("3"))
			continue
		}
		if strings.Contains(msg, `42["X",`) {
			msg = msg[7 : len(msg)-1]
			var i map[string]interface{}
			json.Unmarshal([]byte(msg), &i)

			movePlayer(c, i)
			continue
		}
	}
}

func movePlayer(c *websocket.Conn, info map[string]interface{}) {
	if info["i"] != playerID {
		return
	}

	x := info["x"].(float64) - 10
	y := info["y"].(float64)

	xStr := strconv.Itoa(int(x))
	yStr := strconv.Itoa(int(y))

	c.WriteMessage(websocket.TextMessage, []byte(`42["moveTo", `+xStr+`, `+yStr+`]`))
}
