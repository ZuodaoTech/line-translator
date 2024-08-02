package cmd

import (
	"github.com/zuodaotech/line-translator/config"
	"github.com/zuodaotech/line-translator/store"

	_ "github.com/zuodaotech/line-translator/store/task"

	"github.com/spf13/cobra"
)

func genCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate database operation code",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			cfg := ctx.Value("config").(*config.Config)
			h := store.MustInit(store.Config{
				Driver: cfg.DB.Driver,
				DSN:    cfg.DB.DSN,
			})
			h.Generate()
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(genCmd())
}
