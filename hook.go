// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build readline

package readline

/*
#cgo LDFLAGS: -lreadline -lhistory

#include <readline/readline.h>
*/
import "C"

import (
	"os"
	"os/signal"
	"syscall"
)

// https://code.google.com/p/go/issues/detail?id=4216
func init() {
	// Program received signal SIGSEGV, Segmentation fault.
	// rl_sigwinch_handler (sig=-136463680) at /tmp/buildd/readline6-6.2+dfsg/signals.c:267
	// 267	  RL_UNSETSTATE(RL_STATE_SIGHANDLER);
	C.rl_catch_sigwinch = 0
	resized := make(chan os.Signal, 1)
	go func() {
		for _ = range resized {
			C.rl_resize_terminal()
		}
	}()
	signal.Notify(resized, syscall.SIGWINCH)
}
