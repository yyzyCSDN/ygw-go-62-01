package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/gc"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
	"columnstore/internal/query"
)

func main() {
	cfg, err := ParseConfig(os.Args[1:])
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		log.Fatalf("data dir: %v", err)
	}
	store := col.NewStore(cfg.DataDir)
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	engine := query.NewEngine(store, parts, metaRegistry, dicts)
	engine.SetChunk(cfg.Chunk)
	reclaimer := gc.New(store, gc.DefaultPolicy().WithTTL(cfg.Retention))
	compressService := compress.NewService(dicts, metaRegistry)
	app := NewApp(store, parts, metaRegistry, loader, engine, reclaimer, compressService, cfg.DataDir)
	if err := seedDemo(loader); err != nil {
		log.Printf("demo seed skipped: %v", err)
	}
	handler, err := NewServer(app, "web")
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	log.Printf("columnstore listening on %s", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
