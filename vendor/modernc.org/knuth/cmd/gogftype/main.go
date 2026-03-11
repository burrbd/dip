// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gogftype is the GFtype program by D. R. Fuchs, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftype.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/gftype package.
package main // modernc.org/knuth/cmd/gogftype

import (
	"flag"
	"fmt"
	"os"

	"modernc.org/knuth/gftype"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program GF_type( gf_file, output,stderr);
// @!gf_file:byte_file; {the stuff we are \.{GF}typing}

// Main executes the gftype program using the supplied arguments.
func main() {
	oImages := flag.Bool("images", false, "show characters as pixels")
	oMnemonics := flag.Bool("mnemonics", false, "translate all GF commands")
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 1 {
		fail(2, "expected 1 argument: gf_file\n")
	}

	gfFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer gfFile.Close()

	if err = gftype.Main(gfFile, os.Stdout, os.Stderr, *oMnemonics, *oImages); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
