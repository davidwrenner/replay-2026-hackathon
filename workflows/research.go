package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
	"go.temporal.io/sdk/workflow"
)

// DataSource represents an external data provider
type DataSource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"apiKey"`
}

// ResearchWorkflowInput is the input for the research workflow.
type ResearchWorkflowInput struct {
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
	DataSources   []DataSource `json:"dataSources"`

	// Individual job configs for each data source (nil = skip)
	BloombergConfig         *activities.ResearchInput `json:"bloombergConfig"`
	DowJonesConfig          *activities.ResearchInput `json:"dowJonesConfig"`
	LexisNexisConfig        *activities.ResearchInput `json:"lexisNexisConfig"`
	NYTimesConfig           *activities.ResearchInput `json:"nyTimesConfig"`
	PolymarketConfig        *activities.ResearchInput `json:"polymarketConfig"`
	RedditConfig            *activities.ResearchInput `json:"redditConfig"`
	RefinitivConfig         *activities.ResearchInput `json:"refinitivConfig"`
	TwitterConfig           *activities.ResearchInput `json:"twitterConfig"`
	WallStreetJournalConfig *activities.ResearchInput `json:"wallStreetJournalConfig"`
	YouTubeConfig           *activities.ResearchInput `json:"youTubeConfig"`
	ResearchConfig          *activities.ResearchInput `json:"researchConfig"`
}

// ResearchWorkflowOutput is the output from the research workflow.
type ResearchWorkflowOutput struct {
	Research string `json:"research"`
}

// ResearchWorkflow orchestrates the research process.
func ResearchWorkflow(ctx workflow.Context, input ResearchWorkflowInput) (*ResearchWorkflowOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Check for prompt injection
	err := workflow.ExecuteActivity(ctx, "CheckInjection", activities.InjectionCheckInput{
		Prompt:        input.Prompt,
		RiskTolerance: input.RiskTolerance,
	}).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Step 2: Fetch data from sources
	err = workflowext.ExecuteOptional(ctx, "Bloomberg", input.BloombergConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "DowJones", input.DowJonesConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "LexisNexis", input.LexisNexisConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "NYTimes", input.NYTimesConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Polymarket", input.PolymarketConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Reddit", input.RedditConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Refinitiv", input.RefinitivConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Twitter", input.TwitterConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "WallStreetJournal", input.WallStreetJournalConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "YouTube", input.YouTubeConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Step 3: Perform research
	var result activities.ResearchOutput
	err = workflowext.ExecuteOptional(ctx, "Research", input.ResearchConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
