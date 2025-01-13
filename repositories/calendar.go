package repositories

import (
	"context"
	"time"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	repository[models.Calendar]
	GetByProgramAndDateRange(ctx context.Context, programId int, startDate, endDate time.Time) ([]models.Calendar, error)
	CreateCalendar(ctx context.Context, calendar *models.Calendar) error
	GetCalendarByProgramID(ctx context.Context, programId int) ([]models.Calendar, error)
	UpdateCalendar(ctx context.Context, updatedCalendar *models.Calendar) (*models.Calendar, error)
	GetCalendarByID(ctx context.Context, id int) (*models.Calendar, error)
	DeleteCalendar(ctx context.Context, id int) error
}

type calendarRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Calendar]
}

func NewCalendarRepository(db *gorm.DB) CalendarRepository {
	return &calendarRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Calendar](db),
	}
}

func (r *calendarRepositoryImpl) GetByProgramAndDateRange(ctx context.Context, programId int, startDate, endDate time.Time) ([]models.Calendar, error) {
	var calendars []models.Calendar

	// Query for overlapping events for the same program_id
	if err := r.db.WithContext(ctx).Debug().
		Where("program_id = ? AND NOT (start_date > ? OR end_date < ?)", programId, endDate, startDate).
		Find(&calendars).Error; err != nil {
		return nil, err // handle errors
	}

	return calendars, nil
}

func (r *calendarRepositoryImpl) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	if err := r.db.WithContext(ctx).Create(calendar).Error; err != nil {
		return err
	}
	return nil
}
func (r *calendarRepositoryImpl) GetCalendarByID(ctx context.Context, id int) (*models.Calendar, error) {

	var calendars *models.Calendar
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&calendars).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}

func (r *calendarRepositoryImpl) GetCalendarByProgramID(ctx context.Context, programId int) ([]models.Calendar, error) {
	var calendars []models.Calendar

	if err := r.db.WithContext(ctx).Where("program_id = ?", programId).Find(&calendars).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}

func (r *calendarRepositoryImpl) UpdateCalendar(ctx context.Context, updatedCalendar *models.Calendar) (*models.Calendar, error) {
	var calendar models.Calendar
	if err := r.db.WithContext(ctx).Model(&models.Calendar{}).Where("id = ?", updatedCalendar.ID).First(&calendar).Error; err != nil {
		return nil, err
	}

	calendar.StartDate = updatedCalendar.StartDate
	calendar.EndDate = updatedCalendar.EndDate
	calendar.Title = updatedCalendar.Title
	calendar.Description = updatedCalendar.Description
	calendar.ProgramID = updatedCalendar.ProgramID

	if err := r.db.WithContext(ctx).Save(&calendar).Error; err != nil {
		return nil, err
	}

	return &calendar, nil
}

func (r *calendarRepositoryImpl) DeleteCalendar(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Calendar{}).Error; err != nil {
		return err
	}
	return nil
}
