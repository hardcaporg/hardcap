package srv

import (
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/hardcaporg/hardcap/internal/model"
	"net/http"

	"github.com/go-chi/render"
	"github.com/hardcaporg/hardcap/internal/ctxval"
)

type RegisterHostPayload struct {
	MacAddresses []string `json:"mac"`
	Serial       string   `json:"serial"`

	CPU struct {
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
		ChassisSerialNumber   string `json:"chassis-serial-number"`
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

	sid := model.NewSystemID(payload.Serial, payload.MacAddresses...)
	logger.Trace().Msgf("System ID: %s %s %s", sid.Long, sid.Short, sid.FriendlyName)

	registration := model.Registration{
		Sid:                   sid.Long,
		Name:                  sid.FriendlyName,
		BiosVendor:            payload.DMI.BIOSVendor,
		BiosVersion:           payload.DMI.BIOSVersion,
		BiosReleaseDate:       payload.DMI.BIOSReleaseDate,
		BiosRevision:          payload.DMI.BIOSRevision,
		FirmwareRevision:      payload.DMI.FirmwareRevision,
		SystemManufacturer:    payload.DMI.SystemManufacturer,
		SystemProductName:     payload.DMI.SystemProductName,
		SystemVersion:         payload.DMI.SystemVersion,
		SystemSerialNumber:    payload.DMI.SystemSerialNumber,
		SystemUUID:            payload.DMI.SystemUUID,
		SystemSkuNumber:       payload.DMI.SystemSKUNumber,
		SystemFamily:          payload.DMI.SystemFamily,
		BaseboardManufacturer: payload.DMI.BaseboardManufacturer,
		BaseboardProductName:  payload.DMI.BaseboardProductName,
		BaseboardVersion:      payload.DMI.BaseboardVersion,
		BaseboardSerialNumber: payload.DMI.BaseboardSerialNumber,
		BaseboardAssetTag:     payload.DMI.BaseboardAssetTag,
		ChassisManufacturer:   payload.DMI.ChassisManufacturer,
		ChassisType:           payload.DMI.ChassisType,
		ChassisVersion:        payload.DMI.ChassisVersion,
		ChassisSerialNumber:   payload.DMI.ChassisSerialNumber,
		ChassisAssetTag:       payload.DMI.ChassisAssetTag,
		ProcessorFamily:       payload.DMI.ProcessorFamily,
		ProcessorManufacturer: payload.DMI.ProcessorManufacturer,
		ProcessorVersion:      payload.DMI.ProcessorVersion,
		ProcessorFrequency:    payload.DMI.ProcessorFrequency,
	}
    err := db.Registration.Create(&registration)
    if err != nil {
        panic(err) 
    }
	logger.Debug().Msgf("Host registered: %+v", payload)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
