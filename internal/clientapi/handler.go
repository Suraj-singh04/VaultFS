package clientapi

import (
	"io"
	"log"
	"net/http"
)

type Handler struct {
	coordinator *Coordinator
}

func NewHandler(coordinator *Coordinator) *Handler {
	return &Handler{
		coordinator: coordinator,
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing file path", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read file data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.coordinator.UploadFile(filePath, data)
	err = h.coordinator.UploadFile(filePath, data)
	if err != nil {
		log.Printf("Upload error: %v", err)
		http.Error(w, "Failed to upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing file path", http.StatusBadRequest)
		return
	}

	data, err := h.coordinator.FetchFile(filePath)
	if err != nil {
		http.Error(w, "Failed to fetch file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filePath+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
