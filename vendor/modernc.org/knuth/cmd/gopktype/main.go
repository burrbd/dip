// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gopktype is the PKtype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/local/mfware/pktype.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/pktype package.
package main // modernc.org/knuth/cmd/gopktype

import (
	"flag"
	"fmt"
	"io"
	"os"

	"modernc.org/knuth/pktype"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program PKtype( pk_file, typ_file, output);
// pk_file:byte_file;  {where the input comes from}
// typ_file:text_file; {where the final output goes}

// Main executes the pktype program using the supplied arguments.
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 2 {
		fail(2, "expected 2 or 3 arguments: pk_file [typ_file]\n")
	}

	pkFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer pkFile.Close()

	typFile := io.Writer(os.Stdout)
	if nArg == 2 {
		typFile, err := os.Create(flag.Arg(1))
		if err != nil {
			fail(1, "creating %s: %v\n", flag.Arg(1), err)
		}

		defer func() {
			if err := typFile.Close(); err != nil {
				fail(1, "closing %s: %v\n", flag.Arg(1), err)
			}
		}()

	}

	if err = pktype.Main(pkFile, typFile, os.Stdout, os.Stderr); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
