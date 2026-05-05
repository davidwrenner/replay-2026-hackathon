import type { ResearchFormState, RiskTolerance, DataSource } from "../types";

interface ResearchFormProps {
  dataSources: DataSource[];
  prompt: string;
  onPromptChange: (value: string) => void;
  riskTolerance: RiskTolerance;
  onRiskToleranceChange: (value: RiskTolerance) => void;
  maxBudget: string;
  onMaxBudgetChange: (value: string) => void;
  enabledMap: Record<string, boolean>;
  onEnabledMapChange: (update: (prev: Record<string, boolean>) => Record<string, boolean>) => void;
  onSubmit: (state: ResearchFormState) => void;
  onNavigateConfig: () => void;
}

export default function ResearchForm({
  dataSources,
  prompt,
  onPromptChange,
  riskTolerance,
  onRiskToleranceChange,
  maxBudget,
  onMaxBudgetChange,
  enabledMap,
  onEnabledMapChange,
  onSubmit,
  onNavigateConfig,
}: ResearchFormProps) {
  const toggleSource = (name: string) => {
    onEnabledMapChange((prev) => ({ ...prev, [name]: !prev[name] }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const sources = dataSources.map((s) => ({
      ...s,
      enabled: !!enabledMap[s.name] && s.apiKey.trim().length > 0,
    }));
    onSubmit({ dataSources: sources, prompt, riskTolerance, maxBudget });
  };

  const hasAnyEnabled = dataSources.some(
    (s) => s.apiKey.trim().length > 0 && enabledMap[s.name]
  );
  const isValid = prompt.trim().length > 0 && hasAnyEnabled;

  const riskOptions: { value: RiskTolerance; label: string }[] = [
    { value: "low", label: "Conservative" },
    { value: "medium", label: "Moderate" },
    { value: "high", label: "Aggressive" },
  ];

  return (
    <form className="card" onSubmit={handleSubmit} id="research-form">
      <div className="card-header">
        <h1 className="card-title">Market Research</h1>
        <p className="card-subtitle">
          Describe your trading objective and bias. The bot will analyze markets
          and return a strategy with specific Kalshi trades.
        </p>
      </div>

      {/* Prompt */}
      <div className="form-group">
        <label className="form-label" htmlFor="prompt-input">
          Objective &amp; Bias
        </label>
        <textarea
          id="prompt-input"
          className="form-textarea"
          placeholder="e.g. I think the Fed will cut rates in September. Find election-related markets that could benefit from this macro view..."
          value={prompt}
          onChange={(e) => onPromptChange(e.target.value)}
        />
      </div>

      {/* Data Sources */}
      <div className="form-group">
        <div className="form-label-row">
          <label className="form-label" style={{ marginBottom: 0 }}>
            Data Sources
          </label>
          <button
            type="button"
            className="config-link"
            onClick={onNavigateConfig}
            id="go-to-config"
          >
            Configure API Keys →
          </button>
        </div>
        <div className="sources-grid">
          {dataSources.map((source) => {
            const hasKey = source.apiKey.trim().length > 0;
            const isEnabled = hasKey && !!enabledMap[source.name];

            return (
              <div
                key={source.name}
                className="source-chip-wrapper"
              >
                <button
                  type="button"
                  className={`source-chip ${isEnabled ? "active" : ""} ${
                    !hasKey ? "disabled" : ""
                  }`}
                  onClick={() => hasKey && toggleSource(source.name)}
                  id={`source-${source.name.toLowerCase().replace(/\s+/g, "-")}`}
                  aria-disabled={!hasKey}
                >
                  <span className="source-chip-dot" />
                  <span className="source-chip-name">{source.name}</span>
                </button>
                {!hasKey && (
                  <div className="source-tooltip">
                    No API key configured.{" "}
                    <button
                      type="button"
                      className="tooltip-link"
                      onClick={onNavigateConfig}
                    >
                      Add key
                    </button>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Risk Tolerance */}
      <div className="form-group">
        <label className="form-label">Risk Tolerance</label>
        <div className="risk-group">
          {riskOptions.map((opt) => (
            <button
              key={opt.value}
              type="button"
              className={`risk-pill ${
                riskTolerance === opt.value ? `active-${opt.value}` : ""
              }`}
              onClick={() => onRiskToleranceChange(opt.value)}
              id={`risk-${opt.value}`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      {/* Budget */}
      <div className="form-group">
        <label className="form-label">Max Budget</label>
        <div className="budget-display">
          <span className="budget-amount">
            ${Number(maxBudget).toLocaleString()}
          </span>
          <span className="budget-currency">USD</span>
        </div>
        <input
          type="range"
          id="budget-slider"
          className="budget-slider"
          min="500"
          max="100000"
          step="500"
          value={maxBudget}
          onChange={(e) => onMaxBudgetChange(e.target.value)}
        />
      </div>

      {/* Submit */}
      <button
        type="submit"
        className="btn-primary"
        disabled={!isValid}
        id="submit-research"
      >
        Run Market Research
      </button>
    </form>
  );
}
