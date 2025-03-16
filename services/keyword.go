package services

import (
	"context"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type KeywordService interface {
	GetAllKeywords(ctx context.Context) ([]models.Keyword, error)
	GetKeywords(ctx context.Context, programID string) ([]models.Keyword, error)
	GetKeyword(ctx context.Context, id string) (*models.Keyword, error)
	CreateKeyword(ctx context.Context, keyword *models.Keyword) error
	UpdateKeyword(ctx context.Context, keyword *models.Keyword) error
	DeleteKeyword(ctx context.Context, id string) error
}

type keywordService struct {
	repo repositories.KeywordRepository
}

func NewKeywordService(repo repositories.KeywordRepository) KeywordService {
	return &keywordService{repo: repo}
}

func (s *keywordService) GetAllKeywords(ctx context.Context) ([]models.Keyword, error) {
	return s.repo.FindAll(ctx)
}

func (s *keywordService) GetKeywords(ctx context.Context, programID string) ([]models.Keyword, error) {
	return s.repo.FindAllByProgramId(ctx, programID)
}

func (s *keywordService) GetKeyword(ctx context.Context, id string) (*models.Keyword, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *keywordService) CreateKeyword(ctx context.Context, keyword *models.Keyword) error {
	return s.repo.Create(ctx, keyword)
}

func (s *keywordService) UpdateKeyword(ctx context.Context, keyword *models.Keyword) error {
	return s.repo.Update(ctx, keyword)
}

func (s *keywordService) DeleteKeyword(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
