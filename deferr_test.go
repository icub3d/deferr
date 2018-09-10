package deferr

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type errCloser struct {
	err error
}

func (c *errCloser) String() string {
	return "err closer"
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

func TestLog(t *testing.T) {
	tests := []struct {
		name   string
		close  error
		exp    string
		logged bool
	}{
		{
			name:   "error",
			close:  errors.New("close"),
			exp:    "closing err closer: close",
			logged: true,
		},
		{
			name:   "no err",
			logged: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logged := false
			f := func(format string, args ...interface{}) {
				buf.WriteString(fmt.Sprintf(format, args...))
				logged = true
			}
			Log(&errCloser{test.close}, f)
			if logged != test.logged {
				t.Errorf("expected a message to be logged, but one wasn't")
			}
			if test.exp != "" && !strings.Contains(buf.String(), test.exp) {
				t.Errorf("didn't find '%v' in '%v'", test.exp, buf.String())
			}
		})
	}
}
