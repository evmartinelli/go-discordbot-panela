package provider

import (
	"github.com/evmartinelli/go-discordbot-panela/internal/voice/usecase"
	"go.uber.org/fx"
)

// UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewVoiceUsecase),
)
