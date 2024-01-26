package deploykit

import (
	"errors"
	"fmt"
)

const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("deploykit error: code=%s message=%s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	switch {
	case err == nil:
		return ""
	case errors.As(err, &e):
		return e.Code
	default:
		return EINTERNAL
	}
}

func ErrorMessage(err error) string {
	var e *Error
	switch {
	case err == nil:
		return ""
	case errors.As(err, &e):
		return e.Message
	default:
		return "Internal error."
	}
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

type DeploymentConfig struct {
	App struct {
		Name string `toml:"name"`
	} `toml:"app"`
}

type NetworkRepository interface {
	Create(name string)
}

type NetworkService struct {
	repo NetworkRepository
}

func (s *NetworkService) Create(name string) {
	s.repo.Create(name)
}
