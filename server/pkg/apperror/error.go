package apperror

import (
	"errors"
	"fmt"
	"runtime"
	"sort"

	"server/pkg/code"
)

const stackDepth = 32

type Frame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Error struct {
	Code    code.Code      `json:"code"`
	Message string         `json:"message"`
	Cause   error          `json:"-"`
	Stack   []Frame        `json:"stack,omitempty"`
	Fields  map[string]any `json:"fields,omitempty"`
}

func New(resultCode code.Code, message string) *Error {
	if message == "" {
		message = resultCode.Msg()
	}

	return &Error{
		Code:    resultCode,
		Message: message,
		Stack:   captureStack(2),
	}
}

func Wrap(resultCode code.Code, cause error, message string) *Error {
	if cause == nil {
		return nil
	}

	if message == "" {
		message = resultCode.Msg()
	}

	return &Error{
		Code:    resultCode,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(2),
		Fields:  cloneFields(FieldsOf(cause)),
	}
}

func Panic(recovered any) *Error {
	return &Error{
		Code:    code.CodeServerBusy,
		Message: fmt.Sprintf("panic recovered: %v", recovered),
		Cause:   fmt.Errorf("%v", recovered),
		Stack:   captureStack(3),
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Cause == nil {
		return fmt.Sprintf("%d:%s", e.Code.Code(), e.Message)
	}

	return fmt.Sprintf("%d:%s: %v", e.Code.Code(), e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *Error) WithField(key string, value any) *Error {
	if e == nil || key == "" || value == nil {
		return e
	}

	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[key] = value
	return e
}

func (e *Error) WithFields(fields map[string]any) *Error {
	if e == nil || len(fields) == 0 {
		return e
	}

	if e.Fields == nil {
		e.Fields = make(map[string]any, len(fields))
	}

	for key, value := range fields {
		if key == "" || value == nil {
			continue
		}
		e.Fields[key] = value
	}

	return e
}

func (e *Error) Clone() *Error {
	if e == nil {
		return nil
	}

	stack := make([]Frame, len(e.Stack))
	copy(stack, e.Stack)

	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Cause:   e.Cause,
		Stack:   stack,
		Fields:  cloneFields(e.Fields),
	}
}

func From(err error) *Error {
	if err == nil {
		return nil
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}

	return &Error{
		Code:    code.CodeServerBusy,
		Message: code.CodeServerBusy.Msg(),
		Cause:   err,
		Stack:   captureStack(2),
	}
}

func CodeOf(err error) code.Code {
	if err == nil {
		return code.CodeSuccess
	}
	return From(err).Code
}

func MessageOf(err error) string {
	if err == nil {
		return code.CodeSuccess.Msg()
	}
	return From(err).Message
}

func StackOf(err error) []Frame {
	if err == nil {
		return nil
	}

	appErr := From(err)
	stack := make([]Frame, len(appErr.Stack))
	copy(stack, appErr.Stack)
	return stack
}

func FieldsOf(err error) map[string]any {
	if err == nil {
		return nil
	}
	return cloneFields(From(err).Fields)
}

func SortedFieldKeys(err error) []string {
	fields := FieldsOf(err)
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func captureStack(skip int) []Frame {
	pcs := make([]uintptr, stackDepth)
	count := runtime.Callers(skip+2, pcs)
	if count == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:count])
	result := make([]Frame, 0, count)

	for {
		frame, more := frames.Next()
		result = append(result, Frame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	return result
}

func cloneFields(input map[string]any) map[string]any {
	if len(input) == 0 {
		return nil
	}

	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
