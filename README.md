# critter-bot
A box critters bot to follow your player. Based off of a Perl script I used in 2010 on Club Penguin.

## Usage

Put your accounts into `accounts/accounts.txt` in the format `email:password`

Run the bot: `go run .`


Type in your playerID (`world.player.playerId` in box critters javascript console to get this)

Move around and watch your critters follow you

### Notes

It appears as if box critters mutes your critter if there's more than 1 active sesion per IP. It also appears that you can only have 4 connections per IP.

