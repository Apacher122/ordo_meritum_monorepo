import "@/assets/styles/pages/InfoPage.css";

import React, { useEffect } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { Link } from "react-router-dom";
import { LlmProvider } from "@/shared/types/index.js";
import { useJobInfoForm } from "../hooks/useJobInfoForm";

const llmProviderOptions: LlmProvider[] = [
  "Gemini",
  "Cohere",
  "OpenAI",
  "Groq",
  "Anthropic",
  "Cerebras",
];

export const InfoPage: React.FC = () => {
  
  
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  const {
    formState,
    loading,
    error,
    successMessage,
    settings,    
    handleChange,
    handleSubmit,
  } = useJobInfoForm();

  
  useEffect(() => {
    setHeaderTitle("Analyze a New Job");
    setHeaderSubtitle("Enter the job details below to begin the analysis.");
    setHeaderControls(null); 

    
    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
    };
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls]); 


    const isApiKeyMissing = !settings?.apiKeys[formState.llmProvider];
  const isSubmitDisabled = loading || isApiKeyMissing;

  return (
    <div className="info-page">
      <form onSubmit={handleSubmit} className="info-form">
        {error && <div className="status-message error-message">{error}</div>}
        {successMessage && (
          <div className="status-message success-message">{successMessage}</div>
        )}

        <div className="form-grid">
          <input
            type="text"
            name="companyName"
            value={formState.companyName}
            onChange={handleChange}
            placeholder="Company Name"
            className="input"
            disabled={loading}
            required
          />
          <input
            type="text"
            name="positionTitle"
            value={formState.positionTitle}
            onChange={handleChange}
            placeholder="Position Title"
            className="input"
            disabled={loading}
            required
          />
        </div>

        <input
          type="url"
          name="url"
          value={formState.url}
          onChange={handleChange}
          placeholder="URL (optional)"
          className="input"
          disabled={loading}
        />

        <div className="form-grid single-column">
          <label>LLM Provider</label>
          <select
            name="llmProvider"
            value={formState.llmProvider}
            onChange={handleChange}
            className="input"
            disabled={loading}
          >
            {llmProviderOptions.map((provider) => (
              <option key={provider} value={provider}>
                {provider}
              </option>
            ))}
          </select>
        </div>

        <textarea
          name="description"
          value={formState.description}
          onChange={handleChange}
          placeholder="Paste job description here..."
          className="textarea job-textarea"
          disabled={loading}
          required
        />

        <div className="submit-wrapper">
          <button
            type="submit"
            className="submit-button"
            disabled={isSubmitDisabled}
          >
            {loading ? "Analyzing..." : "Analyze Job"}
          </button>
          {isApiKeyMissing && (
            <div className="api-key-warning">
              API key for {formState.llmProvider} is missing. Please add it in{" "}
              <Link to="/settings">Settings</Link>.
            </div>
          )}
        </div>
      </form>
    </div>
  );
};
