package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

// Repository interface
//
//go:generate mockgen -source=messages_repository.go -destination=../mock/mock_repo.go
type Repository interface {
	GetBadWordList() (ReplyWordStruct, error)
	GetPlayersURL() (Players, error)
	GetPlayersStats(playerID string, data *Response) error
}

type messageRepository struct{}

// NewMessageRepository new message repository
func NewMessageRepository() Repository {
	return &messageRepository{}
}

// ReplyWordStruct structure
type ReplyWordStruct struct {
	BadwordReply []string `json:"badwordReply"`
	KuyReply     []string `json:"kuyReply"`
}

type PlayersURL struct {
	Nick string `json:"nick"`
	Url  string `json:"url"`
}

type Players struct {
	Players []PlayersURL `json:"players"`
}

type Response struct {
	Matches struct {
		Wins    int `json:"wins"`
		Loss    int `json:"loss"`
		Matches int `json:"matches"`
	} `json:"matches"`
}

// GetBadWordList return list of bad word
func (messageRepository) GetBadWordList() (ReplyWordStruct, error) {
	// need to injection config
	messagesFile, err := os.Open("./data/messages.json")
	if err != nil {
		fmt.Println("Error at HandleService: opening messages.json,\nMsg: ", err)
		return ReplyWordStruct{}, err
	}
	defer messagesFile.Close()
	replyWordByteValue, _ := io.ReadAll(messagesFile)
	var replyWord ReplyWordStruct
	json.Unmarshal(replyWordByteValue, &replyWord)
	return replyWord, nil
}

// GetPlayersURL return list of players
func (messageRepository) GetPlayersURL() (Players, error) {
	// need to injection config
	playersFile, err := os.Open("./data/panela.json")
	if err != nil {
		log.Println("Error at HandleService: opening messages.json,\nMsg: ", err)
		return Players{}, err
	}

	defer playersFile.Close()
	playersByteValue, _ := io.ReadAll(playersFile)
	if err != nil {
		log.Println("Error at ReadAll Players: \nMsg: ", err)
		return Players{}, err
	}
	var players Players
	json.Unmarshal(playersByteValue, &players)
	return players, nil
}

func (messageRepository) GetPlayersStats(playerID string, data *Response) error {
	client := resty.New()

	gclubsess := "gclubsess=" + os.Getenv("GCLUB_SESS")

	resp, err := client.R().
		SetResult(data).
		SetPathParams(map[string]string{"playerID": playerID}).
		SetHeader("Cookie", gclubsess).
		Get("https://gamersclub.com.br/api/box/historyFilterDate/{playerID}/2022-11")
	if err != nil {
		return fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)
	}

	log.Println(resp.Result())

	return nil
}
