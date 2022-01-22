package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func init() {
	funcs = append(funcs, printWithoutUnwrap, checkErrorType, isAsEOFError, testGoString)
}

func printWithoutUnwrap() {
	fmt.Printf("不带unwrap的打印:%s", createErrorList())
}

func checkErrorType() {
	err := createErrorList()
	if merr, ok := err.(*multierror.Error); ok {
		fmt.Printf("类型断言:%s", merr)
	}
}

func isAsEOFError() {
	err := createErrorList()
	var eof eofError
	if errors.As(err, &eof) {
		fmt.Printf("是eofError")
	}
}

func testGoString() {
	err := createErrorList()
	merr := err.(*multierror.Error)
	fmt.Println(merr.GoString())
}

type eofError struct{}

func (e eofError) Error() string {
	return "eof"
}

type readError struct{}

func (e *readError) Error() string {
	return "read error"
}

func createErrorList() (err error) {
	err = multierror.Append(err, errors.New("error1"))
	err = multierror.Append(err, errors.New("error2"))
	err = multierror.Append(err, errors.New("error3"))
	var eof eofError
	err = multierror.Append(err, eof)
	err = multierror.Append(err, &readError{})

	return
}
