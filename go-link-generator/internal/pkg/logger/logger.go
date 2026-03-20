// Package logger provides structured logging with wide events support.
// Following loggingsucks.com principles: structured, contextual, and searchable.
//
// This package implements the wide events pattern where each request generates
// a single canonical log line containing all relevant context. Business logic
// can enrich this log event throughout the request lifecycle.
package logger

import (
	"context"
)

// Logger defines the interface for structured logging.
// All methods accept a context for automatic extraction of request metadata
// (request_id, user_id, trace_id) which are automatically included in log output.
type Logger interface {
	// Core logging methods with structured fields and context extraction
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)

	// With methods for building loggers with preset fields
	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

// Field represents a structured logging field.
// Use the field constructor functions (String, Int, Bool, etc.) to create fields.
type Field struct {
	Key   string
	Value interface{}
	Type  FieldType
}

// FieldType represents the type of a logging field for type-safe conversion.
type FieldType int

const (
	FieldTypeAny FieldType = iota
	FieldTypeString
	FieldTypeInt
	FieldTypeInt64
	FieldTypeBool
	FieldTypeDuration
	FieldTypeError
)

// ============================================================================
// Field Constructors
// ============================================================================

// String creates a string field.
func String(key, value string) Field {
	return Field{Key: key, Value: value, Type: FieldTypeString}
}

// Int creates an int field.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value, Type: FieldTypeInt}
}

// Int64 creates an int64 field.
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value, Type: FieldTypeInt64}
}

// Bool creates a boolean field.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value, Type: FieldTypeBool}
}

// Any creates a field with any type (uses reflection).
// Prefer typed constructors (String, Int, etc.) for better performance.
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Type: FieldTypeAny}
}

// Error creates an error field.
func Error(err error) Field {
	return Field{Key: "error", Value: err, Type: FieldTypeError}
}

// Duration creates a duration field.
func Duration(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Type: FieldTypeDuration}
}

// Instance is the global logger instance that implements Logger interface.
// This is initialized by the Initialize function in the framework wrapper.
var Instance Logger
