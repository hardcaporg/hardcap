package srv

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hardcaporg/hardcap/internal/ctxval"
)

type RegisterHostPayload struct {
	MacAddresses []string `json:"mac"`
	CPU          struct {
		Count int `json:"count"`
	} `json:"cpu"`
	DMI struct {
		SystemSerialNumber int `json:"system-serial-number"`
	} `json:"dmi"`
}

func (p *RegisterHostPayload) Bind(_ *http.Request) error {
	return nil
}

func RegisterHostService(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())

	payload := &RegisterHostPayload{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, NewInvalidRequestError(r.Context(), err))
		return
	}

	logger.Debug().Msgf("Host registered: %+v", payload)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
