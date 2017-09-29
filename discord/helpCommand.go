package discord

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
)

func helpCommand(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil { LogPMCreateError(message.Author); return; }

  LogChatCommand(message.Author, "Help")

  msg := "Hi!\nHere's some help:\n"
  for _, key := range ( commandList ) {
    var command = commands[key]
    msg += fmt.Sprintf("%s - %s\n", command.name, command.help)
  }

  session.ChannelMessageSend(channel.ID, msg)
}
