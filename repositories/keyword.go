package repositories

import (
	"context"
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type KeywordRepository interface {
	FindAll(ctx context.Context) ([]models.Keyword, error)
	FindAllByProgramId(ctx context.Context, programID string) ([]models.Keyword, error)
	FindByID(ctx context.Context, id string) (*models.Keyword, error)
	FindByKeywordAndProgramId(ctx context.Context, keyword string, programId int) (*models.Keyword, error)
	Create(ctx context.Context, keyword *models.Keyword) error
	Update(ctx context.Context, keyword *models.Keyword) error
	Delete(ctx context.Context, id string) error
}

type keywordRepository struct {
	DB *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) KeywordRepository {
	return &keywordRepository{DB: db}
}

func (r *keywordRepository) FindByKeywordAndProgramId(ctx context.Context, keyword string, programId int) (*models.Keyword, error) {
	var keywordData *models.Keyword
	err := r.DB.WithContext(ctx).Where("keyword = ? AND program_id = ?", keyword, programId).First(keywordData).Error
	return keywordData, err
}

func (r *keywordRepository) FindAll(ctx context.Context) ([]models.Keyword, error) {
	var keywords []models.Keyword
	err := r.DB.WithContext(ctx).Preload("Program").Find(&keywords).Error
	return keywords, err
}

func (r *keywordRepository) FindAllByProgramId(ctx context.Context, programID string) ([]models.Keyword, error) {
	var keywords []models.Keyword
	err := r.DB.WithContext(ctx).Where("program_id = ?", programID).Preload("Program").Find(&keywords).Error
	return keywords, err
}

func (r *keywordRepository) FindByID(ctx context.Context, id string) (*models.Keyword, error) {
	var keyword *models.Keyword
	err := r.DB.WithContext(ctx).Preload("Program").First(keyword, id).Error
	return keyword, err
}

func (r *keywordRepository) Create(ctx context.Context, keyword *models.Keyword) error {
	return r.DB.WithContext(ctx).Create(keyword).Error
}

func (r *keywordRepository) Update(ctx context.Context, keyword *models.Keyword) error {
	return r.DB.WithContext(ctx).Save(keyword).Error
}

func (r *keywordRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&models.Keyword{}, id).Error
}
