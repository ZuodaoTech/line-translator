package handler

import (
	"errors"
	"net/http"

	"github.com/zuodaotech/line-translator/common/assistant"
	"github.com/zuodaotech/line-translator/config"
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/handler/line"
	"github.com/zuodaotech/line-translator/handler/render"
	"github.com/zuodaotech/line-translator/handler/sys"
	"github.com/zuodaotech/line-translator/session"

	"github.com/go-chi/chi"
)

func New(
	cfg Config,
	syscfg *config.Config,
	se *session.Session,
	composeAssistant *assistant.Assistant,

	conversations core.ConversationStore,

	taskz core.TaskService,
) Server {

	return Server{
		cfg:       cfg,
		syscfg:    syscfg,
		session:   se,
		assistant: composeAssistant,

		conversations: conversations,

		taskz: taskz,
	}
}

type (
	Config struct {
	}

	Server struct {
		cfg       Config
		syscfg    *config.Config
		session   *session.Session
		assistant *assistant.Assistant

		conversations core.ConversationStore

		taskz core.TaskService
	}
)

func (s Server) HandleRest() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", sys.RenderRoot())
	})

	r.Route("/line", func(r chi.Router) {
		r.Post("/webhook", line.HandleWebhook(s.syscfg, s.conversations, s.taskz))
		r.Get("/webhook", sys.RenderRoot())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.NotFound(w, http.StatusNotFound, errors.New("not found"))
	})

	return r
}
