package usecase

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/evmartinelli/go-discordbot-panela/internal/discord"
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
)

var stopChannel chan bool

// Usecase interface
type Usecase interface {
	PlayAudioFile(string, *discordgo.VoiceConnection)
	JoinAndPlayAudioFile(string, *discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild, bool)
	ConnectToVoiceChannel(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild, bool) (*discordgo.VoiceConnection, error)
	StopVoice()
	VoiceCommands() *repository.ResponseCMS
}

type voiceUsecase struct {
	discord            discord.Discord
	messagesRepository repository.Repository
}

// NewVoiceUsecase new voice usecase
func NewVoiceUsecase(mr repository.Repository) Usecase {
	return &voiceUsecase{
		messagesRepository: mr,
	}
}

// StopVoice stop voice channel
func (voiceUsecase) StopVoice() {
	if discord.GetVoiceStatus() {
		stopChannel <- true
	}
}

// PlayAudioFile return youtube download url
func (voiceUsecase) PlayAudioFile(file string, voiceConnection *discordgo.VoiceConnection) {
	if !discord.GetVoiceStatus() {
		stopChannel = make(chan bool)
		discord.UpdateVoiceStatus(true)
		dgvoice.PlayAudioFile(voiceConnection, file, stopChannel)
		close(stopChannel)
		discord.UpdateVoiceStatus(false)
	}
}

// JoinAndPlayAudioFile return youtube download url
func (vu voiceUsecase) JoinAndPlayAudioFile(content string, s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild, isMusicPlaying bool) {
	data, err := vu.messagesRepository.GetAudioItems()
	if err != nil {
		log.Printf("Error: connect to voice channel, Message: '%s'", err)
	}

	for _, v := range data.Items {
		if strings.Contains(content, v.Title) {
			fileToPlay := v.Title + filepath.Ext(v.Attachments[0].Url)
			if _, err := os.Stat(fileToPlay); err == nil {
				fmt.Printf("File exists\n")
			} else {
				err := DownloadFile(fileToPlay, v.Attachments[0].Url)
				if err != nil {
					panic(err)
				}
				fmt.Println("Downloaded: " + v.Attachments[0].Url)
			}
			// if _, err := filenameUsed(v.Title); err == nil {
			// 	fmt.Printf("File exists\n")
			// } else {
			// 	err := DownloadFile(v.Title+ext, v.Attachments[0].Url)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	fmt.Println("Downloaded: " + v.Attachments[0].Url)
			// }
			voiceConnection, err := connectToVoiceChannel(vu.discord, s, m, guild, isMusicPlaying)
			if err != nil {
				log.Printf("Error: connect to voice channel, Message: '%s'", err)
			}
			if !discord.GetVoiceStatus() {
				stopChannel = make(chan bool)
				discord.UpdateVoiceStatus(true)
				dgvoice.PlayAudioFile(voiceConnection, fileToPlay, stopChannel)
				close(stopChannel)
				discord.UpdateVoiceStatus(false)
			}
		}
	}
}

// ConnectToVoiceChannel connect to user voice channelId
func (vu voiceUsecase) ConnectToVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild, isMusicPlaying bool) (*discordgo.VoiceConnection, error) {
	return connectToVoiceChannel(vu.discord, s, m, guild, isMusicPlaying)
}

func (vu voiceUsecase) VoiceCommands() *repository.ResponseCMS {
	data, err := vu.messagesRepository.GetAudioItems()
	if err != nil {
		fmt.Print("error getting itens")
	}

	return data
}

func findVoiceChannelID(guild *discordgo.Guild, m *discordgo.MessageCreate) string {
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == m.Author.ID {
			return voiceState.ChannelID
		}
	}
	return ""
}

func connectToVoiceChannel(discord discord.Discord, s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild, isMustJoin bool) (voiceConnection *discordgo.VoiceConnection, err error) {
	voiceChannelID := findVoiceChannelID(guild, m)
	if voiceChannelID == "" && isMustJoin {
		if err := discord.SendMessageToChannel(m.ChannelID, "กรุณาเข้าห้องก่อนนะค้าบ"); err != nil {
			return nil, err
		}
	}
	voiceConnection, err = s.ChannelVoiceJoin(guild.ID, voiceChannelID, false, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
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

func filenameUsed(name string) (bool, error) {
	matches, err := filepath.Glob(name + ".*")
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}
