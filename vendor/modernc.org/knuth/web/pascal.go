// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web // modernc.org/knuth/web

import (
	"fmt"
	"go/token"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"modernc.org/mathutil"
	mtoken "modernc.org/token"
)

var (
	_ node = (*pasToken)(nil)

	pasKeywords = map[string]token.Token{
		"and":       tokAnd,
		"array":     tokArray,
		"begin":     tokBegin,
		"case":      tokCase,
		"const":     tokConst,
		"div":       tokDiv,
		"do":        tokDo,
		"downto":    tokDownto,
		"else":      tokElse,
		"end":       tokEnd,
		"endif":     tokEndif,
		"endifn":    tokEndifn,
		"file":      tokFile,
		"for":       tokFor,
		"function":  tokFunction,
		"goto":      tokGoto,
		"if":        tokIf,
		"ifdef":     tokIfdef,
		"ifndef":    tokIfndef,
		"in":        tokIn,
		"label":     tokLabel,
		"mod":       tokMod,
		"nil":       tokNil,
		"noreturn":  tokNoreturn,
		"not":       tokNot,
		"of":        tokOf,
		"or":        tokOr,
		"packed":    tokPacked,
		"procedure": tokProcedure,
		"program":   tokProgram,
		"record":    tokRecord,
		"repeat":    tokRepeat,
		"set":       tokSet,
		"then":      tokThen,
		"to":        tokTo,
		"type":      tokType,
		"until":     tokUntil,
		"var":       tokVar,
		"while":     tokWhile,
		"with":      tokWith,
	}

	tokStr = map[rune]string{
		tokEOF:       "tokEOF",
		tokAssign:    "tokAssign",
		tokEllipsis:  "tokEllipsis",
		tokFloat:     "tokFloat",
		tokGeq:       "tokGeq",
		tokIdent:     "tokIdent",
		tokInt:       "tokInt",
		tokInvalid:   "tokInvalid",
		tokLeq:       "tokLeq",
		tokNeq:       "tokNeq",
		tokString:    "tokString",
		tokZero:      "tokZero",
		tokAnd:       "tokAnd",
		tokArray:     "tokArray",
		tokBegin:     "tokBegin",
		tokCase:      "tokCase",
		tokConst:     "tokConst",
		tokDiv:       "tokDiv",
		tokDo:        "tokDo",
		tokDownto:    "tokDownto",
		tokElse:      "tokElse",
		tokEnd:       "tokEnd",
		tokEndif:     "tokEndif",
		tokEndifn:    "tokEndifn",
		tokFile:      "tokFile",
		tokFor:       "tokFor",
		tokFunction:  "tokFunction",
		tokGoto:      "tokGoto",
		tokIf:        "tokIf",
		tokIfdef:     "tokIfdef",
		tokIfndef:    "tokIfndef",
		tokIn:        "tokIn",
		tokLabel:     "tokLabel",
		tokMod:       "tokMod",
		tokNil:       "tokNil",
		tokNoreturn:  "tokNoreturn",
		tokNot:       "tokNot",
		tokOf:        "tokOf",
		tokOr:        "tokOr",
		tokPacked:    "tokPacked",
		tokProcedure: "tokProcedure",
		tokProgram:   "tokProgram",
		tokRecord:    "tokRecord",
		tokRepeat:    "tokRepeat",
		tokSet:       "tokSet",
		tokThen:      "tokThen",
		tokTo:        "tokTo",
		tokType:      "tokType",
		tokUntil:     "tokUntil",
		tokVar:       "tokVar",
		tokWhile:     "tokWhile",
		tokWith:      "tokWith",
	}
)

const (
	tokEOF = -iota - 1

	tokAssign
	tokEllipsis
	tokFloat
	tokGeq
	tokIdent
	tokInt
	tokInvalid
	tokLeq
	tokNeq
	tokString
	tokZero

	tokAnd
	tokArray
	tokBegin
	tokCase
	tokConst
	tokDiv
	tokDo
	tokDownto
	tokElse
	tokEnd
	tokEndif
	tokEndifn
	tokFile
	tokFor
	tokFunction
	tokGoto
	tokIf
	tokIfdef
	tokIfndef
	tokIn
	tokLabel
	tokMod
	tokNil
	tokNoreturn
	tokNot
	tokOf
	tokOr
	tokPacked
	tokProcedure
	tokProgram
	tokRecord
	tokRepeat
	tokSet
	tokThen
	tokTo
	tokType
	tokUntil
	tokVar
	tokWhile
	tokWith
)

const (
	charAlign    = 1
	charSize     = 1
	integerAlign = 4
	integerSize  = 4
	real32Align  = 4
	real32Size   = 4
	realAlign    = 8
	realSize     = 8
)

// node is an item of the CST tree.
type node interface {
	Position() token.Position
}

// pasToken represents a lexeme, its position and its semantic value.
type pasToken struct { // 16 bytes on 64 bit arch
	source *pasSource

	ch    int32
	index int32
}

func (t pasToken) String() string { return fmt.Sprintf("%v: %q %#U", t.Position(), t.Src(), t.ch) }

// isValid reports t is a valid token. Zero value reports false.
func (t pasToken) isValid() bool { return t.source != nil }

// Next returns the token following t or a zero value if no such token exists.
func (t pasToken) Next() (r pasToken) {
	if t.source == nil {
		return r
	}

	if index := t.index + 1; index < int32(len(t.source.toks)) {
		s := t.source
		return pasToken{source: s, ch: s.toks[index].ch, index: index}
	}

	return r
}

// Sep returns any separators, combined, preceding t.
func (t pasToken) Sep() string {
	if t.source == nil {
		return ""
	}

	s := t.source
	if p, ok := s.sepPatches[t.index]; ok {
		return p
	}

	return string(s.buf[s.toks[t.index].sep:s.toks[t.index].src])
}

// Src returns t's source form.
func (t pasToken) Src() string {
	if t.source == nil {
		return ""
	}

	s := t.source
	if p, ok := s.srcPatches[t.index]; ok {
		return p
	}

	if t.ch != int32(tokEOF) {
		next := t.source.off
		if t.index < int32(len(s.toks))-1 {
			next = s.toks[t.index+1].sep
		}
		return string(s.buf[s.toks[t.index].src:next])
	}

	return ""
}

// Positions implements Node.
func (t pasToken) Position() (r token.Position) {
	if t.source == nil {
		return r
	}

	s := t.source
	off := mathutil.MinInt32(int32(len(s.buf)), s.toks[t.index].src)
	return token.Position(s.file.PositionFor(mtoken.Pos(s.base+off), true))
}

type tok struct { // 12 bytes
	ch  int32
	sep int32
	src int32
}

// pasSource represents a single Go pasSource file, editor text buffer etc.
type pasSource struct {
	buf        []byte
	file       *mtoken.File
	name       string
	sepPatches map[int32]string
	srcPatches map[int32]string
	toks       []tok

	base int32
	off  int32
}

// 'buf' becomes owned by the result and must not be modified afterwards.
func newPasSource(name string, buf []byte) *pasSource {
	file := mtoken.NewFile(name, len(buf))
	return &pasSource{
		buf:  buf,
		file: file,
		name: name,
		base: int32(file.Base()),
	}
}

type errWithPosition struct {
	pos token.Position
	err error
}

func (e errWithPosition) String() string {
	switch {
	case e.pos.IsValid():
		return fmt.Sprintf("%v: %v", e.pos, e.err)
	default:
		return fmt.Sprintf("%v", e.err)
	}
}

type errList []errWithPosition

func (e errList) Err() (r error) {
	if len(e) == 0 {
		return nil
	}

	return e
}

func (e errList) Error() string {
	w := 0
	prev := errWithPosition{pos: token.Position{Offset: -1}}
	for _, v := range e {
		if v.pos.Line == 0 || v.pos.Offset != prev.pos.Offset || v.err.Error() != prev.err.Error() {
			e[w] = v
			w++
			prev = v
		}
	}

	var a []string
	for _, v := range e {
		a = append(a, fmt.Sprint(v))
	}
	return strings.Join(a, "\n")
}

func (e *errList) err(pos token.Position, msg string, args ...interface{}) {
	switch {
	case len(args) == 0:
		*e = append(*e, errWithPosition{pos, fmt.Errorf("%s", msg)})
	default:
		*e = append(*e, errWithPosition{pos, fmt.Errorf(msg, args...)})
	}
}

type pasScanner struct {
	*pasSource
	dir  string
	errs errList
	opts *genOpts
	skip int // < 0: skip tokens due to ifdef/ifndef
	tok  tok

	last int32

	errBudget int

	c byte // Look ahead byte.

	eof      bool
	isClosed bool
}

func newPasScanner(name string, buf []byte, opts *genOpts) *pasScanner {
	dir, _ := filepath.Split(name)
	r := &pasScanner{
		dir:       dir,
		errBudget: 10,
		pasSource: newPasSource(name, buf),
		opts:      opts,
	}
	switch {
	case len(buf) == 0:
		r.eof = true
	default:
		r.c = buf[0]
		if r.c == '\n' {
			r.file.AddLine(int(r.base + r.off))
		}
	}
	return r
}

func (s *pasScanner) isDigit(c byte) bool      { return c >= '0' && c <= '9' }
func (s *pasScanner) isIDNext(c byte) bool     { return s.isIDFirst(c) || s.isDigit(c) }
func (s *pasScanner) isOctalDigit(c byte) bool { return c >= '0' && c <= '7' }

