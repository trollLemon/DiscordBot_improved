package Commands

import (
	"bot/Factories"
	"bot/Services/Audio"
	"bot/Services/Database"
	"bot/util"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	audioPlayer    = audio.NewAudioPlayer(factories.CreateStreamService(), factories.CreateVoiceService(), factories.CreateNotificationService())
	searchDatabase = database.NewRepository(factories.CreateDatabaseService())
)

func vcHelper(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guild := os.Getenv("GUILD_ID")
	author := i.Member.User.ID

	dgv, _ := util.JoinVoiceChannel(s, author, guild)
	voiceConn := audio.Voice{
		Vc: dgv,
	}

	notif := audio.Notifier{
		Session: s,
		Channel: i.ChannelID,
	}
	audioPlayer.SetConnection(&voiceConn)
	audioPlayer.UpdateNotifier(&notif)
}

func inVc(s *discordgo.Session, i *discordgo.InteractionCreate) bool {

	guild := os.Getenv("GUILD_ID")
	author := i.Member.User.ID

	voiceState, _ := s.State.VoiceState(guild, author)

	return voiceState != nil && voiceState.ChannelID != ""

}
func Play(s *discordgo.Session, i *discordgo.InteractionCreate) {

	query := i.ApplicationCommandData().Options[0].StringValue()
	url := query
	if !util.IsURL(query) {
		u, err := util.GetURLFromQuery(query)
		if err != nil {

			response := util.GetBasicReply(fmt.Sprintf("Error searching Youtube... '%s'", err.Error()))
			if err := s.InteractionRespond(i.Interaction, response); err != nil {
				log.Printf("error responding to interaction: %v", err)
			}

			return // if we cant get a url dont continue
		}
		url = u
	}

	//vcHelper(s, i)

	guild := os.Getenv("GUILD_ID")
	author := i.Member.User.ID

	dgv, _ := util.JoinVoiceChannel(s, author, guild)
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
	}
}

func Stop(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if !inVc(s, i) {
		response := util.GetBasicReply(fmt.Sprint("You must be in a VC to use this command."))
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}

		return
	}

	audioPlayer.Stop()
	response := util.GetBasicReply(fmt.Sprint("Stopped, leaving VC."))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func Skip(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if !inVc(s, i) {
		response := util.GetBasicReply(fmt.Sprint("You must be in a VC to use this command."))
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}

	}

	err := audioPlayer.Skip()
	content := "Skipping..."
	if err != nil {
		content = "error skipping " + err.Error()

	}
	response := util.GetBasicReply(content)
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func Shuffle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !inVc(s, i) {
		response := util.GetBasicReply(fmt.Sprint("You must be in a VC to use this command."))
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
			return
		}

	}

	audioPlayer.Shuffle()
	response := util.GetBasicReply(fmt.Sprintf("Shuffled."))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}

func RandomPlay(s *discordgo.Session, i *discordgo.InteractionCreate) {

	num := i.ApplicationCommandData().Options[0].IntValue()
	terms, err := searchDatabase.GetRandN(int(num))
	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
			return
		}
	}

	url, err := util.GetURLFromQuery(strings.Join(terms, " "))

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}

		return
	}

	vcHelper(s, i)
	audioPlayer.Play(url)
	response := util.GetBasicReply(fmt.Sprintf("Added %s to the queue. Searched with: %s", url, terms))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}

}
func Add(s *discordgo.Session, i *discordgo.InteractionCreate) {
	term := i.ApplicationCommandData().Options[0].StringValue()
	err := searchDatabase.Add(term)
	content := fmt.Sprintf("Added %s to database", term)
	if err != nil {
		content = err.Error()
	}

	response := util.GetBasicReply(fmt.Sprintf(content))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}

}
func Remove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	term := i.ApplicationCommandData().Options[0].StringValue()
	err := searchDatabase.Remove(term)
	content := fmt.Sprintf("Removed %s from database", term)
	if err != nil {
		content = err.Error()
	}

	response := util.GetBasicReply(fmt.Sprintf(content))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}
func Show(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data, err := searchDatabase.GetAllItems()
	content := fmt.Sprintf(strings.Join(data, "\n"))
	if err != nil {
		content = err.Error()
	}

	response := util.GetBasicReply(fmt.Sprintf(content))
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}

}
