import { NewJobPayload, sendJobInfo } from "../api/sendJobInfo";
import { useEffect, useState } from "react";

import { LlmProvider } from "@/shared/types/index.js";
import { encryptData } from "@/shared/lib/encryption";
import { useSettings } from "../../settings/hooks/useSettings";

interface FormState extends Omit<NewJobPayload, "apiKey"> {
  llmProvider: LlmProvider;
}

const initialState: FormState = {
  companyName: "",
  positionTitle: "",
  url: "",
  description: "",
  llmProvider: "Gemini",
};

/**
 * A hook that provides state and functions for handling the job posting form.
 * It handles form state, validation, and submission to the backend.
 * @returns An object containing the form state, loading state, error message, success message, settings, change handler, and submit handler.
 */
export const useJobInfoForm = () => {
  const [formState, setFormState] = useState<FormState>(initialState);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const { settings } = useSettings();
  
  useEffect(() => {
    if (settings) {
      setFormState((prev) => ({
        ...prev,
        llmProvider: settings.featureAssignments.matchSummary.provider === "None" 
          ? "Gemini" 
          : settings.featureAssignments.matchSummary.provider,
      }));
    }
  }, [settings]);

/**
 * Handles form state changes.
 * @param {React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>} e - The event containing the form element.
 * @returns void
 */
  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormState((prevState) => ({
      ...prevState,
      [name]: value,
    }));
  };

/**
 * Handles form submission by validating the form state, encrypting the API key, and sending the job posting details to the backend.
 * If the submission is successful, it resets the form state to its initial state and displays a success message.
 * If an error occurs during submission, it displays the error message.
 * @param {React.FormEvent} e - The event containing the form element.
 * @returns {Promise<void>} - A promise that resolves when the submission is complete.
 */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !formState.companyName ||
      !formState.positionTitle ||
      !formState.description
    ) {
      setError("Company Name, Position Title, and Description are required.");
      return;
    }

    setLoading(true);
    setError(null);
    setSuccessMessage(null);

    try {
      const providerKeys = settings?.apiKeys[formState.llmProvider] || [];
      const assignment = settings?.featureAssignments.matchSummary;
      const apiKey = providerKeys[assignment?.keyIndex ?? 0] || providerKeys[0];

      if (!apiKey) {
        throw new Error(
          `API key for ${formState.llmProvider} is not set in Settings.`
        );
      }

      const encryptedApiKey = await encryptData(apiKey);
      const { llmProvider, ...jobPayload } = formState;

      const message = await sendJobInfo(jobPayload, {
        llmProvider,
        encryptedApiKey,
      });

      setSuccessMessage(message);
      setFormState((prev) => ({
        ...initialState,
        llmProvider: prev.llmProvider,
      }));
    } catch (err: any) {
      setError(err.message || "An unknown error occurred during submission.");
    } finally {
      setLoading(false);
    }
  };

  return {
    formState,
    loading,
    error,
    successMessage,
    settings,
    handleChange,
    handleSubmit,
  };
};