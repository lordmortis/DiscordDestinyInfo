package discord

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
)

func discordGetChannelInfo(session *discordgo.Session, id string) channelInfo {
  info, exists := discordChannels[id]
  if (exists) { return info }

  discordInfo, err := session.Channel(id)
  if (err != nil) {
    fmt.Println("Could not retrieve info for: " + id)
    return info
  }

  var isDM = discordInfo.Type == discordgo.ChannelTypeDM
  var isGroupDM = discordInfo.Type == discordgo.ChannelTypeGroupDM

  var newInfo channelInfo
  newInfo.DM = isDM || isGroupDM
  discordChannels[id] = newInfo
  return newInfo
}
