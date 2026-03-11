// Copyright 2023 The Knuth Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web // modernc.org/knuth/web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"modernc.org/knuth"
)

var (
	untyped = &notype{}
)

const (
	lib = "prg"
	rtl = "knuth"
)

type mode int

const (
	modeDefault mode = iota
	modeLValue
)

type notype struct{}

type ns struct {
	reg map[string]string
	cnt map[string]int
	g   *gen
}

func newNS() *ns {
	return &ns{
		reg: map[string]string{},
		cnt: map[string]int{
			// Names we cannot use or we want to keep for the backend's use.
			"bool":        1,
			"break":       1,
			"byte":        1,
			"case":        1,
			"chan":        1,
			"char":        1,
			"const":       1,
			"continue":    1,
			"default":     1,
			"defer":       1,
			"else":        1,
			"error":       1,
			"fallthrough": 1,
			"float32":     1,
			"float64":     1,
			"for":         1,
			"func":        1,
			"go":          1,
			"goto":        1,
			"if":          1,
			"import":      1,
			"init":        1,
			"int16":       1,
			"int32":       1,
			"int64":       1,
			"int8":        1,
			"interface":   1,
			"map":         1,
			"nil":         1,
			"package":     1,
			"range":       1,
			"return":      1,
			"rune":        1,
			"select":      1,
			"string":      1,
			"struct":      1,
			"switch":      1,
			"type":        1,
			"uint16":      1,
			"uint32":      1,
			"uint64":      1,
			"uint8":       1,
			"uintptr":     1,
			"var":         1,

			"abs":      1,
			"arraystr": 1,
			"close":    1,
			"fabs":     1,
			"ii":       1,
			"main":     1,
			"math":     1,
			"panic":    1,
			"r":        1,
			"round":    1,
			"signal":   1,
			"stderr":   1,
			"stdin":    1,
			"stdout":   1,
			"strcopy":  1,
			"unsafe":   1,
			lib:        1,
			rtl:        1,
		},
	}
}

func (n *ns) xid0(nm string) (r string) {
	nm = strings.ToLower(nm)
	a := strings.Split(nm, "_")
	for i, v := range a {
		if i != 0 {
			a[i] = strings.ToUpper(v[:1]) + v[1:]
		}
	}
	if r := n.reg[nm]; r != "" {
		return r
	}

	nm2 := strings.Join(a, "")
	c := n.cnt[nm2]
	if c == 0 {
		n.reg[nm] = nm2
		n.cnt[nm2] = 1
		return nm2
	}
	for {
		nm3 := fmt.Sprintf("%s%d", nm2, c)
		if n.cnt[nm3] == 0 {
			n.reg[nm] = nm3
			n.cnt[nm3] = c
			return nm3
		}

		c++
	}
}

func (n *ns) xid(id *identifier) (r string) {
	nm := strings.ToLower(id.Src())
	if nm == "byte" {
		return nm
	}

	if id.scope != nil {
		switch {
		case id.scope.isTLD():
			switch x := id.resolvedTo.(type) {
			case *constantDefinition:
				// ok
			case *variable, *procedureDeclaration, *functionDeclaration:
				return fmt.Sprintf("%s.%s", n.g.rcvrName(), n.xid0(nm))
			default:
				panic(todo("%v %T", id.ident, x))
			}
		}
	}

	return n.xid0(nm)
}

// Option adjusts the Pascal to Go transpiler behavior.
type Option func(p *genOpts) error

type genOpts struct {
	defs           map[string]struct{}
	positionalArgs []string // Substitutions for `@$`, see https://gitlab.com/cznic/knuth/-/issues/6#note_2451638009
}

func (g *genOpts) fillPositionalArgs(src []byte) (r []byte, err error) {
	if len(g.positionalArgs) == 0 {
		return src, nil
	}

	a := bytes.Split(src, []byte("@$"))
	if len(a)-1 != len(g.positionalArgs) {
		return nil, fmt.Errorf("need %v positionalArgs, got %v", len(a)-1, len(g.positionalArgs))
	}

	for i, v := range a {
		r = append(r, v...)
		if i != len(a)-1 {
			r = append(r, []byte(g.positionalArgs[i])...)
		}
	}
	return r, nil
}

func WithPositionalArgs(args ...string) Option {
	return func(p *genOpts) error {
		p.positionalArgs = append(p.positionalArgs, args...)
		return nil
	}
}

func WithDefines(defs ...string) Option {
	return func(p *genOpts) error {
		if p.defs == nil {
			p.defs = map[string]struct{}{}
		}
		for _, v := range defs {
			p.defs[v] = struct{}{}
		}
		return nil
	}
}

type gen struct {
	ast            *ast
	dest           io.Writer
	filesToClose   []string
	generatedBy    string
	ns             *ns
	opts           *genOpts
	packageName    string
	variantRecords map[string]*recordType
}

func newGen(dest io.Writer, ast *ast, packageName, generatedBy string, opts *genOpts) (*gen, error) {
	r := &gen{
		ast:            ast,
		dest:           dest,
		generatedBy:    generatedBy,
		ns:             newNS(),
		opts:           opts,
		packageName:    packageName,
		variantRecords: map[string]*recordType{},
	}
	r.ns.g = r
	return r, nil
}

func (g *gen) xid(id *identifier) string { return g.ns.xid(id) }
func (g *gen) xid0(nm string) string     { return g.ns.xid0(nm) }

func (g *gen) w(s string, args ...interface{}) {
	if _, err := fmt.Fprintf(g.dest, s, args...); err != nil {
		panic(abort(fmt.Errorf("writing backend result: %v", err)))
	}
}

func (g *gen) gen() (err error) {
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

	g.program(g.ast.program)
	return nil
}

func (g *gen) commentNL(s string) string {
	return g.comment(s)
}

// Pascal separator(s) -> Go comment(s)
func (g *gen) comment(s string) string {
	a, cnt := g.paserSeparators(s)
	w := 0
	for _, v := range a {
		switch v[0] {
		case '{':
			a[w] = v
			w++
		default:
			switch b := strings.Split(v, "\n"); len(b) {
			case 0, 1:
				a[w] = v
				w++
			case 2, 3:
				a[w] = strings.Repeat("\n", len(b)-1)
				w++
			default:
				a[w] = "\n\n"
				w++
			}
		}
	}
	a = a[:w]
	for i, v := range a {
		if v[0] != '{' {
			continue
		}

		v = v[1 : len(v)-1]
		if v[0] == ' ' {
			v = v[1:]
		}
		if strings.Index(v, "\n") < 0 {
			if i < len(a)-1 {
				if strings.Contains(a[i+1], "\n") {
					a[i] = "// " + v
					continue
				}

				if cnt == 1 {
					a[i] = fmt.Sprintf("/* %s */", v)
					continue
				}
			}
		}

		b := strings.Split(v, "\n")
		a[i] = "// " + strings.Join(b, "\n// ") + "\n"
		continue
	}
	return strings.Join(a, "")
}

