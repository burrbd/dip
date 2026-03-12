// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package pltotf is the PLtoTF program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/pltotf.web
package pltotf // modernc.org/knuth/pltotf

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program PLtoTF( pl_file, tfm_file, output, stderr);
// pl_file:text;
// @!tfm_file:packed file of 0..255;

// Main executes the vptovf program using the supplied arguments.
func Main(plFile io.Reader, tfmFile, stdout, stderr io.Writer) (mainErr error) {
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
		plFile:  knuth.NewTextFile(plFile, nil, nil),
		tfmFile: knuth.NewBinaryFile(nil, tfmFile, 1, nil),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
