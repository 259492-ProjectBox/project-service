package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository[M any] interface {
	Get(ctx context.Context, id int) (*M, error)
	GetMany(ctx context.Context) ([]M, error)
	Create(ctx context.Context, ent *M) (*M, error)
	CreateMany(ctx context.Context, ent []*M) error
	Update(ctx context.Context, id int, ent *M) (*M, error)
	Upsert(ctx context.Context, ent *M) (*M, error)
	Delete(ctx context.Context, id int) error
}

type repositoryImpl[M any] struct {
	db *gorm.DB
}

func newRepository[M any](db *gorm.DB) *repositoryImpl[M] {
	return &repositoryImpl[M]{
		db: db,
	}
}
func (r *repositoryImpl[M]) Get(ctx context.Context, id int) (*M, error) {
	ent := new(M)
	if err := r.db.First(ent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ent, nil
}

func (r *repositoryImpl[M]) GetMany(ctx context.Context) ([]M, error) {
	var ents []M
	if err := r.db.Find(&ents).Error; err != nil {
		return nil, err
	}
	return ents, nil
}

func (r *repositoryImpl[M]) Create(ctx context.Context, ent *M) (*M, error) {
	query := r.db.WithContext(ctx)

	if err := query.Create(ent).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *repositoryImpl[M]) CreateMany(ctx context.Context, ent []*M) error {
	query := r.db.WithContext(ctx)

	if err := query.Create(ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl[M]) Update(ctx context.Context, id int, ent *M) (*M, error) {
	query := r.db.WithContext(ctx)
	if err := query.Model(ent).Where("id = ?", id).Updates(ent).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *repositoryImpl[M]) Upsert(ctx context.Context, ent *M) (*M, error) {
	query := r.db.WithContext(ctx)
	if err := query.Preload(clause.Associations).Save(ent).First(ent).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *repositoryImpl[M]) Delete(ctx context.Context, id int) error {
	ent := new(M)
	query := r.db.WithContext(ctx)

	if err := query.Where("id = ?", id).Delete(ent).Error; err != nil {
		return err
	}

	return nil
}
