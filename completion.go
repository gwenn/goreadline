// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline

/*
#include <stdlib.h>
#include "goreadline.h"

extern char *goCompletionEntryFunction(char *text, int state);

static char *c_entry_func(const char *text, int state) { // cgo doesn't support const keyword
	return goCompletionEntryFunction((char *)text, state);
}

static char **my_attempted_completion_function(const char *text, int start, int end) {
	return rl_completion_matches(text, c_entry_func);
}

static void register_attempted_completion_function() {
	rl_attempted_completion_function = my_attempted_completion_function;
}
*/
import "C"

import (
	"unsafe"
)

//export goCompletionEntryFunction
func goCompletionEntryFunction(text *C.char, state C.int) *C.char {
	match := completionEntryFunction(C.GoString(text), int(state))
	if match == "" {
		return nil
	}
	return C.CString(match) // freed by readline
}

// CompletionEntryFunction is the generator function.
// It is called repeatedly, returning a string each time.
// The arguments to the generator function are text and state:
//  * text is the partial word to be completed.
//  * state is zero the first time the function is called, allowing the generator to perform any necessary initialization,
//    and a positive non-zero integer for each subsequent call.
// Usually the generator function computes the list of possible completions when state is zero,
// and returns them one at a time on subsequent calls.
// The generator function returns an empty string to the caller when there are no more possibilities left.
// If all completion entries share a common prefix, it is automatically appended to the current line.
type CompletionEntryFunction func(text string, state int) string

var completionEntryFunction CompletionEntryFunction

// SetCompletionEntryFunction registers the specified generator function.
// (See rl_attempted_completion_function http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX361)
func SetCompletionEntryFunction(f CompletionEntryFunction) {
	if f == nil {
		if completionEntryFunction != nil {
			C.rl_attempted_completion_function = nil
		}
	} else if completionEntryFunction == nil {
		C.register_attempted_completion_function()
	}
	completionEntryFunction = f
}

// If an application-specific completion function calls this function with a true value,
// Readline will not perform its default filename completion even if the application's completion function returns no matches.
// It should be call only by an application's completion function.
// (See rl_attempted_completion_over http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX369)
func SetAttemptedCompletionOver(b bool) {
	if b {
		C.rl_attempted_completion_over = 1
	} else {
		C.rl_attempted_completion_over = 0
	}
}

// SetCompleterWordBreakChars sets the list of characters that signal a break between words for completion.
// (See rl_completer_word_break_characters http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX354)
func SetCompleterWordBreakChars(s string) {
	cs := C.CString(s)
	C.free(unsafe.Pointer(C.rl_completer_word_break_characters))
	C.rl_completer_word_break_characters = cs
}

// CompleterWordBreakChars returns the list of characters that signal a break between words for completion.
// The default list is " \t\n\"\\'`@$><=;|&{(".
// (See rl_completer_word_break_characters http://cnswww.cns.cwru.edu/php/chet/readline/readline.html#IDX354)
func CompleterWordBreakChars() string {
	return C.GoString(C.rl_completer_word_break_characters)
}
