package models

type Config struct {
	ConfigName string `json:"config_name" gorm:"unique"`
	Value      string `json:"value"`
}
