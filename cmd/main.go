package main

import (
	"context"
	application "helo-suster"
	"helo-suster/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	application.Start(cfg)
}
