package people

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type Repo interface {
	List(ctx context.Context) ([]Person, error)
	Get(ctx context.Context, id uint) (Person, error)
	Create(ctx context.Context, p *Person) error
	Update(ctx context.Context, p *Person) error
	Delete(ctx context.Context, id uint) error
	EmailExists(ctx context.Context, email string, excludeID uint) (bool, error)
}

type GormRepo struct{ db *gorm.DB }

func NewGormRepo(db *gorm.DB) *GormRepo { return &GormRepo{db: db} }

func (r *GormRepo) List(ctx context.Context) ([]Person, error) {
	var out []Person
	err := r.db.WithContext(ctx).Order("id desc").Find(&out).Error
	return out, err
}

func (r *GormRepo) Get(ctx context.Context, id uint) (Person, error) {
	var p Person
	err := r.db.WithContext(ctx).First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Person{}, ErrNotFound
	}
	return p, err
}

func (r *GormRepo) Create(ctx context.Context, p *Person) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *GormRepo) Update(ctx context.Context, p *Person) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *GormRepo) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&Person{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormRepo) EmailExists(ctx context.Context, email string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&Person{}).Where("email = ?", email)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