func (s *pasScanner) isHexDigit(c byte) bool {
	return s.isDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func (s *pasScanner) isIDFirst(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (s *pasScanner) position() token.Position {
	return token.Position(s.pasSource.file.PositionFor(mtoken.Pos(s.base+s.off), true))
}

func (s *pasScanner) pos(off int32) token.Position {
	return token.Position(s.file.PositionFor(mtoken.Pos(s.base+off), true))
}

func (s *pasScanner) close() {
	if s.isClosed {
		return
	}

	s.tok.ch = int32(tokInvalid)
	s.eof = true
	s.isClosed = true
}

func (s *pasScanner) token() pasToken {
	return pasToken{source: s.pasSource, ch: s.tok.ch, index: int32(len(s.toks) - 1)}
}

func (s *pasScanner) peek(off int) rune {
	if n := s.off + int32(off); n < int32(len(s.buf)) {
		return rune(s.buf[n])
	}

	return -1
}

func (s *pasScanner) next() {
	if s.eof {
		return
	}

	s.off++
	if int(s.off) == len(s.buf) {
		s.c = 0
		s.eof = true
		return
	}

	s.c = s.buf[s.off]
	if s.c == '\n' {
		s.file.AddLine(int(s.base + s.off))
	}
}

func (s *pasScanner) err(off int32, msg string, args ...interface{}) {
	if s.errBudget <= 0 {
		s.close()
		return
	}

	s.errBudget--
	if n := int32(len(s.buf)); off >= n {
		off = n
	}
	s.errs.err(s.pos(off), msg, args...)
}

func (s *pasScanner) scan() (r bool) {
	if s.isClosed {
		return false
	}

	s.last = s.tok.ch
	r = s.preproc()
	s.toks = append(s.toks, s.tok)
	return r
}

func (s *pasScanner) must(ch rune) (r bool) {
again:
	s.tok.ch = tokZero
	r = s.scan0()
	if !r {
		s.err(s.tok.src, "unexpected EOF")
		return false
	}

	if s.tok.ch == tokZero {
		goto again
	}

	if s.tok.ch != ch {
		s.err(s.tok.src, "unexpected token %#U, expected %#U", s.tok.ch, ch)
		s.close()
		return false
	}

	return true
}

func (s *pasScanner) mustParenString() (str string, r bool) {
	if s.must('(') && s.must(tokString) {
		str = string(s.buf[s.tok.src+1 : s.off-1])
		return str, s.must(')')
	}

	return "", false
}

func (s *pasScanner) preproc() (r bool) {
again0:
	s.tok.sep = s.off
	s.tok.ch = tokZero
again1:
	switch r = s.scan0(); {
	case !r: // eof
		return r
	case s.tok.ch == tokZero:
		goto again1
	}

	switch s.tok.ch {
	case tokEndif:
		nm, ok := s.mustParenString()
		if !ok {
			return false
		}

		if _, ok := s.opts.defs[nm]; !ok {
			s.skip++
		}
		goto again0
	case tokEndifn:
		nm, ok := s.mustParenString()
		if !ok {
			return false
		}

		if _, ok := s.opts.defs[nm]; ok {
			s.skip++
		}
		goto again0
	case tokIfdef:
		nm, ok := s.mustParenString()
		if !ok {
			return false
		}

		if _, ok := s.opts.defs[nm]; !ok {
			s.skip--
		}
		goto again0
	case tokIfndef:
		nm, ok := s.mustParenString()
		if !ok {
			return false
		}

		if _, ok := s.opts.defs[nm]; ok {
			s.skip--
		}
		goto again0
	default:
		if s.skip < 0 {
			goto again0
		}

		return r
	}
}

func (s *pasScanner) scan0() (r bool) {
	s.tok.src = mathutil.MinInt32(s.off, int32(len(s.buf)))
	switch s.c {
	case '{':
		off := s.off
		s.next()
		s.generalComment(off)
	case ' ', '\t', '\r', '\n':
		s.next()
	case '(', ',', ')', ';', '=', '[', ']', '+', '-', '*', '/', '^':
		s.tok.ch = rune(s.c)
		s.next()
	case '.':
		s.next()
		if s.c != '.' {
			s.tok.ch = '.'
			return true
		}

		s.next()
		s.tok.ch = int32(tokEllipsis)
	case ':':
		s.tok.ch = ':'
		s.next()
		if s.c == '=' {
			s.next()
			s.tok.ch = int32(tokAssign)
		}
	case '\'', '"':
		s.stringLiteral()
	case '>':
		s.next()
		switch s.c {
		case '=':
			s.next()
			s.tok.ch = int32(tokGeq)
		default:
			s.tok.ch = '>'
		}
	case '<':
		s.next()
		switch s.c {
		case '>':
			s.next()
			s.tok.ch = int32(tokNeq)
		case '=':
			s.next()
			s.tok.ch = int32(tokLeq)
		default:
			s.tok.ch = '<'
		}
	default:
		switch {
		case s.isIDFirst(s.c):
			s.next()
			s.identifierOrKeyword()
		case s.isDigit(s.c):
			s.numericLiteral()
		case s.eof:
			s.close()
			s.tok.ch = int32(tokEOF)
			s.tok.sep = mathutil.MinInt32(s.tok.sep, s.tok.src)
			return false
		default:
			s.err(s.off, "unexpected rune %#U", s.c)
			s.next()
		}
	}
	return true
}

func (s *pasScanner) stringLiteral() {
	// Leadind " or ' not consumed.
	off := s.off
	ch := s.c
	s.next()
	s.tok.ch = int32(tokString)
	for {
		switch {
		case s.c == ch:
			s.next()
			if s.c != ch {
				return
			}
		case s.c == '\n':
			fallthrough
		case s.eof:
			s.err(off, "string literal not terminated")
			return
		}
		s.next()
	}
}

func (s *pasScanner) numericLiteral() {
	// Leading decimal digit not consumed.
more:
	switch s.c {
	case '0':
		s.next()
		switch s.c {
		case 'x', 'X':
			s.tok.ch = int32(tokInt)
			s.next()
			if s.hexadecimals() == 0 {
				s.err(s.base+s.off, "hexadecimal literal has no digits")
				return
			}
		case '.':
			s.tok.ch = int32(tokInt)
			if s.isDigit(byte(s.peek(1))) {
				s.next()
				s.decimals()
				s.tok.ch = int32(tokFloat)
			}
			return
		default:
			invalidOff := int32(-1)
			var invalidDigit byte
			for {
				if s.isOctalDigit(s.c) {
					s.next()
					continue
				}

				if s.isDigit(s.c) {
					if invalidOff < 0 {
						invalidOff = s.off
						invalidDigit = s.c
					}
					s.next()
					continue
				}

				break
			}
			switch s.c {
			case '.', 'e', 'E', 'i':
				break more
			}
			if s.isDigit(s.c) {
				break more
			}
			if invalidOff > 0 {
				s.err(invalidOff, "invalid digit '%c' in octal literal", invalidDigit)
			}
			s.tok.ch = int32(tokInt)
			return
		}
	default:
		s.decimals()
		if s.c == '.' && s.isDigit(byte(s.peek(1))) {
			s.next()
			s.decimals()
			s.tok.ch = int32(tokFloat)
			return
		}
	}
	s.tok.ch = int32(tokInt)

}

func (s *pasScanner) hexadecimals() (r int) {
	for {
		switch {
		case s.isHexDigit(s.c):
			s.next()
			r++
		case s.c == '_':
			for n := 0; s.c == '_'; n++ {
				if n == 1 {
					s.err(s.off, "'_' must separate successive digits")
				}
				s.next()
			}
			if !s.isHexDigit(s.c) {
				s.err(s.off-1, "'_' must separate successive digits")
			}
		default:
			return r
		}
	}
}

func (s *pasScanner) decimals() (r int) {
	for {
		switch {
		case s.isDigit(s.c):
			s.next()
			r++
		default:
			return r
		}
	}
}

func (s *pasScanner) generalComment(off int32) {
	// Leading "{" consumed
	for {
		switch {
		case s.c == '}':
			s.next()
			return
		case s.eof:
			s.tok.ch = 0
			s.err(off, "comment not terminated")
			return
		default:
			s.next()
		}
	}
}

func (s *pasScanner) identifierOrKeyword() {
	for {
		switch {
		case s.isIDNext(s.c):
			s.next()
		default:
			if s.tok.ch = int32(pasKeywords[strings.ToLower(string(s.buf[s.tok.src:s.off]))]); s.tok.ch == 0 {
				s.tok.ch = int32(tokIdent)
			}
			return
		}
	}
}

func pasParse(s *pasScanner) (r *ast, err error) {
	p := newPasParser(s)
	p.scan()
	program := p.program()
	if p.c() != eof {
		panic(todo("", p.token()))
	}

	if err = p.errs.Err(); err != nil {
		return nil, err
	}

	return &ast{
		program: program,
		eof:     p.token(),
	}, nil
}

type pasParser struct {
	*pasScanner
	scope *scope
}

func newPasParser(s *pasScanner) (r *pasParser) {
	return &pasParser{
		pasScanner: s,
		scope:      newScope(nil),
	}
}

func (p *pasParser) c() rune                           { return p.tok.ch }
func (p *pasParser) mustIdent(ch rune) (r *identifier) { return &identifier{ident: p.must(ch)} }
func (p *pasParser) shift() pasToken                   { r := p.token(); p.scan(); return r }

func (p *pasParser) must(ch rune) (r pasToken) {
	if p.c() == ch {
		r = p.token()
	} else {
		panic(todo("", p.token()))
	}

	p.shift()
	return r
}

func (p *pasParser) opt(ch rune) (r pasToken) {
	if p.c() == ch {
		r = p.shift()
	}
	return r
}

type program struct {
	programHeading *programHeading
	semi           pasToken
	block          *block
	dot            pasToken
}

// Program = ProgramHeading ";" Block "." .
func (p *pasParser) program() (r *program) {
	switch p.c() {
	case tokProgram:
		r = &program{
			p.programHeading(),
			p.must(';'),
			p.block(),
			pasToken{},
		}
		if p.c() != tokEOF {
			r.dot = p.must('.')
		}
		return r
	default:
		panic(todo("", p.token()))
	}
}

type identifier struct {
	ident pasToken
	scope *scope

	resolvedTo node
	typ        interface{}
}

func (n identifier) Position() token.Position { return n.ident.Position() }
func (n identifier) Sep() string              { return n.ident.Sep() }
func (n identifier) Src() string              { return n.ident.Src() }

type programHeading struct {
	program              pasToken
	ident                *identifier
	programParameterList *programParameterList
}

// ProgramHeading = "program" Identifier [ ProgramParameterList ] .
func (p *pasParser) programHeading() (r *programHeading) {
	return &programHeading{
		p.must(tokProgram),
		p.mustIdent(tokIdent),
		p.programParameterList(),
	}
}

type programParameterList struct {
	lparen         pasToken
	identifierList []*identifierList
	rparen         pasToken

	idents []*identifier
}

// ProgramParameterList = "(" IdentifierList ")" .
func (p *pasParser) programParameterList() (r *programParameterList) {
	if p.c() != '(' {
		return nil
	}

	return &programParameterList{
		lparen:         p.must('('),
		identifierList: p.identifierList(),
		rparen:         p.must(')'),
	}
}

type identifierList struct {
	comma pasToken
	ident *identifier
}

// IdentifierList = Identifier { "," Identifier } .
func (p *pasParser) identifierList() (r []*identifierList) {
	r = []*identifierList{{ident: p.mustIdent(tokIdent)}}
	for p.c() == ',' {
		r = append(r, &identifierList{p.shift(), p.mustIdent(tokIdent)})
	}
	return r
}

type block struct {
	labelDeclarationPart                *labelDeclarationPart
	constantDefinitionPart              *constantDefinitionPart
	typeDefinitionPart                  *typeDefinitionPart
	variableDeclarationPart             *variableDeclarationPart
	procedureAndFunctionDeclarationPart []*procedureAndFunctionDeclarationPart
	statementPart                       *compoundStatement
}

func (n *block) Position() (r token.Position) {
	if n == nil {
		return r
	}

	if x := n.labelDeclarationPart; x != nil {
		return x.Position()
	}

	if x := n.constantDefinitionPart; x != nil {
		return x.Position()
	}

	if x := n.typeDefinitionPart; x != nil {
		return x.Position()
	}

	if x := n.variableDeclarationPart; x != nil {
		return x.Position()
	}

	if len(n.procedureAndFunctionDeclarationPart) != 0 {
		return n.procedureAndFunctionDeclarationPart[0].Position()
	}

	return n.statementPart.Position()
}

// Block = LabelDeclarationPart
//
//	ConstantDefinitionPart
//	TypeDefinitionPart
//	VariableDeclarationPart
//	ProcedureAndFunctionDeclarationPart
//	StatementPart .
func (p *pasParser) block() (r *block) {
	return &block{
		p.labelDeclarationPart(),
		p.constantDefinitionPart(),
		p.typeDefinitionPart(),
		p.variableDeclarationPart(),
		p.procedureAndFunctionDeclarationPart(),
		p.compoundStatement(),
	}
}

type typeDefinitionPart struct {
	type1              pasToken
	typeDefinitionList []*typeDefinitionList
}

func (n *typeDefinitionPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.type1.Position()
}

type typeDefinitionList struct {
	typeDefinition *typeDefinition
	semi           pasToken
}

// TypeDefinitionPart = [ "type" TypeDefinition ";" { TypeDefinition ";" } ] .
func (p *pasParser) typeDefinitionPart() (r *typeDefinitionPart) {
	if p.c() != tokType {
		return nil
	}

	r = &typeDefinitionPart{
		p.shift(),
		[]*typeDefinitionList{{p.typeDefinition(), p.must(';')}},
	}
	for {
		switch p.c() {
		case tokVar, tokBegin:
			return r
		case tokIdent:
			r.typeDefinitionList = append(r.typeDefinitionList, &typeDefinitionList{p.typeDefinition(), p.must(';')})
		default:
			panic(todo("", p.token()))
		}
	}
}

type typeDefinition struct {
	ident *identifier
	eq    pasToken
	type1 node

	typ interface{}
}

func (n *typeDefinition) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.ident.Position()
}

// TypeDefinition = Identifier "=" Type .
func (p *pasParser) typeDefinition() (r *typeDefinition) {
	r = &typeDefinition{
		ident: p.mustIdent(tokIdent),
		eq:    p.must('='),
		type1: p.type1(),
	}
	if err := p.scope.add(r.ident.Src(), r); err != nil {
		p.errs.err(r.ident.Position(), "%s", err)
	}
	return r
}

// Type = SimpleType | StructuredType | PointerType .
func (p *pasParser) type1() (r node) {
	switch p.c() {
	case tokInt, tokIdent, '-', tokString:
		return p.simpleType()
	case tokPacked, tokArray, tokRecord, tokFile:
		return p.structuredType()
	case '^':
		return &pointerType{p.shift(), p.type1()}
	default:
		panic(todo("", p.token()))
	}
}

type pointerType struct {
	carret pasToken
	type1  node
}

func (n *pointerType) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.carret.Position()
}

// PointerType = "^" DomainType | PointerTypeldentifier .

// StructuredType = [ "packed" ] UnpackedStructuredType | StructuredTypeIdentifier .
func (p *pasParser) structuredType() (r node) {
	packed := p.opt(tokPacked)
	switch p.c() {
	case tokArray, tokRecord, tokFile:
		return p.unpackedStructuredType(packed)
	default:
		panic(todo("", p.token()))
	}
}

// UnpackedStructuredType = ArrayType | RecordType | SetType | FileType .
func (p *pasParser) unpackedStructuredType(packed pasToken) (r node) {
	switch p.c() {
	case tokFile:
		return p.fileType(packed)
	case tokArray:
		return p.arrayType(packed)
	case tokRecord:
		return p.recordType(packed)
	default:
		panic(todo("", p.token()))
	}
}

type field struct {
	off int64
	typ interface{}

	inVariantPart bool
}

type recordType struct {
	packedTok pasToken
	record    pasToken
	fieldList *fieldList
	end       pasToken

	align int64
	size  int64

	fields map[string]*field

	checked     bool
	hasVariants bool
	packed      bool
}

func (n *recordType) isPacked() bool { return n.packed }

func (n *recordType) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.record.Position()
}

func (n *recordType) String() string {
	return fmt.Sprintf("record %s end", n.fieldList)
}

func (n *recordType) canAssignFrom(c *ctx, m *recordType) (r bool) {
	if n == nil || m == nil {
		return false
	}

	return n.fieldList.canAssignFrom(c, m.fieldList)
}

func (n *recordType) field(nm pasToken) *field {
	r := n.fields[nm.Src()]
	if r == nil {
		panic(todo("", nm))
	}

	return r
}

func (n *recordType) check(c *ctx) interface{} {
	if n == nil || n.checked {
		return nil
	}

	n.checked = true
	n.fields = map[string]*field{}
	n.size = n.fieldList.check(c, n, false, 0)
	if !n.packed {
		n.size = roundup(n.size, n.align)
	}
	return n
}

// RecordType = "record" FieldList "end" .
func (p *pasParser) recordType(packed pasToken) (r *recordType) {
	return &recordType{
		packedTok: packed,
		record:    p.must(tokRecord),
		fieldList: p.fieldList(),
		end:       p.must(tokEnd),
		size:      -1,
	}
}

type fieldList struct {
	fixedPart   []*fixedPart
	semi        pasToken
	variantPart *variantPart
	semi2       pasToken
}

func (n *fieldList) Position() (r token.Position) {
	if n == nil {
		return r
	}

	if len(n.fixedPart) != 0 {
		return n.fixedPart[0].Position()
	}

	return n.variantPart.Position()
}

