package drivers

import (
	"context"
	"fmt"

	"github.com/Foxpunk/courseforge/internal/interfaces"
	"gorm.io/gorm"
)

// Базовый репозиторий теперь работает только с типами, которые реализуют Validator
type baseRepository[T interfaces.Validator] struct {
	db *gorm.DB
}

func NewBaseRepository[T interfaces.Validator](db *gorm.DB) interfaces.BaseRepository[T] {
	return &baseRepository[T]{
		db: db,
	}
}

func (r *baseRepository[T]) Create(ctx context.Context, e *T) error {
	// Валидация перед созданием
	if err := (*e).Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return r.db.WithContext(ctx).Create(e).Error
}

func (r *baseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	res := r.db.WithContext(ctx).First(&entity, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &entity, nil
}

func (r *baseRepository[T]) Update(ctx context.Context, e *T) error {
	// Валидация перед обновлением
	if err := (*e).Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return r.db.WithContext(ctx).Save(e).Error
}

func (r *baseRepository[T]) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}

func (r *baseRepository[T]) List(ctx context.Context, limit, offset int) ([]T, error) {
	// Валидация параметров пагинации
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var entities []T
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}
