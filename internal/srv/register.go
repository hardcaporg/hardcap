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
		BIOSVendor            string `json:"bios-vendor"`
		BIOSVersion           string `json:"bios-version"`
		BIOSReleaseDate       string `json:"bios-release-date"`
		BIOSRevision          string `json:"bios-revision"`
		FirmwareRevision      string `json:"firmware-revision"`
		SystemManufacturer    string `json:"system-manufacturer"`
		SystemProductName     string `json:"system-product-name"`
		SystemVersion         string `json:"system-version"`
		SystemSerialNumber    string `json:"system-serial-number"`
		SystemUUID            string `json:"system-uuid"`
		SystemSKUNumber       string `json:"system-sku-number"`
		SystemFamily          string `json:"system-family"`
		BaseboardManufacturer string `json:"baseboard-manufacturer"`
		BaseboardProductName  string `json:"baseboard-product-name"`
		BaseboardVersion      string `json:"baseboard-version"`
		BaseboardSerialNumber string `json:"baseboard-serial-number"`
		BaseboardAssetTag     string `json:"baseboard-asset-tag"`
		ChassisManufacturer   string `json:"chassis-manufacturer"`
		ChassisType           string `json:"chassis-type"`
		ChassisVersion        string `json:"chassis-version"`
		ChassisNumber         string `json:"chassis-serial-number"`
		ChassisAssetTag       string `json:"chassis-asset-tag"`
		ProcessorFamily       string `json:"processor-family"`
		ProcessorManufacturer string `json:"processor-manufacturer"`
		ProcessorVersion      string `json:"processor-version"`
		ProcessorFrequency    string `json:"processor-frequency"`
	} `json:"dmi"`
	Log [][]string `json:"log"`
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
