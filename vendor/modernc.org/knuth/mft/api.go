// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package mft is the MFT program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/mft.web
package mft // modernc.org/knuth/mft

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program MFT( mf_file, change_file, style_file, tex_file, output);
// mf_file:text_file; {primary input}
// change_file:text_file; {updates}
// style_file:text_file; {formatting bootstrap}
// tex_file: text_file;

// Main executes the mft program using the supplied arguments.
func Main(mfFile, changeFile, styleFile io.Reader, texFile, stdout, stderr io.Writer) (mainErr error) {
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
		mfFile:     knuth.NewTextFile(mfFile, nil, nil),
		changeFile: knuth.NewTextFile(changeFile, nil, nil),
		styleFile:  knuth.NewTextFile(styleFile, nil, nil),
		texFile:    knuth.NewTextFile(nil, texFile, nil),
		stdout:     knuth.NewTextFile(nil, stdout, nil),
		stderr:     knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
