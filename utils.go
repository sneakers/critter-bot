package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func askForInput(input string) string {
	fmt.Println(input)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")

	return text
}

func (c *critter) printMessage(msg string) {
	fmt.Printf("[%s] %s \n", c.Email, msg)
}

func socketMessageToMap(msg string) map[string]interface{} {
	msg = msg[7 : len(msg)-1]
	var info map[string]interface{}
	json.Unmarshal([]byte(msg), &info)

	return info
}

func getAccounts() []critter {
	bs, err := ioutil.ReadFile("accounts/accounts.txt")
	if err != nil {
		panic(err)
	}

	critters := make([]critter, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(bs))
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		if len(splits) == 2 {
			critters = append(critters, critter{Email: splits[0], Password: splits[1]})
			continue
		}
	}

	return critters
}
