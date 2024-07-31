package assistant

import (
	"github.com/zuodaotech/line-translator/common/ai"
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
