package dtos

type ConfigResponse struct {
	ConfigName string `json:"config_name"`
	Value      string `json:"value"`
	ProgramID  int    `json:"program_id"`
}
