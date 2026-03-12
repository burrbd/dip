// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package web deals with .web files.
package web // modernc.org/knuth/web

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"modernc.org/knuth"
	"modernc.org/sortutil"
)

var (
	oTrcw bool // testing
)

const (
	blankSet          = " \t\n\r"
	incompleteNameTag = "..."
	poolSumInit       = 271828 // tangle.pdf, p.138
	stackLimit        = 250    // min 124
)

const (
	eof                = -iota - 1
	ccAt               // "@@"
	ccBeginMetaComment // "@{"
	ccBeginPascal      // "@p"
	ccBigLineBreak     // "@#"
	ccCheckSum         // "@$"
	ccDefinition       // "@d"
	ccEnd              // "@>"
	ccEndMetaComment   // "@}"
	ccForceLine        // "@\\"
	ccFormat           // "@f"
	ccHex              // "@\""
	ccJoin             // "@&"
	ccLineBreak        // "@/"
	ccMathBreak        // "@|"
	ccModuleName       // "@<"
	ccNewModule        // "@ "
	ccNewStarredModule // "@*"
	ccNoLineBreak      // "@+"
	ccNoUnderline      // "@?"
	ccOctal            // "@'"
	ccPopMacroArg      // "@Z" internal
	ccPseudoSemi       // "@;"
	ccTeXString        // "@t"
	ccThinSpace        // "@,"
	ccUnderline        // "@!"
	ccVerbatim         // "@="
	ccXrefRoman        // "@^"
	ccXrefTypewriter   // "@."
	ccXrefWildcard     // "@:"
)

type abort error

type webScanner struct {
	controlCodePos token.Position
	src            knuth.RuneSource
	stack          []knuth.RuneSource

	c2          rune
	controlCode rune

	controlCodeValid bool
}

func newWebScanner(src knuth.RuneSource) *webScanner { return &webScanner{src: src} }

func (s *webScanner) srcStack() (r []string) {
	for i, v := range s.stack {
		p := v.Position()
		p.Filename = filepath.Base(p.Filename)
		r = append(r, fmt.Sprintf("%3d: %p %v:", i, v, p))
	}
	if v := s.src; v != nil {
		p := v.Position()
		p.Filename = filepath.Base(p.Filename)
		r = append(r, fmt.Sprintf("TOS: %p %v:", v, p))
	} else {
		r = append(r, fmt.Sprintf("TOS: <nil>"))
	}
	return r
}

