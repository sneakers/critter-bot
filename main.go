package main

var playerID string

func main() {
	prevID := askForInput("Player ID to follow:")
	critters := getAccounts()

	for _, crit := range critters {
		crit := crit
		crit.login()

		if crit.LoginSuccess == false {
			continue
		}

		crit.FollowingID = prevID
		crit.joinGame()
		go crit.listener()

		prevID = crit.PlayerID
	}

	select {}
}
