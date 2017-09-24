# Destiny Discord Bot!

This is a destiny discord bot. Primarily made for our personal destiny discord server.

# Developing

## Prerequisites

You will need a [properly setup](https://golang.org/doc/install) Go development environment to use this.

## Getting code and all dependencies

Once you have your GOPATH setup as above you should:

  1. `go get github.com/bwmarrin/discordgo`
  2. `go get github.com/lordmortis/goBungieNet`

once you have all those, you should create a config file from the `config.sample.json` template.
Ensure you have your Client ID handy to add your bot to a server

 1. `cd $GOPATH/src/github.com/lordmortis/DiscordDestinyInfo`
 2. `go run *.go`
 3. To add the bot to a server visit `https://discordapp.com/oauth2/authorize?scope=bot&permissions=o&client_id=<YOUR_CLIENT_ID>` (the correct URL will also be displayed when you run the bot)

