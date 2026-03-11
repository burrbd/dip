// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command goweave is the WEAVE program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/web/weave.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/weave package.
package main // modernc.org/knuth/cmd/goweave

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"modernc.org/knuth/weave"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program WEAVE( web_file, change_file, tex_file);
// web_file:text_file; {primary input}
// change_file:text_file; {updates}
// tex_file: text_file;

// Main executes the weave program using the supplied arguments.
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 3 {
		fail(2, "expected 1 to 3 argument(s): web_file.web [change_file.ch] [tex_file.tex]\n")
	}

	args := flag.Args()
	root := args[0]
	root = root[:len(root)-len(filepath.Ext(root))]
	webFile, err := os.Open(args[0])
	args = args[1:]
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer webFile.Close()

	changeFile := io.Reader(bytes.NewBuffer(nil))
	if len(args) != 0 {
		if filepath.Ext(filepath.Base(args[0])) == ".ch" {
			changeFile, err = os.Open(args[0])
			args = args[1:]
			if err != nil {
				fail(1, "%s\n", err)
			}
		}
	}

	texNm := root + ".tex"
	if len(args) != 0 {
		if filepath.Ext(filepath.Base(args[0])) == ".tex" {
			texNm = args[0]
			args = args[1:]
		}
	}
	texFile, err := os.Create(texNm)
	if err != nil {
		fail(1, "%s\n", err)
	}

	if err = weave.Main(webFile, changeFile, texFile, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
