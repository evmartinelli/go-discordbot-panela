package delivery

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/evmartinelli/go-discordbot-panela/internal/discord"
	messagesUsecase "github.com/evmartinelli/go-discordbot-panela/internal/messages/usecase"
	voiceUsecase "github.com/evmartinelli/go-discordbot-panela/internal/voice/usecase"
)

// Delivery interface
type Delivery interface {
	GetMessageHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type messageDelivery struct {
	voiceUsecase    voiceUsecase.Usecase
	discord         discord.Discord
	messagesUsecase messagesUsecase.Usecase
}

// NewMessageDelivery new message delivery
func NewMessageDelivery(discord discord.Discord, vu voiceUsecase.Usecase, mu messagesUsecase.Usecase) Delivery {
	return &messageDelivery{
		voiceUsecase:    vu,
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

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Println(err)
	}

	if strings.Contains(m.Content, "malvadao") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/kakule_malvadao.mp3", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "é o kakule malvadão ooo")
	} else if strings.Contains(m.Content, "terra") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/terra_alben.mp3", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "era só ouvir o doutor...")
	} else if strings.Contains(m.Content, "reeday") {
		content, err := md.messagesUsecase.GetPanelaMatches()
		if err != nil {
			log.Println(err)
		}
		md.discord.SendMessageToChannel(m.ChannelID, content)
	} else if strings.Contains(m.Content, "dumb") {
		content, err := md.messagesUsecase.GetPanelaLoss()
		if err != nil {
			log.Println(err)
		}
		md.discord.SendMessageToChannel(m.ChannelID, content)
	} else if strings.Contains(m.Content, "kakule") {
		content, err := md.messagesUsecase.GetPanelaKAST()
		if err != nil {
			log.Println(err)
		}
		md.discord.SendMessageToChannel(m.ChannelID, content)
	}
}
