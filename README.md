# deferr

Package deferr solves the pesky issue of not handling errors when
deferring io.Closer's.

Handling errors is important but some errors can be lost when
deferring a call to Close(). For some cases, the result of a Close()
may not matter. In many cases though, this information may be
important. For example, if a Close() fails on an open file handle,
it's possible that some of the data was not written.

Consider the following simple program:

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := sayHello(); err != nil {
		fmt.Printf("saying hello: %v\n", err)
	}
}

func sayHello() (err error) {
	f, err := os.Create("hello-world.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewBufferString("Hello, world!"))
	return err
}
```

If you check this against ```gometalinter``` you should see this:

```sh
$ gometalinter main.go
main.go:21:15:warning: error return value not checked (defer f.Close()) (errcheck)
```

This package can help fix this. It will do error checking and set the
error value if the close fails. You can give this error higher or
lower priority of other errors that might be returned based on the
call you make. ```Error()``` gives the ```Close()``` lower priority
and ```Closer()``` gives it higher priority. To fix our code, we
change the defer line to:

```go
defer deferr.Error(f, &err)
```

You'll now notice that ```gometalinter``` doesn't complain
anymore. Problem solved! Note that we pass the error by reference and
we use the a named return value for the error. This allows us to
modify the error in this package as the deferred function is running. 

It may not be feasible to do this everywhere, but it's certainly
worthwhile for cases where a failed close may be important. If you
just want to log a possible error message, you can use the ```Log()```
function instead. You can change the previous example to:

```go
deferr.Log(f, log.Printf)
// OR
deferr.Log(f, logrus.Errorf)
// etc.
```

instead of setting a error, it just logs the error using the given log
function.
