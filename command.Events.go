package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/lordmortis/DiscordDestinyInfo/discord"
	"github.com/lordmortis/goBungieNet"

	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.RegisterCommand("Events", "Tell me info about today/this week's events", handleEvents)
	discord.RegisterCommand("WhatsOn", "Tell me info about today/this week's events", handleEvents)
	discord.RegisterCommand("Xur", "Tell me info about Xur", handleXur)
	discord.RegisterCommand("Nightfall", "Tell me info about this week's nightfall", handleNightfall)
	discord.RegisterCommand("Flashpoint", "Tell me info about this week's flashpoint", handleFlashpoint)
	discord.RegisterCommand("Meditations", "Tell me info about this week's meditations", handleMeditations)
	discord.RegisterCommand("FactionRally", "Tell me info about the faction rally", handleFactionRally)
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

	msg := ""

	_ = milestones

	for key, defn := range *defns {
		err = nil
		var newMsg *string
		milestone := (*milestones)[key]
		switch defn.FriendlyName {
		case "FactionRallyPledge", "FactionRallyWinAnnouncement":
			// These just specify if the pledge is available and/or if the winner has been announced.
			fallthrough
		case "ClanProgress", "ClanObjectives":
			// We ignore clan progress because it doesn't tell us much.
			fallthrough
		case "Raid":
			// We ignore raid because it doesn't tell us much.
			fallthrough
		case "CallToArms", "Trials":
			// We're ignoring these because they are there all the time and provide no info
			//			fallthrough
			//		case "FactionRally":
			// I can't identify what the current faction rally is yet :/
		case "Meditations":
			newMsg, err = meditationsMessage(milestone, true)
		case "Nightfall":
			newMsg, err = nightfallMessage(milestone)
		case "FactionRally":
			fallthrough
		case "Hotspot":
			newMsg, err = hotspotMessage(milestone, defn)
		case "":
			// which one of these is it?
			switch defn.DisplayProperties.Name {
			case "Xûr":
				newMsg = xurMessage(milestone.EndDate)
			}
		default:
			err = dumpMilestoneData(milestone, defn)
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

func handleXur(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "Xur")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	var milestone goBungieNet.DestinyMilestone
	found := false
	for key, defn := range *defns {
		if defn.FriendlyName == "" && defn.DisplayProperties.Name == "Xûr" {
			milestone = ((*milestones)[key])
			found = true
			break
		}
	}

	if !found {
		session.ChannelMessageSend(message.ChannelID, "Xur's not around?")
	}

	eventMsg := xurMessage(milestone.EndDate)
	session.ChannelMessageSend(message.ChannelID, *eventMsg)
}

func handleMeditations(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "Meditations")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	var milestone goBungieNet.DestinyMilestone
	found := false
	for key, defn := range *defns {
		if defn.FriendlyName == "Meditations" {
			milestone = ((*milestones)[key])
			found = true
			break
		}
	}

	if !found {
		session.ChannelMessageSend(message.ChannelID, "No Meditations found")
	}

	var eventMsg *string
	eventMsg, err = meditationsMessage(milestone, false)
	if err != nil {
		discord.LogPMError(session, message.Author, channel, err.Error())
		return
	}

	msg := "This week's meditations are:\n"
	msg += *eventMsg

	session.ChannelMessageSend(message.ChannelID, msg)
}

func handleFactionRally(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "FactionRally")

	var milestones *map[uint32]goBungieNet.DestinyMilestone
	var defns *map[uint32]goBungieNet.DestinyMilestoneDefinition
	milestones, defns, err = fetchEventsAndDefinitions("en")
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get milestone data: %s", err.Error())
		return
	}

	var dailyMilestone goBungieNet.DestinyMilestone
	var dailyDefn goBungieNet.DestinyMilestoneDefinition
	dailyFound := false
	for key, defn := range *defns {
		if defn.FriendlyName == "FactionRally" {
			dailyMilestone = (*milestones)[key]
			dailyDefn = defn
			dailyFound = true
		}
	}

	if !dailyFound {
		session.ChannelMessageSend(message.ChannelID, "Faction Rally Events not running")
		return
	}

	if dailyFound {
		dailyMsg, err := hotspotMessage(dailyMilestone, dailyDefn)
		if err == nil {
			msg := "Faction rally is currently on, and the daily is: " + *dailyMsg
			session.ChannelMessageSend(message.ChannelID, msg)
		} else {
			discord.LogPMError(session, message.Author, channel, "Couldn't set daily message: %s", err.Error())
		}
	}
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

	msg = strings.TrimRight(msg, "\n")

	return &msg, nil
}

func meditationsMessage(milestone goBungieNet.DestinyMilestone, preamble bool) (*string, error) {
	activityHashes := make([]uint32, len(milestone.AvailableQuests))
	for index, quest := range milestone.AvailableQuests {
		activityHashes[index] = quest.Activity.ActivityHash
	}

	definitions, err := goBungieNet.ActivityDefinitions("en", activityHashes)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get activity definitions: %s", err.Error())
	}

	if len(*definitions) == 0 {
		msg := "No meditations found"
		return &msg, nil
	}

	msg := ""

	if preamble {
		msg += "Meditations:\n"
	}

	for _, definition := range *definitions {
		msg += fmt.Sprintf(
			"%s - %s\n",
			definition.DisplayProperties.Name,
			definition.DisplayProperties.Description,
		)
	}

	msg = strings.TrimRight(msg, "\n")

	return &msg, nil
}

func xurMessage(xurVanish time.Time) *string {
	amsg := "Xûr will be around for another "
	now := time.Now().UTC()
	amsg += friendlyDuration(xurVanish.Sub(now))
	return &amsg
}

func hotspotMessage(milestone goBungieNet.DestinyMilestone, definition goBungieNet.DestinyMilestoneDefinition) (*string, error) {
	questDefn := definition.Quests[milestone.AvailableQuests[0].QuestItemHash]
	msg := fmt.Sprintf("%s - %s", questDefn.DisplayProperties.Name, questDefn.DisplayProperties.Description)
	return &msg, nil
}

func dumpMilestoneData(milestone goBungieNet.DestinyMilestone, definition goBungieNet.DestinyMilestoneDefinition) error {
	fmt.Printf("Dump of %s\n", definition.FriendlyName)
	fmt.Print("Milestone:\n")
	spew.Dump(milestone)
	fmt.Print("\nDefinition:\n")
	spew.Dump(definition)
	fmt.Printf("\n")
	return nil
}
