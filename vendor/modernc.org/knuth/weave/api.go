// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package weave is the WEAVE program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/web/weave.web
package weave // modernc.org/knuth/weave

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program WEAVE( web_file, change_file, tex_file);
// web_file:text_file; {primary input}
// change_file:text_file; {updates}
// tex_file: text_file;

// Main executes the weave program using the supplied arguments.
func Main(webFile io.ReadSeeker, changeFile io.Reader, texFile, stdout, stderr io.Writer) (mainErr error) {
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
		texFile:    knuth.NewTextFile(nil, texFile, nil),
		stdout:     knuth.NewTextFile(nil, stdout, nil),
		stderr:     knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
