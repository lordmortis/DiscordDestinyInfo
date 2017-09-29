package main

import (
  "fmt"
//  "strings"
//  "strconv"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
//  "github.com/lordmortis/goBungieNet"

  "github.com/bwmarrin/discordgo"
)

func init() {
	discord.RegisterCommand("WhosOn", "Tell me who is on and what they are doing", handleWhosOn)
}

func handleWhosOn(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil { discord.LogPMCreateError(message.Author); return; }

	discord.LogChatCommand(message.Author, "WhosOn")

	var regos *[]Registration
	regos, err = loadRegos()
	if err != nil {	
		discord.LogPMError(session, message.Author, channel, "Couldn't get registrations: %s", err.Error())
		return
	}
	
	for _, rego := range( *regos ) {
		
	}

	session.ChannelMessageSend(message.ChannelID, msg)
}
