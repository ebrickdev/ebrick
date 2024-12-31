package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	field := Any("username", "johndoe")
	assert.Equal(t, "username", field.Key)
	assert.Equal(t, "johndoe", field.Value)
}

func TestString(t *testing.T) {
	field := String("email", "johndoe@example.com")
	assert.Equal(t, "email", field.Key)
	assert.Equal(t, "johndoe@example.com", field.Value)
}

func TestInt(t *testing.T) {
	field := Int("age", 30)
	assert.Equal(t, "age", field.Key)
	assert.Equal(t, 30, field.Value)
}

func TestBool(t *testing.T) {
	fieldTrue := Bool("is_active", true)
	assert.Equal(t, "is_active", fieldTrue.Key)
	assert.Equal(t, true, fieldTrue.Value)

	fieldFalse := Bool("is_active", false)
	assert.Equal(t, "is_active", fieldFalse.Key)
	assert.Equal(t, false, fieldFalse.Value)
}

func TestError(t *testing.T) {
	err := errors.New("file not found")
	field := Error(err)
	assert.Equal(t, "error", field.Key)
	assert.Equal(t, err, field.Value)
}

func TestFloat(t *testing.T) {
	field := Float("pi", 3.14159)
	assert.Equal(t, "pi", field.Key)
	assert.Equal(t, 3.14159, field.Value)
}

func TestTime(t *testing.T) {
	testTime := time.Date(2023, 12, 26, 15, 0, 0, 0, time.UTC)
	field := Time("created_at", testTime)
	assert.Equal(t, "created_at", field.Key)
	assert.Equal(t, testTime.Format(time.RFC3339), field.Value)
}
func TestUserField(t *testing.T) {
	field := UserField("user123")
	assert.Equal(t, "userID", field.Key)
	assert.Equal(t, "user123", field.Value)
}

func TestRequestIDField(t *testing.T) {
	field := RequestIDField("550e8400-e29b-41d4-a716-446655440000")
	assert.Equal(t, "requestID", field.Key)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", field.Value)
}

func TestTenantIDField(t *testing.T) {
	field := TenantIDField("550e8400-e29b-41d4-a716-446655440001")
	assert.Equal(t, "tenantID", field.Key)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", field.Value)
}