func (g *gen) paserSeparators(s string) (r []string, cnt int) {
	for len(s) != 0 {
		switch s[0] {
		case ' ', '\t', '\n', '\r':
			if x := strings.IndexByte(s, '{'); x > 0 {
				r = append(r, s[:x])
				s = s[x:]
				break
			}

			return append(r, s), cnt
		case '{':
			x := strings.IndexByte(s, '}')
			r = append(r, s[:x+1])
			s = s[x+1:]
			cnt++
		default:
			panic(todo("%q", s))
		}
	}
	return r, cnt
}

func (g *gen) program(n *program) {
	g.w(`// Code generated by '%s', DO NOT EDIT.

%s

package %s

import (
	"math"
	"unsafe"

	"modernc.org/knuth"
)

var (
	_ = math.MaxInt32
	_ unsafe.Pointer
)

type (
	char = byte
	signal int
)

func strcopy(dst []char, src string) {
	for i := 0; i < len(dst) && i < len(src); i++ {
		dst[i] = char(src[i])
	}
}

func arraystr(a []char) string {
	b := make([]byte, len(a))
	for i, c := range a {
		b[i] = byte(c)
	}
	return string(b)
}

func abs(n int32) int32 {
	if n >= 0 {
		return n
	}

	return -n
}

func fabs(f float64) float64 {
	if f >= 0 {
		return f
	}

	return -f
}

func round(f float64) int32 {
	if f >= 0 {
		return int32(f+0.5)
	}

	return int32(f-0.5)
}

`,
		g.generatedBy, g.comment(n.programHeading.program.Sep()), g.packageName,
	)
	g.constants(n.block.constantDefinitionPart)
	g.types(n.block.typeDefinitionPart)
	if n.block.variableDeclarationPart != nil {
		g.w("\n\ntype %s struct{", lib)
		g.w("\nstdin, stdout, stderr knuth.File")
		g.filesToClose = []string{"stdin", "stdout", "stderr"}
		g.variables(n.block.variableDeclarationPart, false, true)
		g.w("\n}")
	}
	for i, v := range n.block.procedureAndFunctionDeclarationPart {
		g.procedureAndFunctionDeclarationPart(i, v)
	}
	g.w("\n\n%s", g.commentNL(n.block.statementPart.begin.Sep()))
	g.w("func %smain()", g.rcvr())
	var b strings.Builder
	if len(g.filesToClose) != 0 {
		b.WriteString("\ndefer func() {\n")
		sort.Strings(g.filesToClose)
		for _, v := range g.filesToClose {
			fmt.Fprintf(&b, "\tif %s.%s != nil { %[1]s.%[2]s.Close() }\n", lib, v)
		}
		b.WriteString("}()\n\n")
	}
	g.compoundStatement0(n.block.statementPart, true, false, b.String())
	g.w("%s", g.comment(g.ast.eof.Sep()))
}

func (g *gen) constants(n *constantDefinitionPart) {
	if n == nil || len(n.constantDefinitionList) == 0 {
		return
	}

	g.w("\n\n%sconst(", g.comment(n.const1.Sep()))
	var semi pasToken
	for _, v := range n.constantDefinitionList {
		semi = v.semi
		g.constantDefinition(v.constantDefinition)
		g.semi(semi)
	}
	c := g.comment(semi.Next().Sep())
	for strings.HasSuffix(c, "\n") {
		c = c[:len(c)-1]
	}
	g.w("%s\n)", c)
}

func (g *gen) constantDefinition(n *constantDefinition) {
	if n == nil {
		return
	}

	g.w("%s%s = %s", g.comment(n.ident.Sep()), g.xid(n.ident), g.constExpr(n.constant, false))
}

func (g *gen) semi(s pasToken) {
	c := g.comment(s.Sep())
	g.w("%s", c)
	if !strings.Contains(c, "\n") {
		g.w(";")
	}
}

func (g *gen) constExpr(n node, flat bool) string {
	var b strings.Builder
	g.constExpr0(&b, n, flat)
	return b.String()

}

