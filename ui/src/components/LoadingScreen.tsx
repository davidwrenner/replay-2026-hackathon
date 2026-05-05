import { useState, useEffect } from "react";

interface LoadingScreenProps {
  prompt: string;
}

const STEPS = [
  "Connecting to data sources…",
  "Scanning Kalshi markets…",
  "Analyzing sentiment signals…",
  "Building trade candidates…",
  "Generating strategy report…",
];

export default function LoadingScreen({ prompt }: LoadingScreenProps) {
  const [activeStep, setActiveStep] = useState(0);

  useEffect(() => {
    let step = 0;
    const interval = setInterval(() => {
      step++;
      if (step < STEPS.length) {
        setActiveStep(step);
      } else {
        clearInterval(interval);
      }
    }, 550);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="card">
      <div className="loading-container" id="loading-screen">
        <div className="loading-orb">
          <div className="loading-orb-ring" />
          <div className="loading-orb-ring ring-2" />
          <div className="loading-orb-core">K</div>
        </div>
        <div className="loading-text">
          <strong>Running Market Research</strong>
          <span>{STEPS[activeStep]}</span>
        </div>
        <div className="loading-steps">
          {STEPS.map((step, idx) => (
            <div
              key={step}
              className={`loading-step ${
                idx < activeStep ? "done" : idx === activeStep ? "active" : ""
              }`}
            >
              <span className="loading-step-dot" />
              {step}
            </div>
          ))}
        </div>
        <p className="loading-prompt-preview">
          "{prompt.length > 120 ? prompt.slice(0, 120) + "…" : prompt}"
        </p>
      </div>
    </div>
  );
}
