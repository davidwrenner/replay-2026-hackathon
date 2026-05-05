import { useState } from "react";
import "./index.css";
import type { AppPage, AppPhase, DataSource, ResearchFormState, RiskTolerance } from "./types";
import ResearchForm from "./components/ResearchForm";
import ResearchReport from "./components/ResearchReport";
import LoadingScreen from "./components/LoadingScreen";
import ConfigPage from "./components/ConfigPage";
import { submitResearch } from "./api";

const DEFAULT_SOURCES: DataSource[] = [
  { name: "Twitter", description: "Twitter API", apiKey: "", enabled: false },
  { name: "Reddit", description: "Reddit API", apiKey: "", enabled: false },
  { name: "New York Times", description: "New York Times API", apiKey: "", enabled: false },
  { name: "Wall Street Journal", description: "Wall Street Journal API", apiKey: "", enabled: false },
  { name: "LexisNexis", description: "LexisNexis API", apiKey: "", enabled: false },
  { name: "Refinitiv", description: "Refinitiv API", apiKey: "", enabled: false },
  { name: "Bloomberg", description: "Bloomberg API", apiKey: "", enabled: false },
  { name: "Dow Jones", description: "Dow Jones API", apiKey: "", enabled: false },
  { name: "YouTube", description: "YouTube API", apiKey: "", enabled: false },
  { name: "Polymarket", description: "Polymarket API", apiKey: "", enabled: false },
];

export default function App() {
  const [page, setPage] = useState<AppPage>("research");
  const [phase, setPhase] = useState<AppPhase>("form");
  const [dataSources, setDataSources] = useState<DataSource[]>(DEFAULT_SOURCES);
  const [reportMdx, setReportMdx] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [lastPrompt, setLastPrompt] = useState("");

  // Lifted form state so it persists across page switches
  const [prompt, setPrompt] = useState("");
  const [riskTolerance, setRiskTolerance] = useState<RiskTolerance>("medium");
  const [maxBudget, setMaxBudget] = useState("10000");
  const [enabledMap, setEnabledMap] = useState<Record<string, boolean>>({});

  const handleSubmit = async (form: ResearchFormState) => {
    setPhase("loading");
    setError(null);
    setLastPrompt(form.prompt);

    const minDelay = new Promise((r) => setTimeout(r, 3000));

    try {
      const [mdx] = await Promise.all([submitResearch(form), minDelay]);
      setReportMdx(mdx);
      setPhase("report");
    } catch (err) {
      await minDelay;
      setError(err instanceof Error ? err.message : "An unexpected error occurred");
      setPhase("error");
    }
  };

  const handleBack = () => {
    setPhase("form");
    setReportMdx(null);
    setError(null);
  };

  const handleSaveConfig = (sources: DataSource[]) => {
    setDataSources(sources);
    // Auto-enable newly configured sources that aren't already in the map
    setEnabledMap((prev) => {
      const next = { ...prev };
      for (const s of sources) {
        if (s.apiKey.trim().length > 0 && !(s.name in next)) {
          next[s.name] = true;
        }
        // Disable sources whose key was removed
        if (s.apiKey.trim().length === 0) {
          next[s.name] = false;
        }
      }
      return next;
    });
    setPage("research");
  };

  const configuredCount = dataSources.filter((s) => s.apiKey.trim().length > 0).length;

  // Show nav only when on the form or config page
  const showNav = phase === "form" || page === "config";

  return (
    <div className="app-shell">
      {/* Header */}
      <header className="app-header">
        <div className="app-logo">
          <div className="app-logo-icon">K</div>
          <span>Kalshi Bot</span>
        </div>
        {showNav ? (
          <nav className="app-nav" id="main-nav">
            <button
              className={`nav-tab ${page === "research" ? "active" : ""}`}
              onClick={() => setPage("research")}
              id="nav-research"
            >
              Research
            </button>
            <button
              className={`nav-tab ${page === "config" ? "active" : ""}`}
              onClick={() => setPage("config")}
              id="nav-config"
            >
              Configuration
              {configuredCount > 0 && (
                <span className="nav-tab-badge">{configuredCount}</span>
              )}
            </button>
          </nav>
        ) : (
          <span className="app-badge">
            {phase === "loading" ? "Analyzing" : phase === "report" ? "Report" : "Research"}
          </span>
        )}
      </header>

      {/* Main */}
      <main className="app-main">
        <div
          className="app-content"
          style={
            phase === "report" && page === "research"
              ? { maxWidth: 900 }
              : undefined
          }
        >
          {page === "config" && (
            <ConfigPage dataSources={dataSources} onSave={handleSaveConfig} />
          )}

          {page === "research" && (
            <>
              {phase === "form" && (
                <ResearchForm
                  dataSources={dataSources}
                  prompt={prompt}
                  onPromptChange={setPrompt}
                  riskTolerance={riskTolerance}
                  onRiskToleranceChange={setRiskTolerance}
                  maxBudget={maxBudget}
                  onMaxBudgetChange={setMaxBudget}
                  enabledMap={enabledMap}
                  onEnabledMapChange={setEnabledMap}
                  onSubmit={handleSubmit}
                  onNavigateConfig={() => setPage("config")}
                />
              )}

              {phase === "loading" && <LoadingScreen prompt={lastPrompt} />}

              {phase === "report" && reportMdx && (
                <ResearchReport mdxContent={reportMdx} onBack={handleBack} />
              )}

              {phase === "error" && (
                <div className="card">
                  <div className="error-container" id="error-screen">
                    <div className="error-icon">⚠️</div>
                    <div className="error-title">Research Failed</div>
                    <div className="error-message">{error}</div>
                    <button
                      className="btn-secondary"
                      onClick={handleBack}
                      id="error-back-btn"
                    >
                      ← Try Again
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  );
}
