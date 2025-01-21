package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
	"github.com/project-box/utils"
)

type CalendarService interface {
	CreateCalendar(ctx context.Context, calendar *dtos.CreateCalendarRequest) (dtos.CalendarResponse, error)
	GetCalendarByProgramId(ctx context.Context, programId int) ([]dtos.CalendarResponse, error)
	UpdateCalendar(ctx context.Context, calendar *dtos.UpdateCalendarRequest) (*dtos.CalendarResponse, error)
	DeleteCalendar(ctx context.Context, id int) error
}

type calendarServiceImpl struct {
	calendarRepo repositories.CalendarRepository
	programRepo  repositories.ProgramRepository
}

func NewCalendarService(calendarRepo repositories.CalendarRepository, programRepo repositories.ProgramRepository) CalendarService {
	return &calendarServiceImpl{
		calendarRepo: calendarRepo,
		programRepo:  programRepo,
	}
}

func (s *calendarServiceImpl) CreateCalendar(ctx context.Context, calendar *dtos.CreateCalendarRequest) (dtos.CalendarResponse, error) {
	// Check if the program ID exists
	program, err := s.programRepo.GetProgramById(ctx, calendar.ProgramID)
	if err != nil {
		return dtos.CalendarResponse{}, errors.New("program ID does not exist")
	}

	startDate, err := utils.ParseDateTime(calendar.StartDate)
	if err != nil {
		return dtos.CalendarResponse{}, fmt.Errorf("failed to parse start_date: %w", err)
	}
	endDate, err := utils.ParseDateTime(calendar.EndDate)
	if err != nil {
		return dtos.CalendarResponse{}, fmt.Errorf("failed to parse end_date: %w", err)
	}
	// Check if the start date already exists for the given program
	existingCalendars, err := s.calendarRepo.GetByProgramAndDateRange(ctx, calendar.ProgramID, startDate, endDate)
	fmt.Print(existingCalendars)
	if err != nil {
		return dtos.CalendarResponse{}, err // handle error
	}

	// If any overlapping events are found, return an error
	if len(existingCalendars) > 0 {
		return dtos.CalendarResponse{}, errors.New("a calendar event with the same date range already exists for this program")
	}

	// Convert DTO to model
	newCalendar := &models.Calendar{
		StartDate:   startDate,
		EndDate:     endDate,
		Title:       calendar.Title,
		Description: calendar.Description,
		ProgramID:   calendar.ProgramID,
	}

	// Create the calendar event
	err = s.calendarRepo.CreateCalendar(ctx, newCalendar)
	if err != nil {
		return dtos.CalendarResponse{}, err
	}

	// Convert model to DTO response
	response := dtos.CalendarResponse{
		ID:          newCalendar.ID,
		StartDate:   utils.FormatDate(newCalendar.StartDate),
		EndDate:     utils.FormatDate(newCalendar.EndDate),
		Title:       newCalendar.Title,
		Description: newCalendar.Description,
		Program:     program.ProgramNameTH,
	}

	return response, nil
}

func (s *calendarServiceImpl) GetCalendarByProgramId(ctx context.Context, programId int) ([]dtos.CalendarResponse, error) {
	calendars, err := s.calendarRepo.GetCalendarByProgramID(ctx, programId)
	if err != nil {
		return nil, err
	}

	var calendarResponses []dtos.CalendarResponse
	for _, calendar := range calendars {
		calendarResponse := dtos.CalendarResponse{
			ID:          calendar.ID,
			StartDate:   utils.FormatDate(calendar.StartDate),
			EndDate:     utils.FormatDate(calendar.EndDate),
			Title:       calendar.Title,
			Description: calendar.Description,
			Program:     fmt.Sprintf("%d", calendar.ProgramID),
		}
		calendarResponses = append(calendarResponses, calendarResponse)
	}

	return calendarResponses, nil
}

func (s *calendarServiceImpl) UpdateCalendar(ctx context.Context, calendar *dtos.UpdateCalendarRequest) (*dtos.CalendarResponse, error) {
	// Check if the calendar ID exists
	_, err := s.calendarRepo.GetCalendarByID(ctx, calendar.ID)
	if err != nil {
		return nil, errors.New("calendar ID does not exist")
	}

	startDate, err := utils.ParseDateTime(calendar.StartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_date: %w", err)
	}
	endDate, err := utils.ParseDateTime(calendar.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_date: %w", err)
	}
	// Convert DTO to model
	updatedCalendar := &models.Calendar{
		ID:          calendar.ID,
		StartDate:   startDate,
		EndDate:     endDate,
		Title:       calendar.Title,
		Description: calendar.Description,
		ProgramID:   calendar.ProgramID,
	}

	// Update the calendar
	updatedCalendar, err = s.calendarRepo.UpdateCalendar(ctx, updatedCalendar)
	if err != nil {
		return nil, err
	}

	// convert to calendar response
	response := dtos.CalendarResponse{
		ID:          updatedCalendar.ID,
		StartDate:   utils.FormatDate(updatedCalendar.StartDate),
		EndDate:     utils.FormatDate(updatedCalendar.EndDate),
		Title:       updatedCalendar.Title,
		Description: updatedCalendar.Description,
		Program:     fmt.Sprintf("%d", updatedCalendar.ProgramID),
	}

	return &response, nil
}

func (s *calendarServiceImpl) DeleteCalendar(ctx context.Context, id int) error {
	// Check if the calendar ID exists
	_, err := s.calendarRepo.GetCalendarByID(ctx, id)
	if err != nil {
		return errors.New("calendar ID does not exist")
	}

	// Delete the calendar
	err = s.calendarRepo.DeleteCalendar(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
