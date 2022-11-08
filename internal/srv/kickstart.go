package srv

import (
	"bytes"
    "github.com/hardcaporg/hardcap/internal/model"
    "net/http"
	"strings"
	"text/template"

	"github.com/hardcaporg/hardcap/internal/config"
	"github.com/hardcaporg/hardcap/internal/ctxval"
	"github.com/hardcaporg/hardcap/internal/srv/snip"
	"github.com/hardcaporg/hardcap/internal/version"
)

const templateRegister = `
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

    var serial string
    var macs []string
    if serials, ok := r.Header["X-System-Serial-Number"]; ok && len(serials) == 1 {
        serial = serials[0]
	} else {
        logger.Warn().Msg("Kickstart request without X-System-Serial-Number, use inst.ks.sendsn kernel option (VM perhaps)")
    }
    seenMAC := false
	for name, values := range r.Header {
		for _, value := range values {
            logger.Trace().Msgf("HTTP header '%s': '%s'", name, value)
            if strings.HasPrefix(name, "X-Rhn-Provisioning-Mac") {
				nameAndMAC := strings.Split(value, " ")
				if len(nameAndMAC) == 2 {
					macs = append(macs, nameAndMAC[1])
                    seenMAC = true
				} else {
                    logger.Warn().Msgf("X-RHN-Provisioning-MAC unexpected format: '%s'", nameAndMAC)
                }
			}
		}
	}
    if !seenMAC {
        logger.Warn().Msg("Kickstart request without X-RHN-Provisioning-MAC, use inst.ks.sendmac kernel option")
    }
    
    logger.Trace().Msgf("Generating ID: %+v %+v", serial, macs)
    sid := model.NewSystemID(serial, macs...)
    logger.Trace().Msgf("System ID: %s %s %s", sid.Long, sid.Short, sid.FriendlyName)

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

	t = template.Must(template.New("ks").Parse(templateRegister))
	buf = bytes.NewBuffer(nil)
	err = t.Execute(buf, vars)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to render kickstart template")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
