package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/middlewares"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	log2 "log"
	"net/http"
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
	if service == nil {
		log2.Panic("service must not be nil")
	}
	if logger == nil {
		log2.Panic("logger must not be nil")
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
	h.router.Use(middlewares.LoggerMiddleware(h.logger))
	h.router.Use(middlewares.RecovererMiddleware(h.logger))

	users := chi.NewRouter()

	users.Method(http.MethodPut, "/{username}/location", http.HandlerFunc(h.setUserLocation))
	users.Method(http.MethodGet, "/radius", http.HandlerFunc(h.listUsersInRadius))

	h.router.Mount("/users", users)
}

func (h *HTTPHandler) setUserLocation(w http.ResponseWriter, r *http.Request) {
	var dto port.UserServiceSetUserLocationRequest

	if err := util.DecodeBody(r, &dto); err != nil {
		status, body := errpack.ErrToHTTP(fmt.Errorf("%w: %v", errpack.ErrInternalError, err))
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
		status, body := errpack.ErrToHTTP(fmt.Errorf("%w: %v", errpack.ErrInternalError, err))
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
