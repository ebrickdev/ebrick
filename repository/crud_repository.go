package repository

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CrudRepository defines a set of generic CRUD operations.
// Note: T is expected to be a struct type.
type CrudRepository[T any] interface {
	Create(ctx context.Context, et *T) (*T, error)
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, et *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context) ([]T, error)
	First(ctx context.Context, et *T) (*T, error)
	FindWithEntity(ctx context.Context, et *T) ([]T, error)
	FindWithConditions(ctx context.Context, conditions map[string]any) ([]T, error)
	FindWithOrConditions(ctx context.Context, conditions map[string]any) ([]T, error)
	CountWithConditions(ctx context.Context, conditions map[string]any) (int64, error)
	CountWithEntity(ctx context.Context, et *T) (int64, error)
	Exists(ctx context.Context, conditions map[string]any) (bool, error)
	// Add paging method
	ListPaged(ctx context.Context, offset int, limit int) ([]T, error)
}

// NewCrudRepository creates a new CrudRepository instance for type T.
// It also initializes a validator instance for reuse.
func NewCrudRepository[T any](db *gorm.DB) CrudRepository[T] {
	return &crudRepository[T]{
		db:       db,
		validate: validator.New(),
	}
}

type crudRepository[T any] struct {
	db       *gorm.DB
	validate *validator.Validate
}

func (r *crudRepository[T]) Create(ctx context.Context, et *T) (*T, error) {
	if err := r.validate.Struct(et); err != nil {
		return nil, err
	}
	err := r.db.WithContext(ctx).Create(et).Error
	return et, err
}

func (r *crudRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var et T
	err := r.db.WithContext(ctx).First(&et, id).Error
	return &et, err
}

func (r *crudRepository[T]) Update(ctx context.Context, et *T) (*T, error) {
	if err := r.validate.Struct(et); err != nil {
		return nil, err
	}
	err := r.db.WithContext(ctx).Save(et).Error
	return et, err
}

func (r *crudRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var et T
	return r.db.WithContext(ctx).Delete(&et, id).Error
}

func (r *crudRepository[T]) ListAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) First(ctx context.Context, et *T) (*T, error) {
	var existed T
	err := r.db.WithContext(ctx).Where(et).First(&existed).Error
	return &existed, err
}

func (r *crudRepository[T]) FindWithEntity(ctx context.Context, et *T) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(et).Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) FindWithConditions(ctx context.Context, conditions map[string]any) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(conditions).Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) FindWithOrConditions(ctx context.Context, conditions map[string]any) ([]T, error) {
	var entities []T
	// Return an empty slice if no conditions are provided.
	if len(conditions) == 0 {
		return entities, nil
	}
	query := r.db.WithContext(ctx).Model(new(T))
	for key, value := range conditions {
		query = query.Or(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) CountWithConditions(ctx context.Context, conditions map[string]any) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error
	return count, err
}

func (r *crudRepository[T]) CountWithEntity(ctx context.Context, et *T) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where(et).Count(&count).Error
	return count, err
}

func (r *crudRepository[T]) Exists(ctx context.Context, conditions map[string]any) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error
	return count > 0, err
}

// ListPaged returns a subset of entities based on offset and limit.
func (r *crudRepository[T]) ListPaged(ctx context.Context, offset int, limit int) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&entities).Error
	return entities, err
}
