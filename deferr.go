// Package deferr solves the pesky issue of not handling errors when
// deferring io.Closer's.
//
// Handling errors is important but some errors can be lost when
// deferring a call to Close(). For some cases, the result of a
// Close() may not matter. In many cases though, this information may
// be important. For example, if a Close() fails on an open file handle, it's possible that some of the data
package deferr

import "io"

// Closer calls Close() on c. The error returned from it has
// preference and will overwrite any value pointed to by err.
func Closer(c io.Closer, err *error) {
	e := c.Close()
	if e != nil {
		*err = e
	}
}

// Error calls Close() on c. The error returned from it has a lower
// preference to err and will only be returned if the value pointed to
// by err is nil.
func Error(c io.Closer, err *error) {
	e := c.Close()
	if *err == nil {
		*err = e
	}
}
