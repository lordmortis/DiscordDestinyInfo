package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/lordmortis/DiscordDestinyInfo/discord"
	"github.com/lordmortis/goBungieNet"

	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.RegisterCommand("Events", "Tell me info about today/this week's events", handleEvents)
	discord.RegisterCommand("Nightfall", "Tell me info about this week's nightfall", handleNightfall)
	discord.RegisterCommand("Flashpoint", "Tell me info about this week's flashpoint", handleFlashpoint)
}

func fetchEventsAndDefinitions(languageCode string) (*map[uint32]goBungieNet.DestinyMilestone, *map[uint32]goBungieNet.DestinyMilestoneDefinition, error) {
	milestones, err := goBungieNet.GetMilestones()
	if err != nil {
		return nil, nil, err
	}

	definitions := make(map[uint32]goBungieNet.DestinyMilestoneDefinition)
	for key, milestone := range *milestones {
		defn, err := milestone.Definition(languageCode)
		if err != nil {
			continue
		}
		definitions[key] = *defn
	}

	return milestones, &definitions, nil
}

func handleEvents(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "Events")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	msg := "Found Events:\n"

	_ = milestones

	for key, defn := range *defns {
		err = nil
		var newMsg *string
		switch defn.FriendlyName {
		case "Nightfall":
			newMsg, err = nightfallMessage((*milestones)[key])
		case "Hotspot":
			newMsg, err = hotspotMessage((*milestones)[key], defn)
		default:
			err = dumpMilestoneData((*milestones)[key], defn)
		}

		if err != nil {
			discord.LogPMError(session, message.Author, channel, err.Error())
			msg += fmt.Sprintf("%s - error fetching data", defn.FriendlyName)
		} else if newMsg != nil {
			msg += fmt.Sprintf("%s\n\n", *newMsg)
		}
	}

	msg = strings.TrimSuffix(msg, "\n")
	session.ChannelMessageSend(message.ChannelID, msg)
}

func handleNightfall(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "Nightfall")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	var nightfallMilestone goBungieNet.DestinyMilestone
	found := false
	for key, defn := range *defns {
		if defn.FriendlyName == "Nightfall" {
			nightfallMilestone = ((*milestones)[key])
			found = true
			break
		}
	}

	if !found {
		session.ChannelMessageSend(message.ChannelID, "No Nightfall found")
	}

	var eventMsg *string
	eventMsg, err = nightfallMessage(nightfallMilestone)
	if err != nil {
		discord.LogPMError(session, message.Author, channel, err.Error())
		return
	}

	msg := "This week's nightfall is:\n"
	msg += *eventMsg

	session.ChannelMessageSend(message.ChannelID, msg)
}

func handleFlashpoint(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "Flashpoint")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	var milestone goBungieNet.DestinyMilestone
	var milestoneDefn goBungieNet.DestinyMilestoneDefinition
	found := false
	for key, defn := range *defns {
		if defn.FriendlyName == "Hotspot" {
			milestone = ((*milestones)[key])
			milestoneDefn = defn
			found = true
			break
		}
	}

	if !found {
		session.ChannelMessageSend(message.ChannelID, "No Nightfall found")
	}

	var eventMsg *string
	eventMsg, err = hotspotMessage(milestone, milestoneDefn)
	if err != nil {
		discord.LogPMError(session, message.Author, channel, err.Error())
		return
	}

	msg := "This week's flashpoint is:\n"
	msg += *eventMsg

	session.ChannelMessageSend(message.ChannelID, msg)
}

func nightfallMessage(milestone goBungieNet.DestinyMilestone) (*string, error) {
	nfQuestActivity := milestone.AvailableQuests[0].Activity

	nfActDefn, err := nfQuestActivity.Definition("en")
	if err != nil {
		return nil, fmt.Errorf("Couldn't get activity definition: %s", err.Error())
	}

	var nfActModDefns *[]goBungieNet.DestinyActivityModifierDefinition
	nfActModDefns, err = nfQuestActivity.Modifiers("en")
	if err != nil {
		return nil, fmt.Errorf("Couldn't get modifier definitions: %s", err.Error())
	}

	msg := fmt.Sprintf("%s - %s\n",
		nfActDefn.DisplayProperties.Name,
		nfActDefn.DisplayProperties.Description)
	if len(*nfActModDefns) > 0 {
		msg += "Modifiers\n"
		for index, modifier := range *nfActModDefns {
			if modifier.Redacted {
				msg += fmt.Sprintf("%d. Redacted by Bungie :(\n", index+1)
			} else {
				msg += fmt.Sprintf("%d. %s %s\n", index+1, modifier.DisplayProperties.Name, modifier.DisplayProperties.Description)
			}
		}
	}

	return &msg, nil
}

func hotspotMessage(milestone goBungieNet.DestinyMilestone, definition goBungieNet.DestinyMilestoneDefinition) (*string, error) {
	questDefn := definition.Quests[milestone.AvailableQuests[0].QuestItemHash]
	msg := fmt.Sprintf("%s - %s", questDefn.DisplayProperties.Name, questDefn.DisplayProperties.Description)
	return &msg, nil
}

func dumpMilestoneData(milestone goBungieNet.DestinyMilestone, definition goBungieNet.DestinyMilestoneDefinition) error {
	fmt.Printf("Dump of %s", definition.FriendlyName)
	fmt.Print("Milestone:\n")
	spew.Dump(milestone)
	fmt.Print("\nDefinition:\n")
	spew.Dump(definition)
	fmt.Printf("\n")
	return nil
}
