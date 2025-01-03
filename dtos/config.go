package dtos

type ConfigReponse struct {
	ConfigName string `json:"config_name"`
	Value      string `json:"value"`
	MajorID    int    `json:"major_id"`
}
