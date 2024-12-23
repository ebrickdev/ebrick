package core

import "time"

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new generic Field.
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// NewStringField creates a new Field with a string value.
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// NewIntField creates a new Field with an integer value.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// NewBoolField creates a new Field with a boolean value.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// NewErrorField creates a new Field with an error value.
func Error(key string, err error) Field {
	return Field{Key: key, Value: err}
}

func Float(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value.Format(time.RFC3339)}
}
