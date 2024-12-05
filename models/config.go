package models

type Config struct {
	ConfigName string `json:"config_name"`
	Value      string `json:"value"`
	MajorID    int    `json:"major_id"`
	Major      Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
}
