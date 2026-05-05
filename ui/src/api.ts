import type { ResearchFormState, ResearchRequest } from "./types";

const API_BASE = import.meta.env.VITE_API_BASE ?? "";

export async function submitResearch(form: ResearchFormState): Promise<string> {
  const body: ResearchRequest = {
    dataSources: form.dataSources
      .filter((s) => s.enabled)
      .map(({ name, description, apiKey }) => ({ name, description, apiKey })),
    prompt: form.prompt,
    riskTolerance: form.riskTolerance,
    maxBudget: form.maxBudget,
  };

  const res = await fetch(`${API_BASE}/research`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `Request failed with status ${res.status}`);
  }

  // The backend returns { research: "<base64 encoded MDX>" }
  const json = await res.json();
  const mdx = atob(json.research);
  return mdx;
}
