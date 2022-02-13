package handler

import (
	"fmt"
	log2 "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/schema"
	"gitlab.com/spacewalker/geotracker/internal/app/location/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
	"gitlab.com/spacewalker/geotracker/internal/pkg/middleware"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
)

var (
	schemaDecoder = schema.NewDecoder()
)

// HTTPHandler serves http requests.
type HTTPHandler struct {
	service port.UserService
	router  *chi.Mux
	logger  log.Logger
}

// NewHTTPHandler creates HTTPHandler and returns its pointer.
func NewHTTPHandler(service port.UserService, logger log.Logger) *HTTPHandler {
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

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
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

	users.Method(http.MethodPut, "/{username}/location", http.HandlerFunc(h.setUserLocation))
	users.Method(http.MethodGet, "/radius", http.HandlerFunc(h.listUsersInRadius))

	h.router.Mount("/users", users)
}

func (h *HTTPHandler) setUserLocation(w http.ResponseWriter, r *http.Request) {
	var dto port.UserServiceSetUserLocationRequest

	if err := util.DecodeBody(r, &dto); err != nil {
		status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
		util.Respond(w, status, body)
		return
	}
	dto.Username = chi.URLParam(r, "username")

	res, err := h.service.SetUserLocation(r.Context(), dto)
	if err != nil {
		status, body := errpack.ErrToHTTP(err)
		util.Respond(w, status, body)
		return
	}

	util.Respond(w, http.StatusOK, res)
}

type listUsersInRadiusDTO struct {
	Radius    float64 `schema:"radius"`
	Longitude float64 `schema:"longitude"`
	Latitude  float64 `schema:"latitude"`
	PageToken string  `schema:"page_token"`
	PageSize  int     `schema:"page_size"`
}

func (h *HTTPHandler) listUsersInRadius(w http.ResponseWriter, r *http.Request) {
	var dto listUsersInRadiusDTO
	var res port.UserServiceListUsersInRadiusResponse

	err := r.ParseForm()
	if err != nil {
		status, body := errpack.ErrToHTTP(errpack.ErrInvalidArgument)
		util.Respond(w, status, body)
		return
	}

	err = schemaDecoder.Decode(&dto, r.URL.Query())
	if err != nil {
		fmt.Println(err)
		// TODO: check whether it makes sense to handle different errors
		status, body := errpack.ErrToHTTP(fmt.Errorf("%w", errpack.ErrInvalidArgument))
		util.Respond(w, status, body)
		return
	}

	req := port.UserServiceListUsersInRadiusRequest{
		Point: geo.Point{
			dto.Longitude,
			dto.Latitude,
		},
		Radius:    dto.Radius,
		PageToken: dto.PageToken,
		PageSize:  dto.PageSize,
	}

	res, err = h.service.ListUsersInRadius(r.Context(), req)
	if err != nil {
		status, body := errpack.ErrToHTTP(err)
		util.Respond(w, status, body)
		return
	}

	util.Respond(w, http.StatusOK, res)
}
