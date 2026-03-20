package logger

import (
	"context"
	"os"
	"sync"

	"github.com/jeremygprawira/go-link-generator/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapFieldPool is a sync.Pool for reusing zap.Field slices.
// This reduces GC pressure in high-throughput logging scenarios.
var zapFieldPool = sync.Pool{
	New: func() interface{} {
		// Pre-allocate with capacity for typical log calls (context + user fields)
		fields := make([]zap.Field, 0, 10)
		return &fields
	},
}

// ZapLogger wraps zap.Logger to implement the Logger interface.
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger instance.
func NewZapLogger(zapLog *zap.Logger) Logger {
	return &ZapLogger{logger: zapLog}
}

// Debug logs a debug message with automatic context extraction.
func (z *ZapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	z.logger.Debug(msg, z.buildFields(ctx, fields)...)
}

// Info logs an info message with automatic context extraction.
func (z *ZapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	z.logger.Info(msg, z.buildFields(ctx, fields)...)
}

// Warn logs a warning message with automatic context extraction.
func (z *ZapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	z.logger.Warn(msg, z.buildFields(ctx, fields)...)
}

// Error logs an error message with automatic context extraction.
func (z *ZapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	z.logger.Error(msg, z.buildFields(ctx, fields)...)
}

// Fatal logs a fatal message and exits with automatic context extraction.
func (z *ZapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	z.logger.Fatal(msg, z.buildFields(ctx, fields)...)
}

// With returns a logger with preset fields.
func (z *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{logger: z.logger.With(convertFields(fields)...)}
}

// WithContext returns a logger with context values extracted and preset.
// This is useful for creating a logger that always includes request metadata.
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	contextFields := extractContextFields(ctx)
	if len(contextFields) == 0 {
		return z
	}
	return &ZapLogger{logger: z.logger.With(contextFields...)}
}

// buildFields combines user-provided fields with context-extracted fields.
// Uses sync.Pool to reduce allocations in high-throughput scenarios.
func (z *ZapLogger) buildFields(ctx context.Context, fields []Field) []zap.Field {
	// Extract context fields (request_id, user_id, trace_id)
	contextFields := extractContextFields(ctx)

	// Convert user fields
	userFields := convertFields(fields)

	// Calculate total size needed
	totalSize := len(contextFields) + len(userFields)

	// Get a slice from the pool
	allFieldsPtr := zapFieldPool.Get().(*[]zap.Field)
	allFields := *allFieldsPtr

	// Reset the slice and ensure capacity
	if cap(allFields) < totalSize {
		allFields = make([]zap.Field, 0, totalSize)
	} else {
		allFields = allFields[:0]
	}

	allFields = append(allFields, contextFields...)
	allFields = append(allFields, userFields...)
	return allFields
}

// extractContextFields extracts common fields from context.
// This automatically includes request_id, user_id, and trace_id in all logs.
func extractContextFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}
	fields := make([]zap.Field, 0, 3) // Pre-allocate for common fields
	if requestID := GetRequestID(ctx); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if userID := GetUserID(ctx); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	return fields
}

// convertFields converts logger.Field to zap.Field with type preservation.
func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch f.Type {
		case FieldTypeString:
			zapFields[i] = zap.String(f.Key, f.Value.(string))
		case FieldTypeInt:
			zapFields[i] = zap.Int(f.Key, f.Value.(int))
		case FieldTypeInt64:
			zapFields[i] = zap.Int64(f.Key, f.Value.(int64))
		case FieldTypeBool:
			zapFields[i] = zap.Bool(f.Key, f.Value.(bool))
		case FieldTypeError:
			if err, ok := f.Value.(error); ok {
				zapFields[i] = zap.Error(err)
			} else {
				zapFields[i] = zap.Any(f.Key, f.Value)
			}
		case FieldTypeDuration:
			zapFields[i] = zap.Any(f.Key, f.Value)
		default:
			zapFields[i] = zap.Any(f.Key, f.Value)
		}
	}
	return zapFields
}

// Log is the raw zap.Logger instance.
var Log *zap.Logger

// Initialize configures and initializes the global logger based on environment.
func Initialize(configuration *config.Configuration) {
	var zapConfig zap.Config
	var encoder zapcore.Encoder

	isProduction := configuration.Application.Environment == "prod" ||
		configuration.Application.Environment == "production"

	zapConfig = zap.NewProductionConfig()
	zapConfig.Encoding = "json"
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	encoder = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.LevelKey = "level"
	zapConfig.EncoderConfig.CallerKey = "caller"
	zapConfig.EncoderConfig.StacktraceKey = "stacktrace"

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapConfig.Level,
	)

	Log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	zap.ReplaceGlobals(Log)
	Instance = NewZapLogger(Log)

	if isProduction {
		Log.Info("Logger initialized in production mode",
			zap.String("environment", configuration.Application.Environment),
			zap.String("encoding", "json"),
			zap.String("level", "info"),
		)
	} else {
		Log.Info("Logger initialized in development mode",
			zap.String("environment", configuration.Application.Environment),
			zap.String("encoding", "console"),
			zap.String("level", "debug"),
		)
	}
}
