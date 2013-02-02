// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func checkNoError(t *testing.T, err error, format string) {
	if err != nil {
		t.Fatalf(format, err)
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
	err = SetOutput(out)
	checkNoError(t, err, "error while setting output to /dev/null file: %s")
	return out
}
func CleanOutput(t *testing.T, output *os.File) {
	checkNoError(t, output.Close(), "error while closing output file: %s")
}

func TestReadline(t *testing.T) {
	input := "Hello, world!"
	in := InitInput(t, input)
	defer CleanInput(t, in)
	err := SetInput(in)
	checkNoError(t, err, "error while setting input to temp file: %s")

	out := InitOutput(t)
	defer CleanOutput(t, out)

	line, err := Readline("> ")
	checkNoError(t, err, "error while reading first line: %s")
	if line != input {
		t.Error("%q expected (got %q)", input, line)
	}
	line, err = Readline("> ")
	if len(line) != 0 {
		t.Errorf("EOF expected (got %q)", line)
	}
	if err != io.EOF {
		t.Errorf("EOF expected (got %v)", err)
	}
}
