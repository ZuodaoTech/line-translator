package task

import (
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/store"
	"github.com/zuodaotech/line-translator/store/task/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/task/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.TaskStore) {}, core.Task{})
		},
	)
}

func New(h *store.Handler) core.TaskStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Task).(core.TaskStore)
	if !ok {
		panic("dao.Task is not core.TaskStore")
	}

	return &storeImpl{
		TaskStore: v,
	}
}

type storeImpl struct {
	core.TaskStore
}
