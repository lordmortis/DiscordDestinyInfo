package discord

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
)

type channelInfo struct {
  DM bool
}

var (
  channels map[string]channelInfo
)

func init() {
  channels = make(map[string]channelInfo)
}

func getChannelInfo(session *discordgo.Session, id string) channelInfo {
  info, exists := channels[id]
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
  channels[id] = newInfo
  return newInfo
}
