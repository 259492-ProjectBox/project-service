package services

import (
	"github.com/project-box/repositories"
)

type Major interface {
}

type majorServiceImpl struct {
	majorRepo repositories.MajorRepository
}

func NewMajorService(majorRepo repositories.MajorRepository) Major {
	return &majorServiceImpl{
		majorRepo: majorRepo,
	}
}
