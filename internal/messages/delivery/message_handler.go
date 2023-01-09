package delivery

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	} else if strings.Contains(m.Content, "corno") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/corno_detectado.mpeg", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "e é VOCÊ...")
	} else if strings.Contains(m.Content, "tr") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/encaixou_tr.mp3", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "encaixooou esse tr ai hein...")
	} else if strings.Contains(m.Content, "ct") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/encaixou_ct.mp3", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "encaixooou esse ct ai hein...")
	} else if strings.Contains(m.Content, "carita") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/bota_la_carita.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "ai botou...")
	} else if strings.Contains(m.Content, "feliz") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/proibido_ser_feliz.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "proibido ser feliz...")
	} else if strings.Contains(m.Content, "pegando") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/pegando.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "eu to pegando uns cara...")
	} else if strings.Contains(m.Content, "horroroso") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/horroroso.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "fale por vc...")
	} else if strings.Contains(m.Content, "bosta") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/seubosta.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "SEU B#$#$...")
	} else if strings.Contains(m.Content, "panela") {
		go md.voiceUsecase.JoiAndPlayAudioFile("./sound/panela.mp4", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "é sim...")
	} else if strings.Contains(m.Content, "testcms") {
		fileUrl := "https://pub-31421060051a4b90b63207767964aab4.r2.dev/panelabot-cms/production/media/audio-80e48d57374b8de3f0d99ee438e992a4.mp3"
		if _, err := os.Stat("teste.mp3"); err == nil {
			fmt.Printf("File exists\n")
		} else {
			err := DownloadFile("teste.mp3", fileUrl)
			if err != nil {
				panic(err)
			}
			fmt.Println("Downloaded: " + fileUrl)
		}
		go md.voiceUsecase.JoiAndPlayAudioFile("teste.mp3", s, m, guild, false)
		md.discord.SendMessageToChannel(m.ChannelID, "é sim...")
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

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
