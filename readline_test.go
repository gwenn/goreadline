// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
)

func checkNoError(t *testing.T, err error, format string) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("\n%s:%d: %s", path.Base(file), line, fmt.Sprintf(format, err))
	}
}

func InitInput(t *testing.T, input string) *os.File {
	in, err := ioutil.TempFile("", "testReadline")
	checkNoError(t, err, "error while creating input temp file: %s")
	_, err = in.WriteString(input)
	checkNoError(t, err, "error while creating input temp file: %s")
	_, err = in.WriteString("\n")
	checkNoError(t, err, "error while creating input temp file: %s")
	err = in.Sync()
	checkNoError(t, err, "error while syncing input temp file: %s")
	return in
}
func CleanInput(t *testing.T, input *os.File) {
	checkNoError(t, input.Close(), "error while closing input temp file: %s")
	checkNoError(t, os.Remove(input.Name()), "error while removing input temp file: %s")
}
func InitOutput(t *testing.T) *os.File {
	out, err := os.Open("/dev/null")
	checkNoError(t, err, "error while opening /dev/null file: %s")
	err = setOutput(out)
	checkNoError(t, err, "error while setting output to /dev/null file: %s")
	return out
}
func CleanOutput(t *testing.T, output *os.File) {
	checkNoError(t, output.Close(), "error while closing output file: %s")
}

func TestReadLine(t *testing.T) {
	input := "Hello, world!"
	in := InitInput(t, input)
	defer CleanInput(t, in)
	err := setInput(in)
	checkNoError(t, err, "error while setting input to temp file: %s")

	out := InitOutput(t)
	defer CleanOutput(t, out)

	// need by editline
	lib := LibraryVersion()
	fmt.Printf("Version: %x\n", Version())
	println("LibraryVersion:", LibraryVersion())
	if strings.HasPrefix(lib, "EditLine") {
		err = Initialize()
		checkNoError(t, err, "error while initializing: %s")
	}

	line, eof := ReadLine("> ")
	if eof {
		t.Error("unexpected EOF")
	}
	if line != input {
		t.Errorf("%q expected (got %q)", input, line)
	}
	line, eof = ReadLine("> ")
	if len(line) != 0 {
		t.Errorf("EOF expected (got %q)", line)
	}
	if !eof {
		t.Error("EOF expected")
	}
}

func TestName(t *testing.T) {
	assert.T(t, "" == Name() || "other" == Name())
	SetName("goreadline")
	assert.Equal(t, "goreadline", Name())
}
