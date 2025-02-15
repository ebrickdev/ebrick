package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCrudRepository is a mock implementation of the generic CrudRepository interface.
type MockCrudRepository[T any] struct {
	mock.Mock
}

func (m *MockCrudRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockCrudRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockCrudRepository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockCrudRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCrudRepository[T]) ListAll(ctx context.Context) ([]T, error) {
	args := m.Called(ctx)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockCrudRepository[T]) First(ctx context.Context, entity *T) (*T, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockCrudRepository[T]) FindWithEntity(ctx context.Context, entity *T) ([]T, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockCrudRepository[T]) FindWithConditions(ctx context.Context, conditions map[string]interface{}) ([]T, error) {
	args := m.Called(ctx, conditions)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockCrudRepository[T]) FindWithOrConditions(ctx context.Context, conditions map[string]interface{}) ([]T, error) {
	args := m.Called(ctx, conditions)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockCrudRepository[T]) CountWithConditions(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	args := m.Called(ctx, conditions)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCrudRepository[T]) CountWithEntity(ctx context.Context, entity *T) (int64, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCrudRepository[T]) Exists(ctx context.Context, conditions map[string]interface{}) (bool, error) {
	args := m.Called(ctx, conditions)
	return args.Get(0).(bool), args.Error(1)
}
