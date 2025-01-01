package Commands

import (
	"bot/Core/Factories"
	"bot/Core/Services/Audio"
	"bot/Core/Services/Database"
	"bot/Core/Services/ImageManip"
	"bot/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	audioPlayer    = audio.NewAudioPlayer(factories.CreateStreamService(), factories.CreateVoiceService(), factories.CreateNotificationService())
	searchDatabase = database.NewRepository(factories.CreateDatabaseService())

	api = imagemanip.NewImageAPIWrapper("http://127.0.0.1:8000/api",20)
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

func RandomImageFilter(s *discordgo.Session, i *discordgo.InteractionCreate) {


	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	encoded_url := url.QueryEscape(attachmentUrl)

	kernel_option := i.ApplicationCommandData().Options[1].IntValue()
	lower_option := i.ApplicationCommandData().Options[2].IntValue()
	upper_option := i.ApplicationCommandData().Options[3].IntValue()
	//normalize_option := i.ApplicationCommandData().Options[4].BoolValue()

	image, err := api.RandomFilter(encoded_url, int(kernel_option), int(lower_option), int(upper_option))

	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}

}

func InvertImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.InvertImage(encoded_url)
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}
func SaturateImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	magnitude := i.ApplicationCommandData().Options[1].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.SaturateImage(encoded_url,int(magnitude))
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}

func EdgeDetection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	lower := i.ApplicationCommandData().Options[1].IntValue()
	upper := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.EdgeDetect(encoded_url,int(lower),int(upper))
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}


func  Dilate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	box_size := i.ApplicationCommandData().Options[1].IntValue()
	iterations := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.DilateImage(encoded_url,int(box_size),int(iterations))
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}



func  Erode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	box_size := i.ApplicationCommandData().Options[1].IntValue()
	iterations := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.ErodeImage(encoded_url,int(box_size),int(iterations))
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}



func  AddText(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	text := i.ApplicationCommandData().Options[1].Value.(string)
	font_size := i.ApplicationCommandData().Options[2].IntValue()
	x := i.ApplicationCommandData().Options[3].IntValue()
	y := i.ApplicationCommandData().Options[3].IntValue()
	
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.AddText(encoded_url,text,float32(font_size),float32(x)/100.0,float32(y)/100.0)
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}


func ReduceImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	
	quality := i.ApplicationCommandData().Options[1].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	
	image, err := api.Reduced(encoded_url,float32(quality)/ 100.0 )
	
	var response *discordgo.InteractionResponse
	if err != nil {
		response = util.GetBasicReply(err.Error())
	} else {
		response = util.GetAttatchmentReply(image)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
	
}