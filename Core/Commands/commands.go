package Commands

import (
	"bot/Core/Factories"
	"bot/Core/Services/Audio"
	"bot/Core/Services/Database"
	"bot/Core/Services/ImageManip"
	"bot/util"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	audioPlayer    = audio.NewAudioPlayer(factories.CreateStreamService(), factories.CreateVoiceService(), factories.CreateNotificationService())
	searchDatabase = database.NewRepository(factories.CreateDatabaseService())

	api = imagemanip.NewImageAPIWrapper("http://image:8080/api")
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

func processImageReply(image []byte, err error, s *discordgo.Session, i *discordgo.InteractionCreate) {

	var response_edit *discordgo.WebhookEdit

	if err != nil {

		err_response := "An error occurred: " + err.Error()
		response_edit = &discordgo.WebhookEdit{
			Content: &err_response,
		}
	} else {
		response_edit = &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:   "processed_image.png",
					Reader: bytes.NewReader(image),
				},
			},
		}
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, response_edit); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}

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

	vcHelper(s, i)

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

	response := util.GetBasicReply(content)
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

	response := util.GetBasicReply(content)
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		log.Printf("error responding to interaction: %v", err)
	}
}
func Show(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data, err := searchDatabase.GetAllItems()
	content := strings.Join(data, "\n")
	if err != nil {
		content = err.Error()
	}

	response := util.GetBasicReply(content)
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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := api.RandomFilter(encoded_url, int(kernel_option), int(lower_option), int(upper_option))

	processImageReply(image, err, s, i)

}

func InvertImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	encoded_url := url.QueryEscape(attachmentUrl)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := api.InvertImage(encoded_url)

	processImageReply(image, err, s, i)

}
func SaturateImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	magnitude := i.ApplicationCommandData().Options[1].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.SaturateImage(encoded_url, int(magnitude))

	processImageReply(image, err, s, i)

}

func EdgeDetection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	lower := i.ApplicationCommandData().Options[1].IntValue()
	upper := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.EdgeDetect(encoded_url, int(lower), int(upper))

	processImageReply(image, err, s, i)

}

func Dilate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	box_size := i.ApplicationCommandData().Options[1].IntValue()
	iterations := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.DilateImage(encoded_url, int(box_size), int(iterations))

	processImageReply(image, err, s, i)

}

func Erode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	box_size := i.ApplicationCommandData().Options[1].IntValue()
	iterations := i.ApplicationCommandData().Options[2].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.ErodeImage(encoded_url, int(box_size), int(iterations))

	processImageReply(image, err, s, i)

}

func AddText(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	text := i.ApplicationCommandData().Options[1].Value.(string)
	font_size := i.ApplicationCommandData().Options[2].IntValue()
	x := i.ApplicationCommandData().Options[3].IntValue()
	y := i.ApplicationCommandData().Options[3].IntValue()

	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.AddText(encoded_url, text, float32(font_size), float32(x)/100.0, float32(y)/100.0)

	processImageReply(image, err, s, i)

}

func ReduceImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	quality := i.ApplicationCommandData().Options[1].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.Reduced(encoded_url, float32(quality)/100.0)

	processImageReply(image, err, s, i)

}

func ShuffleImage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	partitions := i.ApplicationCommandData().Options[1].IntValue()
	encoded_url := url.QueryEscape(attachmentUrl)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}
	image, err := api.Shuffle(encoded_url, int(partitions))

	processImageReply(image, err, s, i)

}
