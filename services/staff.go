package services

import (
	"context"
	"errors"
	"strings"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"gorm.io/gorm"
)

type StaffService interface {
	CreateStaff(ctx context.Context, staff *dtos.CreateStaffRequest) (*dtos.StaffResponse, error)
	GetStaffById(ctx context.Context, id int) (*dtos.StaffResponse, error)
	UpdateStaff(ctx context.Context, staff *dtos.UpdateStaffRequest) (*dtos.StaffResponse, error)
	DeleteStaff(ctx context.Context, id int) error
	CreateStaffs(ctx context.Context, staffs []models.Staff) error
	GetStaffByProgramId(ctx context.Context, programId int) ([]dtos.StaffResponse, error)
	GetAllStaff(ctx context.Context) ([]dtos.StaffResponse, error)
	GetAllStaffByProgramId(ctx context.Context, programId int) ([]models.Staff, error)
	GetStaffByEmail(ctx context.Context, email string) (*models.Staff, error)
	GetStaffByName(ctx context.Context, name string) (*models.Staff, error)
}

type staffServiceImpl struct {
	staffRepo repositories.StaffRepository
}

func NewStaffService(staffRepo repositories.StaffRepository) StaffService {
	return &staffServiceImpl{
		staffRepo: staffRepo,
	}
}

func (s *staffServiceImpl) GetStaffByName(ctx context.Context, name string) (*models.Staff, error) {
	nameParts := strings.Split(name, " ")
	if len(nameParts) < 2 {
		return nil, errors.New("invalid name format")
	}

	firstName := nameParts[0]
	lastName := nameParts[1]

	staff, err := s.staffRepo.GetStaffByFirstNameAndLastName(ctx, firstName, lastName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found " + firstName + " " + lastName)
		}
		return nil, err
	}

	return staff, nil
}

func (s *staffServiceImpl) GetAllStaffByProgramId(ctx context.Context, programId int) ([]models.Staff, error) {
	staffs, err := s.staffRepo.GetAllStaffByProgramId(ctx, programId)
	if err != nil {
		return nil, err
	}

	return staffs, nil
}

func (s *staffServiceImpl) GetAllStaff(ctx context.Context) ([]dtos.StaffResponse, error) {
	staffs, err := s.staffRepo.GetAllStaffs(ctx)
	if err != nil {
		return nil, err
	}

	// convert to response
	var StaffResponses []dtos.StaffResponse
	for _, staff := range staffs {
		StaffResponses = append(StaffResponses, dtos.StaffResponse{
			ID:          staff.ID,
			PrefixTH:    staff.PrefixTH,
			PrefixEN:    staff.PrefixEN,
			FirstNameTH: staff.FirstNameTH,
			LastNameTH:  staff.LastNameTH,
			FirstNameEN: staff.FirstNameEN,
			LastNameEN:  staff.LastNameEN,
			Email:       staff.Email,
			ProgramID:   staff.ProgramID,
			Program:     staff.Program,
		})
	}

	return StaffResponses, nil
}

func (s *staffServiceImpl) CreateStaff(ctx context.Context, staffBody *dtos.CreateStaffRequest) (*dtos.StaffResponse, error) {

	staff := &models.Staff{
		PrefixTH:    staffBody.PrefixTH,
		PrefixEN:    staffBody.PrefixEN,
		FirstNameTH: staffBody.FirstNameTH,
		LastNameTH:  staffBody.LastNameTH,
		FirstNameEN: staffBody.FirstNameEN,
		LastNameEN:  staffBody.LastNameEN,
		Email:       staffBody.Email,
		IsActive:    staffBody.IsActive,
		ProgramID:   staffBody.ProgramID,
	}

	staff, err := s.staffRepo.CreateStaff(ctx, staff)
	if err != nil {
		return nil, err
	}

	return &dtos.StaffResponse{
		ID:          staff.ID,
		PrefixTH:    staff.PrefixTH,
		PrefixEN:    staff.PrefixEN,
		FirstNameTH: staff.FirstNameTH,
		LastNameTH:  staff.LastNameTH,
		FirstNameEN: staff.FirstNameEN,
		LastNameEN:  staff.LastNameEN,
		IsActive:    staff.IsActive,
		Email:       staff.Email,
		ProgramID:   staff.ProgramID,
		Program:     staff.Program,
	}, nil
}

