// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline

/*
#include <stdlib.h>
#include <readline/history.h>
*/
import "C"

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

// UsingHistory begins a session in which the history functions might be used. This initializes the interactive variables.
// (See using_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX2)
func UsingHistory() {
	C.using_history()
}

// AddHistory places string at the end of the history list.
// Blank lines are discarded.
// (See add_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX5)
func AddHistory(line string) {
	if len(line) == 0 || len(strings.TrimSpace(line)) == 0 {
		return
	}
	if unicode.IsSpace(rune(line[0])) { // ignorespace
		return
	}
	if prev, err := GetHistory(-1); err == nil && prev == line { // ignore consecutive dups
		return
	}
	cline := C.CString(line)
	C.add_history(cline)
	C.free(unsafe.Pointer(cline))
}

// ReadHistory adds the content of filename to the history list, a line at a time.
// If filename is "", then read from '~/.history'.
// (See read_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX27)
func ReadHistory(filename string) (bool, error) {
	var cfilename *C.char
	if len(filename) != 0 {
		cfilename = C.CString(filename)
	}
	err := C.read_history(cfilename)
	if cfilename != nil {
		C.free(unsafe.Pointer(cfilename))
	}
	if err != 0 {
		e := syscall.Errno(err)
		if e == syscall.ENOENT { // ignored when the file doesn't exist.
			return false, nil
		}
		return false, e
	}
	return true, nil
}

// WriteHistory writes the current history to filename, overwriting filename if necessary.
// If filename is "", then write the history list to `~/.history'.
// (See write_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX29)
func WriteHistory(filename string) error {
	var cfilename *C.char
	if len(filename) != 0 {
		cfilename = C.CString(filename)
	}
	err := C.write_history(cfilename)
	if cfilename != nil {
		C.free(unsafe.Pointer(cfilename))
	}
	if err != 0 {
		return syscall.Errno(err)
	}
	return nil
}

// AppendHistory appends the last nelements of the history list to filename.
// If filename is "", then append to `~/.history'.
// FIXME seems to be unsupported by editline library.
// (See append_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX30)
func AppendHistory(nelements int, filename string) error {
	if HistoryLength() == 0 {
		return nil
	}
	// Checks if the file exists. If not, creates it.
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
	}
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	var cfilename *C.char
	if len(filename) != 0 {
		cfilename = C.CString(filename)
	}
	cerr := C.append_history(C.int(nelements), cfilename)
	if cfilename != nil {
		C.free(unsafe.Pointer(cfilename))
	}
	if cerr != 0 {
		return syscall.Errno(cerr)
	}
	return nil
}

// TruncateHistoryFile truncates the history file filename, leaving only the last nlines lines.
// If filename is "", then `~/.history' is truncated.
// (See history_truncate_file http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX31)
func TruncateHistoryFile(filename string, nlines int) error {
	var cfilename *C.char
	if len(filename) != 0 {
		cfilename = C.CString(filename)
	}
	err := C.history_truncate_file(cfilename, C.int(nlines))
	if cfilename != nil {
		C.free(unsafe.Pointer(cfilename))
	}
	if err != 0 {
		return syscall.Errno(err)
	}
	return nil
}

// ClearHistory clears the history list by deleting all the entries.
// (See clear_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX10)
func ClearHistory() {
	C.clear_history()
}

// StifleHistory cuts off the history list, remembering only the last max entries.
// (See stifle_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX11)
func StifleHistory(max int) {
	C.stifle_history(C.int(max))
}

// UnstifleHistory stops stifling the history.
// This returns the previously-set maximum number of history entries (as set by StifleHistory()).
// The value is positive if the history was stifled, negative if it wasn't.
// (See unstifle_history http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX12)
func UnstifleHistory() int {
	return int(C.unstifle_history())
}

// IsHistoryStifled says if the history is stifled.
// (See history_is_stifled http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX13)
func IsHistoryStifled() bool {
	return C.history_is_stifled() != 0
}

// HistoryLength returns the number of entries currently stored in the history list.
// (See history_length http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX37)
func HistoryLength() int {
	return int(C.history_length)
}

// HistoryBase returns the logical offset of the first entry in the history list.
// (See history_base http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX36)
func HistoryBase() int {
	return int(C.history_base)
}

/*
// Return a nil terminated array of HIST_ENTRY * which is the current input history.
// Element 0 of this list is the beginning of time. If there is no history, return nil.
// (See history_list http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX17)
func HistoryList(offset int) string {
	C.history_list()
}
*/

// GetHistory returns the history entry at position index, starting from 0.
// If there is no entry there, or if index is greater than the history length, return an error.
// (See history_get http://cnswww.cns.cwru.edu/php/chet/readline/history.html#IDX17)
func GetHistory(index int) (string, error) {
	length := HistoryLength()
	if index < 0 {
		index += length
	}
	if index < 0 || index >= length {
		return "", fmt.Errorf("invalid index", index)
	}
	index += HistoryBase() // TODO
	entry := C.history_get(C.int(index))
	if entry == nil {
		return "", fmt.Errorf("invalid index", index)
	}
	return C.GoString(entry.line), nil
}