func (n *fieldList) String() string {
	var a []string
	for _, v := range n.fixedPart {
		a = append(a, v.String())
	}
	return strings.Join(a, "; ")
}

func (n *fieldList) check(c *ctx, r *recordType, inVariantPart bool, off int64) int64 {
	if n == nil {
		return off
	}

	for _, v := range n.fixedPart {
		off = v.check(c, r, inVariantPart, off)
	}
	return n.variantPart.check(c, r, off)
}

func (n *fieldList) canAssignFrom(c *ctx, m *fieldList) (r bool) {
	if n == nil || m == nil || len(n.fixedPart) != len(m.fixedPart) {
		return false
	}

	for i, v := range n.fixedPart {
		if !v.canAssignFrom(c, m.fixedPart[i]) {
			return false
		}
	}

	return n.variantPart == nil && m.variantPart == nil || n.variantPart.canAssignFrom(c, m.variantPart)
}

// FieldList = [ ( FixedPart [ ";" VariantPart ] | VariantPart ) [ ";" ] ] .
func (p *pasParser) fieldList() (r *fieldList) {
	switch p.c() {
	case tokIdent:
		fp, semi := p.fixedPart()
		r = &fieldList{
			fixedPart: fp,
			semi:      semi,
		}
		switch p.c() {
		case tokEnd, ')':
			return r
		case tokCase:
			r.variantPart, r.semi2 = p.variantPart()
			return r
		default:
			panic(todo("", p.token()))
		}
	case tokCase:
		r = &fieldList{}
		r.variantPart, r.semi2 = p.variantPart()
		return r
	default:
		panic(todo("", p.token()))
	}
}

type recordSection struct {
	identifierList []*identifierList
	colon          pasToken
	type1          node

	typ interface{}
}

func (n *recordSection) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.identifierList[0].ident.Position()
}

func (n *recordSection) String() string {
	var a []string
	for _, v := range n.identifierList {
		a = append(a, strings.ToLower(v.ident.Src()))
	}
	return strings.Join(a, ", ") + ": " + typeStr(n.typ)
}

func (n *recordSection) canAssignFrom(c *ctx, m *recordSection) (r bool) {
	if n == nil || m == nil || len(n.identifierList) != len(m.identifierList) {
		return false
	}

	for i, v := range n.identifierList {
		if strings.ToLower(v.ident.Src()) != strings.ToLower(m.identifierList[i].ident.Src()) {
			return false
		}
	}

	return c.checkAssign(n, -1, n.typ, m.typ, nil)
}

func (n *recordSection) check(c *ctx, r *recordType, inVariantPart bool, off int64) int64 {
	if n == nil {
		return off
	}

	n.typ = c.checkType(n.type1)
	if inVariantPart {
		switch underlyingType(n.typ).(type) {
		case *pascalReal:
			n.typ = &pascalReal32{}
		}
	}
	for _, v := range n.identifierList {
		nm := v.ident.Src()
		if _, ok := r.fields[nm]; ok {
			panic(todo("", v.ident))
		}

		align, sz := c.sizeof(n.typ)
		r.align = mathutil.MaxInt64(r.align, align)
		if !r.packed {
			off = roundup(off, align)
		}
		r.fields[nm] = &field{typ: n.typ, inVariantPart: inVariantPart, off: off}
		off += sz
	}
	return off
}

func underlyingType(t interface{}) interface{} {
	switch x := t.(type) {
	case *typeDefinition:
		return x.typ
	default:
		return t
	}
}

func (c *ctx) sizeof(t interface{}) (align, size int64) {
	switch x := underlyingType(t).(type) {
	case *subrangeType:
		switch {
		case x.lo >= 0 && x.hi <= math.MaxUint8:
			return 1, 1
		case x.lo >= 0 && x.hi <= math.MaxUint16:
			return 2, 2
		case x.lo >= 0 && x.hi <= math.MaxUint32:
			return 4, 4
		case x.lo >= math.MinInt8 && x.hi <= math.MaxInt8:
			return 1, 1
		case x.lo >= math.MinInt16 && x.hi <= math.MaxInt16:
			return 2, 2
		case x.lo >= math.MinInt32 && x.hi <= math.MaxInt32:
			return 4, 4
		default:
			panic(todo("", x.lo, x.hi))
		}
	case *pascalInteger:
		return integerAlign, integerSize
	case *recordType:
		if !x.checked || x.size <= 0 || x.align <= 0 {
			panic(todo("", x, x.checked, x.size))
		}

		return x.align, x.size
	case *pascalReal:
		return realAlign, realSize
	case *pascalReal32:
		return real32Align, real32Size
	case *pascalChar:
		return charAlign, charSize
	default:
		panic(todo("%T", x))
	}
}

// RecordSection = IdentifierList ":" Type.
func (p *pasParser) recordSection() (r *recordSection) {
	return &recordSection{
		identifierList: p.identifierList(),
		colon:          p.must(':'),
		type1:          p.type1(),
	}
}

type variantPart struct {
	case1           pasToken
	variantSelector *variantSelector
	of              pasToken
	variants        []*variant
}

func (n *variantPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.case1.Position()
}

func (n *variantPart) check(c *ctx, r *recordType, off int64) int64 {
	if n == nil {
		return off
	}

	r.hasVariants = true
	n.variantSelector.check(c)
	switch x := underlyingType(n.variantSelector.typ).(type) {
	case *subrangeType, *pascalBoolean:
		// ok
	default:
		panic(todo("%v: %T", n.variantSelector.Position(), x))
	}
	off0 := off
	for _, v := range n.variants {
		off = mathutil.MaxInt64(off, v.check(c, n.variantSelector, r, off0))
	}
	return off
}

func (n *variantPart) canAssignFrom(c *ctx, m *variantPart) (r bool) {
	if n == nil || m == nil || len(n.variants) != len(m.variants) {
		return false
	}

	if !c.checkAssign(n, -1, n.variantSelector.typ, m.variantSelector.typ, nil) {
		return false
	}

	for i, v := range n.variants {
		if !v.canAssignFrom(c, m.variants[i]) {
			return false
		}
	}

	return true
}

// VariantPart = "case" VariantSelector "of" Variant { ";" Variant } .
func (p *pasParser) variantPart() (r *variantPart, semi pasToken) {
	r = &variantPart{
		case1:           p.must(tokCase),
		variantSelector: p.variantSelector(),
		of:              p.must(tokOf),
	}
	r.variants, semi = p.variants()
	return r, semi
}

type variant struct {
	semi      pasToken
	constList []*constList
	comma     pasToken
	lparen    pasToken
	fieldList *fieldList
	rparen    pasToken
}

func (n *variant) canAssignFrom(c *ctx, m *variant) (r bool) {
	if n == nil || m == nil || len(n.constList) != len(m.constList) {
		return false
	}

	return n.fieldList.canAssignFrom(c, m.fieldList)
}

func (n *variant) check(c *ctx, selector *variantSelector, r *recordType, off int64) int64 {
	if n == nil {
		return off
	}

	for _, v := range n.constList {
		v.check(c)
		c.checkAssign(v, -1, selector.typ, v.typ, v.val)
	}
	return n.fieldList.check(c, r, true, off)
}

// Variant = Constant { "," Constant } ":" "(" FieIdList ")" .
func (p *pasParser) variants() (r []*variant, semi pasToken) {
	r = []*variant{{
		pasToken{},
		p.constList(),
		p.must(':'),
		p.must('('),
		p.fieldList(),
		p.must(')'),
	}}
	for p.c() == ';' {
		semi2 := p.shift()
		switch p.c() {
		case tokInt, tokIdent:
			r = append(r, &variant{
				semi2,
				p.constList(),
				p.must(':'),
				p.must('('),
				p.fieldList(),
				p.must(')'),
			})
		case tokEnd:
			return r, semi2
		default:
			panic(todo("", p.token()))
		}
	}
	return r, semi
}

type variantSelector struct {
	tagField pasToken
	comma    pasToken
	tagType  node

	typ interface{}
}

func (n *variantSelector) Position() (r token.Position) {
	if n == nil {
		return r
	}

	if n.tagField.isValid() {
		return n.tagField.Position()
	}

	return n.tagType.Position()
}

func (n *variantSelector) check(c *ctx) {
	n.typ = c.checkType(n.tagType)
}

// VariantSelector = [ TagField ":"] TagType .
func (p *pasParser) variantSelector() (r *variantSelector) {
	switch p.c() {
	case tokIdent:
		id := &identifier{ident: p.shift()}
		switch p.c() {
		case tokOf:
			return &variantSelector{tagType: id}
		default:
			panic(todo("", p.token()))
		}
	default:
		panic(todo("", p.token()))
	}
}

type fixedPart struct {
	semi          pasToken
	recordSection *recordSection
}

func (n *fixedPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.recordSection.Position()
}

func (n *fixedPart) String() string {
	return n.recordSection.String()
}

func (n *fixedPart) check(c *ctx, r *recordType, inVariantPart bool, off int64) int64 {
	if n == nil {
		return off
	}

	return n.recordSection.check(c, r, inVariantPart, off)
}

func (n *fixedPart) canAssignFrom(c *ctx, m *fixedPart) (r bool) {
	if n == nil || m == nil {
		return false
	}

	return n.recordSection.canAssignFrom(c, m.recordSection)
}

// FixedPart = RecordSection { ";" RecordSection } .
func (p *pasParser) fixedPart() (r []*fixedPart, semi pasToken) {
	r = []*fixedPart{{recordSection: p.recordSection()}}
	for p.c() == ';' {
		semi := p.shift()
		switch p.c() {
		case tokIdent:
			r = append(r, &fixedPart{semi, p.recordSection()})
		case tokEnd, tokCase:
			return r, semi
		default:
			panic(todo("", p.token()))
		}
	}
	return r, semi
}

type arrayType struct {
	packedTok pasToken
	array     pasToken
	lbracket  pasToken
	typeList  []*typeList
	rbracket  pasToken
	of        pasToken
	elemType  node

	indexTyp []interface{}
	elemTyp  interface{}

	packed bool
}

func (n *arrayType) isPacked() bool { return n.packed }

func (n *arrayType) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.array.Position()
}

func (n *arrayType) canAssignFrom(c *ctx, m *arrayType) bool {
	if n == nil || m == nil || len(n.indexTyp) != len(m.indexTyp) || !c.checkAssign(n, -1, n.elemTyp, m.elemTyp, nil) {
		return false
	}

	for i, v := range n.indexTyp {
		if !c.checkAssign(n, -1, v, m.indexTyp[i], nil) {
			return false
		}
	}

	return true
}

// ArrayType = "array" "[" IndexType { "," IndexType } "]" "of" ComponentType .
func (p *pasParser) arrayType(packed pasToken) (r *arrayType) {
	return &arrayType{
		packedTok: packed,
		array:     p.must(tokArray),
		lbracket:  p.must('['),
		typeList:  p.typeList(),
		rbracket:  p.must(']'),
		of:        p.must(tokOf),
		elemType:  p.type1(),
	}
}

type typeList struct {
	comma pasToken
	type1 node
}

func (p *pasParser) typeList() (r []*typeList) {
	r = []*typeList{{type1: p.type1()}}
	for p.c() == ',' {
		r = append(r, &typeList{p.shift(), p.type1()})
	}
	return r
}

type fileType struct {
	packedTok pasToken
	file      pasToken
	of        pasToken
	type1     node

	elemType interface{}

	packed bool
}

func (n *fileType) isPacked() bool { return n.packed }

func (n *fileType) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.file.Position()
}

// FileType = "file" "of" Type .
func (p *pasParser) fileType(packed pasToken) (r *fileType) {
	return &fileType{
		packedTok: packed,
		file:      p.must(tokFile),
		of:        p.must(tokOf),
		type1:     p.type1(),
	}
}

// SimpleType = OrdinalType | RealTypeldentifier.
func (p *pasParser) simpleType() (r node) {
	switch p.c() {
	case tokInt, '-', tokString:
		return p.ordinalType()
	case tokIdent:
		id := &identifier{ident: p.shift()}
		switch p.c() {
		case tokEllipsis:
			return &subrangeType{
				constant:  id,
				ellipsis:  p.shift(),
				constant2: p.expression(),
			}
		default:
			return id
		}
	default:
		panic(todo("", p.token()))
	}
}

// OrdinalType = EnumeratedType | SubrangeType | OrdinalTypeldentifier .
func (p *pasParser) ordinalType() (r node) {
	switch p.c() {
	case tokInt, '-', tokString:
		return p.subrangeType()
	default:
		panic(todo("", p.token()))
	}
}

type subrangeType struct {
	constant  node
	ellipsis  pasToken
	constant2 node

	lo, hi int64
}

func (n *subrangeType) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.constant.Position()
}

// SubrangeType = Constant ".." Constant .
func (p *pasParser) subrangeType() (r *subrangeType) {
	return &subrangeType{
		constant:  p.expression(),
		ellipsis:  p.must(tokEllipsis),
		constant2: p.expression(),
	}
}

type compoundStatement struct {
	begin             pasToken
	statementSequence []*statementSequence
	end               pasToken
}

func (n *compoundStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.begin.Position()
}

// StatementPart = CompoundStatement .
// CompoundStatement = "begin" StatementSequence "end" .
func (p *pasParser) compoundStatement() (r *compoundStatement) {
	if p.c() == tokEOF {
		return nil
	}

	return &compoundStatement{
		p.must(tokBegin),
		p.statementSequence(),
		p.must(tokEnd),
	}
}

