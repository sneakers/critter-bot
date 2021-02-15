package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

type critter struct {
	Email        string
	Password     string
	SessionID    string
	PlayerID     string
	FollowingID  string
	LoginSuccess bool
	Socket       *websocket.Conn
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TitleID  string `json:"titleId"`
}

func (c *critter) login() {
	data, _ := json.Marshal(login{
		Email:    c.Email,
		Password: c.Password,
		TitleID:  "5417",
	})

	resp, err := http.Post("https://null.playfabapi.com/Client/LoginWithEmailAddress",
		"application/json", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.printMessage("failed login")
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var loginResponse map[string]map[string]interface{}
	json.Unmarshal(body, &loginResponse)

	c.LoginSuccess = true
	c.SessionID = loginResponse["data"]["SessionTicket"].(string)
	c.PlayerID = loginResponse["data"]["PlayFabId"].(string)

	c.printMessage("successful login")
}

func (c *critter) joinGame() {
	socketURL := "wss://boxcritters.herokuapp.com/socket.io/?EIO=4&transport=websocket"
	u, _ := url.Parse(socketURL)

	socket, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	c.Socket = socket
	c.Socket.ReadMessage() // Initial SID Message

	c.Socket.WriteMessage(websocket.TextMessage, []byte(`40`))
	c.Socket.ReadMessage()

	c.Socket.WriteMessage(websocket.TextMessage, []byte(`42["login","`+c.SessionID+`"]`))
	c.printMessage("logging into server...")
	c.Socket.ReadMessage()

	c.Socket.WriteMessage(websocket.TextMessage, []byte(`42["joinRoom", "port"]`))
	c.printMessage("joining port...")
	c.Socket.ReadMessage()
	c.printMessage("joined port")
}

func (c *critter) listener() {
	c.printMessage("starting listener")

	for {
		_, message, _ := c.Socket.ReadMessage()
		msg := string(message)

		if msg == "" {
			continue
		}

		if msg == "2" { // Server ping
			c.Socket.WriteMessage(websocket.TextMessage, []byte("3"))
			continue
		}

		switch string(msg[4]) {
		case "X":
			c.move(msg)
			continue
		case "M":
			c.sendMessage(msg)
			continue
		}
	}
}

func (c *critter) move(msg string) {
	info := socketMessageToMap(msg)

	if info["i"] != c.FollowingID {
		return
	}

	x := info["x"].(float64) - 30
	y := info["y"].(float64)

	xStr := strconv.Itoa(int(x))
	yStr := strconv.Itoa(int(y))

	c.Socket.WriteMessage(websocket.TextMessage, []byte(`42["moveTo", `+xStr+`, `+yStr+`]`))
}

func (c *critter) sendMessage(msg string) {
	info := socketMessageToMap(msg)

	if info["i"] != playerID {
		return
	}

	msg = info["m"].(string)
	c.Socket.WriteMessage(websocket.TextMessage, []byte(`42["message", "`+msg+`"]`))
}
