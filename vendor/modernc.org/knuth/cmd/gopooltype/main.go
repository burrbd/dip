// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gopooltype is the POOLtype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/pooltype.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/pooltype package.
package main // modernc.org/knuth/cmd/gopooltype

import (
	"flag"
	"fmt"
	"os"

	"modernc.org/knuth/pooltype"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program POOLtype( pool_file, output, stderr);
// pool_file:packed file of  char ;
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 1 {
		fail(2, "expected 1: poolFile\n")
	}

	poolFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer poolFile.Close()

	if err = pooltype.Main(poolFile, os.Stdout, os.Stderr); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
