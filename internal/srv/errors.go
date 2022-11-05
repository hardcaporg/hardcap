package srv

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hardcaporg/hardcap/internal/ctxval"
)

// ResponseError implements Go standard error interface as well as Wrapper and Renderer
type ResponseError struct {
	// HTTP status code
	HTTPStatusCode int `json:"-"`

	// user facing error message
	Message string `json:"msg"`

	// trace id from context (if provided)
	TraceId string `json:"trace_id,omitempty"`

	// full root cause
	Error error `json:"error"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func writeErrorBody(w http.ResponseWriter, _ *http.Request, msg, traceId, err string) {
	_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s", "trace_id": "%s", "error": "%s"}`, msg, traceId, err)))
}

func writeBasicError(w http.ResponseWriter, r *http.Request, err error) {
	if logger := ctxval.Logger(r.Context()); logger != nil {
		logger.Error().Msgf("unable to render error %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	wrappedMessage := ""
	if errors.Unwrap(err) != nil {
		wrappedMessage = errors.Unwrap(err).Error()
	}
	traceId := ctxval.TraceId(r.Context())
	writeErrorBody(w, r, err.Error(), traceId, wrappedMessage)
}

func renderError(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	errRender := render.Render(w, r, renderer)
	if errRender != nil {
		writeBasicError(w, r, errRender)
	}
}

func NewInvalidRequestError(ctx context.Context, err error) *ResponseError {
	msg := "invalid request"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Err(err).Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}
