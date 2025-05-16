package Commands

import (
	"bot/Application"
	"bot/Core/Common"
	"bot/Core/Interfaces"
	"bot/util"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func JoinVC(s Interfaces.DiscordSession, voiceState *discordgo.VoiceState, guildId string) (*discordgo.VoiceConnection, error) {

	dgv, err := s.ChannelVoiceJoin(guildId, voiceState.ChannelID, false, false)

	if err != nil {
		return nil, errors.New("error joining voice channel")
	}

	return dgv, nil
}

func Play(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {

	author := i.GetInteractionAuthor()
	guild := a.GuildID
	channel := i.GetChannel()
	voiceState, err := s.VoiceState(guild, author)
	searchQuery := i.ApplicationCommandData().Options[0].StringValue()

	if err != nil {
		Common.Reply(s, i, "Failed to get voice state")
		return err
	}

	if voiceState.ChannelID == "" {
		Common.Reply(s, i, "You must be in a voice channel to use this command")
		return errors.New("user not in a voice channel")
	}

	if !util.IsURL(searchQuery) {
		videoUrl, err := a.Search.SearchWithQuery(searchQuery)
		if err != nil {
			Common.Reply(s, i, fmt.Sprintf("Error searching %s. %s", searchQuery, err.Error()))
			return errors.New("search failed with error : " + err.Error())
		}

		searchQuery = videoUrl

	}

	dgv, err := JoinVC(s, voiceState, guild)
	if err != nil {
		Common.Reply(s, i, "Error joining voice channel")
		return err
	}

	voiceConn := a.ServiceFactory.CreateVoiceService(dgv)
	notifier := a.ServiceFactory.CreateNotificationService(s, channel)

	a.AudioPlayer.UpdateConnection(voiceConn)
	a.AudioPlayer.UpdateNotifier(notifier)
	a.AudioPlayer.Play(searchQuery)

	Common.Reply(s, i, fmt.Sprintf("Added %s to the queue", searchQuery))

	return nil
}

func Stop(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	author := i.GetInteractionAuthor()
	guild := a.GuildID
	voiceState, err := s.VoiceState(guild, author)

	if err != nil {
		Common.Reply(s, i, "Failed to get voice state")
		return err
	}

	if voiceState.ChannelID == "" {
		Common.Reply(s, i, "You must be in a voice channel to use this command")
		return errors.New("user not in a voice channel")
	}

	a.AudioPlayer.Stop()

	Common.Reply(s, i, "Stopped")

	return nil
}

func Skip(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	author := i.GetInteractionAuthor()
	guild := a.GuildID
	voiceState, err := s.VoiceState(guild, author)

	if err != nil {
		Common.Reply(s, i, "Failed to get voice state")
		return err
	}

	if voiceState.ChannelID == "" {
		Common.Reply(s, i, "You must be in a voice channel to use this command")
		return errors.New("user not in a voice channel")
	}

	err = a.AudioPlayer.Skip()

	if err != nil {
		Common.Reply(s, i, "Error skipping : "+err.Error())
		return err
	}

	Common.Reply(s, i, "Skipped")
	return nil

}

func Shuffle(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	author := i.GetInteractionAuthor()
	guild := a.GuildID
	voiceState, err := s.VoiceState(guild, author)

	if err != nil {
		Common.Reply(s, i, "Failed to get voice state")
		return err
	}

	if voiceState.ChannelID == "" {
		Common.Reply(s, i, "You must be in a voice channel to use this command")
		return errors.New("user not in a voice channel")
	}

	err = a.AudioPlayer.Shuffle()

	if err != nil {
		Common.Reply(s, i, "Error shuffling : "+err.Error())
		return err
	}
	Common.Reply(s, i, "Shuffled")
	return nil

}

func RandomPlay(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {

	author := i.GetInteractionAuthor()
	guild := a.GuildID
	channel := i.GetChannel()
	voiceState, err := s.VoiceState(guild, author)
	numTerms := i.ApplicationCommandData().Options[0].IntValue()

	if err != nil {
		Common.Reply(s, i, "Failed to get voice state")
		return err
	}

	if voiceState.ChannelID == "" {
		Common.Reply(s, i, "You must be in a voice channel to use this command")
		return errors.New("user not in a voice channel")
	}

	randomTerms, err := a.WordDatabase.FetchRandom(int(numTerms))
	if err != nil {
		Common.Reply(s, i, "Error fetching random words")
		return err
	}

	searchQuery := strings.Join(randomTerms, " ")

	videoUrl, err := a.Search.SearchWithQuery(searchQuery)
	if err != nil {
		Common.Reply(s, i, fmt.Sprintf("Error searching %s. %s", searchQuery, err.Error()))
		return errors.New("search failed with error : " + err.Error())
	}

	dgv, err := JoinVC(s, voiceState, guild)
	if err != nil {
		Common.Reply(s, i, "Error joining voice channel")
		return err
	}

	voiceConn := a.ServiceFactory.CreateVoiceService(dgv)
	notifier := a.ServiceFactory.CreateNotificationService(s, channel)

	a.AudioPlayer.UpdateConnection(voiceConn)
	a.AudioPlayer.UpdateNotifier(notifier)
	a.AudioPlayer.Play(videoUrl)

	Common.Reply(s, i, fmt.Sprintf("Added %s to the queue. Searched with %s", videoUrl, searchQuery))

	return nil
}
