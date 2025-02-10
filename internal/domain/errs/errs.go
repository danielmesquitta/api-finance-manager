package errs

import (
	"encoding/json"
	"fmt"
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

// NewErr creates a new Err instance from either an error or a string,
// and sets the Type flag to unknown. This is useful when you want to
// create an error that is not expected to happen, and you want to
// log it with stack tracing.
func New(err any) *Err {
	return NewWithType(err, ErrTypeUnknown)
}

func NewWithType(err any, t Type) *Err {
	if err == nil {
		return nil
	}
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

	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return &Err{
				Message:    fmt.Sprintf("unsupported err type %T: %+v", v, err),
				StackTrace: string(debug.Stack()),
				Type:       t,
			}
		}
		return &Err{
			Message:    string(jsonData),
			StackTrace: string(debug.Stack()),
			Type:       t,
		}
	}
}

func (e *Err) Error() string {
	return e.Message
}

var _ error = (*Err)(nil)
