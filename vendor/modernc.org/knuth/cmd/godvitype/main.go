// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command godvitype is the DVItype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/dvitype.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/dvitype package.
package main // modernc.org/knuth/cmd/godvitype

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"modernc.org/knuth"
	"modernc.org/knuth/dvitype"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program DVI_type( dvi_file, output);
// dvi_file:byte_file; {the stuff we are \.[DVI]typing}

// Main executes the dvitype program using the supplied arguments.
func main() {
	sep := string(os.PathListSeparator)
	oDPI := flag.Float64("dpi", 0, "device resolution in pixels per inch; default 300.0")
	oMagnification := flag.Int("magnification", 0, "magnification")
	oMaxPages := flag.Int("max-pages", -1, "process max <arg> pages; default 1e6")
	oOutputMode := flag.Int("output-level", -1, "verbosity level, from 0 to 4; default 4")
	oPageStart := flag.String("page-start", "", "start at <arg>, for example '2' or '5.*.-2'")
	oTexfonts := flag.String("texfonts", "", fmt.Sprintf("a list of search paths separated by %q", sep))
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 1 {
		fail(2, "expected 1 argument: gf_file\n")
	}

	dviFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer dviFile.Close()

	var fontPaths []string
	if s := *oTexfonts; s != "" {
		fontPaths = strings.Split(s, sep)
	}
	if err = dvitype.Main(
		dviFile, os.Stdout, os.Stderr, *oOutputMode, *oPageStart, *oMaxPages, *oDPI, *oMagnification,
		func(fileName string) (r io.Reader, err error) { return knuth.Open(fileName, fontPaths) },
	); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
