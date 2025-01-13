package services

import (
	"context"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type StaffService interface {
	CreateStaff(ctx context.Context, staff *dtos.CreateStaffRequest) (*dtos.StaffResponse, error)
	GetStaffById(ctx context.Context, id int) (*dtos.StaffResponse, error)
	UpdateStaff(ctx context.Context, staff *dtos.UpdateStaffRequest) (*dtos.StaffResponse, error)
	DeleteStaff(ctx context.Context, id int) error
	GetStaffByProgramId(ctx context.Context, programId int) ([]dtos.StaffResponse, error)
}

type staffServiceImpl struct {
	staffRepo repositories.StaffRepository
}

func NewStaffService(staffRepo repositories.StaffRepository) StaffService {
	return &staffServiceImpl{
		staffRepo: staffRepo,
	}
}

func (s *staffServiceImpl) CreateStaff(ctx context.Context, staffBody *dtos.CreateStaffRequest) (*dtos.StaffResponse, error) {

	staff := &models.Staff{
		Prefix:    staffBody.Prefix,
		FirstName: staffBody.FirstName,
		LastName:  staffBody.LastName,
		Email:     staffBody.Email,
		ProgramID: staffBody.ProgramID,
	}

	staff, err := s.staffRepo.Create(ctx, staff)
	if err != nil {
		return nil, err
	}

	return &dtos.StaffResponse{
		ID:        staff.ID,
		Prefix:    staff.Prefix,
		FirstName: staff.FirstName,
		LastName:  staff.LastName,
		Email:     staff.Email,
		ProgramID: staff.ProgramID,
	}, nil

}

func (s *staffServiceImpl) GetStaffById(ctx context.Context, id int) (*dtos.StaffResponse, error) {
	staff, err := s.staffRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dtos.StaffResponse{
		ID:        staff.ID,
		Prefix:    staff.Prefix,
		FirstName: staff.FirstName,
		LastName:  staff.LastName,
		Email:     staff.Email,
		ProgramID: staff.ProgramID,
	}, nil
}

func (s *staffServiceImpl) UpdateStaff(ctx context.Context, staff *dtos.UpdateStaffRequest) (*dtos.StaffResponse, error) {
	// convert from dto to model
	updatedStaff := &models.Staff{
		ID:        staff.ID,
		Prefix:    staff.Prefix,
		FirstName: staff.FirstName,
		LastName:  staff.LastName,
		Email:     staff.Email,
		ProgramID: staff.ProgramID,
	}
	updatedStaff, err := s.staffRepo.Update(ctx, staff.ID, updatedStaff)
	if err != nil {
		return nil, err
	}

	// convert to response
	return &dtos.StaffResponse{
		ID:        updatedStaff.ID,
		Prefix:    updatedStaff.Prefix,
		FirstName: updatedStaff.FirstName,
		LastName:  updatedStaff.LastName,
		Email:     updatedStaff.Email,
		ProgramID: updatedStaff.ProgramID,
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
			ID:        staff.ID,
			Prefix:    staff.Prefix,
			FirstName: staff.FirstName,
			LastName:  staff.LastName,
			Email:     staff.Email,
			ProgramID: staff.ProgramID,
		})
	}

	return StaffResponses, nil
}
