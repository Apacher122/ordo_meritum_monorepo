package web

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/ordo_meritum/config"
	"github.com/ordo_meritum/shared/middleware"
	ordows "github.com/ordo_meritum/websocket"
	"github.com/rs/zerolog/log"
)

type DownloadRequest struct {
	DownloadURL string `json:"download_url"`
	ChangesURL  string `json:"changes_url"`
}

// ServeWs handles WebSocket upgrade requests.
// It verifies the Firebase ID token provided in the "token" query parameter.
// If the token is valid, it upgrades the HTTP connection to a WebSocket connection
// using the gorilla/websocket library and registers the new client with the provided Hub.
// It then starts the client's read and write pumps as separate goroutines.
// If token verification fails or the upgrade fails, it responds with an appropriate
// HTTP error status.
func ServeWs(
	hub *ordows.Hub,
	w http.ResponseWriter,
	r *http.Request,
) {
	authClient, err := config.AuthClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get auth client")
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get auth client"})
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		log.Warn().Msg("Unauthorized access attempt to websocket")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: No token provided"})
		return
	}

	verifiedToken, err := authClient.VerifyIDToken(context.Background(), tokenStr)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid WebSocket token")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Invalid token"})
		return
	}
	userID := verifiedToken.UID

	upgrader := websocket.Upgrader{
		// Allow all origins for WebSocket connections.
		// Consider restricting this in production for security.
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}

	client := &ordows.Client{Hub: hub, UserID: userID, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

// HandleDownload handles requests to download both a PDF document and its associated JSON changes file.
// It expects a JSON request body conforming to the DownloadRequest struct, containing URLs for both files.
// It reads both files from the local filesystem (assuming URLs map to local paths after stripping the prefix)
// and streams them back to the client as a multipart/form-data response with appropriate Content-Disposition headers
// for attachment downloads.
// Returns HTTP errors if the request body is invalid or if either file cannot be opened.
func HandleDownload(
	w http.ResponseWriter,
	r *http.Request,
) {
	defer r.Body.Close()

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pdfPath := strings.TrimPrefix(req.DownloadURL, "http://localhost:8080/")
	jsonPath := strings.TrimPrefix(req.ChangesURL, "http://localhost:8080/")

	pdfPath = filepath.Join("./shared_pdfs", filepath.Base(pdfPath))
	jsonPath = filepath.Join("./shared_pdfs", filepath.Base(jsonPath))

	pdfFile, err := os.Open(req.DownloadURL)
	if err != nil {
		http.Error(w, "Failed to open PDF file", http.StatusNotFound)
		return
	}
	defer pdfFile.Close()

	jsonFile, err := os.Open(req.ChangesURL)
	if err != nil {
		http.Error(w, "Failed to open JSON file", http.StatusNotFound)
	}
	defer jsonFile.Close()

	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Disposition", "attachment; filename="+mw.Boundary())
	w.WriteHeader(http.StatusOK)

	pdfHeader := textproto.MIMEHeader{}
	pdfHeader.Set("Content-Disposition", `attachment; filename="`+filepath.Base(pdfPath)+`"`)
	pdfHeader.Set("Content-Type", "application/pdf")
	pdfPart, err := mw.CreatePart(pdfHeader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create PDF multipart part")
		return
	}
	if _, err := io.Copy(pdfPart, pdfFile); err != nil {
		log.Error().Err(err).Msg("Failed to copy PDF content to multipart")
		return
	}

	jsonHeader := textproto.MIMEHeader{}
	jsonHeader.Set("Content-Disposition", `attachment; filename="`+filepath.Base(jsonPath)+`"`)
	jsonHeader.Set("Content-Type", "application/json")
	jsonPart, err := mw.CreatePart(jsonHeader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create JSON multipart part")
		return
	}
	if _, err := io.Copy(jsonPart, jsonFile); err != nil {
		log.Error().Err(err).Msg("Failed to copy JSON content to multipart")
		return
	}

	if err := mw.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close multipart writer")
	}
}
