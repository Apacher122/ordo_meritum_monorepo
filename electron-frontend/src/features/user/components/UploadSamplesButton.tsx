import React from "react";

/**
 * A button component that triggers the file system to upload writing samples.
 * It handles the full upload and save flow, showing alerts on success or failure.
 * @returns {React.FC}
 */
export const UploadSamplesButton: React.FC = () => {
  const handleUpload = async () => {
    const result = await window.appAPI.writingSamples.upload();
    if (result.success && result.samples) {
      const saveResult = await window.appAPI.writingSamples.save(result.samples);
      if (saveResult.success) {
        alert(`${result.samples.length} writing sample(s) saved successfully!`);
      } else {
        alert(`Error saving samples: ${saveResult.error}`);
      }
    } else if (result.error) {
      alert(`Error opening files: ${result.error}`);
    }
  };

  return (
    <button onClick={handleUpload} className="button">
      Upload Writing Samples
    </button>
  );
};