type statementSequence struct {
	semi      pasToken
	statement node
}

// StatementSequence = Statement { ";" Statement} .
func (p *pasParser) statementSequence() (r []*statementSequence) {
	switch p.c() {
	case
		tokIdent, tokIf, tokRepeat, tokWhile, tokBegin, tokCase, tokFor, tokGoto,
		tokInt, ';', tokEnd:

		r = []*statementSequence{{statement: p.statement(true)}}
	default:
		panic(todo("", p.token()))
	}
	for p.c() == ';' {
		r = append(r, &statementSequence{p.shift(), p.statement(true)})
	}
	return r
}

// Statement = [ Label ":" ] ( SimpleStatement | StructuredStatement ) .
func (p *pasParser) statement(acceptLabel bool) (r node) {
	switch p.c() {
	case tokIdent, tokEnd, tokGoto:
		return p.simpleStatement()
	case tokFor, tokBegin, tokIf, tokWhile, tokRepeat, tokCase:
		return p.structuredStatement()
	case tokUntil, ';':
		return &emptyStatement{p.token()}
	case tokInt:
		if !acceptLabel {
			panic(todo("", p.token()))
		}

		return &labeled{
			label: p.expression(),
			colon: p.must(':'),
			stmt:  p.statement(false),
		}
	case tokElse:
		return nil
	default:
		panic(todo("", p.token()))
	}
}

type labeled struct {
	label node
	plus  pasToken
	off   node
	colon pasToken
	stmt  node
}

func (n *labeled) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.label.Position()
}

func (n *labeled) check(c *ctx) {
	if n == nil {
		return
	}

	if n.off != nil {
		c.checkExpr(n.off)
	}
	c.checkStatement(n.stmt)
}

// StructuredStatement = CompoundStatement | ConditionalStatement
//
//	| RepetitiveStatement | WithStatement .
func (p *pasParser) structuredStatement() (r node) {
	switch p.c() {
	case tokFor, tokWhile, tokRepeat:
		return p.repetitiveStatement()
	case tokBegin:
		return p.compoundStatement()
	case tokIf:
		return p.ifStatement()
	case tokCase:
		return p.caseStatement()
	default:
		panic(todo("", p.token()))
	}
}

type caseStatement struct {
	case1    pasToken
	expr     node
	of       pasToken
	cases    []*caseList
	semi     pasToken
	else1    pasToken
	elseStmt node
	end      pasToken
}

func (n *caseStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.case1.Position()
}

func (n *caseStatement) check(c *ctx) {
	if n == nil {
		return
	}

	t, _ := c.checkExpr(n.expr)
	switch x := underlyingType(t).(type) {
	case *pascalInteger, *subrangeType:
		// ok
	default:
		panic(todo("%v: %T", n.expr.Position(), x))
	}
	for _, v := range n.cases {
		v.case1.check(c)
	}
	c.checkStatement(n.elseStmt)
}

// CaseStatement = "case" Expression "of" Case { ";" Case } [ ";" ] "end" .
func (p *pasParser) caseStatement() (r *caseStatement) {
	r = &caseStatement{
		case1: p.must(tokCase),
		expr:  p.expression(),
		of:    p.must(tokOf),
	}
	var semi pasToken
	r.cases, semi = p.caseList()
	if p.c() == tokElse {
		semi = pasToken{}
		r.else1 = p.shift()
		r.elseStmt = p.statement(true)
	}
	if semi.isValid() {
		r.semi = semi
	} else {
		r.semi = p.opt(';')
	}
	r.end = p.must(tokEnd)
	return r
}

type caseList struct {
	semi  pasToken
	case1 *case1
}

func (p *pasParser) caseList() (r []*caseList, semi pasToken) {
	r = []*caseList{{case1: p.case1()}}
	for p.c() == ';' {
		semi2 := p.shift()
		switch p.c() {
		case tokEnd, tokElse:
			return r, semi2
		default:
			r = append(r, &caseList{semi2, p.case1()})
		}
	}
	return r, semi
}

type case1 struct {
	constList []*constList
	comma     pasToken
	stmt      node
}

// Case = Constant { "," Constant } ":" Statement .
func (p *pasParser) case1() (r *case1) {
	return &case1{
		p.constList(),
		p.must(':'),
		p.statement(true),
	}
}

func (n *case1) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.constList {
		t, val := c.checkExpr(v.const1)
		switch x := t.(type) {
		case *pascalInteger, *subrangeType:
			// ok
		default:
			panic(todo("%v: %T", v.const1.Position(), x))
		}
		switch x := val.(type) {
		case int64:
			// ok
		default:
			panic(todo("%v: %T", v.const1.Position(), x))
		}
	}
	c.checkStatement(n.stmt)
}

type constList struct {
	comma  pasToken
	const1 node

	typ, val interface{}
}

func (n *constList) Position() (r token.Position) {
	if n == nil {
		return r
	}

	if n.comma.isValid() {
		return n.comma.Position()
	}

	return n.const1.Position()
}

func (n *constList) check(c *ctx) {
	n.typ, n.val = c.checkExpr(n.const1)
}

func (p *pasParser) constList() (r []*constList) {
	r = []*constList{{const1: p.expression()}}
	for p.c() == ',' {
		r = append(r, &constList{comma: p.shift(), const1: p.expression()})
	}
	return r
}

type ifStatement struct {
	if1    pasToken
	expr   node
	then   pasToken
	ifStmt node
}

func (n *ifStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.if1.Position()
}

func (n *ifStatement) check(c *ctx) {
	if n == nil {
		return
	}

	t, _ := c.checkExpr2(n.expr)
	switch x := t.(type) {
	case *pascalBoolean:
		// ok
	default:
		panic(todo("%v: %T", n.expr.Position(), x))
	}
	c.checkStatement(n.ifStmt)
}

type ifElseStatement struct {
	if1      pasToken
	expr     node
	then     pasToken
	ifStmt   node
	else1    pasToken
	elseStmt node
}

func (n *ifElseStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.if1.Position()
}

func (n *ifElseStatement) check(c *ctx) {
	if n == nil {
		return
	}

	t, _ := c.checkExpr2(n.expr)
	switch x := t.(type) {
	case *pascalBoolean:
		// ok
	default:
		panic(todo("%v: %T", n.expr.Position(), x))
	}
	c.checkStatement(n.ifStmt)
	c.checkStatement(n.elseStmt)
}

// IfStatement = "if" BooleanExpression "then" Statement
//
//	[ "else" Statement ] .
func (p *pasParser) ifStatement() (r node) {
	x := &ifStatement{
		p.must(tokIf),
		p.expression(),
		p.must(tokThen),
		p.statement(true),
	}
	if p.c() != tokElse {
		return x
	}

	return &ifElseStatement{
		x.if1,
		x.expr,
		x.then,
		x.ifStmt,
		p.shift(),
		p.statement(true),
	}
}

// RepetitiveStatement = WhileStatement | RepeatStatement | ForStatement .
func (p *pasParser) repetitiveStatement() (r node) {
	switch p.c() {
	case tokFor:
		return p.forStatement()
	case tokWhile:
		return p.whileStatement()
	case tokRepeat:
		return p.repeatStatement()
	default:
		panic(todo("", p.token()))
	}
}

type repeatStatement struct {
	repeat pasToken
	stmt   []*statementSequence
	until  pasToken
	expr   node
}

func (n *repeatStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.repeat.Position()
}

func (n *repeatStatement) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.stmt {
		c.checkStatement(v.statement)
	}

	t, _ := c.checkExpr(n.expr)
	switch x := t.(type) {
	case *pascalBoolean:
		// ok
	default:
		panic(todo("%v: %T", n.until.Position(), x))
	}
}

// RepeatStatement = "repeat" StatementSequence "until" Expression .
func (p *pasParser) repeatStatement() (r *repeatStatement) {
	return &repeatStatement{
		p.must(tokRepeat),
		p.statementSequence(),
		p.must(tokUntil),
		p.expression(),
	}
}

type whileStatement struct {
	while pasToken
	expr  node
	do    pasToken
	stmt  node
}

func (n *whileStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.while.Position()
}

func (n *whileStatement) check(c *ctx) {
	if n == nil {
		return
	}

	t, _ := c.checkExpr(n.expr)
	switch x := t.(type) {
	case *pascalBoolean:
		// ok
	default:
		panic(todo("%v: %T", n.expr.Position(), x))
	}
	c.checkStatement(n.stmt)
}

// WhileStatement = "while" BooleanExpression "do" Statement .
func (p *pasParser) whileStatement() (r *whileStatement) {
	return &whileStatement{
		p.must(tokWhile),
		p.expression(),
		p.must(tokDo),
		p.statement(true),
	}
}

type forStatement struct {
	for1         pasToken
	variable     *identifier
	assing       pasToken
	initialValue node
	direction    pasToken
	finalValue   node
	do           pasToken
	statement    node
}

func (n *forStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.for1.Position()
}

func (n *forStatement) check(c *ctx) {
	if n == nil {
		return
	}

	vt, _ := c.checkExpr(n.variable)
	switch x := underlyingType(vt).(type) {
	case *pascalInteger, *subrangeType:
		// ok
	default:
		panic(todo("%v: %T", n.variable.Position(), x))
	}
	ivt, _ := c.checkExpr(n.initialValue)
	switch x := underlyingType(ivt).(type) {
	case *pascalInteger, *subrangeType:
		// ok
	default:
		panic(todo("%v: %T", n.variable.Position(), x))
	}
	fvt, _ := c.checkExpr(n.finalValue)
	switch x := underlyingType(fvt).(type) {
	case *pascalInteger, *subrangeType:
		// ok
	default:
		panic(todo("%v: %T", n.variable.Position(), x))
	}
	c.checkStatement(n.statement)
}

// ForStatement = "for" ControlVariable ":=" InitialValue ( "to" | "downto" ) FinalValue "do" Statement .
func (p *pasParser) forStatement() (r *forStatement) {
	return &forStatement{
		p.must(tokFor),
		p.mustIdent(tokIdent),
		p.must(tokAssign),
		p.expression(),
		p.forDir(),
		p.expression(),
		p.must(tokDo),
		p.statement(true),
	}
}

func (p *pasParser) forDir() (r pasToken) {
	switch p.c() {
	case tokTo, tokDownto:
		return p.shift()
	default:
		panic(todo("", p.token()))
	}
}

type emptyStatement struct {
	t pasToken
}

func (n *emptyStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.t.Position()
}

// SimpleStatement = EmptyStatement | AssignmentStatement | ProcedureStatement | GotoStatement .
func (p *pasParser) simpleStatement() (r node) {
	switch p.c() {
	case tokIdent:
		id := &identifier{ident: p.shift()}
		switch p.c() {
		case '(':
			return p.procedureStatement(id)
		case ';', tokElse, tokUntil, tokEnd:
			switch p.scope.lookup(id).(type) {
			case *procedureDeclaration:
				return p.procedureStatement(id)
			case nil:
				p.errs.err(id.Position(), "undefined: %s", id.Src())
				return id
			default:
				return id
			}
		case tokAssign, '[', '.', '^':
			return p.assignmentStatement(id)
		case ':':
			return &labeled{
				label: id,
				colon: p.shift(),
				stmt:  p.statement(false),
			}
		case '+':
			plus := p.shift()
			expr := p.expression()
			return &labeled{
				label: id,
				plus:  plus,
				off:   expr,
				colon: p.must(':'),
				stmt:  p.statement(false),
			}
		default:
			panic(todo("", p.token()))
		}
	case tokEnd:
		return &emptyStatement{p.token()}
	case tokGoto:
		return p.gotoStatement()
	default:
		panic(todo("", p.token()))
	}
}

type gotoStatement struct {
	goto1 pasToken
	label node
}

func (n *gotoStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.goto1.Position()
}

func (n *gotoStatement) check(c *ctx) {
	if n == nil {
		return
	}

	switch x := n.label.(type) {
	case pasToken, *binaryExpression, *identifier:
		// ok
	default:
		panic(todo("%v: %T", n.label.Position(), x))
	}
}

// GotoStatement = "goto" Label .
func (p *pasParser) gotoStatement() (r *gotoStatement) {
	return &gotoStatement{
		p.must(tokGoto),
		p.expression(),
	}
}

type assignmentStatement struct {
	lhs    node
	assign pasToken
	expr   node
}

func (n *assignmentStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.lhs.Position()
}

func (n *assignmentStatement) check(c *ctx) {
	if n == nil {
		return
	}

	lt, _ := c.checkExpr(n.lhs)
	rt, v := c.checkExpr(n.expr)
	c.mustAssign(n.assign, -1, lt, rt, v)
}

// AssignmentStatement = ( Variable | FunctionIdentifier ) ":=" Expression .
func (p *pasParser) assignmentStatement(lhs node) (r *assignmentStatement) {
	switch x := lhs.(type) {
	case *identifier:
		lhs = p.variable(x)
	default:
		panic(todo("%T %v:", x, lhs.Position()))
	}

	return &assignmentStatement{
		lhs,
		p.must(tokAssign),
		p.expression(),
	}
}

// Variable = EntireVariable | ComponentVariable |
//
//	IdentifiedVariable | BufferVariable .
func (p *pasParser) variable(n node) (r node) {
	for {
		switch p.c() {
		case
			tokAssign, ';', ')', ':', ']', tokNeq, tokGeq, tokMod, tokDiv,
			'*', tokDo, '>', tokLeq, '=', '<', tokOf, ',', tokThen, '+',
			'-', tokElse, tokTo, tokDownto, '/', tokAnd, tokOr:

			return n
		case '[', '.':
			n = p.componentVariable(n)
		case '^':
			n = &deref{n, p.shift()}
		default:
			panic(todo("", p.token()))
		}
	}
}

