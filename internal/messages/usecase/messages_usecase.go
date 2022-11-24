package usecase

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
	"golang.org/x/exp/constraints"
)

// Usecase interface
type Usecase interface {
	GetPanelaMatches() (string, error)
	GetPanelaKAST() (string, error)
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
	matches := make(map[string]int)

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		return "", err
	}

	for _, v := range players.Players {
		var stats repository.Response
		err = mu.messagesRepository.GetPlayersStats(v.Url, &stats)
		if err != nil {
			return "", err
		}

		matches[v.Nick] = stats.Matches.Matches

	}

	return createKeyValuePairs(SortKeys(matches), matches, "%v jogou \"%v\"\n"), nil
}

func (mu messagesUsecase) GetPanelaKAST() (string, error) {
	matches := make(map[string]string)

	players, err := mu.messagesRepository.GetPlayersURL()
	if err != nil {
		return "", err
	}

	for _, v := range players.Players {
		var stats repository.Response
		err = mu.messagesRepository.GetPlayersStats(v.Url, &stats)
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
