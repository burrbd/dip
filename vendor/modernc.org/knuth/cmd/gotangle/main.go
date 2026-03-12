// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gotangle is the TANGLE program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/web/tangle.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/tangle package.
package main // modernc.org/knuth/cmd/gotangle

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"modernc.org/knuth/tangle"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program TANGLE( web_file, change_file, Pascal_file, pool);
// web_file:text_file; {primary input}
// change_file:text_file; {updates}
// Pascal_file: text_file;
// pool: text_file;

// Main executes the tangle program using the supplied arguments.
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg < 1 || nArg > 4 {
		fail(2, "expected 1 to 4 argument(s): web_file.web [change_file.ch] [pascal_file.p] [pool_file.pool]\n")
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

	pascalNm := root + ".p"
	if len(args) != 0 {
		if filepath.Ext(filepath.Base(args[0])) == ".p" {
			pascalNm = args[0]
			args = args[1:]
		}
	}
	pascalFile, err := os.Create(pascalNm)
	if err != nil {
		fail(1, "%s\n", err)
	}

	poolNm := root + ".pool"
	if len(args) != 0 {
		if filepath.Ext(filepath.Base(args[0])) == ".p" {
			poolNm = args[0]
			args = args[1:]
		}
	}
	poolFile, err := os.Create(poolNm)
	if err != nil {
		fail(1, "%s\n", err)
	}

	if err = tangle.Main(webFile, changeFile, pascalFile, poolFile, os.Stdout, os.Stderr); err != nil {
		fail(1, "FAIL: %s\n", err)
	}
}