func (s *webScanner) c() (r rune) {
	if s.controlCodeValid {
		return s.controlCode
	}

pop:
	if s.src == nil {
		n := len(s.stack)
		if n == 0 {
			return eof
		}

		s.src = s.stack[n-1]
		s.stack = s.stack[:n-1]
	}
	c, err := s.src.C()
	switch {
	case err == io.EOF:
		s.src = nil
		if len(s.stack) != 0 {
			goto pop
		}

		return eof
	case err != nil:
		panic(abort(fmt.Errorf("%v: %v", s.src.Position(), err)))
	case c == '@':
		s.controlCodePos = s.position()
		s.consume()
		// The letters L, T , P , M , C, and/or S following each code indicate whether
		// or not that code is allowable in limbo, in TEX text, in Pascal text, in
		// module names, in comments, and/or in strings.
		if s.c2, err = s.src.C(); err != nil {
			panic(abort(fmt.Errorf("%v: %v", s.src.Position(), err)))
		}

		switch s.c2 {
		case '*':
			// @* [!L, !P , !T] This denotes the beginning of a new starred module, i.e., a
			// module that begins a new major group. The title of the new group should
			// appear after the @*, followed by a period. As explained above, TEX control
			// sequences should be avoided in such titles unless they are quite simple.
			// When WEAVE and TANGLE read a @*, they print an asterisk on the terminal
			// followed by the current module number, so that the user can see some
			// indication of progress. The very first module should be starred.
			s.controlCode = ccNewStarredModule
		case 'd', 'D':
			// @d [!P , !T] Macro definitions begin with @d (or @D), followed by the Pascal
			// text for one of the three kinds of macros, as explained earlier.
			s.controlCode = ccDefinition
		case ' ', '\t', '\n':
			// @ [!L, !P , !T] This denotes the beginning of a new (unstarred) module. A
			// tab mark or end-of-line (carriage return) is equivalent to a space when it
			// follows an @ sign.
			s.controlCode = ccNewModule
		case 'p', 'P':
			// @p [!P , !T] The Pascal part of an unnamed module begins with @p (or @P).
			// This causes TANGLE to append the following Pascal code to the initial
			// program text T 0 as explained above. The WEAVE processor does not cause a
			// ‘@p’ to appear explicitly in the TEX output, so if you are creating a WEB
			// file based on a TEX-printed WEB documentation you have to remember to insert
			// @p in the appropriate places of the unnamed modules.
			s.controlCode = ccBeginPascal
		case '<':
			// @< [P, !T] A module name begins with @< followed by TEX text followed by @>;
			// the TEX text should not contain any WEB control codes except @@, unless
			// these control codes appear in Pascal text that is delimited by |...|. The
			// module name may be abbreviated, after its first appearance in a WEB file, by
			// giving any unique prefix followed by ..., where the three dots immediately
			// precede the closing @>. No module name should be a prefix of another. Module
			// names may not appear in Pascal text that is enclosed in |...|, nor may they
			// appear in the definition part of a module (since the appearance of a module
			// name ends the definition part and begins the Pascal part).
			s.controlCode = ccModuleName
		case 'f', 'F':
			// @f [!P , !T] Format definitions begin with @f (or @F); they cause WEAVE to
			// treat identifiers in a special way when they appear in Pascal text. The
			// general form of a format definition is ‘@f l == r’, followed by an optional
			// comment enclosed in braces, where l and r are identifiers; WEAVE will
			// subsequently treat identifier l as it currently treats r. This feature
			// allows a WEB programmer to invent new reserved words and/or to unreserve
			// some of Pascal’s reserved identifiers. The definition part of each module
			// consists of any number of macro definitions (beginning with @d) and format
			// definitions (beginning with @f), intermixed in any order.
			s.controlCode = ccFormat
		case '^':
			// @^ [P, T] The “control text” that follows, up to the next ‘@>’, will be
			// entered into the index together with the identifiers of the Pascal program;
			// this text will appear in roman type. For example, to put the phrase “system
			// dependencies” into the index, you can type ‘@^system dependencies@>’ in each
			// module that you want to index as system dependent. A control text, like a
			// string, must end on the same line of the WEB file as it began. Furthermore,
			// no WEB control codes are allowed in a control text, not even @@. (If you
			// need an @ sign you can get around this restriction by typing ‘\AT!’.)
			s.controlCode = ccXrefRoman
		case '>':
			s.controlCode = ccEnd
		case 'Z':
			s.controlCode = ccPopMacroArg
		case '/':
			// @/ [P] This control code causes a line break to occur within a Pascal
			// program formatted by WEAVE; it is ignored by TANGLE. Line breaks are chosen
			// automatically by TEX according to a scheme that works 99% of the time, but
			// sometimes you will prefer to force a line break so that the program is
			// segmented according to logical rather than visual criteria. Caution: ‘@/’
			// should be used only after statements or clauses, not in the middle of an
			// expression; use @| in the middle of expressions, in order to keep WEAVE’s
			// parser happy.
			s.controlCode = ccLineBreak
		case 't', 'T':
			// @t [P] The “control text” that follows, up to the next ‘@>’, will be put
			// into a TEX \hbox and formatted along with the neighboring Pascal program.
			// This text is ignored by TANGLE, but it can be used for various purposes
			// within WEAVE. For example, you can make comments that mix Pascal and
			// classical mathematics, as in ‘size < 2 15 ’, by typing ‘|size <
			// @t$2^{15}$@>|’. A control text must end on the same line of the WEB file as
			// it began, and it may not contain any WEB control codes.
			s.controlCode = ccTeXString
		case '\'':
			// @ ́ [P, T] This denotes an octal constant, to be formed from the succeeding
			// digits. For example, if the WEB file contains ‘@ ́100’, the TANGLE processor
			// will treat this an equivalent to ‘64’; the constant will be formatted as
			// “ ́100 ” in the TEX output produced via WEAVE. You should use octal notation
			// only for positive constants; don’t try to get, e.g., −1 by saying
			// ‘@ ́777777777777’.
			s.controlCode = ccOctal
		case '.':
			// @. [P, T] The “control text” that follows will be entered into the index in
			// typewriter type; see the rules for ‘@^’, which is analogous.
			s.controlCode = ccXrefTypewriter
		case '@':
			// @@ [C, L, M, P, S, T] A double @ denotes the single character ‘@’. This is
			// the only control code that is legal in limbo, in comments, and in strings.
			s.controlCode = ccAt
		case '#':
			// @# [P] This control code forces a line break, like @/ does, and it also
			// causes a little extra white space to appear between the lines at this break.
			// You might use it, for example, between procedure definitions or between
			// groups of macro definitions that are logically separate but within the same
			// module.
			s.controlCode = ccBigLineBreak
		case ',':
			// @, [P] This control code inserts a thin space in WEAVE’s output; it is
			// ignored by TANGLE. Sometimes you need this extra space if you are using
			// macros in an unusual way, e.g., if two identifiers are adjacent.
			s.controlCode = ccThinSpace
		case ':':
			// @: [P, T] The “control text” that follows will be entered into the index in
			// a format controlled by the TEX macro ‘\9’, which the user should define as
			// desired; see the rules for ‘@^’, which is analogous.
			s.controlCode = ccXrefWildcard
		case '&':
			// @& [P] The @& operation causes whatever is on its left to be adjacent to
			// whatever is on its right, in the Pascal output. No spaces or line breaks
			// will separate these two items. However, the thing on the left should not be
			// a semicolon, since a line break might occur after a semicolon.
			s.controlCode = ccJoin
		case '{':
			// @{ [P] The beginning of a “meta comment,” i.e., a comment that is supposed
			// to appear in the Pascal code, is indicated by @{ in the WEB file. Such
			// delimiters can be used as isolated symbols in macros or modules, but they
			// should be properly nested in the final Pascal program. The TANGLE processor
			// will convert ‘@{’ into ‘{’ in the Pascal output file, unless the output is
			// already part of a meta-comment; in the latter case ‘@{’ is converted into
			// ‘[’, since Pascal does not allow nested comments. The WEAVE processor
			// outputs ‘@{’. Incidentally, module numbers are automatically inserted as
			// meta-comments into the Pascal program, in order to help correlate the
			// outputs of WEAVE and TANGLE (see Appendix C) Meta-comments can be used to
			// put conditional text into a Pascal program; this helps to overcome one of
			// the limitations of WEB, since the simple macro processing routines of TANGLE
			// do not include the dynamic evaluation of boolean expressions.
			s.controlCode = ccBeginMetaComment
		case '}':
			// @} [P] The end of a “meta comment” is indicated by ‘@}’; this is converted
			// either into ‘}’ or ‘]’ in the Pascal output, according to the conventions
			// explained for @{ above. The WEAVE processor outputs ‘@}’.
			s.controlCode = ccEndMetaComment
		case '$':
			// @$ [P] This denotes the string pool check sum.
			s.controlCode = ccCheckSum
		case '?':
			// @? [P, T] This cancels an implicit (or explicit) ‘@!’, so that the next
			// index entry will not be underlined.
			s.controlCode = ccNoUnderline
		case '=':
			// @= [P] The “control text” that follows, up to the next ‘@>’, will be passed
			// verbatim to the Pascal program.
			s.controlCode = ccVerbatim
		case '"':
			// @" [P, T] A hexadecimal constant; ‘@"D0D0’ tangles to 53456 and weaves to
			// ‘ ̋D0D0’.
			s.controlCode = ccHex
		case '\\':
			// @\ [P] Force end-of-line here in the Pascal program file.
			s.controlCode = ccForceLine
		case '!':
			// @! [P, T] The module number in an index entry will be underlined if ‘@!’
			// immediately precedes the identifier or control text being indexed. This
			// convention is used to distinguish the modules where an identifier is
			// defined, or where it is explained in some special way, from the modules
			// where it is used. A reserved word or an identifier of length one will not be
			// indexed except for underlined entries. An ‘@!’ is implicitly inserted by
			// WEAVE just after the reserved words function, procedure, program, and var,
			// and just after @d and @f. But you should insert your own ‘@!’ before the
			// definitions of types, constants, variables, parameters, and components of
			// records and enumerated types that are not covered by this implicit
			// convention, if you want to improve the quality of the index that you get.
			s.controlCode = ccUnderline
		case '+':
			// @+ [P] This control code cancels a line break that might otherwise be
			// inserted by WEAVE, e.g., before the word ‘else’, if you want to put a short
			// if-then-else construction on a single line. It is ignored by TANGLE.
			s.controlCode = ccNoLineBreak
		case ';':
			// @; [P] This control code is treated like a semicolon, for formatting
			// purposes, except that it is invisible.  You can use it, for example, after a
			// module name when the Pascal text represented by that module name ends with a
			// semicolon.
			s.controlCode = ccPseudoSemi
		case '|':
			// @| [P] This control code specifies an optional line break in the midst of an
			// expression. For example, if you have a long condition between if and then,
			// or a long expression on the right-hand side of an assignment statement, you
			// can use ‘@|’ to specify breakpoints more logical than the ones that TEX
			// might choose on visual grounds.
			s.controlCode = ccMathBreak
		default:
			panic(todo("%v: %#U", s.controlCodePos, s.c2))
		}

		s.controlCodeValid = true
		return s.controlCode
	default:
		return c
	}
}