type deref struct {
	n      node
	carret pasToken
}

func (n *deref) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.n.Position()
}

// ComponentVariable = IndexedVariable | FieldDesignator .
func (p *pasParser) componentVariable(n node) (r node) {
	switch p.c() {
	case '[':
		return p.indexedVariable(n)
	case '.':
		return &fieldDesignator{
			variable: n,
			dot:      p.shift(),
			ident:    p.mustIdent(tokIdent),
		}
	default:
		panic(todo("", p.token()))
	}
}

type fieldDesignator struct {
	variable node
	dot      pasToken
	ident    *identifier

	typ interface{}
}

func (n *fieldDesignator) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.dot.Position()
}

type indexedVariable struct {
	variable  node
	lbracket  pasToken
	indexList []*indexList
	rbracket  pasToken

	typ, varTyp interface{}
}

func (n *indexedVariable) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.variable.Position()
}

// IndexedVariable = ArrayVariahle "[" Index { "," Index ] "]" .
func (p *pasParser) indexedVariable(n node) (r node) {
	return &indexedVariable{
		variable:  n,
		lbracket:  p.must('['),
		indexList: p.indexList(),
		rbracket:  p.must(']'),
	}
}

type indexList struct {
	comma pasToken
	index node
}

func (p *pasParser) indexList() (r []*indexList) {
	r = []*indexList{{index: p.expression()}}
	for p.c() == ',' {
		r = append(r, &indexList{p.shift(), p.expression()})
	}
	return r
}

type procedureStatement struct {
	ident      *identifier
	parameters *parameters

	typ interface{}
}

func (n *procedureStatement) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.ident.Position()
}

func (n *procedureStatement) check(c *ctx) {
	if n == nil {
		return
	}

	switch x := c.scope.mustLookup(n.ident).(type) {
	case
		*pascalWrite, *pascalWriteln, *pascalReset, *pascalRewrite, *pascalRead, *pascalReadln, *procedureDeclaration,
		*knuthBreak, *knuthClose, *pascalGet, *knuthBreakIn, *pascalPut, *knuthSetPos, *knuthPanic:

		n.typ = x
		n.parameters.check(c, x)
	default:
		panic(todo("%v %T", n.ident, x))
	}
}

// ProcedureStatement = ProcedureIdentifier [ ActualParameterList | WriteParameterList ] .
func (p *pasParser) procedureStatement(id *identifier) (r *procedureStatement) {
	return &procedureStatement{
		ident:      id,
		parameters: p.parameters(),
	}
}

type parameters struct {
	lparen     pasToken
	parameters []*parameter
	rparen     pasToken
}

func (n *parameters) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.lparen.Position()
}

func (n *parameters) check(c *ctx, t interface{}) {
	if n == nil {
		return
	}

	switch x := t.(type) {
	case *pascalAbs:
		if len(n.parameters) != 1 {
			panic(todo("%v:", n.lparen.Position()))
		}

		p := n.parameters[0]
		p.typ, p.val = c.checkExpr(p.parameter)
		switch y := underlyingType(p.typ).(type) {
		case *pascalInteger, *subrangeType, *pascalReal:
			// ok
		default:
			panic(todo("%v: %T", n.parameters[0].parameter.Position(), y))
		}
	case *pascalChr:
		if len(n.parameters) != 1 {
			panic(todo("%v:", n.lparen.Position()))
		}

		p := n.parameters[0]
		p.typ, p.val = c.checkExpr2(p.parameter)
		p.typ = underlyingType(p.typ)
		switch y := p.typ.(type) {
		case *pascalInteger, *subrangeType, *pascalReal:
			// ok
		default:
			panic(todo("%v: %T", n.parameters[0].parameter.Position(), y))
		}
	case *pascalRound, *pascalTrunc:
		if len(n.parameters) != 1 {
			panic(todo("%v:", n.lparen.Position()))
		}

		p := n.parameters[0]
		p.typ, p.val = c.checkExpr(p.parameter)
		switch y := p.typ.(type) {
		case *pascalReal, *pascalReal32:
			// ok
		default:
			panic(todo("%v: %T", n.parameters[0].parameter.Position(), y))
		}
	case *pascalWrite, *pascalWriteln, *pascalReset, *pascalRewrite, *pascalRead, *pascalReadln:
		for i, v := range n.parameters {
			v.check(c, i, t)
		}
	case *pascalEOF, *pascalEOLN, *pascalOdd, *knuthBreak, *knuthErstat, *knuthClose, *pascalGet, *pascalPut, *knuthCurPos, *pascalOrd:
		if len(n.parameters) != 1 {
			panic(todo("%v:", n.lparen.Position()))
		}

		n.parameters[0].check(c, 0, t)
	case *knuthBreakIn, *knuthSetPos:
		if len(n.parameters) != 2 {
			panic(todo("%v:", n.lparen.Position()))
		}

		for i, v := range n.parameters {
			v.check(c, i, t)
		}
	case *knuthPanic:
		if len(n.parameters) != 1 {
			panic(todo("%v:", n.lparen.Position()))
		}

		for i, v := range n.parameters {
			v.check(c, i, t)
		}
	case *procedureDeclaration:
		fp := x.procedureHeading.fp
		for i, v := range n.parameters {
			if i >= len(fp) {
				panic(todo(""))
			}

			v.check(c, i, fp[i])
		}
	case *functionDeclaration:
		fp := x.functionHeading.fp
		for i, v := range n.parameters {
			if i >= len(fp) {
				panic(todo(""))
			}

			v.check(c, i, fp[i])
		}
	default:
		panic(todo("%v: %T", n.lparen.Position(), x))
	}
}

// ActualParameterList = "(" ActualParameter { "," ActualParameter } ")" .
// ActualParameter = Expression | Variable | Procedureldentifier | FunctionIdentifier .
// WriteParameterList = "(" ( FileVariable | WriteParameter) { "," WriteParameter } ")" .
// WriteParameter = Expression [ ":" IntegerExpression [ ":" IntegerExpression ] ] .
func (p *pasParser) parameters() (r *parameters) {
	if p.c() != '(' {
		return nil
	}

	return &parameters{
		p.shift(),
		p.parameterList(),
		p.must(')'),
	}
}

type parameter struct {
	comma     pasToken
	parameter node

	typ, val interface{}
}

func (n *parameter) check(c *ctx, ix int, t interface{}) {
	if n == nil {
		return
	}

	n.typ, n.val = c.checkExpr2(n.parameter)
	c.mustAssign(n.parameter, ix, t, n.typ, n.val)
}

func (c *ctx) mustAssign(n node, ix int, dt, st, sv interface{}) {
	if !c.checkAssign(n, ix, dt, st, sv) {
		panic(todo("%v: %T <- %T", n.Position(), dt, st))
	}
}

