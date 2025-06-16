package admin

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Identifier string `gorm:"unique" json:"identifier"`
	Password   string `json:"-"`
	Creator    string `json:"creator"`
}
