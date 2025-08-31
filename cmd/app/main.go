package main

import (
	"AggregationService/internal/app"
	"AggregationService/internal/config"
	"AggregationService/internal/pkg/logger"
	"context"
	"os"
)

func main() {
	ctx := app.InitContextWithLogger(context.Background())
	cfg := config.MustLoad()
	provider := app.NewAppProvider()
	application := app.New(ctx, cfg, provider)

	if err := application.Run(ctx); err != nil {
		logger.FromContext(ctx).
			Error("Server stopped with error", "error", err)
		os.Exit(1)
	}
}
