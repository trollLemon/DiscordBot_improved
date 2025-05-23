package Commands

import (
	"bot/Application"
	"bot/Core/Common"
	"bot/Core/Interfaces"
	imagemanip "bot/Core/Services/ImageManip"
	"bot/util"
)

func RandomImageFilter(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	kernelOption := applicationData.Options[1].IntValue()
	lowerOption := applicationData.Options[2].IntValue()
	higherOption := applicationData.Options[3].IntValue()
	normalizeOption := applicationData.Options[4].BoolValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.RandomFilter(a.ImageApi, imgBytes, format, kernelOption, lowerOption, higherOption, normalizeOption)

	Common.ReplyImage(img, err, s, i)

	return err
}

func InvertImage(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {

	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.InvertImage(a.ImageApi, imgBytes, format)

	Common.ReplyImage(img, err, s, i)

	return err

}

func SaturateImage(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	saturationMagnitude := applicationData.Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.SaturateImage(a.ImageApi, imgBytes, format, saturationMagnitude)

	Common.ReplyImage(img, err, s, i)

	return err

}

func EdgeDetection(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	lowerBound := applicationData.Options[1].IntValue()
	upperBound := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.EdgeDetect(a.ImageApi, imgBytes, format, lowerBound, upperBound)

	Common.ReplyImage(img, err, s, i)

	return err

}

func Dilate(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	boxSize := applicationData.Options[1].IntValue()
	iterations := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.DilateImage(a.ImageApi, imgBytes, format, boxSize, iterations)

	Common.ReplyImage(img, err, s, i)

	return err

}

func Erode(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	boxSize := applicationData.Options[1].IntValue()
	iterations := applicationData.Options[2].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.ErodeImage(a.ImageApi, imgBytes, format, boxSize, iterations)

	Common.ReplyImage(img, err, s, i)

	return err

}

func AddText(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	text := applicationData.Options[1].Value.(string)
	fontScaleOption := applicationData.Options[2].IntValue()
	xOption := applicationData.Options[3].IntValue()
	yOption := applicationData.Options[4].IntValue()

	fontScale := float32(fontScaleOption) / 100.0
	x := float32(xOption) / 100.0
	y := float32(yOption) / 100.0

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.AddText(a.ImageApi, imgBytes, format, text, fontScale, x, y)

	Common.ReplyImage(img, err, s, i)

	return err

}

func ReduceImage(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	qualityOption := applicationData.Options[1].IntValue()

	quality := float32(qualityOption) / 100.0

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.Reduced(a.ImageApi, imgBytes, format, quality)

	Common.ReplyImage(img, err, s, i)

	return err

}

func ShuffleImage(s Interfaces.DiscordSession, i Interfaces.DiscordInteraction, a *application.Application) error {
	applicationData := i.ApplicationCommandData()
	attachmentID := applicationData.Options[0].Value.(string)
	attachmentURL := i.GetImageURLFromAttachmentID(attachmentID)
	partitionsOption := applicationData.Options[1].IntValue()

	imgBytes, format, err := util.GetImageFromURL(attachmentURL)

	if err != nil {
		Common.Reply(s, i, "Error downloading given attachment")
		return err
	}
	Common.DeferReply(s, i)

	img, err := imagemanip.Shuffle(a.ImageApi, imgBytes, format, partitionsOption)

	Common.ReplyImage(img, err, s, i)

	return err

}
