import "@/assets/styles/pages/SettingsPage.css";

import { ApiKeyConfig, AssignableFeature, RateLimitConfig, Settings } from "../types/types";
import React, { useEffect, useState } from "react";
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { LlmProvider } from "@/shared/types/index.js";
import { useSettings } from "../hooks/useSettings";

const llmProviderOptions: LlmProvider[] = [ "Gemini", "Cohere", "OpenAI", "Groq", "Anthropic", "Cerebras" ];
const assignableFeatures: { key: AssignableFeature, label: string }[] = [
    { key: 'matchSummary', label: 'Job Match Summary' },
    { key: 'resumeGeneration', label: 'Resume Generation' },
    { key: 'coverLetterGeneration', label: 'Cover Letter Generation' },
];

/**
 * A settings page for managing API keys and feature configurations.
 * @returns {JSX.Element} A JSX element representing the settings page.
 */
export const SettingsPage: React.FC = () => {
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();

  const { settings, loading, error, saveSettings, resetSettings } = useSettings();
  const [formState, setFormState] = useState<Settings | null>(null);
  const [selectedProvider, setSelectedProvider] = useState<LlmProvider | "None">("None");
  
  useEffect(() => {
    setHeaderTitle("Settings");
    setHeaderSubtitle("Manage API keys and feature configurations.");
    setHeaderControls(null); 

    return () => {
      setHeaderTitle("No Job Selected");
      setHeaderSubtitle("Select or analyze a job to begin");
    };
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls]); 

  const getDefaultRateLimit = (): RateLimitConfig => ({
    callsPerDay: 1000,
    callsPerMinute: 60,
    currentDayCount: 0,
    currentMinuteCount: 0,
    lastDayReset: Date.now(),
    lastMinuteReset: Date.now(),
  });

  useEffect(() => {
    if (settings) {
      const sanitizedApiKeys: Partial<Record<LlmProvider, ApiKeyConfig[]>> = {};
      
      llmProviderOptions.forEach(provider => {
        const value = settings.apiKeys[provider];
        if (Array.isArray(value)) {
          sanitizedApiKeys[provider] = value;
        } else {
          sanitizedApiKeys[provider] = [{ key: "", rateLimit: getDefaultRateLimit() }]; 
        }
      });

      setFormState({
        ...settings,
        apiKeys: sanitizedApiKeys
      });
    }
  }, [settings]);

/**
 * Handles changes to the API key associated with a given LLM provider.
 * @param {LlmProvider} provider - The LLM provider whose API key is being changed.
 * @param {number} index - The index of the API key being changed.
 * @param {string} value - The new value of the API key.
 */
  const handleApiKeyChange = (provider: LlmProvider, index: number, field: 'key' | keyof RateLimitConfig, value: any) => {
    setFormState(prev => {
        if (!prev) return null;
        const rawKeys = prev.apiKeys[provider];
        const currentKeys = Array.isArray(rawKeys) ? [...rawKeys] : [{ key: "", rateLimit: getDefaultRateLimit() }];
        
        if (!currentKeys[index]) {
            currentKeys[index] = { key: "", rateLimit: getDefaultRateLimit() };
        }

        if (field === 'key') {
            currentKeys[index] = { ...currentKeys[index], key: value };
        } else {
             const limitField = field;
             const numValue = Number.parseInt(value, 10) || 0;
             currentKeys[index] = {
                 ...currentKeys[index],
                 rateLimit: {
                     ...currentKeys[index].rateLimit,
                     [limitField]: numValue
                 }
             };
        }

        return {
            ...prev,
            apiKeys: { ...prev.apiKeys, [provider]: currentKeys }
        };
    });
  };

  const handleAddNewKey = (provider: LlmProvider) => {
    setFormState(prev => {
        if (!prev) return null;
        const rawKeys = prev.apiKeys[provider];
        const currentKeys = Array.isArray(rawKeys) ? [...rawKeys] : [];
        return {
            ...prev,
            apiKeys: { 
                ...prev.apiKeys, 
                [provider]: [...currentKeys, { key: "", rateLimit: getDefaultRateLimit() }] 
            }
        };
    });
  };

/**
 * Handles changes to the assignment of a given feature.
 * @param {AssignableFeature} feature - The feature whose assignment is being changed.
 * @param {'provider' | 'keyIndex'} field - The field of the feature assignment being changed.
 * @param {*} value - The new value of the feature assignment field.
 */
  const handleAssignmentChange = (feature: AssignableFeature, field: 'provider' | 'keyIndex', value: any) => {
    setFormState(prev => {
        if (!prev) return null;
        return {
            ...prev,
            featureAssignments: { 
                ...prev.featureAssignments, 
                [feature]: { 
                    ...prev.featureAssignments[feature], 
                    [field]: field === 'keyIndex' ? (Number.parseInt(value, 10) || 0) : value 
                } 
            }
        };
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (formState) {
      const plainSettings = structuredClone(formState);
      await saveSettings(plainSettings);
      alert("Settings Saved!");
    }
  };

  const handleReset = async () => {
    if (window.confirm("Are you sure you want to reset all settings? All keys will be removed.")) {
        await resetSettings();
        setSelectedProvider("None");
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
					<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
						<h2>API Keys</h2>
						<button 
							type="button" 
							onClick={handleReset} 
							className="secondary-button" 
							style={{ color: 'var(--error-color)', borderColor: 'var(--error-color)', fontSize: '0.8rem' }}
						>
							Reset to Defaults
						</button>
					</div>
					<p className="page-subtitle">Your API keys are encrypted and stored locally on your machine.</p>
					
					<div className="settings-grid" style={{ marginBottom: '1rem' }}>
						<div>
							<label htmlFor="provider-selector">LLM Provider</label>
							<select
									id="provider-selector"
									value={selectedProvider}
									onChange={(e) => setSelectedProvider(e.target.value as LlmProvider | "None")}
									className="input"
							>
									<option value="None">None</option>
									{llmProviderOptions.map(provider => (
											<option key={provider} value={provider}>{provider}</option>
									))}
							</select>
						</div>
					</div>
					{selectedProvider === "None" ? (
						<p className="hint-text" style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', fontStyle: 'italic' }}>
							Please select a provider
						</p>
					) : (
						<div className="provider-keys-section" style={{ borderTop: '1px solid var(--border-color)', paddingTop: '1rem' }}>
							<div className="settings-grid">
								{(Array.isArray(formState.apiKeys[selectedProvider]) 
									? formState.apiKeys[selectedProvider]
									: []
								).map((keyConfig, index) => (
									<div key={`${selectedProvider}-key-${index}`} style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', marginBottom: '1rem', borderBottom: '1px dashed #444', paddingBottom: '1rem' }}>
										<label>Key {index + 1}</label>
										<input
											type="password"
											value={keyConfig.key}
											onChange={(e) => handleApiKeyChange(selectedProvider, index, 'key', e.target.value)}
											placeholder={`Enter ${selectedProvider} Key ${index + 1}`}
											className="input"
                    />
                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <div style={{ flex: 1 }}>
                            <label style={{ fontSize: '0.8rem' }}>Daily Limit</label>
                            <input
                                type="number"
                                value={keyConfig.rateLimit.callsPerDay}
                                onChange={(e) => handleApiKeyChange(selectedProvider, index, 'callsPerDay', e.target.value)}
                                className="input"
                            />
                        </div>
                        <div style={{ flex: 1 }}>
                            <label style={{ fontSize: '0.8rem' }}>Minute Limit</label>
                            <input
                                type="number"
                                value={keyConfig.rateLimit.callsPerMinute}
                                onChange={(e) => handleApiKeyChange(selectedProvider, index, 'callsPerMinute', e.target.value)}
                                className="input"
                            />
                        </div>
                    </div>
                    <div style={{ fontSize: '0.75rem', color: '#888' }}>
                        Usage: {keyConfig.rateLimit.currentDayCount} today / {keyConfig.rateLimit.currentMinuteCount} this minute
                    </div>
									</div>
								))}
							</div>
							<button 
								type="button" 
								className="secondary-button" 
								style={{ marginTop: '1rem' }}
								onClick={() => handleAddNewKey(selectedProvider)}
							>
								+ Add New Key
							</button>
						</div>
					)}
				</div>

				<div className="card">
					<h2>Feature Assignments</h2>
					<p className="page-subtitle">Assign a default LLM provider and specific key to each feature.</p>
					<div className="settings-grid">
					{assignableFeatures.map(feature => {
						const assignment = formState.featureAssignments[feature.key];
						const rawKeys = assignment.provider === "None"
							? []
							: (formState.apiKeys[assignment.provider] || []);
						
						const fallbackKeys = rawKeys && !Array.isArray(rawKeys) ? [rawKeys] : [];
						const availableKeys = (Array.isArray(rawKeys) ? rawKeys : fallbackKeys);

						return (
							<div key={feature.key}>
								<label htmlFor={`${feature.key}-provider-select`}>{feature.label}</label>
								<div style={{ display: 'flex', gap: '10px' }}>
									<select
										id={`${feature.key}-provider-select`}
										value={assignment.provider}
										onChange={(e) => handleAssignmentChange(feature.key, 'provider', e.target.value)}
										className="input"
										style={{ flex: 2 }}
									>
										<option value="None">None</option>
										{llmProviderOptions.map(provider => (
											<option key={provider} value={provider}>{provider}</option>
										))}
									</select>

									<select
										disabled={assignment.provider === "None"}
										aria-label={`${feature.label} Key Selection`}
										value={assignment.keyIndex}
										onChange={(e) => handleAssignmentChange(feature.key, 'keyIndex', e.target.value)}
										className="input"
										style={{ flex: 1 }}
									>
										{availableKeys.length > 0 ? (
											availableKeys.map((_, idx) => (
												<option key={`${assignment.provider}-key-option-${idx}`} value={idx}>
													Key {idx + 1}
												</option>
											))
										) : (
											<option value={0}>No Keys</option>
										)}
									</select>
								</div>
							</div>
						);
					})}
					</div>
				</div>
				<button type="submit" className="submit-button" disabled={loading}>
					{loading ? "Saving..." : "Save Settings"}
				</button>
			</form>
    </div>
  );
};