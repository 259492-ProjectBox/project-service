package repositories

import (
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type StaffRepository interface {
	repository[models.Staff]
	GetStaffById(id int) (*models.Staff, error)
	GetStaffByProgramId(programId int) ([]models.Staff, error)
	CreateStaff(staff *models.Staff) error
	UpdateStaff(updatedStaff *models.Staff) (*models.Staff, error)
	GetAllStaffs() ([]models.Staff, error)
}

type staffRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Staff]
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Staff](db),
	}
}

// get all staff
func (r *staffRepositoryImpl) GetAllStaffs() ([]models.Staff, error) {
	var staffs []models.Staff
	if err := r.db.Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

func (r *staffRepositoryImpl) GetStaffById(id int) (*models.Staff, error) {
	var staff *models.Staff
	if err := r.db.Where("id = ?", id).First(&staff).Error; err != nil {
		return nil, err
	}

	return staff, nil
}

// get staff by program id
func (r *staffRepositoryImpl) GetStaffByProgramId(programId int) ([]models.Staff, error) {
	var staffs []models.Staff

	if err := r.db.Where("program_id = ?", programId).Find(&staffs).Error; err != nil {
		return nil, err
	}

	return staffs, nil
}

// create staff
func (r *staffRepositoryImpl) CreateStaff(staff *models.Staff) error {
	if err := r.db.Create(staff).Error; err != nil {
		return err
	}
	return nil
}

// update staff
func (r *staffRepositoryImpl) UpdateStaff(updatedStaff *models.Staff) (*models.Staff, error) {
	if err := r.db.Save(updatedStaff).Error; err != nil {
		return nil, err
	}
	return updatedStaff, nil
}
