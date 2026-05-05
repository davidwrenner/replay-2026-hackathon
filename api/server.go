package api

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

//go:embed mock/*.json
var mockData embed.FS

var validSources = map[string]string{
	"twitter":             "mock/twitter.json",
	"reddit":              "mock/reddit.json",
	"nytimes":             "mock/nytimes.json",
	"wall_street_journal": "mock/wall_street_journal.json",
	"lexisnexis":          "mock/lexisnexis.json",
	"refinitiv":           "mock/refinitiv.json",
	"bloomberg":           "mock/bloomberg.json",
	"dow_jones":           "mock/dow_jones.json",
	"youtube":             "mock/youtube.json",
	"polymarket":          "mock/polymarket.json",
}

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /data/{name}", s.handleData)
	mux.HandleFunc("GET /sources", s.handleSources)

	log.Printf("API server starting on %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	filePath, ok := validSources[name]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "source not found",
			"message": "use GET /sources to see available data sources",
		})
		return
	}

	file, err := mockData.Open(filePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read data"})
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, file)
}

func (s *Server) handleSources(w http.ResponseWriter, r *http.Request) {
	sources := make([]string, 0, len(validSources))
	for name := range validSources {
		sources = append(sources, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sources": sources,
		"usage":   "GET /data/{source_name}",
	})
}
