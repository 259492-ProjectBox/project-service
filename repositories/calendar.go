package repositories

import (
	"context"
	"time"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type CalendarRepository interface {
	repository[models.Calendar]
	GetByTitleAndDateRange(ctx context.Context, title string, startDate, endDate time.Time) (*models.Calendar, error)
	CreateCalendar(ctx context.Context, calendar *models.Calendar) error
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

func (r *calendarRepositoryImpl) GetByTitleAndDateRange(ctx context.Context, title string, startDate, endDate time.Time) (*models.Calendar, error) {
	filters := map[string]interface{}{"title": title, "start_date": startDate, "end_date": endDate}
	var calendar models.Calendar
	if err := r.db.WithContext(ctx).Where(filters).First(&calendar).Error; err != nil {
		return nil, err
	}
	return &calendar, nil
}

func (r *calendarRepositoryImpl) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	if err := r.db.WithContext(ctx).Create(calendar).Error; err != nil {
		return err
	}
	return nil
}
