package commands

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/common"
	"github.com/trollLemon/DiscordBot/internal/gomanip"
	"github.com/trollLemon/DiscordBot/internal/util"
)


var (
  errFailedDownload = errors.New("Failed to download the provided image. Please try again.")
)

func RandomImageFilter(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	kernelOption := applicationData.Options[1].IntValue()
	lowerOption := applicationData.Options[2].IntValue()
	higherOption := applicationData.Options[3].IntValue()
	normalizeOption := applicationData.Options[4].BoolValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)
	
	if err != nil {
		common.GomanipError(s, i, "RandomImageFilter failed", errFailedDownload)
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.RandomFilter(a.Gomanip, imgBytes, format, kernelOption, lowerOption, higherOption, normalizeOption)

	if err != nil {
		common.GomanipError(s, i, "Random image filter failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func InvertImage(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {

	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.InvertImage(a.Gomanip, imgBytes, format)

	if err != nil {
		common.GomanipError(s, i, "Invert image failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func SaturateImage(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	saturationMagnitude := applicationData.Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.SaturateImage(a.Gomanip, imgBytes, format, saturationMagnitude)

	if err != nil {
		common.GomanipError(s, i, "Saturate image failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func EdgeDetection(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL

	lowerBound := applicationData.Options[1].IntValue()
	upperBound := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.EdgeDetect(a.Gomanip, imgBytes, format, lowerBound, upperBound)

	if err != nil {
		common.GomanipError(s, i, "Edge detection failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func Dilate(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	boxSize := applicationData.Options[1].IntValue()
	iterations := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.DilateImage(a.Gomanip, imgBytes, format, boxSize, iterations)

	if err != nil {
		common.GomanipError(s, i, "Dilating image failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func Erode(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	boxSize := applicationData.Options[1].IntValue()
	iterations := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.ErodeImage(a.Gomanip, imgBytes, format, boxSize, iterations)

	if err != nil {
		common.GomanipError(s, i, "Eroding image failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func AddText(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	text := applicationData.Options[1].Value.(string)
	fontScaleOption := applicationData.Options[2].IntValue()
	xOption := applicationData.Options[3].IntValue()
	yOption := applicationData.Options[4].IntValue()

	fontScale := float32(fontScaleOption)
	x := float32(xOption) / 100.0
	y := float32(yOption) / 100.0

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.AddText(a.Gomanip, imgBytes, format, text, fontScale, x, y)

	if err != nil {
		common.GomanipError(s, i, "Adding text failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func RandomText(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	numTerms := applicationData.Options[1].IntValue()
	fontScaleOption := applicationData.Options[2].IntValue()
	xOption := applicationData.Options[3].IntValue()
	yOption := applicationData.Options[4].IntValue()

	fontScale := float32(fontScaleOption)
	x := float32(xOption) / 100.0
	y := float32(yOption) / 100.0

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}

	terms, err := a.RandomWords.GetRandom(int(numTerms))
	if err != nil {
		common.Reply(s, i, "Error fetching random words")
		return err
	}

	text := strings.Join(terms, " ")

	common.DeferReply(s, i)

	img, err := gomanip.AddText(a.Gomanip, imgBytes, format, text, fontScale, x, y)

	if err != nil {
		common.GomanipError(s, i, "Adding random text failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func ReduceImage(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	qualityOption := applicationData.Options[1].IntValue()

	quality := float32(qualityOption) / 100.0

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.Reduced(a.Gomanip, imgBytes, format, quality)

	if err != nil {
		common.GomanipError(s, i, "Image quality reduction failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}

func ShuffleImage(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.ApplicationCommandData().Resolved.Attachments[attachmentID].URL
	partitionsOption := applicationData.Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	common.DeferReply(s, i)

	img, err := gomanip.Shuffle(a.Gomanip, imgBytes, format, partitionsOption)

	if err != nil {
		common.GomanipError(s, i, "Shuffling image failed", err)
	} else {
		common.ReplyGomanip(img, s, i)
	}

	return err
}
