package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TitleID  string `json:"titleId"`
}

func doAuth() interface{} {
	email := askForInput("Email:")
	password := askForInput("Password:")

	data, _ := json.Marshal(login{
		Email:    email,
		Password: password,
		TitleID:  "5417",
	})

	resp, err := http.Post("https://null.playfabapi.com/Client/LoginWithEmailAddress",
		"application/json", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Couldn't login - maybe the password was wrong?")
		doAuth()
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var loginResponse map[string]map[string]interface{}
	json.Unmarshal(body, &loginResponse)

	return loginResponse["data"]["SessionTicket"]
}
