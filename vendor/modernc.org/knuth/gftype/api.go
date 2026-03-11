// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package gftype is the GFtype program by D. R. Fuchs, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftype.web
package gftype // modernc.org/knuth/gftype

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program GF_type( gf_file, output,stderr);
// @!gf_file:byte_file; {the stuff we are \.{GF}typing}

// Main executes the gftyype program using the supplied arguments.
func Main(gfFile io.Reader, stdout, stderr io.Writer, wantsMnemonics, wantsPixels bool) (mainErr error) {
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
		gfFile:         knuth.NewBinaryFile(gfFile, nil, 1, nil),
		stdout:         knuth.NewTextFile(nil, stdout, nil),
		stderr:         knuth.NewTextFile(nil, stderr, nil),
		wantsMnemonics: wantsMnemonics,
		wantsPixels:    wantsPixels,
	}).main()
	return nil
}
