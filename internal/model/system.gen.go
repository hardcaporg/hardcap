// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameSystem = "system"

// System mapped from table <system>
type System struct {
	ID          int32  `gorm:"column:id;primaryKey" json:"id"`
	ApplianceID int32  `gorm:"column:appliance_id;not null" json:"appliance_id"`
	Name        string `gorm:"column:name;not null" json:"name"`
	UID         string `gorm:"column:uid;not null" json:"uid"`
	Sid         string `gorm:"column:sid;not null" json:"sid"`
}

// TableName System's table name
func (*System) TableName() string {
	return TableNameSystem
}
