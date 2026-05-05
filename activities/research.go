package activities

import (
	"context"
	"encoding/base64"
)

// ResearchInput is the input for the research activity.
type ResearchInput struct {
	Query string `json:"query"`
}

// ResearchOutput is the output from the research activity.
type ResearchOutput struct {
	Research string `json:"research"`
}

// Activities holds dependencies for all activities.
type Activities struct {
	ReportData []byte
}

// Research is a Temporal activity that generates a research report.
func (a *Activities) Research(ctx context.Context, input ResearchInput) (*ResearchOutput, error) {
	encoded := base64.StdEncoding.EncodeToString(a.ReportData)
	return &ResearchOutput{
		Research: encoded,
	}, nil
}
