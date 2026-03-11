// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package pktype is the PKtype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/local/mfware/pktype.web
package pktype // modernc.org/knuth/pktype

import (
	"fmt"
	"io"
	"runtime/debug"

	"modernc.org/knuth"
)

// program PKtype( pk_file, typ_file, output);
// pk_file:byte_file;  {where the input comes from}
// typ_file:text_file; {where the final output goes}

// Main executes the pktyype program using the supplied arguments.
func Main(pkFile io.Reader, typFile, stdout, stderr io.Writer) (mainErr error) {
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
		pkFile:  knuth.NewBinaryFile(pkFile, nil, 1, nil),
		typFile: knuth.NewTextFile(nil, typFile, nil),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}).main()
	return nil
}
