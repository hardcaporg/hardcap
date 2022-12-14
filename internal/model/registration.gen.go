// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRegistration = "registration"

// Registration mapped from table <registration>
type Registration struct {
	ID                    int32  `gorm:"column:id;primaryKey" json:"id"`
	Sid                   string `gorm:"column:sid;not null" json:"sid"`
	Name                  string `gorm:"column:name;not null" json:"name"`
	BiosVendor            string `gorm:"column:bios_vendor;not null" json:"bios_vendor"`
	BiosVersion           string `gorm:"column:bios_version;not null" json:"bios_version"`
	BiosReleaseDate       string `gorm:"column:bios_release_date;not null" json:"bios_release_date"`
	BiosRevision          string `gorm:"column:bios_revision;not null" json:"bios_revision"`
	FirmwareRevision      string `gorm:"column:firmware_revision;not null" json:"firmware_revision"`
	SystemManufacturer    string `gorm:"column:system_manufacturer;not null" json:"system_manufacturer"`
	SystemProductName     string `gorm:"column:system_product_name;not null" json:"system_product_name"`
	SystemVersion         string `gorm:"column:system_version;not null" json:"system_version"`
	SystemSerialNumber    string `gorm:"column:system_serial_number;not null" json:"system_serial_number"`
	SystemUUID            string `gorm:"column:system_uuid;not null" json:"system_uuid"`
	SystemSkuNumber       string `gorm:"column:system_sku_number;not null" json:"system_sku_number"`
	SystemFamily          string `gorm:"column:system_family;not null" json:"system_family"`
	BaseboardManufacturer string `gorm:"column:baseboard_manufacturer;not null" json:"baseboard_manufacturer"`
	BaseboardProductName  string `gorm:"column:baseboard_product_name;not null" json:"baseboard_product_name"`
	BaseboardVersion      string `gorm:"column:baseboard_version;not null" json:"baseboard_version"`
	BaseboardSerialNumber string `gorm:"column:baseboard_serial_number;not null" json:"baseboard_serial_number"`
	BaseboardAssetTag     string `gorm:"column:baseboard_asset_tag;not null" json:"baseboard_asset_tag"`
	ChassisManufacturer   string `gorm:"column:chassis_manufacturer;not null" json:"chassis_manufacturer"`
	ChassisType           string `gorm:"column:chassis_type;not null" json:"chassis_type"`
	ChassisVersion        string `gorm:"column:chassis_version;not null" json:"chassis_version"`
	ChassisSerialNumber   string `gorm:"column:chassis_serial_number;not null" json:"chassis_serial_number"`
	ChassisAssetTag       string `gorm:"column:chassis_asset_tag;not null" json:"chassis_asset_tag"`
	ProcessorFamily       string `gorm:"column:processor_family;not null" json:"processor_family"`
	ProcessorManufacturer string `gorm:"column:processor_manufacturer;not null" json:"processor_manufacturer"`
	ProcessorVersion      string `gorm:"column:processor_version;not null" json:"processor_version"`
	ProcessorFrequency    string `gorm:"column:processor_frequency;not null" json:"processor_frequency"`
}

// TableName Registration's table name
func (*Registration) TableName() string {
	return TableNameRegistration
}
