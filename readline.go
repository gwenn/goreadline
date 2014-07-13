// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package readline provides access to the editline/readline C library.
package readline

/*
#cgo readline LDFLAGS: -lreadline -lhistory
#cgo !readline LDFLAGS: -ledit

#include <stdio.h>
#include <stdlib.h>
//#include <readline/readline.h>
#include <editline/readline.h>
*/
import "C"

import (
	"os"
	"syscall"
	"unsafe"
)

// ReadLine prints a prompt and then reads and returns a single line of text from the user.
// If ReadLine encounters an EOF while reading the line, and the line is empty at that point, then true is returned.
// Otherwise, the line is ended just as if a newline had been typed.
// (See readline http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX190)
func ReadLine(prompt string) (string, bool) {
	var cprompt *C.char
	if len(prompt) != 0 {
		cprompt = C.CString(prompt)
	}
	cline := C.readline(cprompt)
	if cprompt != nil {
		C.free(unsafe.Pointer(cprompt))
	}
	if cline == nil {
		return "", true
	}
	line := C.GoString(cline)
	C.free(unsafe.Pointer(cline))
	return line, false
}

// Buffer returns the line gathered so far.
// (See rl_line_buffer http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX192)
// TODO Validate String versus []byte
func Buffer() string {
	return C.GoString(C.rl_line_buffer)
}

// Point returns the offset of the current cursor position in Buffer (the point).
// (See rl_point http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX192)
func Point() int { // int32
	return int(C.rl_point)
}

// setInput changes the default input stream (stdin by default)
// (See rl_instream http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX209)
func setInput(in *os.File) error {
	return setStream(in, &C.rl_instream, "r", syscall.Stdin)
}

// setOutput changes the default output stream (stdout by default)
// (See rl_outstream http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX210)
func setOutput(out *os.File) error {
	return setStream(out, &C.rl_outstream, "w", syscall.Stdout)
}

func setStream(f *os.File, cstream **C.FILE, mode string, def int) error {
	fd := def
	if f != nil {
		fd = int(f.Fd())
	}
	cfd := def
	if *cstream != nil {
		cfd = int(C.fileno(*cstream))
	}
	if fd == cfd { // TODO Validate
		return nil
	}
	cname := C.CString(f.Name())
	cmode := C.CString(mode)
	cf, err := C.fopen(cname, cmode) // TODO How to close?
	C.free(unsafe.Pointer(cname))
	C.free(unsafe.Pointer(cmode))
	if err != nil {
		return err
	}
	*cstream = cf
	return nil
}

// Initialize or re-initialize Readline's internal state. It's not strictly necessary to call this; Readline() calls it before reading any input.
// (See rl_initialize http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX316)
func Initialize() error {
	err := C.rl_initialize()
	if err != 0 {
		return syscall.Errno(err)
	}
	return nil
}

// LibraryVersion returns the version number of this revision of the library.
// (See rl_library_version http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX214)
func LibraryVersion() string {
	return C.GoString(C.rl_library_version)
}

// Version returns an integer encoding the current version of the library.
// (See rl_readline_version http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX214)
func Version() int32 {
	return int32(C.rl_readline_version)
}

// Name is set to a unique name by each application using Readline. The value allows conditional parsing of the inputrc file.
// (See rl_readline_name http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX218)
func Name() string {
	return C.GoString(C.rl_readline_name)
}

// SetName set to a unique name by each application using Readline. The value allows conditional parsing of the inputrc file.
// (See rl_readline_name http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX218)
func SetName(name string) {
	cname := C.CString(name)
	/*if Name() != "" {
		C.free(unsafe.Pointer(C.rl_readline_name))
	}*/
	C.rl_readline_name = cname
}
