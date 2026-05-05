import { useState } from "react";
import type { DataSource } from "../types";

interface ConfigPageProps {
  dataSources: DataSource[];
  onSave: (sources: DataSource[]) => void;
}

export default function ConfigPage({ dataSources, onSave }: ConfigPageProps) {
  const [sources, setSources] = useState<DataSource[]>(dataSources);

  const updateKey = (idx: number, apiKey: string) => {
    setSources((prev) =>
      prev.map((s, i) => (i === idx ? { ...s, apiKey } : s))
    );
  };

  const handleSave = () => {
    onSave(sources);
  };

  const configuredCount = sources.filter((s) => s.apiKey.trim().length > 0).length;

  return (
    <div className="card" id="config-page">
      <div className="card-header">
        <h1 className="card-title">Data Source Configuration</h1>
        <p className="card-subtitle">
          Enter API keys for each data source you want the bot to access during
          market research. Sources without keys will be unavailable on the
          research page.
        </p>
      </div>

      <div className="config-status">
        <span className="config-status-count">{configuredCount}</span>
        <span className="config-status-label">
          of {sources.length} sources configured
        </span>
      </div>

      <div className="config-list">
        {sources.map((source, idx) => {
          const hasKey = source.apiKey.trim().length > 0;
          return (
            <div
              key={source.name}
              className={`config-item ${hasKey ? "configured" : ""}`}
              id={`config-${source.name.toLowerCase().replace(/\s+/g, "-")}`}
            >
              <div className="config-item-header">
                <div className="config-item-info">
                  <span
                    className={`config-item-dot ${hasKey ? "active" : ""}`}
                  />
                  <div>
                    <span className="config-item-name">{source.name}</span>
                    <span className="config-item-desc">
                      {source.description}
                    </span>
                  </div>
                </div>
                {hasKey && <span className="config-item-badge">Connected</span>}
              </div>
              <input
                type="password"
                className="form-input config-key-input"
                placeholder={`Enter ${source.name} API key…`}
                value={source.apiKey}
                onChange={(e) => updateKey(idx, e.target.value)}
                id={`apikey-${source.name.toLowerCase().replace(/\s+/g, "-")}`}
              />
            </div>
          );
        })}
      </div>

      <button
        type="button"
        className="btn-primary"
        onClick={handleSave}
        id="save-config"
      >
        Save Configuration
      </button>
    </div>
  );
}
