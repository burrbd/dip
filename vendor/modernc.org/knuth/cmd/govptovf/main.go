// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command govftovp is the VPtoVF program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/etc/vptovf.web
//
// For more details about the original Pascal program and its usage please see
// the modernc.org/knuth/vftovp package.
package main // modernc.org/knuth/cmd/govptovf

import (
	"flag"
	"fmt"
	"os"

	"modernc.org/knuth/vptovf"
)

func fail(rc int, s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(rc)
}

// program VPtoVF( vpl_file, vf_file, tfm_file, output, stderr);
// vpl_file:text;
// vf_file:packed file of 0..255;
// tfm_file:packed file of 0..255;
func main() {
	flag.Parse()
	nArg := flag.NArg()
	if nArg != 3 {
		fail(2, "expected 3 arguments: vpl_file vf_file tfm_file\n")
	}

	vplFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer vplFile.Close()

	vfFile, err := os.Create(flag.Arg(1))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer func() {
		if err := vfFile.Close(); err != nil {
			fail(1, "closing %s: %v", flag.Arg(1), err)
		}
	}()

	tfmFile, err := os.Create(flag.Arg(2))
	if err != nil {
		fail(1, "%s\n", err)
	}

	defer func() {
		if err := tfmFile.Close(); err != nil {
			fail(1, "closing %s: %v", flag.Arg(2), err)
		}
	}()

	if err = vptovf.Main(vplFile, vfFile, tfmFile, os.Stdout, os.Stderr); err != nil {
		fail(2, "FAIL: %s\n", err)
	}
}
