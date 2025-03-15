package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type KeywordRepository interface {
	FindAll(programID string) ([]models.Keyword, error)
	FindByID(id string) (models.Keyword, error)
	Create(keyword models.Keyword) error
	Update(keyword models.Keyword) error
	Delete(id string) error
}

type keywordRepository struct {
	DB *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) KeywordRepository {
	return &keywordRepository{DB: db}
}

func (r *keywordRepository) FindAll(programID string) ([]models.Keyword, error) {
	var keywords []models.Keyword
	err := r.DB.Where("program_id = ?", programID).Find(&keywords).Error
	return keywords, err
}

func (r *keywordRepository) FindByID(id string) (models.Keyword, error) {
	var keyword models.Keyword
	err := r.DB.First(&keyword, id).Error
	return keyword, err
}

func (r *keywordRepository) Create(keyword models.Keyword) error {
	return r.DB.Create(&keyword).Error
}

func (r *keywordRepository) Update(keyword models.Keyword) error {
	return r.DB.Save(&keyword).Error
}

func (r *keywordRepository) Delete(id string) error {
	return r.DB.Delete(&models.Keyword{}, id).Error
}
