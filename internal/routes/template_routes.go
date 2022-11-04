package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/hardcaporg/hardcap/internal/tsrv"
)

func MountTemplateEndpoint(r *chi.Mux) {
	r.Route("/ks", func(r chi.Router) {
		r.Get("/", tsrv.KickstartTemplateService)
	})
}
