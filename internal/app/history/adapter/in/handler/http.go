package handler

import (
	log2 "log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/schema"
	"gitlab.com/spacewalker/geotracker/internal/app/history/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
	"gitlab.com/spacewalker/geotracker/internal/pkg/middleware"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
)

var (
	schemaDecoder = schema.NewDecoder()
)

// HTTPHandler is a handler that serves http requests.
type HTTPHandler struct {
	service port.HistoryService
	router  *chi.Mux
	logger  log.Logger
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHTTPHandler(service port.HistoryService, logger log.Logger) *HTTPHandler {
	if logger == nil {
		log2.Panic("logger must not be nil")
	}
	if service == nil {
		logger.Panic("service must not be nil", nil)
	}

	router := chi.NewRouter()

	handler := &HTTPHandler{
		service: service,
		router:  router,
		logger:  logger,
	}

	handler.setupRoutes()

	return handler
}

func (h *HTTPHandler) setupRoutes() {
	h.router.Use(
		middleware.LoggerMiddleware(h.logger),
		middleware.RecovererMiddleware(h.logger),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
		middleware2.AllowContentType("application/json"),
		middleware2.SetHeader("Content-Type", "application/json"),
	)

	users := chi.NewRouter()

	users.Method(http.MethodGet, "/{username}/distance", http.HandlerFunc(h.getDistance))

	h.router.Mount("/users", users)
}

type getDistanceDTO struct {
	From string `schema:"from"`
	To   string `schema:"to"`
}

func (h *HTTPHandler) getDistance(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	err := r.ParseForm()
	if err != nil {
		status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
		util.Respond(w, status, body)
		return
	}

	var dto getDistanceDTO
	err = schemaDecoder.Decode(&dto, r.URL.Query())
	if err != nil {
		status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
		util.Respond(w, status, body)
		return
	}

	var fromPtr, toPtr *time.Time

	if dto.From != "" {
		from, err := time.Parse(time.RFC3339, dto.From)
		if err != nil {
			status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
			util.Respond(w, status, body)
			return
		}
		fromPtr = &from
	}

	if dto.To != "" {
		to, err := time.Parse(time.RFC3339, dto.To)
		if err != nil {
			status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
			util.Respond(w, status, body)
			return
		}
		toPtr = &to
	}

	var res port.HistoryServiceGetDistanceByUsernameResponse
	res, err = h.service.GetDistanceByUsername(r.Context(), port.HistoryServiceGetDistanceByUsernameRequest{
		Username: username,
		From:     fromPtr,
		To:       toPtr,
	})
	if err != nil {
		status, body := errpack.ErrToHTTP(err)
		util.Respond(w, status, body)
		return
	}

	util.Respond(w, http.StatusOK, res)
}
