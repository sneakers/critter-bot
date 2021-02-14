package main

import (
	"bufio"
	"fmt"
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
