package repositories

import (
	"context"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type StaffRepository interface {
	repository[models.Staff]
	GetStaffById(id int) (*models.Staff, error)
	GetStaffByProgramId(programId int) ([]models.Staff, error)
	CreateStaff(staff *models.Staff) error
	UpsertStaffs(ctx context.Context, staffs []models.Staff) error
	UpdateStaff(updatedStaff *models.Staff) (*models.Staff, error)
	GetAllStaffs(ctx context.Context) ([]models.Staff, error)
	GetByEmail(ctx context.Context, email string) (*models.Staff, error)
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

func (r *staffRepositoryImpl) GetAllStaffs(ctx context.Context) ([]models.Staff, error) {
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

func (r *staffRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).Where("email = ?", email).Preload("Program").First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepositoryImpl) UpsertStaffs(ctx context.Context, staffs []models.Staff) error {
	if len(staffs) == 0 {
		return nil
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, staff := range staffs {
		if err := r.upsertStaff(ctx, tx, staff); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *staffRepositoryImpl) upsertStaff(ctx context.Context, tx *gorm.DB, staff models.Staff) error {
	existingStaff, err := s.GetByEmail(ctx, staff.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingStaff != nil {
		existingStaff.PrefixTH = staff.PrefixTH
		existingStaff.PrefixEN = staff.PrefixEN
		existingStaff.FirstNameTH = staff.FirstNameTH
		existingStaff.LastNameTH = staff.LastNameTH
		existingStaff.FirstNameEN = staff.FirstNameEN
		existingStaff.LastNameEN = staff.LastNameEN
		existingStaff.ProgramID = staff.ProgramID
		existingStaff.IsResigned = staff.IsResigned
		return tx.Save(existingStaff).Error
	}

	return tx.Create(&staff).Error
}
