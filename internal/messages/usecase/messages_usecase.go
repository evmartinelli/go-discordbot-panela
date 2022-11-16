package usecase

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
)

// Usecase interface
type Usecase interface {
	GetPanelaMatches() (string, error)
	// GetPanelaKAST() (string, error)
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

	return createKeyValuePairs(matches), nil
}

// GetPlayersURL return list of players
// func (mu messagesUsecase) GetPanelaKAST() (string, error) {
// 	matches := make(map[string]string)

// 	players, err := mu.messagesRepository.GetPlayersURL()
// 	if err != nil {
// 		return "", err
// 	}

// 	for i, v := range players.Players {
// 		var stats repository.Response
// 		err = mu.messagesRepository.GetPlayersStats(v.Url, &stats)
// 		if err != nil {
// 			return "", err
// 		}

// 		if stats.Stat[i].Stat == "KAST%" {
// 			matches[v.Nick] = stats.Stat[i].Value + " - " + stats.Stat[i].Stat
// 		}

// 	}

// 	return createKeyValuePairs(matches), nil
// }

func createKeyValuePairs(m map[string]int) string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	b := new(bytes.Buffer)

	for _, v := range keys {
		fmt.Fprintf(b, "%v jogou \"%v\"\n", v, m[v])
	}

	return b.String()
}
