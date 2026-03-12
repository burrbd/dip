// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gomft is the MFT program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mfware/mft.web
//
// For more details about the original Pascal program and its usage please see
// the PDF documentation included in the modernc.org/knuth/mft package.
package main // modernc.org/knuth/cmd/gomft

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"modernc.org/knuth"
	"modernc.org/knuth/mft"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program MFT( mf_file, change_file, style_file, tex_file, output);
// mf_file:text_file; {primary input}
// change_file:text_file; {updates}
// style_file:text_file; {formatting bootstrap}
// tex_file: text_file;

// Main executes the mft program using the supplied arguments.
func main() {
	oChange := flag.String("change", "", "apply a change file change file")
	oStyle := flag.String("style", "", "use a named style file instead of the default plain.mft asset")
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 2 {
		fail(2, "expected 1 or 2 arguments: mf_file [tex_file]\n")
	}

	mfFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer mfFile.Close()

	texFile := io.Writer(os.Stdout)
	if nArg == 2 {
		texFile, err := os.Create(flag.Arg(1))
		if err != nil {
			fail(1, "creating %s: %v\n", flag.Arg(1), err)
		}

		defer func() {
			if err := texFile.Close(); err != nil {
				fail(1, "closing %s: %v\n", flag.Arg(1), err)
			}
		}()

	}

	var changeFile, styleFile io.Reader
	switch nm := *oChange; {
	case nm != "":
		f, err := os.Open(nm)
		if err != nil {
			fail(1, "change file: %v", err)
		}

		defer f.Close()

		changeFile = f
	default:
		changeFile = bytes.NewBuffer(nil)
	}
	switch nm := *oStyle; {
	case nm != "":
		f, err := os.Open(nm)
		if err != nil {
			fail(1, "style file: %v", err)
		}

		defer f.Close()

		styleFile = f
	default:
		styleFile, err = knuth.Open("TeXinputs:plain.mft", nil)
		if err != nil {
			fail(1, "style file: %v", err)
		}
	}

	if err = mft.Main(mfFile, changeFile, styleFile, texFile, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
