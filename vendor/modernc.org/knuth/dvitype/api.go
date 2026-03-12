// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate make generate

// Package dvitype is the DVItype program by D. E. Knuth, transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/texware/dvitype.web
package dvitype // modernc.org/knuth/dvitype

import (
	"fmt"
	"io"
	"math"
	"runtime/debug"
	"strings"

	"modernc.org/knuth"
)

// program DVI_type( dvi_file, output);
// dvi_file:byte_file; {the stuff we are \.[DVI]typing}

// Main executes the dvitype program using the supplied arguments.
func Main(dviFile io.ReadSeeker, stdout, stderr io.Writer, outMode int, startingPage string, maxPages int, resolution float64, magnification int, open func(string) (io.Reader, error)) (mainErr error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// ok
		case signal:
			switch {
			case mainErr == nil:
				mainErr = fmt.Errorf("aborted")
			default:
				mainErr = fmt.Errorf("aborted: %v", mainErr)
			}
		case knuth.Error:
			mainErr = x
		default:
			mainErr = fmt.Errorf("PANIC %T: %[1]v, error: %v\n%s", x, mainErr, debug.Stack())
		}
	}()

	prg := &prg{
		dviFile: knuth.NewBinaryFile(dviFile, nil, 1, nil),
		tfmFile: knuth.NewBinaryFile(nil, nil, 1, open),
		stdout:  knuth.NewTextFile(nil, stdout, nil),
		stderr:  knuth.NewTextFile(nil, stderr, nil),
	}

	var (
		k int32
	)

	//  Determine the desired |out_mode|
	prg.outMode = byte(theWorks) // default
	switch {
	case outMode < 0:
		// ok
	case outMode >= 0 && outMode <= 4:
		prg.outMode = byte(outMode)
	default:
		return fmt.Errorf("outMode must be in [0, 4] or < 0 for default: %v", outMode)
	}

	//  Determine the desired |start_count| values
	prg.startVals = 0
	prg.startThere[0] = false
	startingPage = strings.TrimSpace(startingPage) + " "
	if len(startingPage) >= len(prg.buffer) {
		return fmt.Errorf("startingPage too long: %q", startingPage)
	}

	for i := 0; i < len(startingPage) && i < len(prg.buffer); i++ {
		prg.buffer[i] = startingPage[i]
	}
	prg.bufPtr = 0
	k = 0
	if int32(prg.buffer[0]) != ' ' {
		for {
			if int32(prg.buffer[prg.bufPtr]) == '*' {
				prg.startThere[k] = false
				prg.bufPtr = byte(int32(prg.bufPtr) + 1)
			} else {
				prg.startThere[k] = true
				prg.startCount[k] = prg.getInteger()
			}
			if k < 9 && int32(prg.buffer[prg.bufPtr]) == '.' {
				k = k + 1
				prg.bufPtr = byte(int32(prg.bufPtr) + 1)
			} else if int32(prg.buffer[prg.bufPtr]) == ' ' {
				prg.startVals = byte(k)
			} else {
				prg.termOut.Write("Type, e.g., 1.*.-5 to specify the ")
				prg.termOut.Writeln("first page with \\count0=1, \\count2=-5.")

				return fmt.Errorf("invalid startingPage: %q", startingPage)
			}
			if int32(prg.startVals) == k {
				break
			}
		}
	}

	//  Determine the desired |max_pages|
	prg.maxPages = 1000000 // default
	switch {
	case maxPages < 0:
		// ok
	case maxPages > 0 && maxPages < math.MaxInt32:
		prg.maxPages = int32(maxPages)
	default:
		return fmt.Errorf("maxPages must be in [1, math.MaxInt32] or < 0 for default: %v", maxPages)
	}

	//  Determine the desired |resolution|
	prg.resolution = 300.0
	if resolution > 0 {
		prg.resolution = resolution
	}

	//  Determine the desired |new_mag|
	prg.newMag = 0
	if magnification > 0 {
		prg.newMag = int32(magnification)
	}

	prg.main()
	return nil
}
