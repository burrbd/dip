// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package trap is the TRAP variant of the METAFONT program by D. E. Knuth,
// transpiled to Go.
//
//	http://mirrors.ctan.org/systems/knuth/dist/mf/mf.web
package trap // modernc.org/knuth/mf/internal/trap

import (
	// Required for go:embed
	_ "embed"
	"fmt"
	"io"
	"runtime/debug"
	"unsafe"

	"modernc.org/knuth"
)

//go:embed mf.pool
var pool string

// program MF; {all file names are defined dynamically}

// Main executes the tangle program using the supplied arguments.
func Main(stdin io.Reader, stdout, stderr io.Writer, opener func(string) (io.Reader, error)) (mainErr error) {
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
		baseFile: knuth.NewBinaryFile(nil, nil, int(unsafe.Sizeof(memoryWord{})), opener),
		gfFile:   knuth.NewBinaryFile(nil, nil, 1, nil),
		logFile:  knuth.NewTextFile(nil, nil, nil),
		poolFile: knuth.NewPoolFile(pool),
		stderr:   knuth.NewTextFile(nil, stderr, nil),
		termIn:   knuth.NewTextFile(stdin, nil, nil),
		termOut:  knuth.NewTextFile(nil, stdout, nil),
		tfmFile:  knuth.NewBinaryFile(nil, nil, 1, nil),
	}
	for i := range prg.inputFile {
		prg.inputFile[i] = knuth.NewTextFile(nil, nil, opener)
	}
	prg.main()
	return nil
}
