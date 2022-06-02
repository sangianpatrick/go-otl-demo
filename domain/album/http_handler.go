package album

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/go-otl-demo/response"
)

type AlbumHTTPHandler struct {
	service AlbumService
}

func NewAlbumHTTPHandler(router *mux.Router, service AlbumService) {
	handler := &AlbumHTTPHandler{
		service: service,
	}

	router.HandleFunc("/v1/albums", handler.Add).Methods(http.MethodPost)
	router.HandleFunc("/v1/albums", handler.GetMany).Methods(http.MethodGet)
}

func (hh *AlbumHTTPHandler) Add(w http.ResponseWriter, r *http.Request) {
	var params CreateAlbumParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		resp := response.ResponseError(response.StatusBadRequest, err, nil, nil, "")
		resp.WriteJSON(w)
		return
	}

	resp := hh.service.Add(r.Context(), params)
	resp.WriteJSON(w)
}

func (hh *AlbumHTTPHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	resp := hh.service.GetMany(r.Context())
	resp.WriteJSON(w)
}
