package deferr

import (
	"errors"
	"reflect"
	"testing"
)

type errCloser struct {
	err error
}

func (c *errCloser) Close() error {
	return c.err
}

func TestCloser(t *testing.T) {
	tests := []struct {
		name  string
		close error
		err   error
		exp   error
	}{
		{
			name:  "closer error, err error",
			close: errors.New("close"),
			err:   errors.New("error"),
			exp:   errors.New("close"),
		},
		{
			name:  "closer error, no err error",
			close: errors.New("close"),
			err:   nil,
			exp:   errors.New("close"),
		},
		{
			name:  "no closer error, err error",
			close: nil,
			err:   errors.New("error"),
			exp:   errors.New("error"),
		},
		{
			name:  "no closer error, no err error",
			close: nil,
			err:   nil,
			exp:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Closer(&errCloser{test.close}, &test.err)
			if !reflect.DeepEqual(test.exp, test.err) {
				t.Errorf("Closer() changed err to '%v', expected '%v'",
					test.err, test.exp)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name  string
		close error
		err   error
		exp   error
	}{
		{
			name:  "closer error, err error",
			close: errors.New("close"),
			err:   errors.New("error"),
			exp:   errors.New("error"),
		},
		{
			name:  "closer error, no err error",
			close: errors.New("close"),
			err:   nil,
			exp:   errors.New("close"),
		},
		{
			name:  "no closer error, err error",
			close: nil,
			err:   errors.New("error"),
			exp:   errors.New("error"),
		},
		{
			name:  "no closer error, no err error",
			close: nil,
			err:   nil,
			exp:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Error(&errCloser{test.close}, &test.err)
			if !reflect.DeepEqual(test.exp, test.err) {
				t.Errorf("Error() changed err to '%v', expected '%v'",
					test.err, test.exp)
			}
		})
	}
}
