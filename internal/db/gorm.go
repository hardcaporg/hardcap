package db

import (
    "github.com/glebarez/sqlite"
    "gorm.io/gorm/logger"
)
import "gorm.io/gorm"

var Pool *gorm.DB

func Initialize() {
	var err error
	Pool, err = gorm.Open(sqlite.Open(":memory:?_pragma=foreign_keys(1)"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
	if err != nil {
		panic(err)
	}

    SetDefault(Pool)

    Pool.Exec(`

create table appliance (
    id integer not null primary key autoincrement,
    name text not null,
    url text not null
);

create table registration (
    id integer not null primary key autoincrement,
    sid text not null,
    name text not null,
    bios_vendor text,
    bios_version text,
    bios_release_date text,
    bios_revision text,
    firmware_revision text,
    system_manufacturer text,
    system_product_name text,
    system_version text,
    system_serial_number text,
    system_uuid text,
    system_sku_number text,
    system_family text,
    baseboard_manufacturer text,
    baseboard_product_name text,
    baseboard_version text,
    baseboard_serial_number text,
    baseboard_asset_tag text,
    chassis_manufacturer text,
    chassis_type text,
    chassis_version text,
    chassis_serial_number text,
    chassis_asset_tag text,
    processor_family text,
    processor_manufacturer text,
    processor_version text,
    processor_frequency text
);

`)
}
