package repositories

import (
	"context"
	"errors"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type StaffRepository interface {
	repository[models.Staff]
	GetStaffByFirstNameAndLastName(ctx context.Context, firstNameTH, lastNameTH string) (*models.Staff, error)
	GetStaffById(id int) (*models.Staff, error)
	GetStaffByProgramId(programId int) ([]models.Staff, error)
	GetAllStaffs(ctx context.Context) ([]models.Staff, error)
	GetAllStaffByProgramId(ctx context.Context, programId int) ([]models.Staff, error)
	CreateStaff(staff *models.Staff) error
	CreateStaffs(ctx context.Context, staffs []models.Staff) error
	UpdateStaff(updatedStaff *models.Staff) (*models.Staff, error)
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

func (r *staffRepositoryImpl) CreateStaffs(ctx context.Context, staffs []models.Staff) error {
	if len(staffs) == 0 {
		return nil
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, staff := range staffs {
		if err := r.createStaff(ctx, tx, staff); err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errors.New("duplicate entry for unique_email_program key: staff " + staff.FirstNameTH + " " + staff.LastNameTH)
			}
			return err
		}
	}

	return tx.Commit().Error
}

func (s *staffRepositoryImpl) createStaff(ctx context.Context, tx *gorm.DB, staff models.Staff) error {
	return tx.Save(&staff).Error
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

func (r *staffRepositoryImpl) GetStaffByFirstNameAndLastName(ctx context.Context, firstName, lastName string) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.Where("(first_name_th = ? AND last_name_th = ?) OR (first_name_en = ? AND last_name_en = ?)",
		firstName, lastName, firstName, lastName).
		First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepositoryImpl) GetAllStaffByProgramId(ctx context.Context, programId int) ([]models.Staff, error) {
	var program1Staffs []models.Staff
	var otherStaffs []models.Staff

	err := r.db.WithContext(ctx).Where("program_id = ?", programId).Preload("Program").Find(&program1Staffs).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Where("email NOT IN (?)", r.db.Model(&models.Staff{}).Select("Email").Where("program_id = ?", programId)).Preload("Program").Find(&otherStaffs).Error
	if err != nil {
		return nil, err
	}
	// Merge results
	staffs := append(program1Staffs, otherStaffs...)
	return staffs, nil
}
