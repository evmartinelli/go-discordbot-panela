package delivery

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/evmartinelli/go-discordbot-panela/internal/discord"
	messagesUsecase "github.com/evmartinelli/go-discordbot-panela/internal/messages/usecase"
)

// Delivery interface
type Delivery interface {
	GetMessageHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type messageDelivery struct {
	discord         discord.Discord
	messagesUsecase messagesUsecase.Usecase
}

// NewMessageDelivery new message delivery
func NewMessageDelivery(discord discord.Discord, mu messagesUsecase.Usecase) Delivery {
	return &messageDelivery{
		discord:         discord,
		messagesUsecase: mu,
	}
}

func (md messageDelivery) GetMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	botPrefix := os.Getenv("BOT_PREFIX")

	if !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	if strings.Contains(m.Content, "first kill") {
		md.discord.SendMessageToChannel(m.ChannelID, md.messagesUsecase.GetRandomKuyReplyWord())
	} else if strings.Contains(m.Content, "kakule") {
		md.discord.SendMessageToChannel(m.ChannelID, "100% KAST")
	} else if strings.Contains(m.Content, "reeday") {
		content, err := md.messagesUsecase.GetPanelaMatches()
		if err != nil {
			log.Println(err)
		}
		md.discord.SendMessageToChannel(m.ChannelID, content)
	}
}
