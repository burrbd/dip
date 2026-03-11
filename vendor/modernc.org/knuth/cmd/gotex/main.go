// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gotex is the TeX program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/tex/tex.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/tex package.
package main // modernc.org/knuth/cmd/gotex

import (
	"fmt"
	"io"
	"os"
	"strings"

	"modernc.org/knuth/tex"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program TEX; {all file names are defined dynamically}

// Main executes the tex program using the supplied arguments.
func main() {
	in := io.Reader(os.Stdin)
	if len(os.Args) > 1 {
		in = io.MultiReader(strings.NewReader(strings.Join(os.Args[1:], " ")+"\n"), in)
	}
	if err := tex.Main(in, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
