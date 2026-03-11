// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gotftopl is the TFtoPL program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/tftopl.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/tftopl package.
package main // modernc.org/knuth/cmd/gotftopl

import (
	"flag"
	"fmt"
	"io"
	"os"

	"modernc.org/knuth/tftopl"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program TFtoPL( tfm_file, pl_file, output, stderr);
// tfm_file:packed file of 0..255;
// pl_file:text;
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 2 {
		fail(2, "expected 1 or 2 arguments: tfm_file [pl_file]\n")
	}

	tfmFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer tfmFile.Close()

	plFile := io.Writer(os.Stdout)
	if nArg == 2 {
		f, err := os.Create(flag.Arg(1))
		if err != nil {
			fail(1, "creating %v: %v\n", flag.Arg(1), err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				fail(1, "closing %v: %v\n", flag.Arg(1), err)
			}
		}()

		plFile = f
	}

	if err = tftopl.Main(tfmFile, plFile, os.Stdout, os.Stderr); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
