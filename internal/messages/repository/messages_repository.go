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
	GetPlayersURL() (Players, error)
	GetPlayerStats(playerID string, data *Response) error
	GetPlayerStatsAsync(playerID string, rchan chan Response)
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
	Stat []struct {
		Stat  string `json:"stat"`
		Value string `json:"value"`
	} `json:"stat"`
}

const GCURL = "https://csgo.gamersclub.gg/api/box/history/{playerID}"

// GetPlayersURL return list of players
func (messageRepository) GetPlayersURL() (Players, error) {
	// need to injection config
	playersFile, err := os.Open("./data/panela.json")
	// playersFile, err := os.Open("/Users/evandrom/Projects/Personal/go-discordbot-panela/data/panela.json")
	if err != nil {
		log.Println("Error at HandleService: opening panela.json,\nMsg: ", err)
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

func (messageRepository) GetPlayerStats(playerID string, data *Response) error {
	client := resty.New()

	gclubsess := "gclubsess=" + os.Getenv("GCLUB_SESS")

	resp, err := client.R().
		SetResult(data).
		SetPathParams(map[string]string{"playerID": playerID}).
		SetHeader("Cookie", gclubsess).
		Get(GCURL)

		// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	if err != nil {
		return fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("TranslationWebAPI - Translate - trans.Translate: %w", err)
	}

	return nil
}

func (messageRepository) GetPlayerStatsAsync(playerID string, rchan chan Response) {
	defer close(rchan)
	client := resty.New()
	data := &Response{}

	gclubsess := "gclubsess=" + os.Getenv("GCLUB_SESS")

	resp, err := client.R().
		SetResult(data).
		SetPathParams(map[string]string{"playerID": playerID}).
		SetHeader("Cookie", gclubsess).
		Get(GCURL)

		// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	if err != nil {
		log.Print(err)
	}

	if resp.IsError() {
		log.Print(err)
	}

	rchan <- *data
}
