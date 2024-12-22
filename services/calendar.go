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

func (s *calendarServiceImpl) CreateCalendar(ctx context.Context, calendar *dtos.CreateCalendarRequest) (dtos.CalendarResponse, error) {
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
		StartDate:   newCalendar.StartDate,
		EndDate:     newCalendar.EndDate,
		Title:       newCalendar.Title,
		Description: newCalendar.Description,
		Major:       major.MajorName,
	}

	return response, nil
}
