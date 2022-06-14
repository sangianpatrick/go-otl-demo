package album

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/go-otl-demo/response"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

type AlbumHTTPHandler struct {
	service AlbumService
}

var counter syncint64.Counter

func NewAlbumHTTPHandler(router *mux.Router, service AlbumService) {
	handler := &AlbumHTTPHandler{
		service: service,
	}

	meter := global.MeterProvider().Meter("visitor.albums")
	counter, _ = meter.SyncInt64().Counter(
		"album.view.counter",
		instrument.WithUnit("1"),
		instrument.WithDescription("Just a test counter"),
	)

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
	ctx := r.Context()
	resp := hh.service.GetMany(ctx)
	if resp.Error() == nil {
		counter.Add(ctx, 1, attribute.String("type", "list"))
	}

	resp.WriteJSON(w)
}
