package errors

import (
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors/code"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
		want    *Error
	}{
		{
			name:    "create new error",
			code:    code.InvalidArgs,
			message: "test error",
			want: &Error{
				Status: Status{
					Code:    code.InvalidArgs,
					Message: "test error",
				},
			},
		},
		{
			name:    "create error with empty message",
			code:    code.InternalServer,
			message: "",
			want: &Error{
				Status: Status{
					Code:    code.InternalServer,
					Message: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.code, tt.message)
			if got.Code != tt.want.Code || got.Message != tt.want.Message {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *Error
		cause   error
		wantStr string
	}{
		{
			name:    "error without cause",
			err:     New(code.InvalidArgs, "test error"),
			wantStr: "error: message = test error  cause = <nil>",
		},
		{
			name:    "error with cause",
			err:     New(code.InvalidArgs, "test error").WithCause(fmt.Errorf("underlying error")),
			wantStr: "error: message = test error  cause = underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantStr {
				t.Errorf("Error() = %v, want %v", got, tt.wantStr)
			}
		})
	}
}

func TestError_WithCause(t *testing.T) {
	cause := fmt.Errorf("test cause")
	tests := []struct {
		name string
		err  *Error
	}{
		{
			name: "add cause to error",
			err:  New(code.InvalidArgs, "test error"),
		},
		{
			name: "override existing cause",
			err:  New(code.InternalServer, "test error").WithCause(fmt.Errorf("old cause")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.WithCause(cause)
			if got == tt.err {
				t.Error("WithCause() returned same error instance, want new instance")
			}
			if got.cause != cause {
				t.Errorf("WithCause() cause = %v, want %v", got.cause, cause)
			}
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
	}{
		{
			name: "clone error without cause",
			err:  New(code.InvalidArgs, "test error"),
		},
		{
			name: "clone error with cause",
			err:  New(code.InternalServer, "test error").WithCause(fmt.Errorf("test cause")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clone(tt.err)
			if got == tt.err {
				t.Error("Clone() returned same error instance, want new instance")
			}
			if got.Code != tt.err.Code {
				t.Errorf("Clone() code = %v, want %v", got.Code, tt.err.Code)
			}
			if got.Message != tt.err.Message {
				t.Errorf("Clone() message = %v, want %v", got.Message, tt.err.Message)
			}
			if got.cause != tt.err.cause {
				t.Errorf("Clone() cause = %v, want %v", got.cause, tt.err.cause)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		wantCode int
	}{
		{
			name:     "InvalidArgs error",
			err:      InvalidArgs,
			wantCode: code.InvalidArgs,
		},
		{
			name:     "InsufficientBalance error",
			err:      InsufficientBalance,
			wantCode: code.InvalidArgs,
		},
		{
			name:     "RecordNotFound error",
			err:      RecordNotFound,
			wantCode: code.NotFound,
		},
		{
			name:     "InternalDB error",
			err:      InternalDB,
			wantCode: code.InternalServer,
		},
		{
			name:     "InternalServer error",
			err:      InternalServer,
			wantCode: code.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("%s code = %v, want %v", tt.name, tt.err.Code, tt.wantCode)
			}
		})
	}
}
