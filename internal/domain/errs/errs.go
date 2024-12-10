package errs

import (
	"encoding/json"
	"log"
	"runtime/debug"
)

type Err struct {
	Message    string `json:"message,omitempty"`
	StackTrace string `json:"stack_trace,omitempty"`
	Type       Type   `json:"type,omitempty"`
}

type Type string

const (
	ErrTypeUnknown      Type = "unknown"
	ErrTypeNotFound     Type = "not_found"
	ErrTypeUnauthorized Type = "unauthorized"
	ErrTypeForbidden    Type = "forbidden"
	ErrTypeValidation   Type = "validation_error"
)

func newErr(err any, t Type) *Err {
	switch v := err.(type) {
	case *Err:
		return v
	case error:
		return &Err{
			Message:    v.Error(),
			StackTrace: string(debug.Stack()),
			Type:       t,
		}
	case string:
		return &Err{
			Message:    v,
			StackTrace: string(debug.Stack()),
			Type:       t,
		}
	case []byte:
		return &Err{
			Message:    string(v),
			StackTrace: string(debug.Stack()),
			Type:       t,
		}
	case nil:
		return nil
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("trying to create an Err with an unsupported type %T\n%+v", v, err)
			return nil
		}
		return &Err{
			Message:    string(jsonData),
			StackTrace: string(debug.Stack()),
			Type:       t,
		}
	}
}

// NewErr creates a new Err instance from either an error or a string,
// and sets the Type flag to unknown. This is useful when you want to
// create an error that is not expected to happen, and you want to
// log it with stack tracing.
func New(err any) *Err {
	return newErr(err, ErrTypeUnknown)
}

func NewWithType(err any, t Type) *Err {
	return newErr(err, t)
}

func (e *Err) Error() string {
	return e.Message
}

var (
	ErrUserNotFound = newErr(
		"user not found",
		ErrTypeNotFound,
	)
	ErrUnauthorized = newErr(
		"unauthorized",
		ErrTypeUnauthorized,
	)
	ErrAccountsAlreadyRegistered = newErr(
		"accounts already registered",
		ErrTypeForbidden,
	)
	ErrInvalidRetirementAge = newErr(
		"age should be less than retirement age",
		ErrTypeValidation,
	)
	ErrInvalidLifeExpectance = newErr(
		"life expectance should be greater than retirement age",
		ErrTypeValidation,
	)
)

var _ error = (*Err)(nil)
