// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gwenn/goreadline"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var comp_entries []string = make([]string, 0, 10)

func completion(text string, state int) string {
	if state == 0 {
		comp_entries = comp_entries[:0]
		_, err := words.Seek(0, 0)
		check(err)
		scanner := bufio.NewScanner(words)
		for scanner.Scan() {
			word := scanner.Text()
			if strings.HasPrefix(word, text) {
				comp_entries = append(comp_entries, word)
			}
		}
		check(scanner.Err())
	}
	if state < len(comp_entries) {
		return comp_entries[state]
	}
	return ""
}

var words *os.File

func main() {
	var err error
	words, err = os.Open("/usr/share/dict/words")
	check(err)
	defer words.Close()
	readline.SetCompletionEntryFunction(completion)
	for {
		line, eof := readline.ReadLine("> ")
		if eof {
			println()
			break
		}

		fmt.Println(line)
	}
}
