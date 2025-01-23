package repository

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CrudRepository[T any] interface {
	Create(et T) (*T, error)
	FindByID(id uuid.UUID) (*T, error)
	Update(et T) (*T, error)
	Delete(id uuid.UUID) error
	ListAll() ([]T, error)
	First(et T) (*T, error)
	FindWithEntity(et T) ([]T, error)
	FindWithConditions(conditions map[string]any) ([]T, error)
	FindWithOrConditions(conditions map[string]any) ([]T, error)
	CountWithConditions(conditions map[string]any) (int64, error)
	CountWithEntity(et T) (int64, error)
}

func NewCrudRepository[T any](db *gorm.DB) CrudRepository[T] {
	return &crudRepository[T]{db: db}
}

type crudRepository[T any] struct {
	db *gorm.DB
}

func (r *crudRepository[T]) Create(et T) (*T, error) {
	v := validator.New()
	if err := v.Struct(et); err != nil {
		return nil, err
	}
	err := r.db.Create(&et).Error
	return &et, err
}

func (r *crudRepository[T]) FindByID(id uuid.UUID) (*T, error) {
	var et T
	err := r.db.First(&et, id).Error
	return &et, err
}

func (r *crudRepository[T]) Update(et T) (*T, error) {
	v := validator.New()
	if err := v.Struct(et); err != nil {
		return nil, err
	}
	err := r.db.Save(&et).Error
	return &et, err
}

func (r *crudRepository[T]) Delete(id uuid.UUID) error {
	var et T
	return r.db.Delete(&et, id).Error
}

func (r *crudRepository[T]) ListAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) First(et T) (*T, error) {
	var existed T
	err := r.db.Where(et).First(&existed).Error
	return &existed, err
}

func (r *crudRepository[T]) FindWithEntity(et T) ([]T, error) {
	var entities []T
	err := r.db.Where(et).Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) FindWithConditions(conditions map[string]any) ([]T, error) {
	var entities []T
	err := r.db.Where(conditions).Find(&entities).Error
	return entities, err
}
func (r *crudRepository[T]) FindWithOrConditions(conditions map[string]any) ([]T, error) {
	var entities []T
	query := r.db.Model(new(T))
	for key, value := range conditions {
		query.Or(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (r *crudRepository[T]) CountWithConditions(conditions map[string]any) (int64, error) {
	var count int64
	err := r.db.Model(new(T)).Where(conditions).Count(&count).Error
	return count, err
}

func (r *crudRepository[T]) CountWithEntity(et T) (int64, error) {
	var count int64
	err := r.db.Model(new(T)).Where(et).Count(&count).Error
	return count, err
}
