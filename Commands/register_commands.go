package Commands

import (
	"bot/Commands/Audio"
	"bot/Commands/Voice"
	"bot/util"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	audioPlayer = audio.NewAudioPlayer(CreateStreamService(),CreateVoiceService(), CreateNotificationService())

	SlashCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Enqueue a youtube video and play it",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The search query",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop",
			Description: "Enqueue a youtube video and play it",
		},
		{
			Name:        "skip",
			Description: "Skip what is playing, and start the next audio in the queue",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			query := i.ApplicationCommandData().Options[0].StringValue()
			url := query
			if !util.IsURL(query) {
				u, err := util.GetURLFromQuery(query)
				if err != nil {

					response := util.GetBasicReply(fmt.Sprintf("Error searching Youtube... '%s'", err.Error()))
					if err := s.InteractionRespond(i.Interaction, response); err != nil {
						log.Printf("error responding to interaction: %v", err)
						return
					}

				}
				url = u
			}

			guild := os.Getenv("GUILD_ID")
			author := i.Member.User.ID

			dgv, _ := voice.JoinVoiceChannel(s, author, guild)
			voiceConn := audio.Voice{
				Vc: dgv,
			}

			notif := audio.Notifier{
				Session: s,
				Channel: i.ChannelID,
			}
			audioPlayer.SetConnection(&voiceConn)
			audioPlayer.UpdateNotifier(&notif)
			audioPlayer.Play(url)
			response := util.GetBasicReply(fmt.Sprintf("Added %s to the queue.", url))
			if err := s.InteractionRespond(i.Interaction, response); err != nil {
				log.Printf("error responding to interaction: %v", err)
				return
			}

		},
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			audioPlayer.Stop()
			response := util.GetBasicReply(fmt.Sprint("Stopped, leaving VC."))
			if err := s.InteractionRespond(i.Interaction, response); err != nil {
				log.Printf("error responding to interaction: %v", err)
				return
			}
		},
		"shuffle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			audioPlayer.Shuffle()
			response := util.GetBasicReply(fmt.Sprintf("Shuffled."))
			if err := s.InteractionRespond(i.Interaction, response); err != nil {
				log.Printf("error responding to interaction: %v", err)
				return
			}
		},
		"skip": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := audioPlayer.Skip()
			content := "Skipping..."
			if err != nil {
				content = "error skipping " + err.Error()

			}
			response := util.GetBasicReply(content)
			if err := s.InteractionRespond(i.Interaction, response); err != nil {
				log.Printf("error responding to interaction: %v", err)
				return
			}

		},
	}
)
