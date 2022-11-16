package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	replyWordByteValue, _ := ioutil.ReadAll(messagesFile)
	var replyWord ReplyWordStruct
	json.Unmarshal(replyWordByteValue, &replyWord)
	return replyWord, nil
}

// GetPlayersURL return list of players
func (messageRepository) GetPlayersURL() (Players, error) {
	// need to injection config
	playersFile, err := os.Open("./data/panela.json")
	if err != nil {
		fmt.Println("Error at HandleService: opening messages.json,\nMsg: ", err)
		return Players{}, err
	}
	defer playersFile.Close()
	playersByteValue, _ := ioutil.ReadAll(playersFile)
	var players Players
	json.Unmarshal(playersByteValue, &players)
	return players, nil
}

func (messageRepository) GetPlayersStats(playerID string, data *Response) error {
	urlAPI := fmt.Sprintf("https://gamersclub.com.br/api/box/historyFilterDate/%s/2022-11", playerID)
	req, err := http.NewRequest(http.MethodGet, urlAPI, nil)
	if err != nil {
		fmt.Println("Error at HandleService: calling GCs API")
	}

	gclubsess := "gclubsess=" + os.Getenv("GCLUB_SESS")

	req.Header.Add("Cookie", gclubsess)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
		return err
	}

	return nil
}
