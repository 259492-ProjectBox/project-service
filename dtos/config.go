package dtos

type ConfigReponse struct {
	ConfigName string `json:"config_name"`
	Value      string `json:"value"`
	ProgramID  int    `json:"program_id"`
}
