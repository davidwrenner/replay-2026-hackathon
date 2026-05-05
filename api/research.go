package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type DataSource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"apiKey"`
}

type ResearchRequest struct {
	DataSources   []DataSource `json:"dataSources"`
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
}

type ResearchResponse struct {
	ResearchID      string   `json:"researchId"`
	Status          string   `json:"status"`
	AcceptedSources []string `json:"acceptedSources"`
	UnknownSources  []string `json:"unknownSources,omitempty"`
	RiskTolerance   string   `json:"riskTolerance"`
	MaxBudget       int64    `json:"maxBudget"`
}

// sourceAliases maps normalized names that don't directly match validSources keys
// to their canonical validSources keys.
var sourceAliases = map[string]string{
	"new_york_times": "nytimes",
}

// normalizeSourceName takes a display-style source name and returns the canonical
// validSources key, plus a bool indicating whether it was resolved.
func normalizeSourceName(name string) (string, bool) {
	normalized := strings.ToLower(strings.ReplaceAll(name, " ", "_"))

	// Direct lookup first
	if _, ok := validSources[normalized]; ok {
		return normalized, true
	}

	// Alias lookup
	if canonical, ok := sourceAliases[normalized]; ok {
		if _, ok := validSources[canonical]; ok {
			return canonical, true
		}
	}

	return normalized, false
}

func (s *Server) handleResearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ResearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid request body",
			"message": "request body must be valid JSON matching the ResearchRequest schema",
		})
		return
	}

	// Validate riskTolerance
	riskTolerance := strings.ToLower(req.RiskTolerance)
	if riskTolerance == "" {
		riskTolerance = "medium"
	} else if riskTolerance != "low" && riskTolerance != "medium" && riskTolerance != "high" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid riskTolerance",
			"message": "riskTolerance must be one of: low, medium, high",
		})
		return
	}

	// Parse maxBudget
	maxBudgetStr := strings.TrimSpace(req.MaxBudget)
	maxBudget, err := strconv.ParseInt(maxBudgetStr, 10, 64)
	if err != nil || maxBudget <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid maxBudget",
			"message": "maxBudget must be a positive integer",
		})
		return
	}

	// Normalize sources and partition into accepted/unknown
	acceptedSources := make([]string, 0)
	unknownSources := make([]string, 0)

	for _, ds := range req.DataSources {
		canonical, ok := normalizeSourceName(ds.Name)
		if ok {
			acceptedSources = append(acceptedSources, canonical)
		} else {
			unknownSources = append(unknownSources, ds.Name)
		}
	}

	researchID := uuid.NewString()

	resp := ResearchResponse{
		ResearchID:      researchID,
		Status:          "pending",
		AcceptedSources: acceptedSources,
		RiskTolerance:   riskTolerance,
		MaxBudget:       maxBudget,
	}
	if len(unknownSources) > 0 {
		resp.UnknownSources = unknownSources
	}

	json.NewEncoder(w).Encode(resp)
}
