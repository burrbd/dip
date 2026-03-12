// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gftodvi is the GFtoDVI program by Tomas Rokicki, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/gftodvi.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/gftodvi package.
package main // modernc.org/knuth/cmd/gogftodvi

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"modernc.org/knuth"
	"modernc.org/knuth/gftodvi"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program GF_to_DVI( output);

// Main executes the gftodvi program using the supplied arguments.
func main() {
	sep := string(os.PathListSeparator)
	oTexfonts := flag.String("texfonts", "", fmt.Sprintf("a list of search paths separated by %q", sep))
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 2 {
		fail(2, "expected 1 or 2 arguments: gf_file [dvi_file]\n")
	}

	gfFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer gfFile.Close()

	var dviFile *os.File
	switch {
	case nArg == 1:
		base := filepath.Base(flag.Arg(0))
		if dviFile, err = os.Create(base[:len(base)-len(filepath.Ext(base))] + ".dvi"); err != nil {
			fail(1, "creating dvi_file: %v", err)
		}
	default:
		if dviFile, err = os.Create(flag.Arg(1)); err != nil {
			fail(1, "creating dvi_file: %v", err)
		}
	}

	defer func() {
		if err := dviFile.Close(); err != nil {
			fail(1, "closing dvi_file: %v", err)
		}
	}()

	var fontPaths []string
	if s := *oTexfonts; s != "" {
		fontPaths = strings.Split(s, sep)
	}
	if err = gftodvi.Main(gfFile, dviFile, os.Stdout, os.Stderr,
		func(fileName string) (r io.Reader, err error) { return knuth.Open(fileName, fontPaths) },
	); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
