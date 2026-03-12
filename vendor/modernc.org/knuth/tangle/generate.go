// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"modernc.org/knuth"
	"modernc.org/knuth/web"
)

func fail(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}

func main() {
	const (
		base = "tangle"

		chFn   = base + ".ch"
		goFn   = base + ".go"
		pasFn  = base + ".pas"
		poolFn = base + ".pool"
		webFn  = base + ".web"
	)

	dest, err := os.Create(filepath.FromSlash(goFn))
	if err != nil {
		fail("creating %s: %v\n", goFn, err)
	}

	defer func() {
		if err = dest.Close(); err != nil {
			fail("closing %s: %v\n", goFn, err)
		}
	}()

	pas, err := os.Create(filepath.FromSlash(pasFn))
	if err != nil {
		fail("creating %s: %v\n", pasFn, err)
	}

	defer func() {
		if err = pas.Close(); err != nil {
			fail("closing %s: %v\n", pasFn, err)
		}
	}()

	pool, err := os.Create(filepath.FromSlash(poolFn))
	if err != nil {
		fail("creating %s: %v\n", poolFn, err)
	}

	defer func() {
		if err = pool.Close(); err != nil {
			fail("closing %s: %v\n", poolFn, err)
		}
	}()

	webSrc, err := os.ReadFile(filepath.FromSlash(webFn))
	if err != nil {
		fail("reading %v: %v\n", webFn, err)
	}

	chSrc, err := os.ReadFile(filepath.FromSlash(chFn))
	if err != nil {
		fail("reading %v: %v\n", chFn, err)
	}

	src, err := knuth.NewChanger(
		knuth.NewRuneSource(webFn, webSrc, knuth.Unicode),
		knuth.NewRuneSource(chFn, chSrc, knuth.Unicode),
	)
	if err != nil {
		fail("processing %s and %s: %v\n", webFn, chFn, err)
	}

	if err := web.Go(dest, pas, pool, src, base); err != nil {
		fail("generate: %v\n", err)
	}
}
