package delivery

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/evmartinelli/go-discordbot-panela/internal/discord"
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
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

	if contains(md.voiceUsecase.VoiceCommands(), m.Content) {
		go md.voiceUsecase.JoinAndPlayAudioFile(m.Content, s, m, guild, false)
	} else if strings.Contains(m.Content, "ajuda") {
		content := md.voiceUsecase.VoiceCommands()
		var b bytes.Buffer
		for _, v := range content.Items {
			b.WriteString(v.Title)
			b.WriteString("\n")
		}
		md.discord.SendMessageToChannel(m.ChannelID, b.String())
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
	} else if strings.Contains(m.Content, "socaforte") {
		content, err := md.messagesUsecase.GetPanelaADR()
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

func contains(s *repository.ResponseCMS, str string) bool {
	for _, v := range s.Items {
		if strings.HasSuffix(str, v.Title) {
			return true
		}
	}

	return false
}
