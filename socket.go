package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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

		if msg == "" {
			continue
		}

		if msg == "2" { // Server ping
			c.WriteMessage(websocket.TextMessage, []byte("3"))
			continue
		}

		switch string(msg[4]) {
		case "X":
			movePlayer(c, msg)
			continue
		case "M":
			sendMessage(c, msg)
			continue
		}
	}
}

func movePlayer(c *websocket.Conn, msg string) {
	msg = msg[7 : len(msg)-1]
	var info map[string]interface{}
	json.Unmarshal([]byte(msg), &info)

	if info["i"] != playerID {
		return
	}

	x := info["x"].(float64) - 30
	y := info["y"].(float64)

	xStr := strconv.Itoa(int(x))
	yStr := strconv.Itoa(int(y))

	c.WriteMessage(websocket.TextMessage, []byte(`42["moveTo", `+xStr+`, `+yStr+`]`))
}

func sendMessage(c *websocket.Conn, msg string) {
	msg = msg[7 : len(msg)-1]
	var info map[string]interface{}
	json.Unmarshal([]byte(msg), &info)

	if info["i"] != playerID {
		return
	}

	msg = info["m"].(string)
	c.WriteMessage(websocket.TextMessage, []byte(`42["message", "`+msg+`"]`))
}
