package main

import "net/http"

// NewServer builds the http.Handler for the column store, serving the
// console page from the given static root directory.
func NewServer(app *App, webRoot string) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", app.health)
	mux.HandleFunc("/api/partitions", app.partitions)
	mux.HandleFunc("/api/blocks", app.blocks)
	mux.HandleFunc("/api/stats", app.scanStats)
	mux.HandleFunc("/api/verify", app.handleVerify)
	mux.HandleFunc("/api/load", app.handleLoad)
	mux.HandleFunc("/api/query", app.handleQuery)
	mux.HandleFunc("/api/compress", app.handleCompress)
	mux.HandleFunc("/api/gc", app.handleGC)
	mux.Handle("/", http.FileServer(http.Dir(webRoot)))
	return mux, nil
}
