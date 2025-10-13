import "@/assets/styles/pages/SettingsPage.css"; 
import { AssignableFeature, Settings } from "../types/types";
import React, { useEffect, useState } from "react";
import { LlmProvider } from "@/shared/types/index.js";

import { useSetHeaderTitle, useSetHeaderSubtitle, useSetHeaderControls } from "@/components/Layouts/providers/HeaderProvider";
import { useSettings } from "../hooks/useSettings";

const llmProviderOptions: LlmProvider[] = [ "Gemini", "Cohere", "OpenAI", "Groq", "Anthropic", "Cerebras" ];
const assignableFeatures: { key: AssignableFeature, label: string }[] = [
    { key: 'matchSummary', label: 'Job Match Summary' },
    { key: 'resumeGeneration', label: 'Resume Generation' },
    { key: 'coverLetterGeneration', label: 'Cover Letter Generation' },
];

export const SettingsPage: React.FC = () => {
  
  
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  const { settings, loading, error, saveSettings } = useSettings();
  const [formState, setFormState] = useState<Settings | null>(null);

  
  useEffect(() => {
    setHeaderTitle("Settings");
    setHeaderSubtitle("Manage API keys and feature configurations.");
    setHeaderControls(null); 

    
    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
    };
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls]); 


  useEffect(() => {
    if (settings) {
      setFormState(settings);
    }
  }, [settings]);

  const handleApiKeyChange = (provider: LlmProvider, value: string) => {
    setFormState(prev => prev && {
        ...prev,
        apiKeys: { ...prev.apiKeys, [provider]: value }
    });
  };

  const handleAssignmentChange = (feature: AssignableFeature, value: LlmProvider) => {
    setFormState(prev => prev && {
        ...prev,
        featureAssignments: { ...prev.featureAssignments, [feature]: value }
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (formState) {
      await saveSettings(formState);
      alert("Settings Saved!");
    }
  };

  if (loading || !formState) {
    return <div>Loading settings...</div>;
  }

  return (
    <div className="settings-page">
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit} className="info-form">
            <div className="card">
                <h2>API Keys</h2>
                <p className="page-subtitle">Your API keys are encrypted and stored locally on your machine.</p>
                <div className="settings-grid">
                    {llmProviderOptions.map(provider => (
                        <div key={provider}>
                            <label>{provider}</label>
                            <input
                                type="password"
                                value={formState.apiKeys[provider] || ''}
                                onChange={(e) => handleApiKeyChange(provider, e.target.value)}
                                placeholder={`${provider} API Key`}
                                className="input"
                            />
                        </div>
                    ))}
                </div>
            </div>

            <div className="card">
                <h2>Feature Assignments</h2>
                <p className="page-subtitle">Assign a default LLM provider to each feature.</p>
                 <div className="settings-grid">
                    {assignableFeatures.map(feature => (
                        <div key={feature.key}>
                            <label>{feature.label}</label>
                            <select
                                value={formState.featureAssignments[feature.key]}
                                onChange={(e) => handleAssignmentChange(feature.key, e.target.value as LlmProvider)}
                                className="input"
                            >
                                {llmProviderOptions.map(provider => (
                                    <option key={provider} value={provider}>{provider}</option>
                                ))}
                            </select>
                        </div>
                    ))}
                </div>
            </div>

             <button type="submit" className="submit-button" disabled={loading}>
                {loading ? "Saving..." : "Save Settings"}
            </button>
        </form>
    </div>
  );
};

