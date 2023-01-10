package repository_test

import (
	"testing"

	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
)

func TestGetAudioItems(t *testing.T) {
	data, err := repository.NewMessageRepository().GetAudioItems()
	if err != nil {
	}
	t.Log(data)
}
