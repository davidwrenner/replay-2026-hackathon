package activities

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"
	"text/template"
	"time"
)

// DataSource represents an external data provider
type DataSource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"apiKey"`
}

// ResearchInput is the input for the research activity.
type ResearchInput struct {
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
	DataSources   []DataSource `json:"dataSources"`
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
	// Parse DataSources into easier to check booleans
	hasSocial := false
	hasMarket := false
	sourceNames := []string{}
	for _, ds := range input.DataSources {
		name := strings.ToLower(ds.Name)
		if name == "twitter" || name == "reddit" {
			hasSocial = true
		} else if name == "bloomberg" || name == "wsj" || name == "refinitiv" || name == "new york times" {
			hasMarket = true
		}
		sourceNames = append(sourceNames, ds.Name)
	}
	
	if len(sourceNames) == 0 {
		sourceNames = append(sourceNames, "None provided")
	}

	// Pick Kalshi trade based on Risk Tolerance
	marketName := "Will the Fed cut rates in June 2026?"
	ticker := "FED-26JUN-T25"
	price := 65
	side := "yes"
	action := "buy"
	contracts := 100

	if input.RiskTolerance == "medium" {
		marketName = "Will S&P 500 close above 5500 this week?"
		ticker = "INX-26MAY10-B5500"
		price = 45
		contracts = 200
	} else if input.RiskTolerance == "high" {
		marketName = "Will Bitcoin hit $100k by EOY?"
		ticker = "BTC-100K-26DEC31"
		price = 28
		contracts = 500
	}

	// Template Data
	tmplData := map[string]interface{}{
		"Date":          time.Now().Format("Jan 2, 2006"),
		"SourceCount":   len(input.DataSources),
		"Sources":       strings.Join(sourceNames, ", "),
		"HasSocial":     hasSocial,
		"HasMarket":     hasMarket,
		"Prompt":        input.Prompt,
		"RiskTolerance": input.RiskTolerance,
		"MaxBudget":     input.MaxBudget,
		"MarketName":    marketName,
		"Ticker":        ticker,
		"Side":          side,
		"Action":        action,
		"Price":         price,
		"Contracts":     contracts,
	}

	tmpl := `
# Market Intelligence Report

**Generated:** {{ .Date }} | **Sources ({{ .SourceCount }}):** {{ .Sources }}

---
{{ if .HasMarket }}
## 1. Market Overview

The market shows dynamic movement based on your prompt: "{{ .Prompt }}". Traditional indicators show mixed performance with a lean towards growth sectors.

### Index Performance

` + "```vega-lite" + `
{
  "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
  "title": "Major Index Performance (%)",
  "width": 400,
  "height": 200,
  "data": {
    "values": [
      {"index": "S&P 500", "change": 0.85, "color": "positive"},
      {"index": "Dow Jones", "change": 0.39, "color": "positive"},
      {"index": "Nasdaq", "change": -0.25, "color": "negative"},
      {"index": "Russell 2000", "change": 0.59, "color": "positive"}
    ]
  },
  "mark": "bar",
  "encoding": {
    "x": {"field": "index", "type": "nominal", "axis": {"labelAngle": 0}, "title": null},
    "y": {"field": "change", "type": "quantitative", "title": "Change (%)"},
    "color": {
      "field": "color",
      "type": "nominal",
      "scale": {"domain": ["positive", "negative"], "range": ["#22c55e", "#ef4444"]},
      "legend": null
    }
  }
}
` + "```" + `
{{ end }}

{{ if .HasSocial }}
## 2. Social Sentiment Analysis

Aggregated sentiment from social sources reveals strong retail positioning and speculative interest around high-growth tech stocks.

### Ticker Mentions & Sentiment

` + "```vega-lite" + `
{
  "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
  "title": "Social Media Ticker Sentiment",
  "width": 400,
  "height": 250,
  "data": {
    "values": [
      {"ticker": "NVDA", "mentions": 45230, "sentiment": 0.78},
      {"ticker": "GME", "mentions": 32100, "sentiment": 0.65},
      {"ticker": "TSLA", "mentions": 28900, "sentiment": 0.42},
      {"ticker": "AAPL", "mentions": 21500, "sentiment": 0.71}
    ]
  },
  "mark": "circle",
  "encoding": {
    "x": {"field": "mentions", "type": "quantitative", "title": "Mentions (24h)", "scale": {"zero": false}},
    "y": {"field": "sentiment", "type": "quantitative", "title": "Sentiment Score", "scale": {"domain": [0, 1]}},
    "size": {"field": "mentions", "type": "quantitative", "legend": null},
    "color": {"field": "sentiment", "type": "quantitative", "scale": {"scheme": "redyellowgreen", "domain": [0, 1]}, "title": "Sentiment"},
    "tooltip": [
      {"field": "ticker", "type": "nominal"},
      {"field": "mentions", "type": "quantitative"},
      {"field": "sentiment", "type": "quantitative", "format": ".2f"}
    ]
  }
}
` + "```" + `
{{ end }}

{{ if and (not .HasMarket) (not .HasSocial) }}
## 1. No Data Sources Enabled

You did not provide API keys for any data sources. The AI model is relying entirely on its baseline training data and your prompt: "{{ .Prompt }}".
{{ end }}

---

## Conclusion & Kalshi Trade Recommendation

Based on your **{{ .RiskTolerance }}** risk tolerance and budget of **${{ .MaxBudget }}**, we recommend the following trade on Kalshi.

### Market: {{ .MarketName }}
**Ticker:** ` + "`{{ .Ticker }}`" + `
**Action:** Buy {{ .Side }}
**Target Price:** {{ .Price }}¢
**Contracts:** {{ .Contracts }}

#### Execute Trade via Kalshi API
` + "```bash" + `
curl -X POST https://api.elections.kalshi.com/trade-api/v2/portfolio/orders \
  -H "Authorization: Bearer <YOUR_API_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "ticker": "{{ .Ticker }}",
    "action": "buy",
    "side": "{{ .Side }}",
    "count": {{ .Contracts }},
    "yes_price": {{ .Price }},
    "client_order_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
` + "```" + `
`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, tmplData); err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &ResearchOutput{
		Research: encoded,
	}, nil
}
