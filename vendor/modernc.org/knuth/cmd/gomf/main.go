// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gomf is the METAFONT program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/mf.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/mft package.
package main // modernc.org/knuth/cmd/gomf

import (
	"fmt"
	"io"
	"os"
	"strings"

	"modernc.org/knuth/mf"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program MF; {all file names are defined dynamically}

// Main executes the mft program using the supplied arguments.
func main() {
	in := io.Reader(os.Stdin)
	if len(os.Args) > 1 {
		in = io.MultiReader(strings.NewReader(strings.Join(os.Args[1:], " ")+"\n"), in)
	}
	if err := mf.Main(in, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
