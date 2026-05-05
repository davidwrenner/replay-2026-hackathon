package api

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/workflows"
	"go.temporal.io/sdk/client"
)

//go:embed mock/*.json
var mockData embed.FS

//go:embed report.mdx
var reportMDX []byte

// GetReportData returns the embedded report data.
func GetReportData() []byte {
	return reportMDX
}

const TaskQueue = "research-task-queue"

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
	addr           string
	temporalClient client.Client
}

func NewServer(addr string, temporalClient client.Client) *Server {
	return &Server{
		addr:           addr,
		temporalClient: temporalClient,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /data/{name}", s.handleData)
	mux.HandleFunc("GET /sources", s.handleSources)
	mux.HandleFunc("POST /research", s.handleResearch)

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

// ResearchRequest matches the UI's ResearchRequest type.
type ResearchRequest struct {
	DataSources   []DataSource `json:"dataSources"`
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
}

// DataSource represents a data source from the request.
type DataSource struct {
	Name string `json:"name"`
}

func (s *Server) handleResearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var req ResearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Convert to workflow input
	dataSources := make([]workflows.DataSource, len(req.DataSources))
	for i, ds := range req.DataSources {
		dataSources[i] = workflows.DataSource{Name: ds.Name}
	}

	// Build base job config
	jobConfig := &activities.ResearchInput{
		Prompt:        req.Prompt,
		RiskTolerance: req.RiskTolerance,
		MaxBudget:     req.MaxBudget,
		DataSources:   convertDataSources(req.DataSources),
	}

	input := workflows.ResearchWorkflowInput{
		DataSources:             dataSources,
		Prompt:                  req.Prompt,
		RiskTolerance:           req.RiskTolerance,
		MaxBudget:               req.MaxBudget,
		BloombergConfig:         configIf(hasDataSource(req.DataSources, "bloomberg"), jobConfig),
		DowJonesConfig:          configIf(hasDataSource(req.DataSources, "dow_jones"), jobConfig),
		LexisNexisConfig:        configIf(hasDataSource(req.DataSources, "lexisnexis"), jobConfig),
		NYTimesConfig:           configIf(hasDataSource(req.DataSources, "nytimes"), jobConfig),
		PolymarketConfig:        configIf(hasDataSource(req.DataSources, "polymarket"), jobConfig),
		RedditConfig:            configIf(hasDataSource(req.DataSources, "reddit"), jobConfig),
		RefinitivConfig:         configIf(hasDataSource(req.DataSources, "refinitiv"), jobConfig),
		TwitterConfig:           configIf(hasDataSource(req.DataSources, "twitter"), jobConfig),
		WallStreetJournalConfig: configIf(hasDataSource(req.DataSources, "wall_street_journal"), jobConfig),
		YouTubeConfig:           configIf(hasDataSource(req.DataSources, "youtube"), jobConfig),
		ResearchConfig:          jobConfig,
	}

	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: TaskQueue,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflows.ResearchWorkflow, input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to start workflow"})
		return
	}

	var result workflows.ResearchWorkflowOutput
	if err := we.Get(ctx, &result); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"research": result.Research,
	})
}

func hasDataSource(ds []DataSource, name string) bool {
	for _, d := range ds {
		if strings.EqualFold(d.Name, name) {
			return true
		}
	}
	return false
}

func configIf(enabled bool, cfg *activities.ResearchInput) *activities.ResearchInput {
	if enabled {
		return cfg
	}
	return nil
}

func convertDataSources(ds []DataSource) []activities.DataSource {
	res := make([]activities.DataSource, len(ds))
	for i, d := range ds {
		res[i] = activities.DataSource{Name: d.Name}
	}
	return res
}
