package requestmeta

import (
	"context"
	"fmt"
	"sort"
	"strconv"
)

const (
	FieldRequestID = "request_id"
	FieldUserID    = "user_id"
	FieldUserName  = "user_name"
	FieldSessionID = "session_id"
	FieldIsAdmin   = "is_admin"
)

type fieldsKey struct{}

type Fields map[string]any

func WithField(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if key == "" || value == nil {
		return ctx
	}

	current := Get(ctx)
	next := make(Fields, len(current)+1)
	for currentKey, currentValue := range current {
		next[currentKey] = currentValue
	}
	next[key] = value
	return context.WithValue(ctx, fieldsKey{}, next)
}

func WithFields(ctx context.Context, fields Fields) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	for key, value := range fields {
		ctx = WithField(ctx, key, value)
	}

	return ctx
}

func Get(ctx context.Context) Fields {
	if ctx == nil {
		return Fields{}
	}

	fields, ok := ctx.Value(fieldsKey{}).(Fields)
	if !ok || len(fields) == 0 {
		return Fields{}
	}

	result := make(Fields, len(fields))
	for key, value := range fields {
		result[key] = value
	}
	return result
}

func String(ctx context.Context, key string) string {
	value, ok := Get(ctx)[key]
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func Int64(ctx context.Context, key string) int64 {
	value, ok := Get(ctx)[key]
	if !ok || value == nil {
		return 0
	}

	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int32:
		return int64(typed)
	case int64:
		return typed
	case uint:
		return int64(typed)
	case uint32:
		return int64(typed)
	case uint64:
		return int64(typed)
	case float64:
		return int64(typed)
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func Bool(ctx context.Context, key string) bool {
	value, ok := Get(ctx)[key]
	if !ok || value == nil {
		return false
	}

	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, err := strconv.ParseBool(typed)
		if err != nil {
			return false
		}
		return parsed
	default:
		return false
	}
}

func SlogArgs(ctx context.Context) []any {
	fields := Get(ctx)
	if len(fields) == 0 {
		return nil
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	args := make([]any, 0, len(keys)*2)
	for _, key := range keys {
		args = append(args, key, fields[key])
	}
	return args
}
