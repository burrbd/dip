// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package gftodvi is the GFtoDVI program by Tomas Rokicki, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftodvi.web
package gftodvi // modernc.org/knuth/gftodvi

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program GF_to_DVI( output);

// Main executes the gftodvi program using the supplied arguments.
func Main(gfFile io.Reader, dviFile, stdout, stderr io.Writer, open func(string) (io.Reader, error)) (mainErr error) {
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
		default:
			mainErr = fmt.Errorf("PANIC %T: %[1]v, error: %v\n%s", x, mainErr, debug.Stack())
		}
	}()

	(&prg{
		gfFile:  knuth.NewBinaryFile(gfFile, nil, 1, nil),
		dviFile: knuth.NewBinaryFile(nil, dviFile, 1, nil),
		tfmFile: knuth.NewBinaryFile(nil, nil, 1, open),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
