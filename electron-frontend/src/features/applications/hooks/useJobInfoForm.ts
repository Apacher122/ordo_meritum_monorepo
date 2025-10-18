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
        llmProvider: settings.featureAssignments.matchSummary,
      }));
    }
  }, [settings]);

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
      const apiKey = settings?.apiKeys[formState.llmProvider];
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
