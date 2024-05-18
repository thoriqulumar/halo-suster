package main

import (
	"context"
	application "halo-suster"
	"halo-suster/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	application.Start(cfg)
}
