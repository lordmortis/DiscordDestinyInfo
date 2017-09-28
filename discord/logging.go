package discord

import (
  "fmt"
  "log"

  "github.com/bwmarrin/discordgo"
)

func LogPMError(session *discordgo.Session, user *discordgo.User, channel *discordgo.Channel, error string, v ...interface{}) {
  errorString := fmt.Sprintf(error, v...)
  log.Printf("ERROR: %s - %s", user.Username, errorString)
  msg := fmt.Sprintf("Error! %s", errorString)
  session.ChannelMessageSend(channel.ID, msg)
}

func LogPMCreateError(user *discordgo.User) {
  LogChatError(user, "Could not create channel for DM!")
}

func LogChatError(user *discordgo.User, error string, v ...interface{}) {
  errorString := fmt.Sprintf(error, v...)
  log.Printf("ERROR: %s - %s", user.Username, errorString)
}

func LogChatMsg(user *discordgo.User, msg string, v ...interface{}) {
  msgString := fmt.Sprintf(msg, v...)
  log.Printf("%s - %s\n", user.Username, msgString)
}

func LogChatCommand(user *discordgo.User, msg string, v ...interface{}) {
  msgString := fmt.Sprintf(msg, v...)
  log.Printf("%s - sent command: %s\n", user.Username, msgString)
}
