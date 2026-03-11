// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package tangle is the TANGLE program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/web/tangle.web
package tangle // modernc.org/knuth/tangle

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program TANGLE( web_file, change_file, Pascal_file, pool);
// web_file:text_file; {primary input}
// change_file:text_file; {updates}
// Pascal_file: text_file;
// pool: text_file;

// Main executes the tangle program using the supplied arguments.
func Main(webFile, changeFile io.Reader, pascalFile, poolFile, stdout, stderr io.Writer) (mainErr error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// ok
		case signal:
			switch {
			case mainErr == nil:
				mainErr = fmt.Errorf("aborted")
			default:
				mainErr = fmt.Errorf("aborted: %v", mainErr)
			}
		case knuth.Error:
			mainErr = x
		default:
			mainErr = fmt.Errorf("PANIC %T: %[1]v, error: %v\n%s", x, mainErr, debug.Stack())
		}
	}()

	(&prg{
		webFile:    knuth.NewTextFile(webFile, nil, nil),
		changeFile: knuth.NewTextFile(changeFile, nil, nil),
		pascalFile: knuth.NewTextFile(nil, pascalFile, nil),
		pool:       knuth.NewTextFile(nil, poolFile, nil),
		stdout:     knuth.NewTextFile(nil, stdout, nil),
		stderr:     knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