func (c *ctx) checkAssign(n node, ix int, dt, st, sv interface{}) bool {
	dt = underlyingType(dt)
	st = underlyingType(st)
	if f, ok := st.(*functionDeclaration); ok {
		st = underlyingType(f.functionHeading.result)
	}
	switch x := dt.(type) {
	case *pascalInteger, *subrangeType:
		switch y := st.(type) {
		case *pascalInteger, *subrangeType:
			return true
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := st.(type) {
		case *pascalReal, *pascalReal32, *pascalInteger, *subrangeType:
			return true
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal32:
		switch y := st.(type) {
		case *pascalReal32, *pascalReal:
			return true
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalWrite, *pascalWriteln:
		switch y := st.(type) {
		case *stringLiteral, *pascalInteger, *subrangeType, *pascalChar, *pascalReal, *recordType:
			return true
		case *pascalText, *fileType, *pascalOutput, *pascalStderr:
			return ix == 0
		case *arrayType:
			switch z := y.elemTyp.(type) {
			case *pascalChar:
				return true
			default:
				panic(todo("%v: %T", n.Position(), z))
			}
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReset, *pascalRewrite:
		switch y := st.(type) {
		case *fileType, *pascalText:
			return ix == 0
		case *arrayType:
			switch z := y.elemTyp.(type) {
			case *pascalChar:
				return ix == 1
			default:
				panic(todo("%v: %T", n.Position(), z))
			}
		case *stringLiteral:
			return ix != 0
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *knuthBreakIn:
		switch y := st.(type) {
		case *fileType:
			return ix == 0
		case *pascalBoolean:
			return ix == 1
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *knuthPanic:
		switch y := st.(type) {
		case *pascalInteger:
			return ix == 0
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *knuthSetPos:
		switch y := st.(type) {
		case *fileType:
			return ix == 0
		case *pascalInteger:
			return ix == 1
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalEOF, *pascalEOLN, *knuthBreak, *knuthErstat, *knuthClose, *pascalGet, *pascalPut, *knuthCurPos:
		switch y := st.(type) {
		case *fileType, *pascalText, *pascalOutput, *pascalInput, *pascalStderr:
			return ix == 0
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalOdd:
		switch y := st.(type) {
		case *pascalInteger, *subrangeType:
			return ix == 0
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalOrd:
		switch y := st.(type) {
		case *pascalChar, *subrangeType, *pascalInteger:
			return ix == 0
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *arrayType:
		switch y := x.elemTyp.(type) {
		case *pascalChar:
			switch z := st.(type) {
			case *stringLiteral:
				if len(x.indexTyp) != 1 {
					panic(todo("", n.Position()))
				}

				return true
			default:
				panic(todo("%v: %T", n.Position(), z))
			}
		default:
			switch z := st.(type) {
			case *arrayType:
				return x.isPacked() == z.isPacked() && x.canAssignFrom(c, z)
			default:
				panic(todo("%v: %T", n.Position(), z))
			}
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalChar:
		switch y := st.(type) {
		case *pascalChar:
			return true
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalBoolean:
		switch y := st.(type) {
		case *pascalBoolean:
			return true
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalRead, *pascalReadln:
		switch y := st.(type) {
		case *fileType, *pascalText, *pascalInput:
			return ix == 0
		case *pascalInteger, *subrangeType, *pascalChar, *recordType:
			return true
		default:
			panic(todo("%v: %d %T", n.Position(), ix, y))
		}
	case *functionDeclaration:
		nm := x.functionHeading.ident.Src()
		if c.inFunc[nm] == 0 {
			panic(todo("%v: %T", n.Position(), x))
		}

		return c.checkAssign(n, -1, x.functionHeading.result, st, sv)
	case *recordType:
		switch y := st.(type) {
		case *recordType:
			return x.isPacked() == y.isPacked() && x.canAssignFrom(c, y)
		default:
			panic(todo("%v: %d %T", n.Position(), ix, y))
		}
	case *fileType:
		switch y := st.(type) {
		case *fileType:
			return x.isPacked() == y.isPacked()
		default:
			panic(todo("%v: %d %T", n.Position(), ix, y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (p *pasParser) parameterList() (r []*parameter) {
	if p.c() == ')' {
		return nil
	}

	r = []*parameter{{parameter: p.parameter()}}
	for {
		switch p.c() {
		case ')':
			return r
		case ',':
			r = append(r, &parameter{comma: p.shift(), parameter: p.parameter()})
		default:
			panic(todo("", p.token()))
		}
	}
}

func (p *pasParser) parameter() (r node) {
	r = p.expression()
	switch p.c() {
	case ')', ',':
		return r
	case ':':
		wp := &writeParameter{
			param:  r,
			comma1: p.shift(),
			expr1:  p.expression(),
		}
		if p.c() == ':' {
			wp.comma2 = p.shift()
			wp.expr2 = p.expression()
		}
		return wp
	default:
		panic(todo("", p.token()))
	}
}

type writeParameter struct {
	param  node // param : expr1 [ : expr2 ]
	comma1 pasToken
	expr1  node
	comma2 pasToken
	expr2  node

	wtyp, wval   interface{}
	wtyp2, wval2 interface{}
}

func (n *writeParameter) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.param.Position()
}

type binaryExpression struct {
	lhs node
	op  pasToken
	rhs node

	typ, value interface{}
	cmpTyp     interface{}
}

func (n *binaryExpression) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.lhs.Position()
}

// Expression = SimpleExression [ RelationalOperator SimpleExression ] .
func (p *pasParser) expression() (r node) {
	r = p.simpleExpression(true)
	switch p.c() {
	case '=', tokNeq, '<', tokLeq, '>', tokGeq, tokIn:
		return &binaryExpression{
			lhs: r,
			op:  p.shift(),
			rhs: p.simpleExpression(true),
		}
	default:
		return r
	}
}

// SimpleExpression = [ Sign ] Term { AddingOperator Term } .
func (p *pasParser) simpleExpression(acceptSign bool) (r node) {
	switch p.c() {
	case '-', '+':
		if !acceptSign {
			panic(todo("", p.token()))
		}

		r = &signed{
			sign: p.shift(),
			node: p.term(),
		}
	default:
		r = p.term()
	}
	for {
		switch p.c() {
		case '+', '-', tokOr:
			r = &binaryExpression{
				lhs: r,
				op:  p.shift(),
				rhs: p.term(),
			}
		default:
			return r
		}
	}
}

// Term = Factor { MultiplyingOperator Factor } .
func (p *pasParser) term() (r node) {
	r = p.factor()
	for {
		switch p.c() {
		case '*', '/', tokDiv, tokMod, tokAnd:
			r = &binaryExpression{
				lhs: r,
				op:  p.shift(),
				rhs: p.factor(),
			}
		default:
			return r
		}
	}
}

// Factor = UnsignedConstant | BoundIdentifier | Variable
//
//	| SetConstructor | FunctionDesignator |
//	"not" factor | "(" Expression ")" .
func (p *pasParser) factor() (r node) {
	switch p.c() {
	case '-':
		sgn := p.shift()
		switch p.c() {
		case tokInt, tokFloat, tokString:
			return &signed{
				sign: sgn,
				node: p.shift(),
			}
		default:
			panic(todo("", p.token()))
		}
	case tokString, tokInt, tokFloat, tokNil:
		return p.shift()
	case tokIdent:
		id := &identifier{ident: p.shift()}
		switch p.c() {
		case '(':
			return p.functionCall(id)
		case '[', '.', '^':
			return p.variable(id)
		default:
			if _, ok := p.scope.lookup(id).(*functionDeclaration); ok {
				return &functionCall{ident: id}
			}

			return id
		}
	case '(':
		return p.parenthesizedExpression()
	case tokNot:
		return &not{
			not:    p.shift(),
			factor: p.factor(),
		}
	default:
		panic(todo("", p.token()))
	}
}

type not struct {
	not    pasToken
	factor node

	typ, val interface{}
}

func (n *not) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.not.Position()
}

type parenthesizedExpression struct {
	lparen pasToken
	expr   node
	rparen pasToken

	typ, val interface{}
}

func (n *parenthesizedExpression) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.lparen.Position()
}

func (p *pasParser) parenthesizedExpression() (r *parenthesizedExpression) {
	return &parenthesizedExpression{
		lparen: p.must('('),
		expr:   p.expression(),
		rparen: p.must(')'),
	}
}

type functionCall struct {
	ident      *identifier
	parameters *parameters

	typ, ft interface{}
}

func (n *functionCall) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.ident.Position()
}

func (p *pasParser) functionCall(id *identifier) (r *functionCall) {
	return &functionCall{
		ident:      id,
		parameters: p.parameters(),
	}
}

type procedureAndFunctionDeclarationPart struct {
	declaration node
	semi        pasToken
}

func (n *procedureAndFunctionDeclarationPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.declaration.Position()
}

// ProcedureAndFunctionDeclarationPart = { ( ProcedureDeclaration | FunctionDeclaration ) ";" } .
func (p *pasParser) procedureAndFunctionDeclarationPart() (r []*procedureAndFunctionDeclarationPart) {
	for {
		switch p.c() {
		case tokNoreturn:
			p.shift()
			fallthrough
		case tokProcedure:
			r = append(r, &procedureAndFunctionDeclarationPart{p.procedureDeclaration(), p.must(';')})
		case tokBegin, tokEOF:
			return r
		case tokFunction:
			r = append(r, &procedureAndFunctionDeclarationPart{p.functionDeclaration(), p.must(';')})
		default:
			panic(todo("", p.token()))
		}
	}
}

// FunctionDeclaration = FunctionHeading ";" Block |
//
//	FunctionHeading ";" Directive |
//	FunctionIdentification ";" Block .
func (p *pasParser) functionDeclaration() (r *functionDeclaration) {
	p.scope = newScope(p.scope)

	defer func() { p.scope = p.scope.parent }()

	r = &functionDeclaration{}
	r.functionHeading = p.functionHeading(r)
	r.semi = p.must(';')
	r.block = p.blockOrForward()
	return r
}

func (p *pasParser) blockOrForward() (r node) {
	if p.c() == tokIdent {
		id := p.shift()
		if strings.ToLower(id.Src()) == "forward" {
			return id
		}

		panic(todo("id.Src=%q", id.Src()))
	}

	return p.block()
}

type functionHeading struct {
	function            pasToken
	ident               *identifier
	formalParameterList *formalParameterList
	comma               pasToken
	type1               node

	fp     []interface{}
	result interface{}
}

func (n *functionHeading) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.function.Position()
}

func (n *functionHeading) check(c *ctx) {
	if n == nil {
		return
	}

	n.fp = n.formalParameterList.check(c)
	n.result = c.checkType(n.type1)
}

// FunctionHeading = "function" Identifier [ FormalParameterList ] ":" Type .
func (p *pasParser) functionHeading(n *functionDeclaration) (r *functionHeading) {
	r = &functionHeading{
		function: p.must(tokFunction),
		ident:    p.mustIdent(tokIdent),
	}
	if nm := r.ident.Src(); p.scope.parent.nodes[nm] == nil {
		if err := p.scope.parent.add(nm, n); err != nil {
			p.errs.err(r.ident.Position(), "%s", err)
		}
	}
	switch p.c() {
	case '(':
		r.formalParameterList = p.formalParameterList()
		r.comma = p.must(':')
		r.type1 = p.type1()
		return r
	case ';':
		return r
	case ':':
		r.comma = p.shift()
		r.type1 = p.type1()
		return r
	default:
		panic(todo("", p.token()))
	}
}

type functionDeclaration struct {
	functionHeading *functionHeading
	semi            pasToken
	block           node

	paramTypes []interface{}
	scope      *scope
}

func (n *functionDeclaration) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.functionHeading.Position()
}

func (n *functionDeclaration) isFwd() bool {
	if n == nil {
		return false
	}

	_, ok := n.block.(pasToken)
	return ok
}

func (n *functionDeclaration) check(c *ctx) {
	if n == nil {
		return
	}

	nm := n.functionHeading.ident.Src()
	switch x := n.block.(type) {
	case *block:
		if x, ok := c.scope.nodes[nm].(*functionDeclaration); ok && x.isFwd() {
			delete(c.scope.nodes, nm)
			n.functionHeading = x.functionHeading
		}
		if err := c.scope.add(nm, n); err != nil {
			c.errs.err(n.functionHeading.ident.Position(), "%s", err)
		}

		defer func() { c.scope = c.scope.parent }()

		scope := n.scope
		if scope == nil {
			c.scope = newScope(c.scope)
			n.scope = c.scope
			n.functionHeading.check(c)
		}
		c.inFunc[nm]++

		defer func() { c.inFunc[nm]-- }()

		x.check(c)
	case pasToken: // forward
		if err := c.scope.add(nm, n); err != nil {
			c.errs.err(n.functionHeading.ident.Position(), "%s", err)
		}
		c.scope = newScope(c.scope)
		n.scope = c.scope

		defer func() { c.scope = c.scope.parent }()

		n.functionHeading.check(c)
	default:
		panic(todo("%v %T", n.semi, x))
	}
}

type procedureDeclaration struct {
	procedureHeading *procedureHeading
	semi             pasToken
	block            node

	scope *scope
}

func (n *procedureDeclaration) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.procedureHeading.Position()
}

func (n *procedureDeclaration) check(c *ctx) {
	if n == nil {
		return
	}

	nm := n.procedureHeading.ident.Src()
	switch x := n.block.(type) {
	case *block:
		if x, ok := c.scope.nodes[nm].(*procedureDeclaration); ok && x.isFwd() {
			delete(c.scope.nodes, nm)
			n.procedureHeading = x.procedureHeading
		}
		if err := c.scope.add(nm, n); err != nil {
			c.errs.err(n.procedureHeading.ident.Position(), "%s", err)
		}

		defer func() { c.scope = c.scope.parent }()

		scope := n.scope
		if scope == nil {
			c.scope = newScope(c.scope)
			n.scope = c.scope
			n.procedureHeading.check(c)
		}
		x.check(c)
	case pasToken: // forward
		if err := c.scope.add(nm, n); err != nil {
			c.errs.err(n.procedureHeading.ident.Position(), "%s", err)
		}
		c.scope = newScope(c.scope)
		n.scope = c.scope

		defer func() { c.scope = c.scope.parent }()

		n.procedureHeading.check(c)
	default:
		panic(todo("%v %T", n.semi, x))
	}
}

func (n *procedureDeclaration) isFwd() bool {
	if n == nil {
		return false
	}

	_, ok := n.block.(pasToken)
	return ok
}

// ProcedureDeclaration = ProcedureHeading ";" Block |
//
//	ProcedureHeading ";" Directive |
//	ProcedureIdentification ";" Block.
func (p *pasParser) procedureDeclaration() (r *procedureDeclaration) {
	p.scope = newScope(p.scope)

	defer func() { p.scope = p.scope.parent }()

	r = &procedureDeclaration{}
	r.procedureHeading = p.procedureHeading(r)
	r.semi = p.must(';')
	r.block = p.blockOrForward()
	return r
}

type formalParameterList struct {
	lparen pasToken
	params []*formalParameterSection
	rparen pasToken
}

func (n *formalParameterList) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.lparen.Position()
}

func (n *formalParameterList) check(c *ctx) (r []interface{}) {
	if n == nil {
		return
	}

	for _, v := range n.params {
		r = append(r, v.check(c)...)
	}
	return r
}

// FormalParameterList = "(" FormalParameterSection { ";" FormalParameterSection } ")" .
func (p *pasParser) formalParameterList() (r *formalParameterList) {
	if p.c() != '(' {
		return nil
	}

	return &formalParameterList{
		p.must('('),
		p.formalParameterSection(),
		p.must(')'),
	}
}

type formalParameterSection struct {
	semi  pasToken
	param node
}

func (n *formalParameterSection) check(c *ctx) (r []interface{}) {
	if n == nil {
		return
	}

	switch x := n.param.(type) {
	case *parameterSpecification:
		r = append(r, x.check(c)...)
	default:
		panic(todo("%v: %T", n.param.Position(), x))
	}
	return r
}

// FormalParameterSection = ValueParameterSpecification |
//
//	VariableParamererSpecification |
//	ProceduralParameterSpecification |
//	FunctionalParameterSpecification .
func (p *pasParser) formalParameterSection() (r []*formalParameterSection) {
	switch p.c() {
	case tokIdent:
		r = []*formalParameterSection{{param: p.parameterSpecification(pasToken{})}}
	case tokVar:
		r = []*formalParameterSection{{param: p.parameterSpecification(p.shift())}}
	default:
		panic(todo("", p.token()))
	}
	for p.c() == ';' {
		semi := p.shift()
		switch p.c() {
		case tokIdent:
			r = append(r, &formalParameterSection{semi, p.parameterSpecification(pasToken{})})
		case tokVar:
			r = append(r, &formalParameterSection{semi, p.parameterSpecification(p.shift())})
		default:
			panic(todo("", p.token()))
		}
	}
	return r
}

type parameterSpecification struct {
	var1           pasToken
	identifierList []*identifierList
	comma          pasToken
	type1          node

	typ interface{}
}

func (n *parameterSpecification) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.identifierList[0].ident.Position()
}

func (n *parameterSpecification) check(c *ctx) (r []interface{}) {
	if n == nil {
		return
	}

	n.typ = c.checkType(n.type1)
	for _, v := range n.identifierList {
		r = append(r, n.typ)
		p := &variable{
			ident:   v.ident,
			typ:     n.typ,
			isParam: true,
		}
		if err := c.scope.add(v.ident.Src(), p); err != nil {
			c.errs.err(v.ident.Position(), "%s", err)
		}
	}
	return r
}

// ValueParameterSpecification = IdentifierList ":" Type .
func (p *pasParser) parameterSpecification(var1 pasToken) (r *parameterSpecification) {
	r = &parameterSpecification{
		var1:           var1,
		identifierList: p.identifierList(),
		comma:          p.must(':'),
		type1:          p.type1(),
	}
	for _, v := range r.identifierList {
		if err := p.scope.add(v.ident.Src(), r); err != nil {
			p.errs.err(v.ident.Position(), "%s", err)
		}
	}
	return r
}

type procedureHeading struct {
	procedure           pasToken
	ident               *identifier
	formalParameterList *formalParameterList

	fp []interface{}
}

func (n *procedureHeading) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.procedure.Position()
}

func (n *procedureHeading) check(c *ctx) {
	if n == nil {
		return
	}

	n.fp = n.formalParameterList.check(c)
}

// ProcedureHeading = "procedure" Identifier [ FormalParameterList ] .
func (p *pasParser) procedureHeading(n *procedureDeclaration) (r *procedureHeading) {
	r = &procedureHeading{
		procedure:           p.must(tokProcedure),
		ident:               p.mustIdent(tokIdent),
		formalParameterList: p.formalParameterList(),
	}
	if nm := r.ident.Src(); p.scope.parent.nodes[nm] == nil {
		if err := p.scope.parent.add(nm, n); err != nil {
			p.errs.err(r.ident.Position(), "%s", err)
		}
	}
	return r
}

type variableDeclarationPart struct {
	var1                    pasToken
	variableDeclarationList []*variableDeclarationList
}

func (n *variableDeclarationPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.var1.Position()
}

// VariableDeclarationPart = [ "var" VariableDeclaration ";" { VariableDeclaration ";" } ] .
func (p *pasParser) variableDeclarationPart() (r *variableDeclarationPart) {
	if p.c() != tokVar {
		return nil
	}

	r = &variableDeclarationPart{
		p.shift(),
		[]*variableDeclarationList{{p.variableDeclaration(), p.must(';')}},
	}
	for {
		switch p.c() {
		case tokProcedure, tokFunction, tokBegin:
			return r
		case tokIdent:
			r.variableDeclarationList = append(r.variableDeclarationList, &variableDeclarationList{p.variableDeclaration(), p.must(';')})
		default:
			panic(todo("", p.token()))
		}
	}
}

type variableDeclaration struct {
	identifierList []*identifierList
	colon          pasToken
	type1          node

	typ interface{}
}

func (n *variableDeclaration) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.identifierList[0].ident.Position()
}

// VariableDeclaration = IdentifierList ":" Type .
func (p *pasParser) variableDeclaration() (r *variableDeclaration) {
	r = &variableDeclaration{
		identifierList: p.identifierList(),
		colon:          p.must(':'),
		type1:          p.type1(),
	}
	for _, v := range r.identifierList {
		if err := p.scope.add(v.ident.Src(), r); err != nil {
			p.errs.err(v.ident.Position(), "%s", err)
		}
	}
	return r
}

type variableDeclarationList struct {
	variableDeclaration *variableDeclaration
	semi                pasToken
}

type constantDefinitionPart struct {
	const1                 pasToken
	constantDefinitionList []*constantDefinitionList
}

func (n *constantDefinitionPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.const1.Position()
}

type constantDefinitionList struct {
	constantDefinition *constantDefinition
	semi               pasToken
}

// ConstantDefinitionPart = [ "const" ConstantDefinition ";" { ConstantDefinition ";" } ] .
func (p *pasParser) constantDefinitionPart() (r *constantDefinitionPart) {
	if p.c() != tokConst {
		return nil
	}

	r = &constantDefinitionPart{
		p.shift(),
		[]*constantDefinitionList{{p.constantDefinition(), p.must(';')}},
	}
	for {
		switch p.c() {
		case tokIdent:
			r.constantDefinitionList = append(r.constantDefinitionList, &constantDefinitionList{p.constantDefinition(), p.must(';')})
		case tokType, tokVar:
			return r
		default:
			panic(todo("", p.token()))
		}
	}
}

type constantDefinition struct {
	ident    *identifier
	eq       pasToken
	constant node

	value, typ interface{}
}

func (n *constantDefinition) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.ident.Position()
}

// ConstantDefinition = Identifier "=" Constant .
func (p *pasParser) constantDefinition() (r *constantDefinition) {
	r = &constantDefinition{
		ident:    p.mustIdent(tokIdent),
		eq:       p.must('='),
		constant: p.expression(),
	}
	if err := p.scope.add(r.ident.Src(), r); err != nil {
		p.errs.err(r.ident.Position(), "%s", err)
	}
	return r
}

type numericConstant struct {
	sign  pasToken
	value node
}

func (n *numericConstant) Position() (r token.Position) {
	if n == nil {
		return r
	}

	if n.sign.isValid() {
		return n.sign.Position()
	}

	return n.value.Position()
}

type signed struct {
	sign pasToken
	node node

	typ interface{}
}

func (n *signed) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.sign.Position()
}

// Constant = [ Sign ] ( UnsignedNumher | ConstantIdentifier) | CharacterString .
func (p *pasParser) constant(acceptSign bool) (r node) {
	switch p.c() {
	case tokInt:
		return &numericConstant{value: p.expression()}
	case tokIdent, tokString:
		return p.shift()
	case '-', '+':
		if !acceptSign {
			panic(todo("", p.token()))
		}

		return &signed{
			sign: p.shift(),
			node: p.constant(false),
		}
	default:
		panic(todo("", p.token()))
	}
}

type labelDeclarationPart struct {
	label             pasToken
	digitSequenceList []*digitSequenceList
	semi              pasToken
}

func (n *labelDeclarationPart) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.label.Position()
}

type digitSequenceList struct {
	comma         pasToken
	digitSequence node
}

// LabelDeclarationPart = [ "label" DigitSequence { "," DigitSequence } ";" ]
func (p *pasParser) labelDeclarationPart() (r *labelDeclarationPart) {
	if p.c() != tokLabel {
		return nil
	}

	r = &labelDeclarationPart{
		label: p.shift(),
	}
	r.digitSequenceList = []*digitSequenceList{{digitSequence: p.expression()}}
	for p.c() == ',' {
		r.digitSequenceList = append(r.digitSequenceList, &digitSequenceList{p.shift(), p.expression()})
	}
	r.semi = p.must(';')
	return r
}

type scope struct {
	nodes  map[string]node
	parent *scope

	setResolution bool
}

func newScope(parent *scope) *scope {
	setResolution := false
	if parent != nil {
		setResolution = parent.setResolution
	}
	return &scope{parent: parent, setResolution: setResolution}
}

func (s *scope) isTLD() bool      { return s.parent != nil && s.parent.isUniverse() }
func (s *scope) isUniverse() bool { return s.parent == nil }

func (s *scope) add(nm string, n node) (err error) {
	nm = strings.ToLower(nm)
	if ex, ok := s.nodes[nm]; ok {
		return fmt.Errorf("%v: %s redefined, previous definition at %v: %T", n.Position(), nm, ex.Position(), ex)
	}

	if s.nodes == nil {
		s.nodes = map[string]node{}
	}
	s.nodes[nm] = n
	return nil
}

func (s *scope) mustLookup(id *identifier) (r node) {
	if r := s.lookup(id); r != nil {
		return r
	}

	// for s != nil {
	// 	trc("tried %p, universe %v, tld %v, setResolution %v", s, s.isUniverse(), s.isTLD(), s.setResolution)
	// 	s = s.parent
	// }
	panic(todo("%v: undefined %q", id.Position(), id.Src()))
}

func (s *scope) lookup(id *identifier) node {
	nm := strings.ToLower(id.Src())
	for s != nil {
		if r := s.nodes[nm]; r != nil {
			if s.setResolution {
				id.scope = s
				id.resolvedTo = r
			}
			return r
		}

		s = s.parent
	}

	return nil
}

type ctx struct {
	errs  errList
	scope *scope

	inFunc map[string]int
}

func newCtx() (r *ctx) {
	universe := newScope(nil)
	universe.setResolution = true
	universe.add("abs", &pascalAbs{})
	universe.add("boolean", &pascalBoolean{})
	universe.add("break", &knuthBreak{})
	universe.add("break_in", &knuthBreakIn{})
	universe.add("char", &pascalChar{})
	universe.add("chr", &pascalChr{})
	universe.add("close", &knuthClose{})
	universe.add("cur_pos", &knuthCurPos{})
	universe.add("eof", &pascalEOF{})
	universe.add("eoln", &pascalEOLN{})
	universe.add("erstat", &knuthErstat{})
	universe.add("false", &pascalFalse{})
	universe.add("get", &pascalGet{})
	universe.add("input", &pascalInput{})
	universe.add("integer", &pascalInteger{})
	universe.add("max_int", &pascalMaxInt{})
	universe.add("odd", &pascalOdd{})
	universe.add("ord", &pascalOrd{})
	universe.add("output", &pascalOutput{})
	universe.add("panic", &knuthPanic{})
	universe.add("put", &pascalPut{})
	universe.add("read", &pascalRead{})
	universe.add("readln", &pascalReadln{})
	universe.add("real", &pascalReal{})
	universe.add("reset", &pascalReset{})
	universe.add("rewrite", &pascalRewrite{})
	universe.add("round", &pascalRound{})
	universe.add("set_pos", &knuthSetPos{})
	universe.add("stderr", &pascalStderr{})
	universe.add("text", &pascalText{})
	universe.add("true", &pascalTrue{})
	universe.add("trunc", &pascalTrunc{})
	universe.add("write", &pascalWrite{})
	universe.add("writeln", &pascalWriteln{})
	return &ctx{
		inFunc: map[string]int{},
		scope:  newScope(universe),
	}
}

type ast struct {
	eof      pasToken
	program  *program
	tldScope *scope
}

func (n *ast) check() error {
	c := newCtx()
	n.tldScope = c.scope
	n.program.check(c)
	return c.errs.Err()
}

func (n *program) check(c *ctx) {
	if n == nil {
		return
	}

	n.programHeading.check(c)
	n.block.check(c)
	var progParams []*identifier
	if n.programHeading.programParameterList != nil {
		progParams = n.programHeading.programParameterList.idents
	}
	for _, v := range progParams {
		if x := c.scope.mustLookup(v); x == nil {
			panic(todo("", v))
		}
	}
}

func (n *programHeading) check(c *ctx) {
	if n == nil {
		return
	}

	n.programParameterList.check(c)
}

func (n *programParameterList) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.identifierList {
		n.idents = append(n.idents, v.ident)
	}
}

func (n *block) check(c *ctx) {
	if n == nil {
		return
	}

	n.labelDeclarationPart.check(c)
	n.constantDefinitionPart.check(c)
	n.typeDefinitionPart.check(c)
	n.variableDeclarationPart.check(c)
	for _, v := range n.procedureAndFunctionDeclarationPart {
		switch x := v.declaration.(type) {
		case *procedureDeclaration:
			x.check(c)
		case *functionDeclaration:
			x.check(c)
		default:
			panic(todo("%v: %T", x.Position(), x))
		}
	}
	n.statementPart.check(c)
}

type label struct {
	node

	nm       string
	lhs, rhs interface{}
}

func (n *labelDeclarationPart) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.digitSequenceList {
		switch x := v.digitSequence.(type) {
		case *identifier:
			y := &label{x, x.Src(), x.Src(), nil}
			if err := c.scope.add(y.nm, y); err != nil {
				c.errs.err(n.Position(), "%s", err)
			}
		case pasToken:
			y := &label{x, x.Src(), x.Src(), nil}
			if err := c.scope.add(y.nm, y); err != nil {
				c.errs.err(y.Position(), "%s", err)
			}
		case *binaryExpression:
			if x.op.ch != '+' {
				panic(todo(""))
			}

			switch y := x.lhs.(type) {
			case *identifier:
				_, rhs := c.checkExpr(x.rhs)
				lbl := &label{x, fmt.Sprintf("%s+%v", y.Src(), rhs), y.Src(), rhs}
				if err := c.scope.add(lbl.nm, lbl); err != nil {
					c.errs.err(lbl.Position(), "%s", err)
				}
			default:
				panic(todo("%T", y))
			}
		default:
			panic(todo("%v: %T", x.Position(), x))
		}
	}
}

func (n *constantDefinitionPart) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.constantDefinitionList {
		v.constantDefinition.check(c)
	}
}

func (n *constantDefinition) check(c *ctx) {
	if n == nil {
		return
	}

	n.typ, n.value = c.checkExpr(n.constant)
	if n.value == nil {
		panic(todo("%v: not a constant expression", n.constant.Position()))
	}

	nm := strings.ToLower(n.ident.Src())
	if _, ok := c.scope.nodes[nm].(*label); ok {
		return
	}

	if err := c.scope.add(nm, n); err != nil {
		c.errs.err(n.Position(), "%s", err)
	}
}

func (c *ctx) checkExpr2(n node) (typ interface{}, val interface{}) {
	typ, val = c.checkExpr(n)
	if x, ok := typ.(*functionDeclaration); ok {
		typ = x.functionHeading.result
	}
	return typ, val
}

func (c *ctx) checkExpr(n node) (typ interface{}, val interface{}) {
	switch x := n.(type) {
	case *identifier:
		defer func() { x.typ = typ }()

		switch y := c.scope.mustLookup(x).(type) {
		case *constantDefinition:
			return y.typ, y.value
		case *variable:
			return y.typ, nil
		case *pascalTrue:
			return &pascalBoolean{}, true
		case *pascalFalse:
			return &pascalBoolean{}, false
		case *pascalEOF, *pascalEOLN, *functionDeclaration, *label, *pascalOutput, *pascalInput, *pascalStderr:
			return y, nil
		case *pascalMaxInt:
			return &pascalInteger{}, int64(math.MaxInt32)
		default:
			panic(todo("%v: %T", x.Position(), y))
		}
	case pasToken:
		switch x.ch {
		case tokInt:
			s := x.Src()
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(todo("%v: %v", x.Position(), err))
			}

			return &pascalInteger{}, n
		case tokFloat:
			s := x.Src()
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(todo("%v: %v", x.Position(), err))
			}

			return &pascalReal{}, n
		case tokString:
			s := x.Src()
			s = s[1 : len(s)-1]                  // remove quotes
			s = strings.ReplaceAll(s, "''", "'") // unquote
			switch a := []rune(s); len(a) {
			case 1:
				return &pascalChar{}, int64(a[0])
			default:
				return &stringLiteral{}, s
			}
		default:
			panic(todo("", x))
		}
	case *binaryExpression:
		t1, v1 := c.checkExpr2(x.lhs)
		t2, v2 := c.checkExpr2(x.rhs)
		switch x.op.ch {
		case '+':
			x.typ, x.value = c.add(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case '-':
			x.typ, x.value = c.sub(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case '*':
			x.typ, x.value = c.mul(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case '/':
			x.typ, x.value = c.div(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case tokDiv:
			x.typ, x.value = c.idiv(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case tokMod:
			x.typ, x.value = c.mod(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case '=':
			x.typ, x.value = c.eq(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case '>':
			x.typ, x.value = c.gt(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case '<':
			x.typ, x.value = c.lt(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case tokNeq:
			x.typ, x.value = c.neq(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case tokGeq:
			x.typ, x.value = c.geq(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case tokLeq:
			x.typ, x.value = c.leq(x.op, t1, v1, t2, v2)
			x.cmpTyp = c.binType(t1, t2)
			return x.typ, x.value
		case tokOr:
			x.typ, x.value = c.or(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		case tokAnd:
			x.typ, x.value = c.and(x.op, t1, v1, t2, v2)
			return x.typ, x.value
		default:
			panic(todo("%v %T %v %T %v", x.op, t1, v1, t2, v2))
		}
	case *not:
		t, v := c.checkExpr2(x.factor)
		switch y := t.(type) {
		case *pascalBoolean:
			if v != nil {
				v = !v.(bool)
			}
			x.typ = y
			x.val = v
			return x.typ, x.val
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *signed:
		t, v := c.checkExpr2(x.node)
		switch y := underlyingType(t).(type) {
		case *pascalInteger, *subrangeType:
			if v != nil {
				v = -v.(int64)
			}
			x.typ = y
			return x.typ, v
		case *pascalReal:
			if v != nil {
				v = -v.(float64)
			}
			x.typ = y
			return x.typ, v
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *indexedVariable:
		x.varTyp, _ = c.checkExpr(x.variable)
		switch y := underlyingType(x.varTyp).(type) {
		case *arrayType:
			switch xDims, yDims := len(x.indexList), len(y.typeList); {
			case xDims == yDims:
				for _, v := range x.indexList {
					c.checkExpr2(v.index)
				}
				x.typ = y.elemTyp
				return x.typ, nil
			default:
				panic(todo("%v: %v %v", n.Position(), xDims, yDims))
			}
		default:
			panic(todo("%v: %T", x.lbracket.Position(), y))
		}
	case *fieldDesignator:
		vt, _ := c.checkExpr(x.variable)
		switch y := underlyingType(vt).(type) {
		case *recordType:
			x.typ = y.field(x.ident.ident).typ
			return x.typ, nil
		default:
			panic(todo("%v: %T", x.dot.Position(), y))
		}
	case *functionCall:
		switch y := c.scope.mustLookup(x.ident).(type) {
		case *pascalAbs:
			x.ft = y
			x.parameters.check(c, y)
			if len(x.parameters.parameters) == 1 {
				x.typ = x.parameters.parameters[0].typ
			}
		case *pascalChr:
			x.typ = &pascalChar{}
			x.ft = y
			x.parameters.check(c, y)
		case *pascalRound, *knuthCurPos, *pascalOrd, *pascalTrunc:
			x.typ = &pascalInteger{}
			x.ft = y
			x.parameters.check(c, y)
		case *pascalEOF, *pascalEOLN, *pascalOdd:
			x.typ = &pascalBoolean{}
			x.ft = y
			x.parameters.check(c, y)
		case *knuthErstat:
			x.typ = &pascalInteger{}
			x.ft = y
			x.parameters.check(c, y)
		case *functionDeclaration:
			x.typ = y.functionHeading.result
			x.ft = y
			x.parameters.check(c, y)
		default:
			panic(todo("%v: %T", x.ident.Position(), y))
		}
		return x.typ, nil
	case *parenthesizedExpression:
		x.typ, x.val = c.checkExpr(x.expr)
		return x.typ, x.val
	case *writeParameter:
		x.wtyp, x.wval = c.checkExpr2(x.expr1)
		x.wval = x.wval.(int64)
		if x.expr2 != nil {
			x.wtyp2, x.wval2 = c.checkExpr2(x.expr2)
			x.wval2 = x.wval2.(int64)
		}
		return c.checkExpr2(x.param)
	case *deref:
		t, _ := c.checkExpr2(x.n)
		switch y := underlyingType(t).(type) {
		case *fileType:
			return y.elemType, nil
		case *pascalInput:
			return &pascalChar{}, nil
		default:
			panic(todo("%v: %T", x.n.Position(), y))
		}
	default:
		panic(todo("%v: %T", x.Position(), x))
	}
}

func (c *ctx) div(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v2 == int64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = float64(v1.(int64)) / float64(v2.(int64))
			}
			return &pascalReal{}, value
		case *pascalReal:
			if v2 == float64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = float64(v1.(int64)) / v2.(float64)
			}
			return &pascalReal{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v2 == float64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = v1.(float64) / v2.(float64)
			}
			return &pascalReal{}, value
		case *pascalInteger, *subrangeType:
			if v2 == int64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = v1.(float64) / float64(v2.(int64))
			}
			return &pascalReal{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) idiv(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v2 == int64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = v1.(int64) / v2.(int64)
			}
			return c.binType(t1, t2), value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) mod(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v2 == int64(0) {
				panic(todo(""))
			}

			if v1 != nil && v2 != nil {
				value = v1.(int64) % v2.(int64)
			}
			return c.binType(t1, t2), value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) mul(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) * v2.(int64)
			}
			return c.binType(t1, t2), value
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = float64(v1.(int64)) * v2.(float64)
			}
			return y, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) * v2.(float64)
			}
			return t1, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) * float64(v2.(int64))
			}
			return t1, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal32:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) * v2.(float64)
			}
			return &pascalReal{}, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) * float64(v2.(int64))
			}
			return &pascalReal{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) add(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) + v2.(int64)
			}
			return c.binType(t1, t2), value
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = float64(v1.(int64)) + v2.(float64)
			}
			return &pascalReal{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) + v2.(float64)
			}
			return t1, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) + float64(v2.(int64))
			}
			return t1, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *label:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			v1 = x.lhs
			if v1 != nil && v2 != nil {
				value = v1.(int64) + v2.(int64)
			}
			return t1, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) sub(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) - v2.(int64)
			}
			return c.binType(t1, t2), value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) - v2.(float64)
			}
			return t1, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) - float64(v2.(int64))
			}
			return t1, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) binType(a, b interface{}) interface{} {
	a = underlyingType(a)
	b = underlyingType(b)
	switch x := a.(type) {
	case *pascalInteger, *subrangeType, *pascalChar:
		switch y := b.(type) {
		case *pascalInteger, *subrangeType, *pascalChar:
			return &pascalInteger{}
		case *pascalReal:
			return y
		default:
			panic(todo("%T", y))
		}
	case *pascalReal, *pascalBoolean:
		return a
	case *pascalReal32:
		return &pascalReal{}
	default:
		panic(todo("%T", x))
	}
}

