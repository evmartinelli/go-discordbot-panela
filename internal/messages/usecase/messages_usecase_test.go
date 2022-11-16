package usecase_test

import (
	"fmt"
	"testing"

	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/usecase"
)

type messageUseCaseFixture struct {
	usecase usecase.Usecase
	repo    repository.Repository
}

func TestTestListPostsUseCase(t *testing.T) {
	setup := func() *messageUseCaseFixture {
		repo := repository.NewMessageRepository()
		usecase := usecase.NewMessagesUsecase(repo)
		return &messageUseCaseFixture{
			usecase: usecase,
			repo:    repo,
		}
	}

	t.Run("Given no post exists, it returns an empty slice", func(t *testing.T) {
		f := setup()

		player, err := f.usecase.GetPanelaMatches()
		if err != nil {
			fmt.Println(player)
		}
		fmt.Println(player)
	})
}
