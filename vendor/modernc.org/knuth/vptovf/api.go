// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package vptovf is the VPtoVF program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/etc/vptovf.web
package vptovf // modernc.org/knuth/vptovf

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program VPtoVF( vpl_file, vf_file, tfm_file, output, stderr);
// vpl_file:text;
// vf_file:packed file of 0..255;
// tfm_file:packed file of 0..255;

// Main executes the vptovf program using the supplied arguments.
func Main(vplFile io.Reader, vfFile, tfmFile, stdout, stderr io.Writer) (mainErr error) {
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
		vplFile: knuth.NewTextFile(vplFile, nil, nil),
		vfFile:  knuth.NewBinaryFile(nil, vfFile, 1, nil),
		tfmFile: knuth.NewBinaryFile(nil, tfmFile, 1, nil),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
