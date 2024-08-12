package assistant

import (
	"github.com/lyricat/goutils/ai"
)

type (
	Assistant struct {
		cfg    Config
		aiInst *ai.Instant
	}
	Config struct {
		BotID uint64
	}
)

func New(cfg Config, aiInst *ai.Instant) *Assistant {
	return &Assistant{
		cfg:    cfg,
		aiInst: aiInst,
	}
}
