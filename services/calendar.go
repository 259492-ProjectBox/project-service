package services

import (
	"context"
	"errors"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
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
	_, err := s.majorRepo.GetByMajorID(ctx, calendar.MajorID)
	if err != nil {
		return dtos.CalendarResponse{}, errors.New("major ID does not exist")
	}

	// Check if the start date already exists for the given major
	existingCalendar, err := s.calendarRepo.GetByTitleAndDateRange(ctx, calendar.Title, calendar.StartDate, calendar.EndDate)
	if err == nil && existingCalendar != nil {
		return dtos.CalendarResponse{}, errors.New("a calendar event with the same start date already exists for this major")
	}

	// Convert DTO to model
	newCalendar := &models.Calendar{
		StartDate:   calendar.StartDate,
		EndDate:     calendar.EndDate,
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
		MajorID:     newCalendar.MajorID,
	}

	return response, nil
}
