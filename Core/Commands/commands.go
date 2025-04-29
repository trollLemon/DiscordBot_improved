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
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	audioPlayer    audio.AudioPlayer
	searchDatabase database.Repository

	imageApi *imagemanip.ImageAPIWrapper
)

func init() {

	if err := InitDependencies(); err != nil {
		log.Fatal("Error setting up dependency services: ", err)
	}
}

func InitDependencies() error {
	streamSvc, err := factories.CreateStreamService(factories.YTDLP)
	if err != nil {
		return err
	}

	voiceSvc, err := factories.CreateVoiceService(factories.DiscordVoice)
	if err != nil {
		return err
	}

	notificationSvc, err := factories.CreateNotificationService(factories.DiscordNotification)
	if err != nil {
		return err
	}

	databaseSvc, err := factories.CreateDatabaseService(factories.Redis)

	if err != nil {
		return err
	}

	imageSvc, err := factories.CreateImageAPIService(factories.Imagemanip)

	if err != nil {
		return err
	}

	audioPlayer = *audio.NewAudioPlayer(streamSvc, voiceSvc, notificationSvc)

	imageApi = imageSvc

	searchDatabase = *database.NewRepository(databaseSvc)

	return nil
}

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

	var responseEdit *discordgo.WebhookEdit

	if err != nil {

		errResponse := "An error occurred: " + err.Error()
		responseEdit = &discordgo.WebhookEdit{
			Content: &errResponse,
		}
	} else {
		responseEdit = &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:   "processed_image.png",
					Reader: bytes.NewReader(image),
				},
			},
		}
	}

	if _, err := s.InteractionResponseEdit(i.Interaction, responseEdit); err != nil {
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

			return // if we cant get a url don't continue
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

	err := audioPlayer.Shuffle()
	response := util.GetBasicReply(fmt.Sprintf("Shuffled."))

	if err != nil {
		response = util.GetBasicReply(err.Error())
	}
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

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	kernelOption := i.ApplicationCommandData().Options[1].IntValue()
	lowerOption := i.ApplicationCommandData().Options[2].IntValue()
	upperOption := i.ApplicationCommandData().Options[3].IntValue()
	normalizeOption := i.ApplicationCommandData().Options[4].BoolValue()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.RandomFilter(imgBytes, format, kernelOption, lowerOption, upperOption, normalizeOption)
	processImageReply(image, err, s, i)

}

func InvertImage(s *discordgo.Session, i *discordgo.InteractionCreate) {

	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.InvertImage(imgBytes, format)
	processImageReply(image, err, s, i)
}
func SaturateImage(s *discordgo.Session, i *discordgo.InteractionCreate) {

	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	saturationVal := i.ApplicationCommandData().Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.SaturateImage(imgBytes, format, saturationVal)
	processImageReply(image, err, s, i)
}

func EdgeDetection(s *discordgo.Session, i *discordgo.InteractionCreate) {

	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	lowerOption := i.ApplicationCommandData().Options[1].IntValue()
	higherOption := i.ApplicationCommandData().Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.EdgeDetect(imgBytes, format, lowerOption, higherOption)
	processImageReply(image, err, s, i)
}

func Dilate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	kernelOption := i.ApplicationCommandData().Options[1].IntValue()
	iterationsOption := i.ApplicationCommandData().Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.DilateImage(imgBytes, format, kernelOption, iterationsOption)
	processImageReply(image, err, s, i)
}

func Erode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	kernelOption := i.ApplicationCommandData().Options[1].IntValue()
	iterationsOption := i.ApplicationCommandData().Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.ErodeImage(imgBytes, format, kernelOption, iterationsOption)
	processImageReply(image, err, s, i)
}

func AddText(s *discordgo.Session, i *discordgo.InteractionCreate) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	textOption := i.ApplicationCommandData().Options[1].Value.(string)
	fontScaleOption := i.ApplicationCommandData().Options[2].IntValue()
	xOption := i.ApplicationCommandData().Options[3].IntValue()
	yOption := i.ApplicationCommandData().Options[4].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	xPerc := float32(xOption) / 100.0
	yPerc := float32(yOption) / 100.0

	image, err := imageApi.AddText(imgBytes, format, textOption, float32(fontScaleOption), xPerc, yPerc)
	processImageReply(image, err, s, i)
}

func ReduceImage(s *discordgo.Session, i *discordgo.InteractionCreate) {

	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	reduceOption := i.ApplicationCommandData().Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	reducePerc := float32(reduceOption) / 100.0

	image, err := imageApi.Reduced(imgBytes, format, reducePerc)
	processImageReply(image, err, s, i)

}

func ShuffleImage(s *discordgo.Session, i *discordgo.InteractionCreate) {

	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachmentUrl := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	partitionsOption := i.ApplicationCommandData().Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentUrl)

	if err != nil {
		response := util.GetBasicReply(err.Error())
		if err := s.InteractionRespond(i.Interaction, response); err != nil {
			log.Printf("error responding to interaction: %v", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error during interaction defer: \n %s", err.Error())
		return
	}

	image, err := imageApi.Shuffle(imgBytes, format, partitionsOption)
	processImageReply(image, err, s, i)
}
