package usecase

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
	"golang.org/x/exp/constraints"
)

// Usecase interface
type Usecase interface {
	GetPanelaMatches() (string, error)
	GetPanelaMatchesAsync(rchan chan repository.Response)
	GetPanelaKAST() (string, error)
	GetPanelaLoss() (string, error)
}

type messagesUsecase struct {
	messagesRepository repository.Repository
}

// NewMessagesUsecase new message delivery
func NewMessagesUsecase(mr repository.Repository) Usecase {
	return &messagesUsecase{
		messagesRepository: mr,
	}
}

// GetPlayersURL return list of players
func (mu messagesUsecase) GetPanelaMatches() (string, error) {
	defer timeTrack(time.Now(), "MATCHES")
	matches := make(map[string]int)

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		return "", err
	}

	for _, v := range players.Players {
		var stats repository.Response
		err = mu.messagesRepository.GetPlayerStats(v.Url, &stats)
		if err != nil {
			return "", err
		}

		matches[v.Nick] = stats.Matches.Matches

	}

	return createKeyValuePairs(SortKeys(matches), matches, "%v jogou \"%v\"\n"), nil
}

// GetPlayersURL return list of players
func (mu messagesUsecase) GetPanelaLoss() (string, error) {
	defer timeTrack(time.Now(), "Loss")
	matches := make(map[string]int)

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		return "", err
	}

	for _, v := range players.Players {
		var stats repository.Response
		err = mu.messagesRepository.GetPlayerStats(v.Url, &stats)
		if err != nil {
			return "", err
		}

		result := float64(stats.Matches.Loss) / float64(stats.Matches.Matches) * 100

		matches[v.Nick] = int(result)

	}

	return createKeyValuePairs(SortKeys(matches), matches, "%v perdeu \"%v\" porcento \n"), nil
}

// GetPlayersURL return list of players
func (mu messagesUsecase) GetPanelaMatchesAsync(rchan chan repository.Response) {
	defer timeTrack(time.Now(), "MATCHESASYNC")

	defer close(rchan)
	results := []chan repository.Response{}

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		log.Println(err)
	}

	for i, player := range players.Players {
		results = append(results, make(chan repository.Response))
		go mu.messagesRepository.GetPlayerStatsAsync(player.Url, results[i])
	}

	for i := range results {
		for r1 := range results[i] {
			rchan <- r1
		}
	}
}

func (mu messagesUsecase) GetPanelaKAST() (string, error) {
	defer timeTrack(time.Now(), "KAST")
	matches := make(map[string]string)

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		return "", err
	}

	for _, v := range players.Players {
		var stats repository.Response
		err = mu.messagesRepository.GetPlayerStats(v.Url, &stats)
		if err != nil {
			return "", err
		}

		for _, stat := range stats.Stat {
			if stat.Stat == "KAST%" {
				matches[v.Nick] = stat.Value
			}
		}
	}

	return createKeyValuePairs(SortKeys(matches), matches, "%v KAST \"%v\"\n"), nil
}

func createKeyValuePairs[K constraints.Ordered, V constraints.Ordered](keys []K, m map[K]V, literal string) string {
	b := new(bytes.Buffer)

	for _, v := range keys {
		fmt.Fprintf(b, literal, v, m[v])
	}

	return b.String()
}

func SortKeys[K constraints.Ordered, V constraints.Ordered](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func scrapListURL(urlToProcess []string, rchan chan repository.Response) {
	defer close(rchan)
	results := []chan repository.Response{}

	for i, url := range urlToProcess {
		results = append(results, make(chan repository.Response))
		go repository.NewMessageRepository().GetPlayerStatsAsync(url, results[i])
	}

	for i := range results {
		for r1 := range results[i] {
			rchan <- r1
		}
	}
}
