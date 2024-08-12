/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/lyricat/goutils/ai"
	"github.com/zuodaotech/line-translator/common/assistant"
	"github.com/zuodaotech/line-translator/config"
	"github.com/zuodaotech/line-translator/handler"
	taskZ "github.com/zuodaotech/line-translator/service/task"
	"github.com/zuodaotech/line-translator/store"
	"github.com/zuodaotech/line-translator/store/task"
	"github.com/zuodaotech/line-translator/worker/tasker"

	"github.com/zuodaotech/line-translator/session"

	"strconv"
	"syscall"
	"time"

	"github.com/zuodaotech/line-translator/worker"

	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	httpdOnly bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		se := session.From(ctx)
		cfg := ctx.Value("config").(*config.Config)

		se.WithJWTSecret([]byte(cfg.Auth.JwtSecret))

		aiInst := ai.New(ai.Config{
			OpenAIApiKey:                     cfg.OpenAI.APIKey,
			AzureOpenAIApiKey:                cfg.Azure.OpenAI.APIKey,
			AzureOpenAIEndpoint:              cfg.Azure.OpenAI.Endpoint,
			AzureOpenAIGptDeploymentID:       cfg.Azure.OpenAI.GptDeploymentID,
			AzureOpenAIEmbeddingDeploymentID: cfg.Azure.OpenAI.EmbeddingDeploymentID,
			Provider:                         "azure",
			Debug:                            true,
		})

		composeAssistant := assistant.New(assistant.Config{}, aiInst)

		// db & stores
		h := store.MustInit(store.Config{
			Driver: cfg.DB.Driver,
			DSN:    cfg.DB.DSN,
		})
		tasks := task.New(h)

		taskz := taskZ.New(taskZ.Config{}, tasks)

		g := errgroup.Group{}

		if !httpdOnly {
			workers := []worker.Worker{
				tasker.New(tasker.Config{
					LineChannelID:     cfg.Line.ChannelID,
					LineChannelKey:    cfg.Line.ChannelKey,
					LineJWTPrivateKey: cfg.Line.JWTPrivateKey,

					AzureAPIKey:   cfg.Azure.Speech.APIKey,
					AzureEndpoint: cfg.Azure.Speech.Endpoint,
				}, composeAssistant, tasks, taskz),
			}

			g, ctx := errgroup.WithContext(ctx)
			for idx := range workers {
				w := workers[idx]
				g.Go(func() error {
					return w.Run(ctx)
				})
			}
		}

		g.Go(func() error {
			var err error
			mux := chi.NewMux()
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.StripSlashes)
			mux.Use(cors.AllowAll().Handler)
			mux.Use(middleware.Logger)
			mux.Use(middleware.NewCompressor(5).Handler)
			{
				restSvr := handler.New(handler.Config{},
					cfg,
					se,
					composeAssistant,
					taskz,
				)
				restHandler := restSvr.HandleRest()

				mux.Mount("/", restHandler)
			}

			port := 8080
			if len(args) > 0 {
				port, err = strconv.Atoi(args[0])
				if err != nil {
					port = 8080
				}
			}

			// launch server
			if err != nil {
				panic(err)
			}
			addr := fmt.Sprintf(":%d", port)

			svr := &http.Server{
				Addr:    addr,
				Handler: mux,
			}

			slog.Info("[server] run httpd server", "addr", addr)
			if err := svr.ListenAndServe(); err != http.ErrServerClosed {
				slog.Error("[server] server aborted", "error", err)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			slog.Info("[server] shutdown server")

			// Create a context with a timeout of 5 seconds to gracefully shutdown the server
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := svr.Shutdown(ctx); err != nil {
				slog.Error("[server] server shutdown failed", "error", err)
			}

			slog.Info("[server] server shutdown complete")
			return nil
		})

		if err := g.Wait(); err != nil {
			slog.Error("[server] run httpd & worker", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&httpdOnly, "httpd", "", false, "only run httpd, no workers")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
