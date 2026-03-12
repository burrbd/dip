// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command govftovp is the VFtoVP program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/etc/vftovp.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/vftovp package.
package main // modernc.org/knuth/cmd/govftovp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"modernc.org/knuth"
	"modernc.org/knuth/vftovp"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program VFtoVP( vf_file, tfm_file, vpl_file, output, stderr);
func main() {
	sep := string(os.PathListSeparator)
	oTexfonts := flag.String("texfonts", "", fmt.Sprintf("a list of search paths separated by %q", sep))
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 2 || nArg > 3 {
		fail(2, "expected 2 or 3 arguments: vf_file tfm_file [vpl_file]\n")
	}

	var fontPaths []string
	if s := *oTexfonts; s != "" {
		fontPaths = strings.Split(s, sep)
	}
	vfFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer vfFile.Close()

	tfmFile, err := os.Open(flag.Arg(1))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer tfmFile.Close()

	vplFile := io.Writer(os.Stdout)
	if nArg == 3 {
		f, err := os.Create(flag.Arg(2))
		if err != nil {
			fail(1, "creating %v: %v\n", flag.Arg(2), err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				fail(1, "closing %v: %v\n", flag.Arg(2), err)
			}
		}()

		vplFile = f
	}
	if err = vftovp.Main(
		vfFile, tfmFile, vplFile, os.Stdout, os.Stderr,
		func(fileName string) (r io.Reader, err error) { return knuth.Open(fileName, fontPaths) },
	); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
