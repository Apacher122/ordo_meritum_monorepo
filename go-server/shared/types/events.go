package types

type DocumentCompletionEvent struct {
	UserID      string `json:"user_id"`
	JobID       string `json:"job_id"`
	Status      string `json:"status"`
	DownloadURL string `json:"download_url,omitempty"`
	Error       string `json:"error,omitempty"`
}
