package main

import (
	"context"

	"github.com/zuodaotech/line-translator/cmd"
	"github.com/zuodaotech/line-translator/session"
)

var (
	Version = "0.0.1"
)

func main() {
	ctx := context.Background()
	s := &session.Session{Version: Version}
	ctx = session.With(ctx, s)

	cmd.ExecuteContext(ctx)
}
