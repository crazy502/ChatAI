package observe

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"sync"

	"server/pkg/apperror"
	"server/pkg/requestmeta"
)

var (
	defaultLogger *slog.Logger
	loggerOnce    sync.Once
)

func Init() {
	loggerOnce.Do(func() {
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}))
		slog.SetDefault(defaultLogger)
	})
}

func Logger(ctx context.Context) *slog.Logger {
	Init()

	args := requestmeta.SlogArgs(ctx)
	if len(args) == 0 {
		return defaultLogger
	}

	return defaultLogger.With(args...)
}

func Info(ctx context.Context, msg string, attrs ...any) {
	Logger(ctx).Info(msg, attrs...)
}

func Warn(ctx context.Context, msg string, attrs ...any) {
	Logger(ctx).Warn(msg, attrs...)
}

func Error(ctx context.Context, msg string, err error, attrs ...any) {
	if err == nil {
		Logger(ctx).Error(msg, attrs...)
		return
	}

	appErr := apperror.From(err)
	logger := Logger(ctx)

	details := []any{
		"error_code", appErr.Code.Code(),
		"error_message", appErr.Message,
	}

	if appErr.Cause != nil {
		details = append(details, "cause", appErr.Cause.Error())
	}

	if len(appErr.Stack) > 0 {
		details = append(details, "stack", appErr.Stack)
	}

	fieldKeys := apperror.SortedFieldKeys(appErr)
	for _, key := range fieldKeys {
		details = append(details, key, appErr.Fields[key])
	}

	logger.Error(msg, append(details, attrs...)...)
}

func ErrorFields(err error) []any {
	appErr := apperror.From(err)
	if appErr == nil {
		return nil
	}

	fields := []any{
		"error_code", appErr.Code.Code(),
		"error_message", appErr.Message,
	}

	if appErr.Cause != nil {
		fields = append(fields, "cause", appErr.Cause.Error())
	}

	keys := apperror.SortedFieldKeys(appErr)
	sort.Strings(keys)
	for _, key := range keys {
		fields = append(fields, key, appErr.Fields[key])
	}

	return fields
}
