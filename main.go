package main

var playerID string

func main() {
	sessionTicket := doAuth()
	playerID = askForInput("Player ID to follow:")
	socket := createConnection()
	socket.ReadMessage() // For the initial SID message

	doLogin(socket, sessionTicket)

	go messageListener(socket)

	for {
	}
}
