package album

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AlbumHTTPHandler struct {
	service AlbumService
}

func NewAlbumHTTPHandler(router *mux.Router, service AlbumService) {
	handler := &AlbumHTTPHandler{
		service: service,
	}

	router.HandleFunc("/v1/albums", handler.GetMany).Methods(http.MethodGet)
}

func (hh *AlbumHTTPHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	resp := hh.service.GetMany(r.Context())
	resp.WriteJSON(w)
}
