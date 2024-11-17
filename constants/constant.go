package constants

type Status string

const (
	StatusInReview Status = "อยู่ในการพิจารณา"
	StatusPass     Status = "ผ่านการพิจารณา"
)

type Role string

const ()

type ResourceType string

const (
	Github     ResourceType = "github"
	Youtube    ResourceType = "youtube"
	Powerpoint ResourceType = "powerpoint"
	PDF        ResourceType = "pdf"
	Download   ResourceType = "download"
	Asset      ResourceType = "asset"
)
