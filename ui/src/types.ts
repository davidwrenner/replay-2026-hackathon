export interface DataSource {
  name: string;
  description: string;
  apiKey: string;
  enabled: boolean;
}

export type RiskTolerance = "low" | "medium" | "high";

export interface ResearchRequest {
  dataSources: Omit<DataSource, "enabled">[];
  prompt: string;
  riskTolerance: RiskTolerance;
  maxBudget: string;
}

export interface ResearchFormState {
  dataSources: DataSource[];
  prompt: string;
  riskTolerance: RiskTolerance;
  maxBudget: string;
}

export type AppPage = "research" | "config";

export type AppPhase = "form" | "loading" | "report" | "error";

export interface AppState {
  phase: AppPhase;
  formState: ResearchFormState;
  reportMdx: string | null;
  error: string | null;
}
