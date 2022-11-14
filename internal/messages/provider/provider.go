package provider

import (
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/delivery"
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/repository"
	"github.com/evmartinelli/go-discordbot-panela/internal/messages/usecase"
	"go.uber.org/fx"
)

// DeliveryModule .
var DeliveryModule = fx.Options(
	fx.Provide(delivery.NewMessageDelivery),
)

// RepositoryModule .
var RepositoryModule = fx.Options(
	fx.Provide(repository.NewMessageRepository),
)

// UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewMessagesUsecase),
)
