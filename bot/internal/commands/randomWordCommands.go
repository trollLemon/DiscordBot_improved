package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/trollLemon/DiscordBot/internal/application"
	"github.com/trollLemon/DiscordBot/internal/common"
)

func Add(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.RandomWords.Insert(term); err != nil {
		common.RandomWordsError(s, i, "Add failed", term, err)
		return err
	}

	common.Reply(s, i, fmt.Sprintf("Added `%s` to the word list.", term))
	return nil
}

func Remove(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	term := i.ApplicationCommandData().Options[0].StringValue()

	if err := a.RandomWords.Delete(term); err != nil {
		common.RandomWordsError(s, i, "Removed failed", term, err)
		return err
	}

	common.Reply(s, i, fmt.Sprintf("Removed `%s` from the word list.", term))
	return nil
}

func Show(s *discordgo.Session, i *discordgo.InteractionCreate, a *application.Application) error {
	terms, err := a.RandomWords.GetAll()

	if err != nil {
	common.RandomWordsError(s, i, "Show failed", "", err)
		return err
	}
	wordsString := strings.Join(terms, "\n")

	common.Reply(s, i, wordsString)

	return nil
}
