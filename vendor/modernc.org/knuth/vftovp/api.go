// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package vftovp is the VFtoVP program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/etc/vftovp.web
package vftovp // modernc.org/knuth/vftovp

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program VFtoVP( vf_file, tfm_file, vpl_file, output, stderr);
// vf_file:packed file of byte;
// tfm_file:packed file of byte;
// vpl_file:text;

// Main executes the vftovp program using the supplied arguments.
func Main(vfFile, tfmFile io.Reader, vplFile, stdout, stderr io.Writer, open func(string) (io.Reader, error)) (mainErr error) {
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
		vfFile:  knuth.NewBinaryFile(vfFile, nil, 1, nil),
		tfmFile: knuth.NewBinaryFile(tfmFile, nil, 1, open),
		vplFile: knuth.NewTextFile(nil, vplFile, nil),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
