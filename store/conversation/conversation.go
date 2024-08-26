package conversation

import (
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/store"
	"github.com/zuodaotech/line-translator/store/conversation/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/conversation/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.ConversationStore) {}, core.Conversation{})
		},
	)
}

func New(h *store.Handler) core.ConversationStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Conversation).(core.ConversationStore)
	if !ok {
		panic("dao.Conversation is not core.ConversationStore")
	}

	return &storeImpl{
		ConversationStore: v,
	}
}

type storeImpl struct {
	core.ConversationStore
}
