package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Karthik-Kondaveeti/blink/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) (*Handler, error) {
	return &Handler{
		service: service,
	}, nil
}

type LinkRequest struct {
	URL string `json:"url"`
}

type LinkResponse struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) AddLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var req LinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link := req.URL
	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortCode, err := h.service.AddLink(ctx, link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(LinkResponse{ShortCode: shortCode})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.PathValue("shortCode")

	originalURL, err := h.service.GetLink(ctx, shortCode)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
