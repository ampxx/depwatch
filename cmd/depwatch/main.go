package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/depwatch/internal/checker"
	"github.com/depwatch/internal/config"
	"github.com/depwatch/internal/notifier"
	"github.com/depwatch/internal/store"
	"github.com/depwatch/internal/watcher"
)

const defaultProxyURL = "https://proxy.golang.org"

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	storePath := flag.String("store", "depwatch.store.json", "path to version store file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	s, err := store.New(*storePath)
	if err != nil {
		log.Fatalf("store: %v", err)
	}

	client := checker.NewClient(defaultProxyURL)
	n := notifier.New(cfg.WebhookURL)
	w := watcher.New(cfg, client, n, s)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("depwatch: starting")
	w.Run(ctx)
	log.Println("depwatch: stopped")
}
