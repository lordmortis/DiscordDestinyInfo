package main

import (
  "fmt"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
  "github.com/lordmortis/goBungieNet"

  "github.com/bwmarrin/discordgo"
)

func init() {
  discord.RegisterCommand("Nightfall", "Tell me info about this week's nightfall!", handleNightfall)
}

func handleNightfall(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil { discord.LogPMCreateError(message.Author); return; }

  discord.LogChatCommand(message.Author, "Nightfall")

  var milestones *map[uint32]goBungieNet.DestinyMilestone
  milestones, err = goBungieNet.GetMilestones()
  if err != nil {
      discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
      return
  }

  var nfMs *goBungieNet.DestinyMilestone
  for _, milestone := range(*milestones) {
    defn, err := milestone.Definition("en")
    if err != nil {
      discord.LogPMError(session, message.Author, channel, "Couldn't get milestone definition: %s", err.Error())
      continue
    }
    if defn.FriendlyName == "Nightfall" { nfMs = &milestone; break; }
  }

  if nfMs == nil {
    session.ChannelMessageSend(message.ChannelID, "No Nightfall found")
  }

  nfQuestActivity := nfMs.AvailableQuests[0].Activity

  var nfActDefn *goBungieNet.DestinyActivity
  nfActDefn, err = nfQuestActivity.Definition("en")
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Couldn't get activity definition: %s", err.Error())
    return
  }

  var nfActModDefns *[]goBungieNet.DestinyActivityModifierDefinition
  nfActModDefns, err = nfQuestActivity.Modifiers("en")
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Couldn't get modifier definitions: %s", err.Error())
    return
  }

  msg := "This week's nightfall is:\n"
  msg += fmt.Sprintf("%s - %s\n",
    nfActDefn.DisplayProperties.Name,
    nfActDefn.DisplayProperties.Description)
  for index, modifier := range (*nfActModDefns) {
    if (modifier.Redacted) {
      msg += fmt.Sprintf("Modifier %d - Redacted by Bungie :(\n", index +1)
    } else {
      msg += fmt.Sprintf("Modifier %s - %s\n", modifier.DisplayProperties.Name, modifier.DisplayProperties.Description)
    }
  }

  session.ChannelMessageSend(message.ChannelID, msg)
}
