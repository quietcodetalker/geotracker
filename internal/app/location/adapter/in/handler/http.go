package handler

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	schemaDecoder = schema.NewDecoder()
)

// HTTPHandler serves http requests.
type HTTPHandler struct {
	service port.UserService
	router  *chi.Mux
}

// NewHTTPHandler creates HTTPHandler and returns its pointer.
func NewHTTPHandler(service port.UserService) *HTTPHandler {
	router := chi.NewRouter()

	handler := &HTTPHandler{
		service: service,
		router:  router,
	}

	handler.setupRoutes()

	return handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *HTTPHandler) setupRoutes() {
	users := chi.NewRouter()

	users.Method(http.MethodPut, "/{username}/location", http.HandlerFunc(h.setUserLocation))
	users.Method(http.MethodGet, "/radius", http.HandlerFunc(h.listUsersInRadius))

	h.router.Mount("/users", users)
}

func (h *HTTPHandler) setUserLocation(w http.ResponseWriter, r *http.Request) {
	var dto port.UserServiceSetUserLocationRequest

	if err := util.DecodeBody(r, &dto); err != nil {
		util.RespondErr(w, http.StatusOK, getHttpError(err))
		return
	}
	dto.Username = chi.URLParam(r, "username")

	res, err := h.service.SetUserLocation(r.Context(), dto)
	if err != nil {
		util.RespondErr(w, http.StatusOK, getHttpError(err))
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
		util.RespondErr(w, http.StatusOK, getHttpError(err))
		return
	}

	err = schemaDecoder.Decode(&dto, r.URL.Query())
	if err != nil {
		fmt.Println(err)
		// TODO: check whether it makes sense to handle different errors
		util.RespondErr(w, http.StatusOK, getHttpError(&port.InvalidArgumentError{}))
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
		util.RespondErr(w, http.StatusOK, getHttpError(err))
		return
	}

	util.Respond(w, http.StatusOK, res)
}

type httpError struct {
	Code    uint32     `json:"code"`
	Message string     `json:"message"`
	Status  codes.Code `json:"status"`
}

func (e httpError) Error() string {
	return e.Message
}

func getHttpError(err error) error {
	var invalidArgumentError *port.InvalidArgumentError

	if err != nil {
		switch {
		case errors.As(err, &invalidArgumentError):
			fallthrough
		case errors.Is(err, port.ErrInvalidUsername):
			return httpError{
				Status:  codes.FailedPrecondition,
				Code:    grpcCodeToHTTP(codes.FailedPrecondition),
				Message: err.Error(),
			}
		default:
			return httpError{
				Status:  codes.Internal,
				Code:    grpcCodeToHTTP(codes.Internal),
				Message: "internal error",
			}
		}
	}

	return nil
}

func grpcCodeToHTTP(code codes.Code) uint32 {
	switch code {
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Internal:
		fallthrough
	case codes.Unknown:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
