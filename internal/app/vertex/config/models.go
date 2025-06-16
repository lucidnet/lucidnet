package config

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Value string
}