func (s *webScanner) consume() {
	s.controlCodeValid = false
	s.src.Consume()
}

func (s *webScanner) position() (r token.Position) {
	if s.controlCodeValid {
		return s.controlCodePos
	}

	if s.src != nil {
		return s.src.Position()
	}

	return r
}

// Tangle processes 'src' and outputs the resulting Pascal code to 'pascal' and
// a string pool to 'pool'. To apply a change file, pass knuth.NewChanger(src,
// changes) as 'src'.
//
// The result is similar, but not compatible, to what the original TANGLE
// outputs. It's also not compatible with the ISO Pascal Standard.
func Tangle(pascal, pool io.Writer, src knuth.RuneSource) (err error) {
	if !doPanic {
		defer func() {
			e := recover()
			switch x := e.(type) {
			case nil:
				return
			case abort:
				err = error(x)
				return
			}

			err = fmt.Errorf("PANIC %T: %[1]s, %s\n%s", e, err, debug.Stack())
		}()
	}

	t := newTangle(pascal, pool, src)
	if err := t.scan(); err != nil {
		return err
	}

	// Sanitize code names.
	t.codeNames = t.codeNames[:sortutil.Dedupe(sort.StringSlice(t.codeNames))]
	for i, v := range t.codeNames {
		if i < len(t.codeNames)-1 && strings.HasPrefix(t.codeNames[i+1], v) {
			panic(todo("", i))
		}
	}
	for _, m := range t.modules {
		for _, c := range m.codes {
			if nm := c.name; nm != "" {
				nm = t.completeName(nm)
				c.name = nm
				t.codesByName[nm] = append(t.codesByName[nm], c)
			}
		}
	}
	t.src = nil
	for _, m := range t.modules {
		m.render(t)
	}
	t.post()
	return nil
}

type tangle struct {
	*webScanner
	codeNames   []string
	codes       []*code
	codesByName map[string][]*code
	definitions map[string]*definition // @d
	formats     map[string]*format     // @f
	macroArgs   []func() knuth.RuneSource
	modules     []*module
	pascal      io.Writer
	pascal0     bytes.Buffer
	pool        io.Writer
	poolSum     int
	strings     map[string]int

	constCount       int
	constInjectState int
	metaCommentLevel int
}

func newTangle(pascal, pool io.Writer, src knuth.RuneSource) *tangle {
	return &tangle{
		codesByName: map[string][]*code{},
		definitions: map[string]*definition{},
		formats:     map[string]*format{},
		pascal:      pascal,
		pool:        pool,
		poolSum:     poolSumInit,
		strings:     map[string]int{},
		webScanner:  newWebScanner(src),
	}
}

func (t *tangle) post() {
	tagJoin := []byte("@&")
	tagCheckSum := []byte("@$")
	nl3 := []byte("\n\n\n")
	b := t.pascal0.Bytes()
	for len(b) != 0 {
		x := bytes.Index(b, tagJoin)
		if x < 0 {
			break
		}

		c := b[:x]
		b = b[x+len(tagJoin):]
		c = bytes.TrimRight(c, blankSet)
		if _, err := t.pascal.Write(c); err != nil {
			panic(todo("", err))
		}

		b = bytes.TrimLeft(b, blankSet)
	}
	for bytes.Index(b, nl3) >= 0 {
		b = bytes.ReplaceAll(b, nl3, nl3[:2])
	}
	for bytes.Index(b, tagCheckSum) >= 0 {
		b = bytes.ReplaceAll(b, tagCheckSum, []byte(fmt.Sprintf(" %d ", t.poolSum)))
	}
	if _, err := t.pascal.Write(b); err != nil {
		panic(todo("", err))
	}
	if _, err := fmt.Fprintf(t.pool, "*%09d\n", t.poolSum); err != nil {
		panic(todo("", err))
	}
}

