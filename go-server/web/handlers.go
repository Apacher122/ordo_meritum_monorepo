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

type DownloadRequest struct {
	DownloadURL string `json:"download_url"`
	ChangesURL  string `json:"changes_url"`
}

func HandleDownload(
	w http.ResponseWriter,
	r *http.Request,
) {
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
	pdfPart, _ := mw.CreatePart(pdfHeader)
	io.Copy(pdfPart, pdfFile)

	jsonHeader := textproto.MIMEHeader{}
	jsonHeader.Set("Content-Disposition", `attachment; filename="`+filepath.Base(jsonPath)+`"`)
	jsonHeader.Set("Content-Type", "application/json")
	jsonPart, _ := mw.CreatePart(jsonHeader)
	io.Copy(jsonPart, jsonFile)

	mw.Close()
}
