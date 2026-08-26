package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/gc"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
	"columnstore/internal/query"
)

func TestHealthEndpoint(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	engine := query.NewEngine(store, parts, metaRegistry, dicts)
	reclaimer := gc.New(store, gc.DefaultPolicy())
	app := NewApp(store, parts, metaRegistry, loader, engine, reclaimer, compress.NewService(dicts, metaRegistry), t.TempDir())
	handler, err := NewServer(app, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("health body must not be empty")
	}
}

func TestConsolePageServed(t *testing.T) {
	store := col.NewStore(t.TempDir())
	parts := part.NewRegistry()
	metaRegistry := meta.NewRegistry()
	dicts := compress.NewDictRegistry()
	loader := load.NewLoader(store, parts, metaRegistry)
	engine := query.NewEngine(store, parts, metaRegistry, dicts)
	reclaimer := gc.New(store, gc.DefaultPolicy())
	app := NewApp(store, parts, metaRegistry, loader, engine, reclaimer, compress.NewService(dicts, metaRegistry), t.TempDir())
	webRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(webRoot, "console.html"), []byte("<html>console</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler, err := NewServer(app, webRoot)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/console.html", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected console page, got %d", rec.Code)
	}
}
