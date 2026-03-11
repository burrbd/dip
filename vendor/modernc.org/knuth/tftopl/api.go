// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package tftopl is the TFtoPL program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/tftopl.web
package tftopl // modernc.org/knuth/tftopl

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program TFtoPL( tfm_file, pl_file, output, stderr);
// tfm_file:packed file of 0..255;
// pl_file:text;

// Main executes the tftopl program using the supplied arguments.
func Main(tfmFile io.Reader, plFile, stdout, stderr io.Writer) (mainErr error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// ok
		case knuth.Error:
			mainErr = x
		default:
			mainErr = fmt.Errorf("PANIC %T: %[1]v, error: %s\n%s", x, mainErr, debug.Stack())
		}
	}()

	(&prg{
		tfmFile: knuth.NewBinaryFile(tfmFile, nil, 1, nil),
		plFile:  knuth.NewTextFile(nil, plFile, nil),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stdout, nil),
	}).main()
	return nil
}
