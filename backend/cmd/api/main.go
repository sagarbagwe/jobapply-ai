package main

import (
	"context"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/ai"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/auth"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/config"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/httpapi"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/store"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, e := config.Load()
	if e != nil {
		log.Fatal(e)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, e := store.Open(ctx, cfg.DatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	origins := map[string]bool{}
	for _, o := range cfg.AllowedOrigins {
		origins[o] = true
	}
	if cfg.Env != "production" {
		origins["chrome-extension://development"] = true
	}
	api := httpapi.API{AI: ai.NewGemini(cfg.GeminiKey, cfg.GeminiModel), Store: db, Tokens: auth.New(cfg.JWTSecret, cfg.AccessTTL), Origins: origins}
	log.Printf("JobApply AI API listening on %s", cfg.HTTPAddr)
	if e = httpapi.Router(api).Run(cfg.HTTPAddr); e != nil {
		log.Fatal(e)
	}
}
