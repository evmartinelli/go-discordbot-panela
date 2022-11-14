package routes

import (
	"github.com/evmartinelli/go-discordbot-panela/internal/discord"
	messageDelivery "github.com/evmartinelli/go-discordbot-panela/internal/messages/delivery"
	"go.uber.org/fx"
)

// NewRoutes new Routes Handler
func NewRoutes(discord discord.Discord, messageDelivery messageDelivery.Delivery) {
	discord.AddHandler(messageDelivery.GetMessageHandler)
}

// Module .
var Module = fx.Options(
	fx.Invoke(NewRoutes),
)