func (t *tangle) w(s string, args ...interface{}) {
	b := []byte(fmt.Sprintf(s, args...))
	if oTrcw {
		os.Stderr.Write(b)
	}
	var w []byte
	for i := 0; i < len(b); {
		c := b[i]
		switch c {
		case '@':
			if i+1 == len(b) {
				break
			}

			switch b[i+1] {
			case '{':
				t.metaCommentLevel++
				if t.metaCommentLevel == 1 {
					w = append(w, '{')
				} else {
					w = append(w, '[')
				}
				i += 2
				continue
			case '}':
				if t.metaCommentLevel == 0 {
					panic(todo("%v:", t.position()))
				}

				t.metaCommentLevel--
				if t.metaCommentLevel == 0 {
					w = append(w, '}')
				} else {
					w = append(w, ']')
				}
				i += 2
				continue
			case '@':
				i++
			case '&':
				w = append(w, "@&"...)
				i += 2
				continue
			}
		case '{':
			if t.metaCommentLevel != 0 {
				c = '['
			}
		case '}':
			if t.metaCommentLevel != 0 {
				c = ']'
			}
		}
		i++
		w = append(w, c)
	}
	b = w
	// fmt.Printf("%s", b)
	if _, err := t.pascal0.Write(b); err != nil {
		panic(abort(fmt.Errorf("%v: writing tangle result: %v", t.src.Position(), err)))
	}
}

func (t *tangle) push(src knuth.RuneSource) {
	if len(t.stack) == stackLimit {
		panic(todo("", t.srcStack()))
	}

	if t.src != nil {
		t.stack = append(t.stack, t.src)
	}
	t.src = src
}

func (t *tangle) pushCode(c *code) {
	src := knuth.NewRuneSource(c.pos.Filename, []byte(c.pascal), knuth.Unicode)
	src.AddLineColumnInfo(0, c.pos.Filename, c.pos.Line, c.pos.Column)
	t.push(src)
}

