// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package readline provides access to the readline C library.
package readline

/*
#cgo LDFLAGS: -lreadline -lhistory

#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
*/
import "C"

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// ReadLine prints a prompt and then reads and returns a single line of text from the user.
// If ReadLine encounters an EOF while reading the line, and the line is empty at that point, then an io.EOF error is returned.
// Otherwise, the line is ended just as if a newline had been typed.
// (See readline http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX190)
func ReadLine(prompt string) (string, error) {
	var cprompt *C.char
	if len(prompt) != 0 {
		cprompt = C.CString(prompt)
	}
	cline := C.readline(cprompt)
	if cprompt != nil {
		C.free(unsafe.Pointer(cprompt))
	}
	if cline == nil {
		return "", io.EOF
	}
	line := C.GoString(cline)
	C.free(unsafe.Pointer(cline))
	return line, nil
}

// Buffer returns the line gathered so far.
// (See rl_line_buffer http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX192)
func Buffer() string {
	return C.GoString(C.rl_line_buffer)
}

// Point returns the offset of the current cursor position in Buffer (the point).
// (See rl_point http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX192)
func Point() int {
	return int(C.rl_point)
}

// SetInput changes the default input stream (stdin by default)
// (See rl_instream http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX209)
func SetInput(in *os.File) error {
	return setStream(in, &C.rl_instream, "r", syscall.Stdin)
}

// SetOutput changes the default output stream (stdout by default)
// (See rl_outstream http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX210)
func SetOutput(out *os.File) error {
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
