package discord

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
)

func helpCommand(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil {
    fmt.Println("Could not create channel for DM!")
    return
  }

  session.ChannelMessageSend(channel.ID, "Hi!")
  session.ChannelMessageSend(channel.ID, "Here's some help:")
  for _, key := range ( commandList ) {
    var command = commands[key]
    var help = fmt.Sprintf("%s - %s", command.name, command.help)
    session.ChannelMessageSend(channel.ID, help)

    fmt.Println(parameters)
  }
}