func (t *tangle) scan() error {
	t.scanLimbo()
	for {
		switch c := t.c(); c {
		case eof:
			return nil
		case ccNewStarredModule:
			t.consume() // "@*"
			t.addModule(t.scanModule(t.scanModuleNameDot()))
		case ccNewModule:
			t.consume() // "@ "
			t.addModule(t.scanModule(""))
		default:
			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) addModule(m *module) {
	t.modules = append(t.modules, m)
	m.number = len(t.modules)
}

func (t *tangle) scanModuleNameDot() string {
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case '.':
			t.consume()
			return strings.TrimSpace(b.String())
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanModule(nm string) *module {
	m := &module{name: nm}
	m.pos, m.tex = t.scanTeX()
	for {
		switch c := t.c(); c {
		case ccDefinition:
			t.addDefinition(t.scanDefinition())
		case ccNewStarredModule, ccNewModule, eof:
			return m
		case ccBeginPascal:
			t.consume()
			pos, s := t.scanPascal(false, true)
			c := &code{pos: pos, pascal: s, inModule: m}
			m.codes = append(m.codes, c)
		case ccModuleName:
			nm := t.scanModuleName()
			t.addCodeName(nm)
			t.scanBlank()
			if t.c() == '+' {
				t.consume() // handle "+=" the same as "=".
			}
			switch c := t.c(); c {
			case '=':
				t.consume()
				pos, s := t.scanPascal(false, true)
				c := &code{name: nm, pos: pos, pascal: s, inModule: m}
				m.codes = append(m.codes, c)
			default:
				panic(todo("%v: %#U %#U", t.position(), c, t.c2))
			}
		case ccFormat:
			t.addFormat(t.scanFormat())
		default:
			if c >= 0 {
				panic(todo("%v: %#U %#U", t.position(), c, t.c2))
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
		t.scanSeparator()
	}
}

func (t *tangle) completeName(nm string) string {
	if strings.HasSuffix(nm, incompleteNameTag) {
		nm0 := nm
		nm = nm[:len(nm)-len(incompleteNameTag)]
		x := sort.SearchStrings(t.codeNames, nm)
		if x == len(t.codeNames) {
			panic(todo("%q %q", nm, t.codeNames))
		}

		nm1 := nm
		nm = t.codeNames[x]
		if !strings.HasPrefix(nm, nm1) {
			for i, v := range t.codeNames {
				trc("%3d: %q", i, v)
			}
			panic(todo("%q -> %q -> %q, x %d", nm0, nm1, nm, x))
		}
	}
	return nm
}

func (t *tangle) findCode(nm string) (r []*code) {
	nm = t.completeName(nm)
	r = t.codesByName[nm]
	if r == nil {
		panic(todo("%q", nm))
	}

	return r
}

func (t *tangle) addCodeName(nm string) {
	if strings.HasSuffix(nm, incompleteNameTag) {
		return
	}

	t.codeNames = append(t.codeNames, nm)
}

func (t *tangle) addFormat(f *format) {
	if ex, ok := t.formats[f.l]; ok {
		panic(todo("%v: %q redefined, previous at %v:", f.pos, f.l, ex.pos))
	}

	t.formats[f.l] = f
}

func (t *tangle) addDefinition(d *definition) {
	if ex, ok := t.definitions[d.name]; ok {
		panic(todo("%v: %q redefined, previous at %v:", d.pos, d.name, ex.pos))
	}

	d.ord = len(t.definitions)
	t.definitions[d.name] = d
	if d.kind == "=" {
		t.constCount++
	}
}

func (t *tangle) scanFormat() *format {
	t.consume()
	f := &format{pos: t.position()}
	_, _, f.l = t.scanIdentifier()
	t.scanBlank()
	switch c := t.c(); c {
	case '=':
		t.consume()
		switch c := t.c(); c {
		case '=':
			t.consume()
			_, _, f.r = t.scanIdentifier()
			f.postSep = t.scanSeparator()
		default:
			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	default:
		panic(todo("%v: %#U %#U", t.position(), c, t.c2))
	}
	return f
}

func (t *tangle) scanDefinition() (r *definition) {
	t.consume()
	d := &definition{pos: t.position()}

	defer func() {
		if r != nil {
			r.replacement = strings.TrimRight(r.replacement, blankSet) + " "
		}
	}()

	_, _, d.name = t.scanIdentifier()
	t.scanBlank()
	switch c := t.c(); c {
	case '=':
		t.consume()
		switch c := t.c(); c {
		case '=':
			t.consume()
			d.kind = "=="
			d.replPos = t.position()
			_, d.replacement = t.scanPascal(true, false)
		default:
			d.kind = "="
			d.replPos = t.position()
			_, d.replacement = t.scanPascal(true, true)
		}
	case '(':
		t.consume()
		t.scanBlank()
		switch c := t.c(); c {
		case '#':
			t.consume()
			t.scanBlank()
			switch c := t.c(); c {
			case ')':
				t.consume()
				t.scanBlank()
				switch c := t.c(); c {
				case '=':
					t.consume()
					switch c := t.c(); c {
					case '=':
						t.consume()
						d.kind = "(#)"
						d.replPos = t.position()
						_, d.replacement = t.scanPascal(true, false)
					default:
						panic(todo("%v: %#U %#U", t.position(), c, t.c2))
					}
				default:
					panic(todo("%v: %#U %#U", t.position(), c, t.c2))
				}
			default:
				panic(todo("%v: %#U %#U", t.position(), c, t.c2))
			}
		default:
			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	default:
		panic(todo("%v: %#U %#U", t.position(), c, t.c2))
	}
	return d
}

func (t *tangle) scanPascal(def, sep bool) (pos token.Position, s string) {
	pos = t.position()
	var b strings.Builder
	for {
		s := t.scanSeparator()
		if s != "" {
			switch {
			case sep:
				b.WriteString(s)
			default:
				b.WriteByte(' ')
			}
		}
		switch c := t.c(); c {
		case '\'':
			_, s := t.scanPascalStringLiteral()
			b.WriteString(s)
		case '"':
			_, s := t.scanQuotedStringLiteral()
			b.WriteString(s)
		case ccNewStarredModule, ccDefinition, ccNewModule, ccBeginPascal, ccFormat, eof:
			return pos, b.String()
		case ccModuleName:
			if def {
				return pos, b.String()
			}

			nm := t.scanModuleName()
			t.addCodeName(nm)
			fmt.Fprintf(&b, "@<%s@>", nm)
		case ccLineBreak:
			t.consume()
			b.WriteRune('\n')
		case ccOctal:
			_, sep, s := t.scanOctal()
			n, err := strconv.ParseUint(s, 8, 64)
			if err != nil {
				panic(todo("", err))
			}

			fmt.Fprintf(&b, "%s{0%o=}%[2]d", sep, n)
		case ccHex:
			_, sep, s := t.scanHex()
			n, err := strconv.ParseUint(s, 16, 64)
			if err != nil {
				panic(todo("", err))
			}

			fmt.Fprintf(&b, "%s{0x%x=}%[2]d", sep, n)
		case ccXrefRoman, ccXrefTypewriter, ccXrefWildcard:
			_, s := t.scanXref()
			fmt.Fprintf(&b, "{ %s }", commentSafe(s))
		case ccTeXString:
			fmt.Fprintf(&b, "{ %s }", commentSafe(t.scanTeXString()))
		case ccJoin:
			t.consume()
			b.WriteString("@&")
		case ccCheckSum:
			t.consume()
			b.WriteString("@$")
		case ccVerbatim:
			b.WriteString(t.scanVerbatim())
		case ccBeginMetaComment:
			t.consume()
			b.WriteString("@{")
		case ccEndMetaComment:
			t.consume()
			b.WriteString("@}")
		default:
			if c >= 0 {
				switch {
				case t.isIdentFirst(c):
					_, _, s := t.scanIdentifier()
					b.WriteString(s)
				case t.isDigit(c):
					_, _, s := t.scanDecimal()
					b.WriteString(s)
				default:
					switch c {
					case
						'(', '#', ')', ',', ';', ':', '=',
						'+', '-', '[', ']', '>', '.', '*',
						'<', '/', '^', '$':
						t.consume()
						b.WriteRune(c)
					default:
						panic(todo("%v: %#U %#U", t.position(), c, t.c2))
					}
				}
				break
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanVerbatim() string {
	t.consume()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ccEnd:
			t.consume()
			return b.String()
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanHex() (pos token.Position, sep, n string) {
	t.consume() // "@\""
	sep = t.scanSeparator()
	pos = t.position()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		default:
			if c >= 0 {
				if t.isDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' {
					b.WriteRune(c)
					t.consume()
					continue
				}

				return pos, sep, b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanOctal() (pos token.Position, sep, n string) {
	t.consume() // "@'"
	sep = t.scanSeparator()
	pos = t.position()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		default:
			if c >= 0 {
				if c >= '0' && c <= '7' {
					b.WriteRune(c)
					t.consume()
					continue
				}

				return pos, sep, b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanNumeric() (pos token.Position, sep, n string) {
	sep = t.scanSeparator()
	pos = t.position()
	var b strings.Builder
	hex := false
	for {
		switch c := t.c(); c {
		case ccNoLineBreak:
			return pos, sep, b.String()
		default:
			if c >= 0 {
				if c >= '0' && c <= '9' || hex && (c >= 'a' && c < 'f') {
					b.WriteRune(c)
					t.consume()
					continue
				}

				if b.String() == "0" && c == 'x' {
					hex = true
					b.WriteRune(c)
					t.consume()
					continue
				}

				return pos, sep, b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanDecimal() (pos token.Position, sep, n string) {
	sep = t.scanSeparator()
	pos = t.position()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ccNoLineBreak:
			return pos, sep, b.String()
		default:
			if c >= 0 {
				if c >= '0' && c <= '9' {
					b.WriteRune(c)
					t.consume()
					continue
				}

				return pos, sep, b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanModuleName() (r string) {
	t.consume() // "@<"
	var b strings.Builder
	var last rune
	for {
		c := t.c()
		switch c {
		case ' ', '\t', '\n', '\r':
			c = ' '
			if last == ' ' {
				t.consume()
				continue
			}
		}

		last = c
		switch c {
		case ccEnd:
			t.consume()
			r = strings.TrimSpace(b.String())
			return r
		case ccThinSpace:
			t.consume()
			b.WriteString("@,")
		case ccOctal:
			t.consume()
			b.WriteString("@'")
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanTeXInComment() (r string) {
	var b strings.Builder
	// pos := t.position()
	t.consume() // "$"
	b.WriteByte('$')
	for {
		switch c := t.c(); c {
		case '$':
			t.consume()
			b.WriteRune(c)
			r = commentSafe(b.String())
			return r
		case ccAt:
			t.consume()
			b.WriteString("@@")
		case ccHex:
			t.consume()
			b.WriteString("@\"")
		default:
			if c >= 0 {
				t.consume()
				b.WriteRune(c)
				break
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanPascalBracedComment() (r string) {
	// pos := t.position()
	var b strings.Builder
	lvl := 0
	for {
		switch c := t.c(); c {
		case '{':
			t.consume()
			lvl++
			b.WriteString("{")
		case '}':
			t.consume()
			lvl--
			b.WriteString("}")
			if lvl == 0 {
				s := b.String()
				return "{" + commentSafe(s[1:len(s)-1]) + "}"
			}
		case '$':
			b.WriteString(t.scanTeXInComment())
		case '\\':
			t.consume()
			b.WriteRune(c)
			switch c := t.c(); c {
			default:
				if c >= 0 {
					t.consume()
					b.WriteRune(c)
					break
				}

				panic(todo("%v: %#U %#U", t.position(), c, t.c2))
			}
		case ccTeXString:
			b.WriteString(t.scanTeXString())
		case ccOctal:
			t.consume()
			b.WriteString("@'")
		case ccAt:
			t.consume()
			b.WriteString("@@")
		case ccUnderline:
			t.consume()
			b.WriteString("@!")
		case ccHex:
			t.consume()
			b.WriteString("@\"")
		case ccBeginMetaComment:
			t.consume()
			b.WriteString(" ")
		case ccEndMetaComment:
			t.consume()
			b.WriteString(" ")
		default:
			if c >= 0 {
				t.consume()
				b.WriteRune(c)
				break
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanQuotedStringLiteral() (pos token.Position, s string) {
	var b strings.Builder
	pos = t.position()
	t.consume() // leading "\""
	b.WriteRune('"')
out:
	for {
		switch c := t.c(); c {
		case '"':
			t.consume()
			b.WriteRune(c)
			if t.c() != '"' {
				break out
			}

			t.consume()
			b.WriteRune(c)
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				t.consume()
				b.WriteRune(c)
				break
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
	s0 := b.String()
	s = s0[1 : len(s0)-1]
	s = strings.ReplaceAll(s, `""`, `"`)
	s = strings.ReplaceAll(s, `@@`, `@`)
	if a := []rune(s); len(a) == 1 {
		return pos, fmt.Sprintf("{%s=}%d", commentSafe(s0), a[0])
	}

	id := t.strings[s]
	if id == 0 {
		id = 256 + len(t.strings)
		t.strings[s] = id
		if _, err := fmt.Fprintf(t.pool, "%02d%s\n", len(s), s); err != nil {
			panic(todo("", err))
		}
		const prime = 03777777667
		t.poolSum += t.poolSum + len(s)
		for t.poolSum > prime {
			t.poolSum -= prime
		}
		for i := 0; i < len(s); i++ {
			t.poolSum += t.poolSum + int(s[i])
			for t.poolSum > prime {
				t.poolSum -= prime
			}
		}
	}
	return pos, fmt.Sprintf("{%s=}%d", commentSafe(s0), id)
}

func (t *tangle) scanPascalStringLiteral() (pos token.Position, s string) {
	var b strings.Builder
	t.consume() // leading "'"
	b.WriteRune('\'')
	for {
		switch c := t.c(); c {
		case '\'':
			t.consume()
			b.WriteRune(c)
			if t.c() != '\'' {
				return pos, b.String()
			}

			t.consume()
			b.WriteRune(c)
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				t.consume()
				b.WriteRune(c)
				break
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) isIdentFirst(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (t *tangle) isIdentNext(c rune) bool {
	return t.isIdentFirst(c) || t.isDigit(c)
}

func (t *tangle) isDigit(c rune) bool { return c >= '0' && c <= '9' }

func (t *tangle) scanIdentifier() (pos token.Position, sep, id string) {
	sep = t.scanSeparator()
	pos = t.position()
	var b strings.Builder
	first := true
	for {
		switch c := t.c(); c {
		default:
			if first && t.isIdentFirst(c) || t.isIdentNext(c) {
				first = false
				b.WriteRune(c)
				t.consume()
				continue
			}

			return pos, sep, b.String()
		}
	}
}

func (t *tangle) scanBlank() string {
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ' ', '\t', '\r':
			b.WriteByte(' ')
			t.consume()
		case '\n':
			b.WriteByte('\n')
			t.consume()
		default:
			return b.String()
		}
	}
}

func (t *tangle) scanSeparator() string {
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ' ', '\t', '\r':
			b.WriteByte(' ')
			t.consume()
		case '\n':
			b.WriteByte('\n')
			t.consume()
		case '{':
			b.WriteString(t.scanPascalBracedComment())
		case
			ccNewStarredModule, ccDefinition, ccNewModule, ccBeginPascal,
			ccFormat, eof, ccModuleName, ccOctal, ccXrefRoman,
			ccXrefTypewriter, ccXrefWildcard, ccTeXString,
			ccJoin, ccCheckSum, ccVerbatim, ccHex, ccBeginMetaComment, ccEndMetaComment:

			return b.String()
		case ccBigLineBreak, ccForceLine, ccLineBreak:
			t.consume()
			b.WriteString("\n")
		case ccThinSpace, ccNoUnderline, ccNoLineBreak, ccMathBreak, ccPseudoSemi, ccUnderline:
			t.consume()
			b.WriteByte(' ')
		default:
			if c >= 0 || c == eof {
				return b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanSeparator2() string {
	var b strings.Builder
outer:
	for {
		switch c := t.c(); c {
		case ' ', '\t', '\r':
			b.WriteByte(' ')
			t.consume()
		case '\n':
			b.WriteByte('\n')
			t.consume()
		case '{':
			t.consume()
			b.WriteRune(c)
			for {
				switch c := t.c(); c {
				case '}':
					t.consume()
					b.WriteRune(c)
					continue outer
				case eof:
					return b.String()
				case ccForceLine:
					t.consume()
					b.WriteByte('\n')
				case ccAt:
					t.consume()
					b.WriteString("@@")
				case ccOctal:
					t.consume()
					b.WriteString("@'")
				case ccUnderline:
					t.consume()
					b.WriteString("@!")
				case ccHex:
					t.consume()
					b.WriteString("@\"")
				default:
					if c >= 0 {
						t.consume()
						b.WriteRune(c)
						continue
					}

					panic(todo("%v: %#U %#U", t.position(), c, t.c2))
				}
			}
		case ccModuleName, ccPopMacroArg, ccBeginMetaComment, ccEndMetaComment:
			return b.String()
		default:
			if c >= 0 || c == eof {
				return b.String()
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanTeX() (pos token.Position, r string) {
	defer func() { r = strings.TrimRight(r, blankSet) }()

	pos = t.position()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case eof:
			return pos, b.String()
		case ccXrefRoman, ccXrefTypewriter, ccXrefWildcard:
			_, s := t.scanXref()
			b.WriteString(s)
		case ccNewStarredModule, ccDefinition, ccNewModule, ccBeginPascal, ccModuleName, ccFormat:
			return pos, b.String()
		case ccTeXString:
			b.WriteString(t.scanTeXString())
		case ccOctal:
			t.consume()
			b.WriteString("@'")
		case ccHex:
			t.consume()
			b.WriteString("@\"")
		case ccThinSpace, ccUnderline:
			t.consume()
			b.WriteString(" ")
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanTeXString() string {
	t.consume()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ccEnd:
			t.consume()
			return b.String()
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanXref() (pos token.Position, s string) {
	pos = t.position()
	t.consume()
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ccEnd:
			t.consume()
			return pos, fmt.Sprintf("\\xref[%s]", b.String())
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), c, t.c2))
		}
	}
}

func (t *tangle) scanLimbo() {
	var b strings.Builder
	for {
		switch c := t.c(); c {
		case ccNewStarredModule, ccNewModule, eof:
			if b.Len() == 0 {
				return
			}

			t.w("{ %s }", commentSafe(b.String()))
			return
		case ccAt:
			t.consume()
			b.WriteString("@@")
		default:
			if c >= 0 {
				b.WriteRune(c)
				t.consume()
				continue
			}

			panic(todo("%v: %#U", t.position(), t.c2))
		}
	}
}

func (t *tangle) injectConstants(appending bool) {
	if !appending {
		t.w("\n\nconst")
	}
	var a []*definition
	for _, d := range t.definitions {
		if d.kind == "=" {
			a = append(a, d)
		}
	}
	sort.Slice(a, func(i, j int) bool { return a[i].ord < a[j].ord })
	for _, d := range a {
		t.w("\n  %s = %s;", d.name, d.replacement)
	}
	t.w("\n")
}

type module struct {
	name  string // Non blank for @* modules only.
	codes []*code
	pos   token.Position
	tex   string

	number int

	teXRendered bool
}

func (m *module) render(t *tangle) {
	if len(m.codes) == 0 {
		m.renderTeX(t)
		return
	}

	for _, c := range m.codes {
		if c.name == "" {
			c.render(t)
		}
	}
}

func (m *module) renderTeX(t *tangle) {
	if m.teXRendered {
		return
	}

	t.w("\n")
	switch s := strings.TrimSpace(m.name); {
	case s != "":
		t.w("\n{ %d. %s }", m.number, commentSafe(s))
	default:
		t.w("\n{ %d. }", m.number)
	}
	if s := strings.TrimSpace(m.tex); s != "" {
		pos := m.pos
		pos.Line--
		t.w("\n\n{tangle:pos %v: }", m.pos)
		t.w("\n\n{ %s }", commentSafe(s))
	}
	m.teXRendered = true
}

type format struct {
	l, r    string // l == r
	pos     token.Position
	postSep string
}

type code struct {
	inModule *module
	name     string
	pos      token.Position
	pascal   string
}

func (c *code) render(t *tangle) {
	c.inModule.renderTeX(t)
	t.pushCode(c)
	c.scan(t)
}

func (c *code) scan(t *tangle) {
	const (
		injZero = iota
		injSeenProgram
		injDone
	)
	for {
		switch ch := t.c(); ch {
		case ' ', '\n', '\t', '\r', '{':
			s := t.scanSeparator2()
			t.w("%s", s)
		case
			'(', ',', ')', ';', '=', '.', ':', '[', ']', '+', '-', '>',
			'*', '<', '/', '$', '^':
			t.consume()
			t.w("%c", ch)
		case '\'':
			_, s := t.scanPascalStringLiteral()
			t.w("%s", s)
		case '"':
			_, s := t.scanQuotedStringLiteral()
			t.w("%s", s)
		case ccModuleName:
			nm := t.scanModuleName()
			codes := t.findCode(nm)
			t.w("\n{ %s }", commentSafe(nm))
			for i := len(codes) - 1; i >= 0; i-- {
				t.pushCode(codes[i])
			}
		case '#':
			if len(t.macroArgs) == 0 {
				panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
			}

			t.consume()
			t.push(t.macroArgs[len(t.macroArgs)-1]())
		case eof:
			return
		case ccPopMacroArg:
			t.consume()
			t.macroArgs = t.macroArgs[:len(t.macroArgs)-1]
		case ccBeginMetaComment:
			t.consume()
			t.w("@{")
		case ccEndMetaComment:
			t.consume()
			t.w("@}")
		case ccJoin:
			t.consume()
			t.w("@&")
		case ccCheckSum:
			t.consume()
			t.w("@$")
		default:
			if ch >= 0 {
				switch {
				case t.isIdentFirst(ch):
					_, _, id := t.scanIdentifier()
					switch d := t.definitions[id]; {
					case d != nil && d.kind != "=":
						c.expand(t, d)
					default:
						if t.constCount != 0 {
							switch t.constInjectState {
							case injZero:
								if id == "program" {
									t.constInjectState = injSeenProgram
								}
							case injSeenProgram:
								switch id {
								case "const":
									t.w("\nconst\n")
									t.injectConstants(true)
									t.constInjectState = injDone
									continue
								case "type", "var", "procedure", "function", "begin":
									t.injectConstants(false)
									t.constInjectState = injDone
									t.w("\n")
								}
							}
						}
						t.w("%s", id)
					}
				case t.isDigit(ch):
					_, _, s := t.scanDecimal()
					t.w("%s", s)
				default:
					panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
				}
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
		}
	}
}

func (c *code) expand(t *tangle, d *definition) {
	switch d.kind {
	case "=", "==":
		repl := d.replacement
		replSrc := knuth.NewRuneSource(d.pos.Filename, []byte(repl), knuth.Unicode)
		p := d.replPos
		replSrc.AddLineColumnInfo(0, p.Filename, p.Line, p.Column)
		t.push(replSrc)
	case "(#)":
	out:
		for {
			switch ch := t.c(); ch {
			case '(':
				break out
			case '\n', ' ':
				t.scanSeparator2()
			case ccPopMacroArg:
				t.consume()
				t.macroArgs = t.macroArgs[:len(t.macroArgs)-1]
			default:
				panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
			}
		}

		p := t.position()
		arg := c.scanMacroArg(t)
		t.macroArgs = append(t.macroArgs, func() knuth.RuneSource {
			argSrc := knuth.NewRuneSource(p.Filename, []byte(arg), knuth.Unicode)
			argSrc.AddLineColumnInfo(0, p.Filename, p.Line, p.Column)
			return argSrc
		})
		replSrc := knuth.NewRuneSource(d.pos.Filename, []byte(d.replacement+"@Z"), knuth.Unicode)
		p = d.replPos
		replSrc.AddLineColumnInfo(0, p.Filename, p.Line, p.Column)
		t.push(replSrc)
	default:
		panic(todo("%v: %q %q", d.pos, d.name, d.kind))
	}
}

func (c *code) scanMacroArg(t *tangle) string {
	var b strings.Builder
	lvl := 0
	for {
		switch ch := t.c(); ch {
		case '(':
			t.consume()
			if lvl != 0 {
				b.WriteRune(ch)
			}
			lvl++
		case ')':
			t.consume()
			lvl--
			if lvl == 0 {
				return b.String()
			}

			b.WriteRune(ch)
		case ccPopMacroArg:
			t.consume()
			b.WriteString("@Z")
		case '#':
			t.consume()
			t.push(t.macroArgs[len(t.macroArgs)-1]())
		case ' ', '\n', '{':
			s := t.scanSeparator2()
			b.WriteString(s)
		case
			',', ';', '=', '.', ':', '[', ']', '+', '-', '>', '*', '<',
			'/', '^':
			t.consume()
			b.WriteRune(ch)
		case '\'':
			_, s := t.scanPascalStringLiteral()
			b.WriteString(s)
		case '"':
			_, s := t.scanQuotedStringLiteral()
			b.WriteString(s)
		case ccCheckSum:
			t.consume()
			b.WriteString("@$")
		default:
			if ch >= 0 {
				switch {
				case t.isIdentFirst(ch):
					_, _, id := t.scanIdentifier()
					b.WriteByte(' ')
					b.WriteString(id)
				case t.isDigit(ch):
					_, _, s := t.scanNumeric()
					b.WriteByte(' ')
					b.WriteString(s)
				default:
					panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
				}
				continue
			}

			panic(todo("%v: %#U %#U", t.position(), ch, t.c2))
		}
	}
}

type definition struct {
	kind        string // "=", "==", "(#)"
	name        string
	ord         int
	pos         token.Position
	replPos     token.Position
	replacement string
}
