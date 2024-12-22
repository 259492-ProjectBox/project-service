package repositories

import (
	"context"
	"time"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	repository[models.Calendar]
	GetByMajorAndDateRange(ctx context.Context, majorID int, startDate, endDate time.Time) ([]models.Calendar, error)
	CreateCalendar(ctx context.Context, calendar *models.Calendar) error
	GetCalendarByMajorID(ctx context.Context, majorID int) ([]models.Calendar, error)
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

func (r *calendarRepositoryImpl) GetByMajorAndDateRange(ctx context.Context, majorID int, startDate, endDate time.Time) ([]models.Calendar, error) {
	var calendars []models.Calendar

	// Query for overlapping events for the same major_id
	if err := r.db.WithContext(ctx).Debug().
		Where("major_id = ? AND NOT (start_date > ? OR end_date < ?)", majorID, endDate, startDate).
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

func (r *calendarRepositoryImpl) GetCalendarByMajorID(ctx context.Context, majorID int) ([]models.Calendar, error) {
	var calendars []models.Calendar

	if err := r.db.WithContext(ctx).Where("major_id = ?", majorID).Find(&calendars).Error; err != nil {
		return nil, err
	}

	return calendars, nil
}
