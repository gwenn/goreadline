// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline

import (
	"io/ioutil"
	"os"
	"testing"
)

func initHistory(t *testing.T) *os.File {
	history, err := ioutil.TempFile("", ".goreadline_history")
	checkNoError(t, err, "error while creating history temp file: %s")
	return history
}
func cleanHistory(t *testing.T, history *os.File) {
	checkNoError(t, history.Close(), "error while closing history temp file: %s")
	checkNoError(t, os.Remove(history.Name()), "error while removing history temp file: %s")
	ClearHistory()
}

func TestAddHistory(t *testing.T) {
	UsingHistory()
	AddHistory("") // empty line ignored
	assertHistoryLength(t, 0)
	AddHistory(" \t") // blank line ignored
	assertHistoryLength(t, 0)
	AddHistory(" line") // line starting with space ignored
	assertHistoryLength(t, 0)

	AddHistory("line")
	assertHistoryLength(t, 1)
	AddHistory("line") // consecutive duplicates ignored
	assertHistoryLength(t, 1)
	assertHistoryLength(t, 1)

	history := initHistory(t)
	defer cleanHistory(t, history)

	err := WriteHistory(history.Name())
	checkNoError(t, err, "error while writting history: %s")

	_, err = ReadHistory(history.Name())
	checkNoError(t, err, "error while reading history: %s")

	// TODO list entries
}

func TestClearHistory(t *testing.T) {
	UsingHistory()
	AddHistory("line")
	assertHistoryLength(t, 1)
	ClearHistory()
	assertHistoryLength(t, 0)
}

func TestStifleHistory(t *testing.T) {
	UsingHistory()
	AddHistory("line1")
	AddHistory("line2")
	assertHistoryLength(t, 2)
	assert(t, "history is not stifled by default", !IsHistoryStifled())
	StifleHistory(1)
	assert(t, "history must be stifled now", IsHistoryStifled())
	assertHistoryLength(t, 1)
	AddHistory("line3")
	AddHistory("line4")
	assertHistoryLength(t, 1)
	assert(t, "msg", 1 == UnstifleHistory())
	assert(t, "history must not be stifled now", !IsHistoryStifled())
}

func assertHistoryLength(t *testing.T, expected int32) {
	actual := HistoryLength()
	if expected != actual {
		t.Errorf("expecting %d line(s) in history but found %d", expected, actual)
	}
}

func assert(t *testing.T, msg string, actual bool) {
	if !actual {
		t.Error(msg)
	}
}
