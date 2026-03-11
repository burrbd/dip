// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package gftopk is the GFtype program by Tomas Rokicki, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftopk.web
package gftopk // modernc.org/knuth/gftopk

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program GFtoPK( gf_file, pk_file, output)
// gf_file:byte_file; {the stuff we are \.[GFtoPK]ing}
// pk_file:byte_file; {the stuff we have \.[GFtoPK]ed}

// Main executes the gftyype program using the supplied arguments.
func Main(gfFile io.ReadSeeker, pkFile, stdout, stderr io.Writer) (mainErr error) {
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
		gfFile: knuth.NewBinaryFile(gfFile, nil, 1, nil),
		pkFile: knuth.NewBinaryFile(nil, pkFile, 1, nil),
		stdout: knuth.NewTextFile(nil, stdout, nil),
		stderr: knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