func (g *gen) constExpr0(b *strings.Builder, n node, flat bool) {
	switch x := n.(type) {
	case *identifier:
		b.WriteString(g.idTok(x, flat))
	case pasToken:
		switch x.ch {
		case tokInt:
			b.WriteString(g.intTok(x, flat))
		case tokString:
			s := x.Src()
			s = s[1 : len(s)-1]
			s = strings.ReplaceAll(s, `""`, `"`)
			fmt.Fprintf(b, "%q", s)
		default:
			panic(todo("", x))
		}
	case *binaryExpression:
		if flat && x.value != nil {
			v := x.value.(int64)
		out:
			for x := x; ; {
				switch y := x.lhs.(type) {
				case *identifier:
					switch z := y.resolvedTo.(type) {
					case *constantDefinition:
						lv := z.value.(int64)
						switch x.op.ch {
						case '+':
							fmt.Fprintf(b, "%s+%v", g.idTok(y, true), v-lv)
							return
						case '-':
							// const x = 42; expr x-5 (37): lv = 42, v = 37
							fmt.Fprintf(b, "%s-%v", g.idTok(y, true), lv-v)
							return
						default:
							panic(todo("", x.op.Src()))
						}
					}
				case *binaryExpression:
					x = y
				default:
					break out
				}
			}
		}

		g.constExpr0(b, x.lhs, flat)
		switch {
		case flat:
			b.WriteString(x.op.Src())
		default:
			fmt.Fprintf(b, "%s%s", g.comment(x.op.Sep()), x.op.Src())
		}
		g.constExpr0(b, x.rhs, flat)
	case *signed:
		b.WriteString(" ")
		b.WriteString(x.sign.Src())
		g.constExpr0(b, x.node, flat)
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (g *gen) idTok(id *identifier, flat bool) string {
	switch {
	case flat:
		return g.xid(id)
	default:
		return fmt.Sprintf("%s%s", g.comment(id.Sep()), g.xid(id))
	}
}

func (g *gen) strTok(t pasToken) (string, interface{}) {
	s := t.Src()
	s = s[1 : len(s)-1]
	s = strings.ReplaceAll(s, "''", "'")
	return fmt.Sprintf("%s%q", g.comment(t.Sep()), s), &pascalString{}
}

func (g *gen) intTok(t pasToken, flat bool) string {
	const tag = "=}"
	sep := t.Sep()
	s := strings.TrimRight(sep, blankSet)
	tail := sep[len(s):]
	if strings.HasSuffix(s, tag) {
		x := strings.LastIndex(s, "{")
		sep = sep[:x] + tail
		switch hint := s[x+1 : len(s)-len(tag)]; {
		case strings.HasPrefix(hint, "0"):
			n, err := strconv.ParseUint(t.Src(), 10, 64)
			if err != nil {
				panic(todo(""))
			}

			if t.Src() == "0" {
				return fmt.Sprintf("%s0", g.comment(sep))
			}

			return fmt.Sprintf("%s0%o", sep, n)
		case strings.HasPrefix(hint, "\""):
			s := hint[1 : len(hint)-1]
			s = strings.ReplaceAll(s, `""`, `"`)
			if a := []rune(s); len(a) == 1 {
				num, err := strconv.ParseUint(t.Src(), 10, 32)
				if err != nil {
					panic(todo("", err))
				}

				return strconv.QuoteRuneToASCII(rune(num))
			}

			return fmt.Sprintf("%s/* %q */%s", g.comment(sep), s, t.Src())
		default:
			panic(todo("%q %v", hint, t))
		}
	}

	switch {
	case flat:
		return t.Src()
	default:
		return fmt.Sprintf("%s%s", g.comment(t.Sep()), t.Src())
	}
}

func (g *gen) types(n *typeDefinitionPart) {
	if n == nil || len(n.typeDefinitionList) == 0 {
		return
	}

	g.w("\n\ntype(")
	var semi pasToken
	for _, v := range n.typeDefinitionList {
		g.typeDefinition(v.typeDefinition)
		semi = v.semi
		g.w("%s", g.comment(semi.Sep()))
	}
	c := g.comment(semi.Next().Sep())
	for strings.HasSuffix(c, "\n") {
		c = c[:len(c)-1]
	}
	g.w("%s\n)", c)
	var a []string
	for k := range g.variantRecords {
		a = append(a, k)
	}
	sort.Strings(a)
	for _, tnm := range a {
		rt := g.variantRecords[tnm]
		var b []string
		for k := range rt.fields {
			b = append(b, k)
		}
		sort.Strings(b)
		for _, fnm := range b {
			f := rt.fields[fnm]
			ft := g.typeLiteral(f.typ, false)
			g.w("\n\nfunc (r *%s) %s() *%s {", tnm, fnm, ft)
			g.w("return (*%s)(unsafe.Add(unsafe.Pointer(&r.data), %d))", ft, f.off)
			g.w("\n}\n")
		}
	}
}

func (g *gen) typeDefinition(n *typeDefinition) {
	if n == nil {
		return
	}

	pnm := strings.ToLower(n.ident.Src())
	switch {
	case pnm == "byte":
		switch x := n.type1.(type) {
		case *subrangeType:
			if x.lo == 0 && x.hi == 255 {
				g.w("%s // %s = %s", g.comment(n.ident.Sep()), pnm, g.typeLiteral(x, true))
				return
			}
		}
	}
	switch x := n.typ.(type) {
	case *arrayType:
		g.w("%s%s %s;", g.comment(n.ident.Sep()), g.xid0(pnm), g.typeLiteral(n.type1, true))
	case *recordType:
		nm := g.xid0(pnm)
		eq := "= "
		if x.hasVariants {
			eq = ""
			g.variantRecords[nm] = x
		}
		g.w("%s%s %s%s;", g.comment(n.ident.Sep()), nm, eq, g.typeLiteral(n.type1, true))
	default:
		g.w("%s%s = %s;", g.comment(n.ident.Sep()), g.xid0(pnm), g.typeLiteral(n.type1, true))
	}
}

func (g *gen) strID(t interface{}) string {
	var b strings.Builder
	g.str0ID(&b, t)
	return b.String()
}

func (g *gen) str0ID(b *strings.Builder, t interface{}) {
	switch x := t.(type) {
	case *fileType:
		if x.isPacked() {
			b.WriteString("packed ")
		}
		b.WriteString("file of ")
		g.str0ID(b, x.elemType)
	case *subrangeType:
		b.WriteString(strings.TrimSpace(g.typeLiteral(x, false)))
	case *pascalChar:
		b.WriteString("char")
	case *recordType:
		if x.isPacked() {
			b.WriteString("packed ")
		}
		b.WriteString("record ")
		g.str0ID(b, x.fieldList)
		b.WriteString(" end;")
	case *fieldList:
		if x == nil {
			break
		}

		for _, v := range x.fixedPart {
			g.str0ID(b, v)
		}
		g.str0ID(b, x.variantPart)
	case *variantPart:
		if x == nil {
			break
		}

		b.WriteString("case ")
		g.str0ID(b, x.variantSelector)
		b.WriteString(" of ")
		for i, v := range x.variants {
			if i != 0 {
				b.WriteString(" ")
			}
			g.str0ID(b, v)
		}
	case *variantSelector:
		if x.tagField.isValid() {
			fmt.Fprintf(b, "%s: ", x.tagField.Src())
		}
		g.str0ID(b, x.typ)
	case *variant:
		for i, v := range x.constList {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(g.constExpr(v.const1, true))
		}
		b.WriteString(": (")
		g.str0ID(b, x.fieldList)
		b.WriteString(");")
	case *fixedPart:
		g.str0ID(b, x.recordSection)
		b.WriteString(";")
	case *recordSection:
		for i, v := range x.identifierList {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(strings.ToLower(v.ident.Src()))
		}
		b.WriteString(" ")
		g.str0ID(b, x.typ)
	case *pascalInteger, *pascalMaxInt:
		b.WriteString("int32")
	case *pascalReal:
		b.WriteString("float64")
	case *pascalReal32:
		b.WriteString("float32")
	case *pascalBoolean:
		b.WriteString("boolean")
	case *pascalString:
		b.WriteString("<string>")
	case *pascalText:
		b.WriteString("<text>")
	case *arrayType:
		b.WriteString("[")
		for i, v := range x.indexTyp {
			if i != 0 {
				b.WriteString(". ")
			}
			g.str0ID(b, v)
		}
		b.WriteString("]")
		g.str0ID(b, x.elemTyp)
	case *identifier:
		b.WriteString(x.Src())
	case *typeDefinition:
		switch y := underlyingType(x).(type) {
		case *subrangeType, *pascalInteger, *pascalReal, *recordType:
			g.str0ID(b, y)
		default:
			panic(todo("%T", y))
		}
	case *notype:
		b.WriteString("<any>")
	default:
		panic(todo("%T", x))
	}
}

func (g *gen) identifierList(w io.Writer, n []*identifierList) (names []string) {
	for i, v := range n {
		if i != 0 {
			fmt.Fprintf(w, "%s,", g.comment(v.comma.Sep()))
		}
		id := g.xid(v.ident)
		names = append(names, id)
		fmt.Fprintf(w, "%s%s", g.comment(v.ident.Sep()), id)
	}
	return names
}

func (g *gen) typeLiteral(t interface{}, comments bool) string {
	var b strings.Builder
	g.typeLiteral0(&b, t, comments)
	return b.String()
}

func (g *gen) typeLiteral0(b *strings.Builder, t interface{}, comments bool) {
	switch x := t.(type) {
	case *subrangeType:
		if comments {
			fmt.Fprintf(b, " /* %s..%s */ ", g.constExpr(x.constant, true), g.constExpr(x.constant2, true))
		}
		switch {
		case x.lo >= 0 && x.hi <= math.MaxUint8:
			b.WriteString(" byte")
		case x.lo >= 0 && x.hi <= math.MaxUint16:
			b.WriteString(" uint16")
		case x.lo >= 0 && x.hi <= math.MaxUint32:
			b.WriteString(" uint32")
		case x.lo >= math.MinInt8 && x.hi <= math.MaxInt8:
			b.WriteString(" int8")
		case x.lo >= math.MinInt16 && x.hi <= math.MaxInt16:
			b.WriteString(" int16")
		default:
			panic(todo("", x.lo, x.hi))
		}
	case *fileType:
		// fmt.Fprintf(b, "/* %s */", g.strID(x))
		b.WriteString("knuth.File")
	case *identifier:
		switch y := g.ast.tldScope.mustLookup(x).(type) {
		case *pascalText:
			fmt.Fprintf(b, "/* %s */ knuth.File", g.strID(x))
		case *pascalInteger, *pascalBoolean, *pascalReal, *pascalChar:
			g.typeLiteral0(b, y, comments)
		case *typeDefinition:
			g.typeLiteral0(b, y, comments)
		default:
			panic(todo("%v: %T", x.Position(), y))
		}
	case *arrayType:
		if x.isPacked() && comments {
			fmt.Fprintf(b, "/* packed */ ")
		}
		for _, v := range x.indexTyp {
			fmt.Fprintf(b, "[%d]", g.cardinality(v))
		}
		g.typeLiteral0(b, x.elemTyp, comments)
	case *recordType:
		if x.isPacked() && comments {
			fmt.Fprintf(b, "/* packed */ ")
		}
		// fmt.Fprintf(b, "/* %s */", g.strID(x))
		g.recordTypeLiteral(b, x)
	case *pascalBoolean:
		b.WriteString("bool")
	case *pascalInteger:
		b.WriteString("int32")
	case *pascalChar:
		b.WriteString("char")
	case *pascalReal:
		b.WriteString("float64")
	case *pascalReal32:
		b.WriteString("float32")
	case *typeDefinition:
		b.WriteString(g.xid(x.ident))
	default:
		panic(todo("%T", x))
	}
}

func (g *gen) recordTypeLiteral(b *strings.Builder, t *recordType) {
	if t.hasVariants || t.isPacked() {
		switch t.size {
		case 1:
			fmt.Fprintf(b, "struct{ data byte }")
		case 2:
			fmt.Fprintf(b, "struct{ data uint16 }")
		case 4:
			fmt.Fprintf(b, "struct{ data uint32 }")
		case 8:
			fmt.Fprintf(b, "struct{ data uint64 }")
		default:
			fmt.Fprintf(b, "struct{ data uint64; _ [%d]byte }", t.size-8)
		}
		return
	}

	b.WriteString("struct{")
	for _, v := range t.fieldList.fixedPart {
		rs := v.recordSection
		g.identifierList(b, rs.identifierList)
		b.WriteString(" ")
		g.typeLiteral0(b, rs.typ, false)
		b.WriteString(";")
	}
	c := g.comment(t.fieldList.semi.Next().Sep())
	for strings.HasSuffix(c, "\n") {
		c = c[:len(c)-1]
	}
	fmt.Fprintf(b, "%s\n}", c)
}

func (g *gen) cardinality(t interface{}) int64 {
	switch x := underlyingType(t).(type) {
	case *subrangeType:
		return x.hi - x.lo + 1
	case *pascalChar:
		return 256
	default:
		panic(todo("%T", x))
	}
}

func (g *gen) variables(n *variableDeclarationPart, wrap, registerFiles bool) {
	if n == nil || len(n.variableDeclarationList) == 0 {
		return
	}

	if wrap {
		g.w("%svar(", g.comment(n.var1.Sep()))
		defer g.w("\n);")
	}
	var semi pasToken
	for _, v := range n.variableDeclarationList {
		names, typ := g.varDeclaration(v.variableDeclaration)
		if g.isFile(typ) {
			g.filesToClose = append(g.filesToClose, names...)
		}
		semi = v.semi
		g.w("%s", g.comment(semi.Sep()))
	}
	c := g.comment(semi.Next().Sep())
	for strings.HasSuffix(c, "\n") {
		c = c[:len(c)-1]
	}
	g.w("%s", c)
}

func (g *gen) isFile(t interface{}) bool {
	switch x := underlyingType(t).(type) {
	case *fileType:
		return true
	case *identifier:
		switch y := g.ast.tldScope.mustLookup(x).(type) {
		case *pascalText:
			return true
		case *typeDefinition:
			return g.isFile(y)
		}
	}
	return false
}

func (g *gen) varDeclaration(n *variableDeclaration) (names []string, typ node) {
	if n == nil {
		return
	}

	names = g.identifierList(g.dest, n.identifierList)
	g.w(" %s;", g.typeLiteral(n.type1, true))
	return names, n.type1
}

func (g *gen) procedureAndFunctionDeclarationPart(ix int, n *procedureAndFunctionDeclarationPart) {
	if n == nil {
		return
	}

	switch x := n.declaration.(type) {
	case *procedureDeclaration:
		g.procedureDeclaration(ix, x)
	case *functionDeclaration:
		g.functionDeclaration(ix, x)
	default:
		panic(todo("%v: %T", n.declaration.Position(), x))
	}
}

func (g *gen) procedureDeclaration(ix int, n *procedureDeclaration) {
	if n == nil || n.isFwd() {
		return
	}

	switch {
	case ix == 0:
		g.w("\n\n")
	default:
		g.w("%s", g.commentNL(n.procedureHeading.procedure.Sep()))
	}
	g.w("func %s %s", g.rcvr(), g.xid(n.procedureHeading.ident))
	g.formalParameterList(n.procedureHeading.formalParameterList)
	g.block(n.block.(*block), false)
}

func (g *gen) rcvr() string     { return fmt.Sprintf("(%s *%s)", g.rcvrName(), lib) }
func (g *gen) rcvrName() string { return lib }

func (g *gen) functionDeclaration(ix int, n *functionDeclaration) {
	if n == nil || n.isFwd() {
		return
	}

	switch {
	case ix == 0:
		g.w("\n\n")
	default:
		g.w("%s", g.commentNL(n.functionHeading.function.Sep()))
	}
	g.w("func %s %s", g.rcvr(), g.xid(n.functionHeading.ident))
	g.formalParameterList(n.functionHeading.formalParameterList)
	g.w("(r %s)", g.typeLiteral(n.functionHeading.result, true))
	g.block(n.block.(*block), true)
}

func (g *gen) formalParameterList(n *formalParameterList) {
	g.w("(")

	defer g.w(")")

	if n == nil {
		return
	}

	for i, v := range n.params {
		if i != 0 {
			g.w("%s,", g.comment(v.semi.Sep()))
		}

		switch x := v.param.(type) {
		case *parameterSpecification:
			g.parameterSpecification(x)
		default:
			panic(todo("%v: %T", v.param.Position(), x))
		}
	}
}

func (g *gen) parameterSpecification(n *parameterSpecification) {
	if n == nil {
		return
	}

	g.identifierList(g.dest, n.identifierList)
	g.w(" %s", g.typeLiteral(n.typ, true))
}

func (g *gen) block(n *block, injectRetval bool) {
	if n == nil {
		return
	}

	g.w("{")
	g.constants(n.constantDefinitionPart)
	g.types(n.typeDefinitionPart)
	g.variables(n.variableDeclarationPart, true, false)
	g.compoundStatement(n.statementPart, false, true)
	if injectRetval {
		g.w("return r")
	}
	g.w("}")
}

func (g *gen) compoundStatement(n *compoundStatement, braced, comment bool) {
	g.compoundStatement0(n, braced, comment, "")
}

func (g *gen) compoundStatement0(n *compoundStatement, braced, comment bool, inject string) {
	if n == nil {
		return
	}

	if braced {
		g.w("{")
		defer g.w("}")
	}

	g.w("%s", inject)
	for _, v := range n.statementSequence {
		g.w("%s", g.comment(v.semi.Next().Sep()))
		g.statement(v.statement)
		g.w(";")
	}
}

func (g *gen) exprStr(n node) (s string, t interface{}) {
	var b strings.Builder
	t = g.exprStr0(&b, n)
	return b.String(), t
}

func (g *gen) exprStr0(b *strings.Builder, n node) (t interface{}) {
	switch x := n.(type) {
	case pasToken:
		switch x.ch {
		case tokString:
			s, t := g.strTok(x)
			b.WriteString(s)
			return t
		case tokInt:
			b.WriteString(g.intTok(x, false))
			return untyped
		case tokFloat:
			b.WriteString(x.Src())
			return &pascalReal{}
		default:
			panic(todo("", x))
		}
	case *identifier:
		switch y := x.resolvedTo.(type) {
		case *variable, *constantDefinition:
			b.WriteString(g.xid(x))
			return x.typ
		case *pascalFalse:
			b.WriteString(g.xid0("false"))
			return &pascalBoolean{}
		case *pascalTrue:
			b.WriteString(g.xid0("true"))
			return &pascalBoolean{}
		case *pascalInput:
			fmt.Fprintf(b, "%s.stdin", g.rcvrName())
			return y
		case *pascalOutput:
			fmt.Fprintf(b, "%s.stdout", g.rcvrName())
			return y
		case *pascalStderr:
			fmt.Fprintf(b, "%s.stderr", g.rcvrName())
			return y
		case *functionDeclaration:
			b.WriteString("r")
			return y.functionHeading.result
		case *pascalMaxInt:
			b.WriteString("math.MaxInt32")
			return y
		default:
			panic(todo("%v: %q in scope %p %T", x.Position(), x.Src(), x.scope, y))
		}
	case *signed:
		b.WriteString(" ")
		b.WriteString(x.sign.Src())
		s, t := g.exprStr(x.node)
		switch y := underlyingType(t).(type) {
		case *subrangeType:
			switch {
			case y.lo >= 0:
				fmt.Fprintf(b, "int32(%s)", s)
				return &pascalInteger{}
			default:
				b.WriteString(s)
				return x.typ
			}
		case *notype, *pascalReal, *pascalInteger, *pascalMaxInt:
			b.WriteString(s)
			return x.typ
		default:
			panic(todo("%v: %T %q", x.Position(), y, s))
		}
	case *binaryExpression:
		b.WriteString("(")
		ct := x.typ
		switch x.op.ch {
		case '=', '<', '>', tokNeq, tokLeq, tokGeq:
			ct = x.cmpTyp
		}
		b.WriteString(g.expr(ct, x.lhs))
		switch x.op.ch {
		case '+', '-', '*', '>', '<':
			b.WriteRune(x.op.ch)
		case '/', tokDiv:
			b.WriteString("/")
		case tokMod:
			b.WriteString("%")
		case '=':
			b.WriteString("==")
		case tokNeq:
			b.WriteString("!=")
		case tokGeq:
			b.WriteString(">=")
		case tokLeq:
			b.WriteString("<=")
		case tokOr:
			b.WriteString("||")
		case tokAnd:
			b.WriteString("&&")
		default:
			panic(todo("", x.op))
		}
		b.WriteString(g.expr(ct, x.rhs))
		b.WriteString(")")
		return x.typ
	case *parenthesizedExpression:
		b.WriteString("(")
		g.exprStr0(b, x.expr)
		b.WriteString(")")
		return x.typ
	case *indexedVariable:
		g.exprStr0(b, x.variable)
		b.WriteString("[")
		at, ok := underlyingType(x.varTyp).(*arrayType)
		if !ok {
			panic(todo("%v: %T", x.variable.Position(), x.varTyp))
		}

		for i, v := range at.indexTyp {
			if i != 0 {
				b.WriteString("][")
			}
			ie := x.indexList[i].index
			switch y := v.(type) {
			case *subrangeType:
				switch {
				case y.lo == 0:
					g.exprStr0(b, ie)
				default:
					b.WriteString("(")
					g.exprStr0(b, ie)
					switch {
					case y.lo < 0:
						fmt.Fprintf(b, ")+%d", -y.lo)
					default:
						fmt.Fprintf(b, ")-%d", y.lo)
					}
				}
			case *pascalChar:
				s, t := g.exprStr(ie)
				switch z := t.(type) {
				case *pascalChar:
					b.WriteString(s)
				case *pascalString:
					s, err := strconv.Unquote(strings.TrimSpace(s))
					if err != nil {
						panic(todo("", err))
					}

					a := []rune(s)
					if len(a) != 1 {
						panic(todo("`%s` -> %v", s, a))
					}

					b.WriteString(strconv.QuoteRuneToASCII(a[0]))
				default:
					panic(todo("%v: %T %q", ie.Position(), z, s))
				}
			default:
				panic(todo("%v: %T", x.lbracket.Position(), y))
			}
		}
		b.WriteString("]")
		return x.typ
	case *fieldDesignator:
		s, t := g.exprStr(x.variable)
		variant := false
		switch y := underlyingType(t).(type) {
		case *recordType:
			if y.hasVariants {
				variant = true
				b.WriteString("(*(")
			}
		default:
			panic(todo("%v: %T", x.dot.Position(), y))
		}
		b.WriteString(s)
		fmt.Fprintf(b, ".%s", g.xid(x.ident))
		if variant {
			b.WriteString("()))")
		}
		return x.typ
	case *deref:
		s, t := g.exprStr(x.n)
		switch y := underlyingType(t).(type) {
		case *fileType:
			switch z := underlyingType(y.elemType).(type) {
			case *recordType:
				switch {
				case z.hasVariants:
					fmt.Fprintf(b, "(*(*%s)(unsafe.Pointer(%s.Data%dP())))", g.typeLiteral(y.elemType, false), s, z.size)
				default:
					panic(todo("%T B", z))
				}
			default:
				fmt.Fprintf(b, "(*(%s.%s()))", s, g.typeLiteralP(y.elemType))
			}
			return y.elemType
		case *pascalInput:
			fmt.Fprintf(b, "(*(%s.ByteP()))", s)
			return &pascalChar{}
		default:
			panic(todo("%v: %T %q", x.n.Position(), y, s))
		}
	case *functionCall:
		switch y := x.ft.(type) {
		case *functionDeclaration:
			fmt.Fprintf(b, "%s", g.xid(x.ident))
			b.WriteString(g.parameters(x.ft, x.parameters, 0))
		case *pascalAbs:
			arg := x.parameters.parameters[0]
			switch y := underlyingType(arg.typ).(type) {
			case *pascalInteger, *subrangeType:
				b.WriteString("abs(")
				b.WriteString(g.expr(&pascalInteger{}, arg.parameter))
				b.WriteString(")")
			case *pascalReal:
				b.WriteString("fabs(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		case *pascalChr:
			arg := x.parameters.parameters[0]
			switch y := arg.typ.(type) {
			case *pascalInteger, *subrangeType:
				b.WriteString("char(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		case *pascalOrd:
			arg := x.parameters.parameters[0]
			switch y := arg.typ.(type) {
			case *pascalInteger, *subrangeType:
				b.WriteString("int32(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")")
			case *pascalChar:
				b.WriteString("int32(")
				b.WriteString(g.expr(&pascalInteger{}, arg.parameter))
				b.WriteString(")")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		case *pascalEOF:
			arg := x.parameters.parameters[0]
			g.exprStr0(b, arg.parameter)
			b.WriteString(".EOF()")
		case *pascalEOLN:
			arg := x.parameters.parameters[0]
			g.exprStr0(b, arg.parameter)
			b.WriteString(".EOLN()")
		case *knuthErstat:
			arg := x.parameters.parameters[0]
			g.exprStr0(b, arg.parameter)
			b.WriteString(".ErStat()")
		case *knuthCurPos:
			arg := x.parameters.parameters[0]
			g.exprStr0(b, arg.parameter)
			b.WriteString(".CurPos()")
		case *pascalOdd:
			arg := x.parameters.parameters[0]
			switch y := underlyingType(arg.typ).(type) {
			case *pascalInteger, *subrangeType:
				b.WriteString("((")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")&1 != 0)")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		case *pascalRound:
			arg := x.parameters.parameters[0]
			switch y := arg.typ.(type) {
			case *pascalReal:
				b.WriteString("round(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")")
			case *pascalReal32:
				b.WriteString("round(float64(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString("))")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		case *pascalTrunc:
			arg := x.parameters.parameters[0]
			switch y := arg.typ.(type) {
			case *pascalReal:
				b.WriteString("int32(")
				b.WriteString(g.expr(untyped, arg.parameter))
				b.WriteString(")")
			default:
				panic(todo("%v: %T", arg.parameter.Position(), y))
			}
		default:
			panic(todo("%v: %T", x.Position(), y))
		}
		return x.typ
	case *not:
		b.WriteString("!(")
		g.exprStr0(b, x.factor)
		b.WriteString(")")
		return x.typ
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (g *gen) typeLiteralP(t interface{}) string {
	switch x := underlyingType(t).(type) {
	case *subrangeType:
		s := strings.TrimSpace(g.typeLiteral(x, false))
		s = strings.ToUpper(s[:1]) + s[1:]
		return s + "P"
	case *pascalChar:
		return "ByteP"
	default:
		panic(todo("%T", x))
	}
}

func (g *gen) expr(t interface{}, n node) string {
	s, st := g.exprStr(n)
	return g.convert(t, st, s)
}

func (g *gen) convert(dt, st interface{}, expr string) string {
	if _, ok := dt.(*notype); ok {
		return expr
	}

	if dt == st {
		return expr
	}

	dts := g.strID(dt)
	sts := g.strID(st)
	if dts == sts {
		return expr
	}

	switch x := dt.(type) {
	case *pascalInteger:
		switch y := st.(type) {
		case *subrangeType, *typeDefinition, *pascalChar:
			return fmt.Sprintf("int32(%s)", expr)
		case *notype:
			return expr
		case *pascalString:
			s, err := strconv.Unquote(strings.TrimSpace(expr))
			if err != nil {
				panic(todo("", err))
			}

			a := []rune(s)
			if len(a) != 1 {
				panic(todo("`%s`.%d", s, len(a)))
			}

			return strconv.QuoteRuneToASCII(a[0])
		default:
			panic(todo("%T %q", y, expr))
		}
	case *pascalChar:
		switch y := st.(type) {
		case *pascalString:
			s, err := strconv.Unquote(strings.TrimSpace(expr))
			if err != nil {
				panic(todo("`%s` %q %v", expr, s, err))
			}

			a := []rune(s)
			if len(a) != 1 {
				panic(todo("`%s`.%d", s, len(a)))
			}

			return strconv.QuoteRuneToASCII(a[0])
		default:
			panic(todo("%T", y))
		}
	case *pascalReal:
		return fmt.Sprintf("float64(%s)", expr)
	case *pascalReal32:
		return fmt.Sprintf("float32(%s)", expr)
	case *typeDefinition:
		switch y := underlyingType(x).(type) {
		case *subrangeType, *pascalInteger:
			return fmt.Sprintf("%s(%s)", g.xid(x.ident), expr)
		default:
			panic(todo("%T", y))
		}
	case *subrangeType:
		switch y := st.(type) {
		case *pascalInteger, *subrangeType:
			return fmt.Sprintf("%s(%s)", g.typeLiteral(x, false), expr)
		case *typeDefinition:
			return fmt.Sprintf("%s(%s)", g.typeLiteral(underlyingType(x), false), expr)
		case *notype:
			return expr
		default:
			panic(todo("%T", y))
		}
	default:
		panic(todo("%s <- %s %T, %s", dts, sts, x, expr))
	}
}

func (g *gen) procedureStatement(n *procedureStatement) {
	switch n.typ.(type) {
	case *pascalWrite, *pascalWriteln:
		switch len(n.parameters.parameters) {
		case 0:
			g.w("%s.stdout.Write", g.rcvrName())
			if _, ok := n.typ.(*pascalWriteln); ok {
				g.w("ln")
			}
			g.w("()")
		default:
			p0 := n.parameters.parameters[0]
			skip := 0
			switch y := underlyingType(p0.typ).(type) {
			case *stringLiteral, *pascalInteger, *subrangeType, *pascalChar, *pascalReal:
				g.w("%s.stdout", g.rcvrName())
			case *pascalText, *fileType, *pascalOutput, *pascalStderr:
				g.w("%s", g.expr(untyped, p0.parameter))
				skip = 1
			case *arrayType:
				switch z := y.elemTyp.(type) {
				case *pascalChar:
					g.w("%s.stdout", g.rcvrName())
				default:
					panic(todo("%v: %T", n.Position(), z))
				}
			default:
				panic(todo("%v: %T", n.Position(), y))
			}
			g.w(".Write")
			if _, ok := n.typ.(*pascalWriteln); ok {
				g.w("ln")
			}
			g.w("%s", g.parameters(n.typ, n.parameters, skip))
		}
	case *pascalReset, *pascalRewrite, *pascalRead, *pascalReadln:
		switch len(n.parameters.parameters) {
		case 0:
			panic(todo("%v:", n.Position()))
		default:
			skip := 0
			p0 := n.parameters.parameters[0]
			switch y := underlyingType(p0.typ).(type) {
			case *fileType, *pascalText, *pascalInput:
				g.w("%s", g.expr(untyped, p0.parameter))
				skip = 1
			case *pascalInteger, *subrangeType, *pascalChar, *pascalReal:
				g.w("%s.stdin", g.rcvrName())
			default:
				panic(todo("%v: %T", n.Position(), y))
			}
			switch n.typ.(type) {
			case *pascalRead:
				g.w(".Read%s", g.parameters(n.typ, n.parameters, skip))
			case *pascalReadln:
				g.w(".Readln%s", g.parameters(n.typ, n.parameters, skip))
			case *pascalReset:
				g.w(".Reset%s", g.parameters(n.typ, n.parameters, skip))
			case *pascalRewrite:
				g.w(".Rewrite%s", g.parameters(n.typ, n.parameters, skip))
			default:
				panic(todo(""))
			}
		}
	case *pascalGet:
		switch len(n.parameters.parameters) {
		case 1:
			g.w("%s.Get()", g.expr(untyped, n.parameters.parameters[0].parameter))
		default:
			panic(todo("%v:", n.Position()))
		}
	case *pascalPut:
		switch len(n.parameters.parameters) {
		case 1:
			g.w("%s.Put()", g.expr(untyped, n.parameters.parameters[0].parameter))
		default:
			panic(todo("%v:", n.Position()))
		}
	case *knuthSetPos:
		switch len(n.parameters.parameters) {
		case 2:
			g.w("%s.SetPos%s", g.expr(untyped, n.parameters.parameters[0].parameter), g.parameters(n.typ, n.parameters, 1))
		default:
			panic(todo("%v:", n.Position()))
		}
	case *knuthClose:
		switch len(n.parameters.parameters) {
		case 1:
			g.w("%s.Close()", g.expr(untyped, n.parameters.parameters[0].parameter))
		default:
			panic(todo("%v:", n.Position()))
		}
	case *knuthPanic:
		s, _ := g.exprStr(n.parameters.parameters[0].parameter)
		g.w("panic(signal(%s))", s)
	default:
		g.w("%s", g.xid(n.ident))
		g.w("%s", g.parameters(n.typ, n.parameters, 0))
	}
}

func (g *gen) parameters(calleeType interface{}, n *parameters, skip int) (r string) {
	var b strings.Builder
	b.WriteString("(")
	defer func() {
		b.WriteString(")")
		r = b.String()
	}()

	if n == nil {
		return
	}

	params := n.parameters[skip:]
	switch x := calleeType.(type) {
	case *pascalReset, *pascalRewrite:
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			switch es, et := g.exprStr(v.parameter); x := et.(type) {
			case *arrayType:
				switch x.elemTyp.(type) {
				case *pascalChar:
					fmt.Fprintf(&b, "arraystr(%s[:])", es)
					continue
				}
			}
			b.WriteString(g.expr(untyped, v.parameter))
		}
	case *knuthClose, *knuthBreak, *knuthBreakIn, *knuthSetPos, *knuthPanic:
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(g.expr(untyped, v.parameter))
		}
	case *pascalWrite, *pascalWriteln:
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			switch y := underlyingType(v.parameter).(type) {
			case *writeParameter:
				es, t := g.exprStr(y.param)
				switch t.(type) {
				case *pascalChar:
					fmt.Fprintf(&b, "string(rune(%s))", es)
				default:
					b.WriteString(es)
				}

				fmt.Fprintf(&b, ", %s.WriteWidth(%d)", rtl, y.wval)
				if y.wval2 != nil {
					fmt.Fprintf(&b, ", %s.WriteWidth(%d)", rtl, y.wval2)
				}
			default:
				es, t := g.exprStr(v.parameter)
				if _, ok := t.(*pascalChar); ok {
					fmt.Fprintf(&b, "string(rune(%s))", es)
					break
				}

				switch x := underlyingType(t).(type) {
				case *recordType:
					switch {
					case x.hasVariants:
						b.WriteString("&")
						b.WriteString(g.expr(untyped, v.parameter))
						b.WriteString(".data")
					default:
						fmt.Fprintf(&b, "(*[%d]byte)(unsafe.Pointer(&%s))", x.size, g.expr(untyped, v.parameter))
					}
				default:
					b.WriteString(g.expr(untyped, v.parameter))
				}
			}
		}
	case *pascalRead, *pascalReadln:
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			switch y := underlyingType(v.typ).(type) {
			case *fileType, *pascalText:
				// ok
			case *pascalInteger, *subrangeType, *pascalChar:
				b.WriteString("&")
			case *recordType:
				switch {
				case y.hasVariants:
					b.WriteString("&")
					b.WriteString(g.expr(untyped, v.parameter))
					b.WriteString(".data")
				default:
					fmt.Fprintf(&b, "(*[%d]byte)(unsafe.Pointer(&%s))", y.size, g.expr(untyped, v.parameter))
				}
				continue
			default:
				panic(todo("%v: %T", v.parameter.Position(), y))
			}
			b.WriteString(g.expr(untyped, v.parameter))
		}
	case *procedureDeclaration:
		fp := x.procedureHeading.fp
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(g.expr(fp[i], v.parameter))
		}
	case *functionDeclaration:
		fp := x.functionHeading.fp
		for i, v := range params {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(g.expr(fp[i], v.parameter))
		}
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
	return
}

func (g *gen) assignmentStatement(n *assignmentStatement) {
	s, t := g.exprStr(n.lhs)
	switch x := t.(type) {
	case *arrayType:
		switch len(x.indexTyp) {
		case 1:
			switch y := x.elemTyp.(type) {
			case *pascalChar:
				switch z := n.expr.(type) {
				case pasToken:
					switch z.ch {
					case tokString:
						s2, _ := g.strTok(z)
						g.w("strcopy(%s[:], %s)", s, s2)
						return
					default:
						panic(todo("", z))
					}
				case *identifier:
					switch a := z.resolvedTo.(type) {
					case *constantDefinition:
						switch b := a.value.(type) {
						case string:
							g.w("strcopy(%s[:], %q)", s, b)
							return
						default:
							panic(todo("%v: %T", z.Position(), b))
						}
					default:
						panic(todo("%v: %T", z.Position(), a))
					}
				default:
					panic(todo("%v: %T", n.Position(), z))
				}
			default:
				panic(todo("%v: %T", n.Position(), y))
			}
		default:
			panic(todo("%v: %d", n.Position(), len(x.indexTyp)))
		}
	case
		*pascalInteger, *subrangeType, *pascalBoolean, *pascalReal, *pascalReal32,
		*pascalChar, *recordType:

		// ok
	case *functionDeclaration:
		t = x.functionHeading.result
	case *typeDefinition:
		t = underlyingType(t)
	default:
		panic(todo("%v: %T", n.Position(), x))
	}

	g.w("%s", s)
	g.w(" = ")
	g.w("%s", g.expr(t, n.expr))
}

func (g *gen) forStatement(n *forStatement) {
	inc := "++"
	cmp := "<="
	switch n.direction.ch {
	case tokTo:
		// ok
	case tokDownto:
		inc = "--"
		cmp = ">="
	default:
		panic(todo(""))
	}
	g.w("for ii := int32(%s); ii %s %s; ii%s {", g.expr(untyped, n.initialValue), cmp, g.expr(&pascalInteger{}, n.finalValue), inc)
	g.w("%s = %s; _ = %[1]s;", g.xid(n.variable), g.convert(n.variable.typ, &pascalInteger{}, "ii"))
	g.w("%s", g.comment(n.do.Next().Sep()))
	g.unbracedStatement(n.statement)
	g.w("}")
}

func (g *gen) unbracedStatement(n node) {
	switch x := n.(type) {
	case *compoundStatement:
		g.compoundStatement(x, false, false)
	default:
		g.statement(n)
	}
}

func (g *gen) ifStatement(n *ifStatement) {
	g.w("if %s", g.expr(untyped, n.expr))
	g.w("{%s", g.comment(n.if1.Next().Sep()))
	g.unbracedStatement(n.ifStmt)
	g.w("}")
}

func (g *gen) ifElseStatement(n *ifElseStatement) {
	g.w("if %s", g.expr(untyped, n.expr))
	g.w("{%s", g.comment(n.if1.Next().Sep()))
	g.unbracedStatement(n.ifStmt)
	g.w("} else ")
	switch x := n.elseStmt.(type) {
	case *ifStatement, *ifElseStatement:
		g.statement(x)
	default:
		g.w("{%s", g.comment(n.else1.Next().Sep()))
		g.unbracedStatement(n.elseStmt)
		g.w("}")
	}
}

func (g *gen) whileStatement(n *whileStatement) {
	g.w("for ")
	v, _ := g.exprStr(n.expr)
	g.w("%s", v)
	g.w("{%s", g.comment(n.do.Next().Sep()))
	g.unbracedStatement(n.stmt)
	g.w("}")
}

func (g *gen) repeatStatement(n *repeatStatement) {
	g.w("for ")
	g.w("{%s", g.comment(n.repeat.Next().Sep()))
	for _, v := range n.stmt {
		g.w("%s", g.comment(v.semi.Next().Sep()))
		g.statement(v.statement)
		g.w(";")
	}
	v, _ := g.exprStr(n.expr)
	g.w("if %s { break }", v)
	g.w("}")
}

func (g *gen) caseStatement(n *caseStatement) {
	g.w("switch ")
	v, _ := g.exprStr(n.expr)
	g.w("%s {%s", v, g.comment(n.of.Sep()))
	for _, v := range n.cases {
		g.w("%s%s%scase ", g.comment(v.semi.Sep()), v.semi.Src(), g.comment(v.semi.Next().Sep()))
		for i, w := range v.case1.constList {
			g.w("%s", w.comma.Src())
			if i != 0 && i%4 == 0 {
				g.w("\n")
			}
			g.w("%s", g.constExpr(w.const1, true))
		}
		g.w(":%s", g.comment(v.case1.comma.Next().Sep()))
		g.unbracedStatement(v.case1.stmt)
	}
	if n.else1.isValid() {
		g.w("\n;%sdefault: ", g.comment(n.else1.Sep()))
		g.unbracedStatement(n.elseStmt)
	}
	g.w("%s}", g.comment(n.end.Sep()))
}

func (g *gen) labeled(n *labeled) {
	s := g.label(n.label)
	switch n.off.(type) {
	case nil:
		// ok
	default:
		var op string
		switch n.plus.ch {
		case '+':
			op = "plus"
		default:
			panic(todo("", n.plus))
		}
		s = fmt.Sprintf("%s_%s_%s", s, op, g.expr(untyped, n.off))
	}
	g.w("%s:", s)
	g.statement(n.stmt)
}

func (g *gen) label(n node) string {
	switch x := n.(type) {
	case *identifier:
		return g.idTok(x, false)
	case pasToken:
		switch x.ch {
		case tokInt:
			return fmt.Sprintf("_%s", x.Src())
		default:
			panic(todo("%v: %q", n.Position(), x.Src()))
		}
	case *binaryExpression:
		var op string
		switch x.op.ch {
		case '+':
			op = "plus"
		default:
			panic(todo("", x.op))
		}
		return fmt.Sprintf("%s_%s_%s", g.label(x.lhs), op, g.expr(untyped, x.rhs))
	default:
		panic(todo("%v: %T", n.Position(), x))
	}
}

func (g *gen) gotoStatement(n *gotoStatement) {
	g.w("%sgoto %s", g.comment(n.goto1.Sep()), g.label(n.label))
}

func (g *gen) statement(n node) {
	switch x := n.(type) {
	case *emptyStatement:
	// nop
	case *procedureStatement:
		g.procedureStatement(x)
	case *assignmentStatement:
		g.assignmentStatement(x)
	case *forStatement:
		g.forStatement(x)
	case *ifStatement:
		g.ifStatement(x)
	case *ifElseStatement:
		g.ifElseStatement(x)
	case *whileStatement:
		g.whileStatement(x)
	case *repeatStatement:
		g.repeatStatement(x)
	case *compoundStatement:
		g.compoundStatement(x, true, true)
	case *caseStatement:
		g.caseStatement(x)
	case *labeled:
		g.labeled(x)
	case *gotoStatement:
		g.gotoStatement(x)
	default:
		g.w("panic(`TODO %v: %T`)", n.Position(), x)
	}
}

// Go processes 'web' and writes the Pascal code to 'pascal', the string pool
// to 'pool' and the Go code in package 'pkg' to 'dest'.  To apply a change
// file, pass knuth.NewChanger(web, changes) as 'web'.
func Go(dest, pascal, pool io.Writer, web knuth.RuneSource, pkg string, options ...Option) (err error) {
	var pascal0, go0 bytes.Buffer
	if err = Tangle(&pascal0, pool, web); err != nil {
		return fmt.Errorf("tangle: %v", err)
	}

	if pascal != nil {
		if _, err := pascal.Write(pascal0.Bytes()); err != nil {
			return fmt.Errorf("writing Pascal: %v", err)
		}
	}

	opts := &genOpts{}

	for _, v := range options {
		if err = v(opts); err != nil {
			return err
		}
	}

	pascalBytes, err := opts.fillPositionalArgs(pascal0.Bytes())
	if err != nil {
		return err
	}

	s := newPasScanner("web.pas", pascalBytes, opts)
	ast, err := pasParse(s)
	if err != nil {
		return fmt.Errorf("parsing Pascal: %v", err)
	}

	if err := ast.check(); err != nil {
		return fmt.Errorf("type checking Pascal: %v", err)
		return nil
	}

	g, err := newGen(&go0, ast, pkg, fmt.Sprint(os.Args), opts)
	if err != nil {
		return fmt.Errorf("generating Go: %v", err)
	}

	if err = g.gen(); err != nil {
		return fmt.Errorf("generating Go: %v", err)
	}

	gofmt, err := exec.LookPath("gofmt")
	if err != nil {
		return fmt.Errorf("searching for gofmt: %v", err)
	}

	tmp, err := ioutil.TempFile("", "web2go-")
	if err != nil {
		return fmt.Errorf("creating temporary file: %v", err)
	}

	fn := tmp.Name()
	if _, err := tmp.Write(go0.Bytes()); err != nil {
		return fmt.Errorf("writing temporary file: %v", err)
	}

	if err = tmp.Close(); err != nil {
		return fmt.Errorf("closing temporary file: %v", err)
	}

	if b, err := exec.Command(gofmt, "-w", "-s", "-r", "(x) -> x", fn).CombinedOutput(); err != nil {
		dest.Write(go0.Bytes())
		return fmt.Errorf("executing gofmt: output: %s\nerror: %v", b, err)
	}

	b, err := os.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("reading temporary file: %v", err)
	}

	if _, err := dest.Write(b); err != nil {
		return fmt.Errorf("writing resulting Go: %v", err)
	}

	return nil
}
