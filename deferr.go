// Package deferr solves the pesky issue of not handling errors when
// deferring io.Closer's.
//
// Handling errors is important but some errors can be lost when
// deferring a call to Close(). For some cases, the result of a
// Close() may not matter. In many cases though, this information may
// be important. For example, if a Close() fails on an open file
// handle, it's possible that some of the data was not written.
package deferr

import "io"

// Closer calls Close() on c. The error returned from Close() has
// preference and will overwrite any value pointed to by err.
func Closer(c io.Closer, err *error) {
	e := c.Close()
	if e != nil {
		*err = e
	}
}

// Error calls Close() on c. The error returned from Close() has a
// lower priority to err and will only be returned if the value
// pointed to by err is nil.
func Error(c io.Closer, err *error) {
	e := c.Close()
	if *err == nil {
		*err = e
	}
}

// LogFunc is the function definition for sending a log message. It is
// used by Log() to send a message to the logger when a Close() fails.
type LogFunc func(format string, args ...interface{})

// Log calls Close() on c. If the error is not nil, then the given
// LogFunc is called using the format "closing %v: %v" with the
// io.Closer and the error. If you want it to look pretty,
// implementing the fmt.Stringer interface would be useful. Otherwise,
// it's going to print the structure.
func Log(c io.Closer, f LogFunc) {
	if err := c.Close(); err != nil {
		f("closing %v: %v", c, err)
	}
}
