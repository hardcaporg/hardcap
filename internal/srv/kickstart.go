package srv

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/hardcaporg/hardcap/internal/config"
	"github.com/hardcaporg/hardcap/internal/ctxval"
	"github.com/hardcaporg/hardcap/internal/srv/snip"
	"github.com/hardcaporg/hardcap/internal/version"
)

const ksTemplate = `
# Generated by Hardcap {{ .BuildCommit }} {{ .BuildTime }}

%pre --interpreter=/bin/python3 --erroronfail --log=/mnt/sysimage/root/ks-pre.log
{{ .Pre }}
%end

%include /tmp/pre-generated.ks
`

type KickstartVars struct {
	BuildCommit string
	BuildTime   string
	Pre         string
	Address     string
}

func KickstartTemplateService(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())

    for name, values := range r.Header {
        // Loop over all values for the name.
        for _, value := range values {
            logger.Debug().Msgf("%s: %s", name, value)
        }
    }
    
	vars := KickstartVars{
		BuildCommit: version.BuildCommit,
		BuildTime:   version.BuildTime,
		Address:     config.Application.AdvertisedAddress,
	}

	preTemplate, err := snip.EmbedFS.ReadFile("pre.py")
	if err != nil {
		logger.Error().Err(err).Msg("Unable to read pre template")
	}

	t := template.Must(template.New("pre").Parse(string(preTemplate)))
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, vars)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to render kickstart template")
	}
	vars.Pre = buf.String()

	t = template.Must(template.New("ks").Parse(ksTemplate))
	buf = bytes.NewBuffer(nil)
	err = t.Execute(buf, vars)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to render kickstart template")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
