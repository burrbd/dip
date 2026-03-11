// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package pooltype is the POOLtype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/pooltype.web
package pooltype // modernc.org/knuth/pooltype

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program POOLtype( pool_file, output, stderr);
// pool_file:packed file of  char ;

// Main executes the pooltype program using the supplied arguments.
func Main(poolFile io.Reader, stdout, stderr io.Writer) (mainErr error) {
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
		stdout:   knuth.NewTextFile(nil, stdout, nil),
		stderr:   knuth.NewTextFile(nil, stderr, nil),
		poolFile: knuth.NewTextFile(poolFile, nil, nil),
	}).main()
	return nil
}
