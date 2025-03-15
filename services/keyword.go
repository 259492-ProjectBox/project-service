package services

import (
	"context"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type KeywordService interface {
	GetAllKeywords(ctx context.Context) ([]models.Keyword, error)
	GetKeywords(programID string) ([]models.Keyword, error)
	GetKeyword(id string) (models.Keyword, error)
	CreateKeyword(keyword models.Keyword) error
	UpdateKeyword(keyword models.Keyword) error
	DeleteKeyword(id string) error
}

type keywordService struct {
	repo repositories.KeywordRepository
}

func NewKeywordService(repo repositories.KeywordRepository) KeywordService {
	return &keywordService{repo: repo}
}

func (s *keywordService) GetAllKeywords(ctx context.Context) ([]models.Keyword, error) {
	return s.repo.FindAll()
}

func (s *keywordService) GetKeywords(programID string) ([]models.Keyword, error) {
	return s.repo.FindAllByProgramId(programID)
}

func (s *keywordService) GetKeyword(id string) (models.Keyword, error) {
	return s.repo.FindByID(id)
}

func (s *keywordService) CreateKeyword(keyword models.Keyword) error {
	return s.repo.Create(keyword)
}

func (s *keywordService) UpdateKeyword(keyword models.Keyword) error {
	return s.repo.Update(keyword)
}

func (s *keywordService) DeleteKeyword(id string) error {
	return s.repo.Delete(id)
}