func (c *ctx) eq(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) == v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) == float64(v2.(int64))
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalChar:
		switch y := t2.(type) {
		case *pascalChar:
			if v1 != nil && v2 != nil {
				value = v1.(int64) == v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalBoolean:
		switch y := t2.(type) {
		case *pascalBoolean:
			if v1 != nil && v2 != nil {
				value = v1.(bool) == v2.(bool)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) lt(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) < v2.(int64)
			}
			return &pascalBoolean{}, value
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = float64(v1.(int64)) < v2.(float64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) < v2.(float64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) gt(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) > v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) > v2.(float64)
			}
			return &pascalBoolean{}, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) > float64(v2.(int64))
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) neq(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) != v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalChar:
		switch y := t2.(type) {
		case *pascalChar:
			if v1 != nil && v2 != nil {
				value = v1.(int64) != v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalBoolean:
		switch y := t2.(type) {
		case *pascalBoolean:
			if v1 != nil && v2 != nil {
				value = v1.(bool) != v2.(bool)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) != v2.(float64)
			}
			return &pascalBoolean{}, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) != float64(v2.(int64))
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) and(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalBoolean:
		switch y := t2.(type) {
		case *pascalBoolean:
			if v1 != nil && v2 != nil {
				value = v1.(bool) && v2.(bool)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) or(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalBoolean:
		switch y := t2.(type) {
		case *pascalBoolean:
			if v1 != nil && v2 != nil {
				value = v1.(bool) || v2.(bool)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) leq(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) <= v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (c *ctx) geq(n node, t1, v1, t2, v2 interface{}) (typ, value interface{}) {
	t1 = underlyingType(t1)
	t2 = underlyingType(t2)
	switch x := t1.(type) {
	case *pascalInteger, *subrangeType:
		switch y := t2.(type) {
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(int64) >= v2.(int64)
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	case *pascalReal:
		switch y := t2.(type) {
		case *pascalReal:
			if v1 != nil && v2 != nil {
				value = v1.(float64) >= v2.(float64)
			}
			return &pascalBoolean{}, value
		case *pascalInteger, *subrangeType:
			if v1 != nil && v2 != nil {
				value = v1.(float64) >= float64(v2.(int64))
			}
			return &pascalBoolean{}, value
		default:
			panic(todo("%v: %T", n.Position(), y))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (n *typeDefinitionPart) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.typeDefinitionList {
		v.check(c)
	}
}

func (n *typeDefinitionList) check(c *ctx) {
	if n == nil {
		return
	}

	n.typeDefinition.check(c)
}

func (n *typeDefinition) check(c *ctx) {
	if n == nil {
		return
	}

	n.typ = c.checkType(n.type1)
	if err := c.scope.add(n.ident.Src(), n); err != nil {
		c.errs.err(n.Position(), "%s", err)
	}
}

func (c *ctx) checkType(t node) (r interface{}) {
	switch x := t.(type) {
	case nil:
		return nil
	case *subrangeType:
		if x.lo != 0 || x.hi != 0 {
			return x
		}

		_, v := c.checkExpr(x.constant)
		switch y := v.(type) {
		case int64:
			x.lo = y
		default:
			panic(todo("%T", y))
		}
		_, v = c.checkExpr(x.constant2)
		switch y := v.(type) {
		case int64:
			x.hi = y
		default:
			panic(todo("%T", y))
		}
		return x
	case *fileType:
		if x.elemType == nil {
			x.elemType = c.checkType(x.type1)
		}
		return x
	case *identifier:
		defer func() { x.typ = r }()

		switch y := c.scope.mustLookup(x).(type) {
		case *typeDefinition:
			return y
		case *pascalText:
			return y
		case *pascalInteger:
			return y
		case *pascalBoolean:
			return y
		case *pascalReal:
			return y
		case *pascalChar:
			return y
		default:
			panic(todo("%v: %v %T", x.Position(), x, y))
		}
	case *arrayType:
		if x.elemTyp == nil {
			x.indexTyp = nil
			for _, v := range x.typeList {
				x.indexTyp = append(x.indexTyp, underlyingType(c.checkType(v.type1)))
			}
			x.elemTyp = c.checkType(x.elemType)
		}
		return x
	case *recordType:
		x.check(c)
		return x
	default:
		panic(todo("%v: %T", t.Position(), x))
	}
}

func (n *variableDeclarationPart) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.variableDeclarationList {
		v.check(c)
	}
}

func (n *variableDeclarationList) check(c *ctx) {
	if n == nil {
		return
	}

	n.variableDeclaration.check(c)
}

type variable struct {
	ident *identifier
	typ   interface{}

	isParam bool
}

func (n *variable) Position() (r token.Position) {
	if n == nil {
		return r
	}

	return n.ident.Position()
}

func (n *variableDeclaration) check(c *ctx) {
	if n == nil {
		return
	}

	n.typ = c.checkType(n.type1)
	for _, v := range n.identifierList {
		if err := c.scope.add(v.ident.Src(), &variable{ident: v.ident, typ: n.typ}); err != nil {
			c.errs.err(n.Position(), "%s", err)
		}
	}
}

func (n *compoundStatement) check(c *ctx) {
	if n == nil {
		return
	}

	for _, v := range n.statementSequence {
		c.checkStatement(v.statement)
	}
}

func (c *ctx) checkStatement(n node) {
	if n == nil {
		return
	}

	switch x := n.(type) {
	case *procedureStatement:
		x.check(c)
	case *assignmentStatement:
		x.check(c)
	case *forStatement:
		x.check(c)
	case *compoundStatement:
		x.check(c)
	case *emptyStatement:
		// nop
	case *ifElseStatement:
		x.check(c)
	case *ifStatement:
		x.check(c)
	case *whileStatement:
		x.check(c)
	case *repeatStatement:
		x.check(c)
	case *identifier:
		switch y := c.scope.mustLookup(x).(type) {
		case *procedureDeclaration:
			x.typ = y
			if len(y.procedureHeading.fp) != 0 {
				panic(todo(""))
			}

			// ok
		default:
			panic(todo("%v: %T", x.Position(), y))
		}
	case *caseStatement:
		x.check(c)
	case *gotoStatement:
		x.check(c)
	case *labeled:
		x.check(c)
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

// ============================================================================

type nopos struct{}

func (nopos) Position() (r token.Position) { return r }

type (
	knuthBreak    struct{ nopos }
	knuthBreakIn  struct{ nopos }
	knuthClose    struct{ nopos }
	knuthCurPos   struct{ nopos }
	knuthErstat   struct{ nopos }
	knuthPanic    struct{ nopos }
	knuthSetPos   struct{ nopos }
	pascalAbs     struct{ nopos }
	pascalBoolean struct{ nopos }
	pascalChar    struct{ nopos }
	pascalChr     struct{ nopos }
	pascalEOF     struct{ nopos }
	pascalEOLN    struct{ nopos }
	pascalFalse   struct{ nopos }
	pascalGet     struct{ nopos }
	pascalInput   struct{ nopos }
	pascalInteger struct{ nopos }
	pascalMaxInt  struct{ nopos }
	pascalOdd     struct{ nopos }
	pascalOrd     struct{ nopos }
	pascalOutput  struct{ nopos }
	pascalPut     struct{ nopos }
	pascalRead    struct{ nopos }
	pascalReadln  struct{ nopos }
	pascalReal    struct{ nopos }
	pascalReal32  struct{ nopos }
	pascalReset   struct{ nopos }
	pascalStderr  struct{ nopos }
	pascalRewrite struct{ nopos }
	pascalRound   struct{ nopos }
	pascalString  struct{ nopos }
	pascalText    struct{ nopos }
	pascalTrue    struct{ nopos }
	pascalTrunc   struct{ nopos }
	pascalWrite   struct{ nopos }
	pascalWriteln struct{ nopos }
	stringLiteral struct{ nopos }
)

func typeStr(t interface{}) string {
	switch x := t.(type) {
	case *subrangeType:
		return fmt.Sprintf("%d..%d", x.lo, x.hi)
	default:
		panic(todo("%T", x))
	}
}
