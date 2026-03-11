// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gopltotf is the PLtoTF program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/pltotf.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/pltotf package.
package main // modernc.org/knuth/cmd/gopltotf

import (
	"flag"
	"fmt"
	"os"

	"modernc.org/knuth/pltotf"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program PLtoTF( pl_file, tfm_file, output, stderr);
// pl_file:text;
// @!tfm_file:packed file of 0..255;

// Main executes the vptovf program using the supplied arguments.
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 2 {
		fail(2, "expected 2 arguments: pl_file tfm_file\n")
	}

	plFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer plFile.Close()

	tfmFile, err := os.Create(flag.Arg(1))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer func() {
		if err := tfmFile.Close(); err != nil {
			fail(1, "closing %s: %v", flag.Arg(1), err)
		}
	}()

	if err = pltotf.Main(plFile, tfmFile, os.Stdout, os.Stderr); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
