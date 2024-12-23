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
	CreateCalendarService(ctx context.Context, calendar *dtos.CreateCalendarRequest) (dtos.CalendarResponse, error)
	GetCalendarByMajorIDService(ctx context.Context, majorID int) ([]dtos.CalendarResponse, error)
	UpdateCalendarService(ctx context.Context, calendar *dtos.UpdateCalendarRequest) (*dtos.CalendarResponse, error)
	DeleteCalendarService(ctx context.Context, id int) error
}

type calendarServiceImpl struct {
	calendarRepo repositories.CalendarRepository
	majorRepo    repositories.MajorRepository
}

func NewCalendarService(calendarRepo repositories.CalendarRepository, majorRepo repositories.MajorRepository) CalendarService {
	return &calendarServiceImpl{
		calendarRepo: calendarRepo,
		majorRepo:    majorRepo,
	}
}

func (s *calendarServiceImpl) CreateCalendarService(ctx context.Context, calendar *dtos.CreateCalendarRequest) (dtos.CalendarResponse, error) {
	// Check if the major ID exists
	major, err := s.majorRepo.GetByMajorID(ctx, calendar.MajorID)
	if err != nil {
		return dtos.CalendarResponse{}, errors.New("major ID does not exist")
	}

	startDate, err := utils.ParseDateTime(calendar.StartDate)
	if err != nil {
		return dtos.CalendarResponse{}, fmt.Errorf("failed to parse start_date: %w", err)
	}
	endDate, err := utils.ParseDateTime(calendar.EndDate)
	if err != nil {
		return dtos.CalendarResponse{}, fmt.Errorf("failed to parse end_date: %w", err)
	}
	// Check if the start date already exists for the given major
	existingCalendars, err := s.calendarRepo.GetByMajorAndDateRange(ctx, calendar.MajorID, startDate, endDate)
	fmt.Print(existingCalendars)
	if err != nil {
		return dtos.CalendarResponse{}, err // handle error
	}

	// If any overlapping events are found, return an error
	if len(existingCalendars) > 0 {
		return dtos.CalendarResponse{}, errors.New("a calendar event with the same date range already exists for this major")
	}

	// Convert DTO to model
	newCalendar := &models.Calendar{
		StartDate:   startDate,
		EndDate:     endDate,
		Title:       calendar.Title,
		Description: calendar.Description,
		MajorID:     calendar.MajorID,
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
		Major:       major.MajorName,
	}

	return response, nil
}

func (s *calendarServiceImpl) GetCalendarByMajorIDService(ctx context.Context, majorID int) ([]dtos.CalendarResponse, error) {
	calendars, err := s.calendarRepo.GetCalendarByMajorID(ctx, majorID)
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
			Major:       fmt.Sprintf("%d", calendar.MajorID),
		}
		calendarResponses = append(calendarResponses, calendarResponse)
	}

	return calendarResponses, nil
}

func (s *calendarServiceImpl) UpdateCalendarService(ctx context.Context, calendar *dtos.UpdateCalendarRequest) (*dtos.CalendarResponse, error) {
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
		MajorID:     calendar.MajorID,
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
		Major:       fmt.Sprintf("%d", updatedCalendar.MajorID),
	}

	return &response, nil
}

func (s *calendarServiceImpl) DeleteCalendarService(ctx context.Context, id int) error {
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
