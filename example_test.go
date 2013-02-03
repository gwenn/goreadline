// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package readline_test

import (
	"github.com/gwenn/goreadline"
	"io"
	"os/user"
	"path"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Example() {
	readline.UsingHistory()
	user, err := user.Current()
	check(err)
	history := path.Join(user.HomeDir, ".goreadline_test.rc")
	// readline.StifleHistory(100) // to limit memory usage
	_, err = readline.ReadHistory(history)
	check(err)
	for {
		line, err := readline.ReadLine("> ")
		if err == io.EOF {
			println()
			break
		}
		check(err)

		// ...

		readline.AddHistory(line)
	}
	// TODO save history in a deferred block?
	readline.StifleHistory(100) // to limit disk usage
	err = readline.WriteHistory(history)
	// err = readline.AppendHistory(100, history) // for multi-sessions handling
	check(err)
}

func ExampleSetCompletionEntryFunction() {
	readline.SetCompletionEntryFunction(func(text string, state int) string {
		if state == 0 {
			return text + "s"
		}
		return ""
	})
	for {
		line, err := readline.ReadLine("> ")
		if err == io.EOF {
			println()
			break
		}
		check(err)

		// ...
		println(line)
	}
}
