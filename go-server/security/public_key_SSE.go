package security

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func PublicKeyStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	publicKeyStr := os.Getenv("PUBLIC_KEY")
	if publicKeyStr == "" {
		log.Println("SSE Error: PUBLIC_KEY environment variable not set.")
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", publicKeyStr)
	flusher.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			log.Println("SSE client disconnected.")
			return
		case <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", publicKeyStr)
			flusher.Flush()
		}
	}
}
