// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gogftopk is the GFtype program by Tomas Rokicki, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftopk.web

// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/gftype package.
package main // modernc.org/knuth/cmd/gogftopk

import (
	"flag"
	"fmt"
	"os"

	"modernc.org/knuth/gftopk"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program GFtoPK( gf_file, pk_file, output)
// gf_file:byte_file; {the stuff we are \.[GFtoPK]ing}
// pk_file:byte_file; {the stuff we have \.[GFtoPK]ed}

// Main executes the gftype program using the supplied arguments.
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 2 {
		fail(2, "expected 2 arguments: gf_file pk_file\n")
	}

	gfFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer gfFile.Close()

	pkFile, err := os.Create(flag.Arg(1))
	if err != nil {
		fail(1, "creating pk_file: %v", err)
	}

	if err = gftopk.Main(gfFile, pkFile, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