func (s *staffServiceImpl) GetStaffById(ctx context.Context, id int) (*dtos.StaffResponse, error) {
	staff, err := s.staffRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dtos.StaffResponse{
		ID:          staff.ID,
		PrefixTH:    staff.PrefixTH,
		PrefixEN:    staff.PrefixEN,
		FirstNameTH: staff.FirstNameTH,
		LastNameTH:  staff.LastNameTH,
		FirstNameEN: staff.FirstNameEN,
		LastNameEN:  staff.LastNameEN,
		IsActive:    staff.IsActive,
		Email:       staff.Email,
		ProgramID:   staff.ProgramID,
	}, nil
}

func (s *staffServiceImpl) UpdateStaff(ctx context.Context, staff *dtos.UpdateStaffRequest) (*dtos.StaffResponse, error) {
	// convert from dto to model
	updatedStaff := &models.Staff{
		ID:          staff.ID,
		PrefixTH:    staff.PrefixTH,
		PrefixEN:    staff.PrefixEN,
		FirstNameTH: staff.FirstNameTH,
		LastNameTH:  staff.LastNameTH,
		FirstNameEN: staff.FirstNameEN,
		LastNameEN:  staff.LastNameEN,
		IsActive:    staff.IsActive,
		Email:       staff.Email,
		ProgramID:   staff.ProgramID,
	}
	updatedStaff, err := s.staffRepo.Update(ctx, staff.ID, updatedStaff)
	if err != nil {
		return nil, err
	}

	// convert to response
	return &dtos.StaffResponse{
		ID:          updatedStaff.ID,
		PrefixTH:    updatedStaff.PrefixTH,
		PrefixEN:    updatedStaff.PrefixEN,
		FirstNameTH: updatedStaff.FirstNameTH,
		LastNameTH:  updatedStaff.LastNameTH,
		FirstNameEN: updatedStaff.FirstNameEN,
		LastNameEN:  updatedStaff.LastNameEN,
		IsActive:    staff.IsActive,
		Email:       updatedStaff.Email,
		ProgramID:   updatedStaff.ProgramID,
	}, nil
}

func (s *staffServiceImpl) DeleteStaff(ctx context.Context, id int) error {
	return s.staffRepo.Delete(ctx, id)
}

func (s *staffServiceImpl) GetStaffByProgramId(ctx context.Context, programId int) ([]dtos.StaffResponse, error) {
	staffs, err := s.staffRepo.GetStaffByProgramId(programId)
	if err != nil {
		return nil, err
	}

	// convert to response
	var StaffResponses []dtos.StaffResponse
	for _, staff := range staffs {
		StaffResponses = append(StaffResponses, dtos.StaffResponse{
			ID:          staff.ID,
			PrefixTH:    staff.PrefixTH,
			PrefixEN:    staff.PrefixEN,
			FirstNameTH: staff.FirstNameTH,
			LastNameTH:  staff.LastNameTH,
			FirstNameEN: staff.FirstNameEN,
			LastNameEN:  staff.LastNameEN,
			Email:       staff.Email,
			IsActive:    staff.IsActive,
			ProgramID:   staff.ProgramID,
		})
	}

	return StaffResponses, nil
}

func (s *staffServiceImpl) GetStaffByEmail(ctx context.Context, email string) (*models.Staff, error) {
	staff, err := s.staffRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffServiceImpl) CreateStaffs(ctx context.Context, staffs []models.Staff) error {
	return s.staffRepo.CreateStaffs(ctx, staffs)
}
