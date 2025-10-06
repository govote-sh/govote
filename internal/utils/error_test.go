package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrMsg_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ErrMsg
		want string
	}{
		{
			name: "with HTTP status code",
			err: ErrMsg{
				Err:            errors.New("not found"),
				HTTPStatusCode: 404,
			},
			want: "HTTP 404: not found",
		},
		{
			name: "with HTTP 500",
			err: ErrMsg{
				Err:            errors.New("internal server error"),
				HTTPStatusCode: 500,
			},
			want: "HTTP 500: internal server error",
		},
		{
			name: "without HTTP status code",
			err: ErrMsg{
				Err:            errors.New("something went wrong"),
				HTTPStatusCode: 0,
			},
			want: "something went wrong",
		},
		{
			name: "with negative status code",
			err: ErrMsg{
				Err:            errors.New("network error"),
				HTTPStatusCode: -1,
			},
			want: "network error",
		},
		{
			name: "with HTTP 200 (unusual but valid)",
			err: ErrMsg{
				Err:            errors.New("unexpected success"),
				HTTPStatusCode: 200,
			},
			want: "HTTP 200: unexpected success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}
