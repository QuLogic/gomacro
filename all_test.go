/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * all_test.go
 *
 *  Created on: Mar 06 2017
 *      Author: Massimiliano Ghilardi
 */

package main

import (
	"go/ast"
	"go/token"
	r "reflect"
	"testing"

	. "github.com/cosmos72/gomacro/ast2"
	. "github.com/cosmos72/gomacro/base"
	fast "github.com/cosmos72/gomacro/fast_interpreter"
	classic "github.com/cosmos72/gomacro/interpreter"
)

type TestFor int

const (
	I TestFor = 1 << iota
	F
	A TestFor = I | F
	B TestFor = I // temporarily disabled compiler test

)

type TestCase struct {
	testfor TestFor
	name    string
	program string
	result0 interface{}
	results []interface{}
}

func TestFastInterp(t *testing.T) {
	env := classic.New()
	comp := fast.New()
	for _, test := range tests {
		if test.testfor&F != 0 {
			c := test
			t.Run(c.name, func(t *testing.T) { c.compile(t, comp, env) })
		}
	}
}

func TestInterpreter(t *testing.T) {
	env := classic.New()
	// env.Options |= OptDebugCallStack | OptDebugPanicRecover
	for _, test := range tests {
		if test.testfor&I != 0 {
			c := test
			t.Run(c.name, func(t *testing.T) { c.interpret(t, env) })
		}
	}
}

func (c *TestCase) compile(t *testing.T, comp *fast.CompEnv, env *classic.Env) {
	// parse + macroexpansion phase
	form := env.ParseAst(c.program)

	// compile phase
	f := comp.CompileAst(form)
	rets := PackValues(comp.Run(f))

	c.compareResults(t, rets)
}

func (c *TestCase) interpret(t *testing.T, env *classic.Env) {
	// parse + macroexpansion phase
	form := env.ParseAst(c.program)
	// eval phase
	rets := PackValues(env.EvalAst(form))

	c.compareResults(t, rets)
}

const sum_s = "func sum(n int) int { total := 0; for i := 1; i <= n; i++ { total += i }; return total }"
const fib_s = "func fibonacci(n uint) uint { if n <= 2 { return 1 }; return fibonacci(n-1) + fibonacci(n-2) }"

var ti = r.StructOf(
	[]r.StructField{
		r.StructField{Name: ReflectGensymPrefix, Type: r.TypeOf((*interface{})(nil)).Elem()},
		r.StructField{Name: "String", Type: r.TypeOf((*func() string)(nil)).Elem()},
	},
)
var si = r.Zero(ti).Interface()

var zeroValues = []r.Value{}

var tests = []TestCase{
	TestCase{A, "1+1", "1+1", 1 + 1, nil},
	TestCase{I, "1+'A'", "1+'A'", 66, nil},  // interpreter is not accurate in this case... returns <int> instead of <int32>
	TestCase{F, "1+'A'", "1+'A'", 'B', nil}, // fast_interpreter instead *IS* accurate
	TestCase{I, "int8+1", "int8(1)+1", int8(2), nil},
	TestCase{I, "int8_overflow", "int8(64)+64", int8(-128), nil},
	TestCase{I, "interface", "type Stringer interface { String() string }; var s Stringer", si, nil},
	TestCase{A, "string", "\"foobar\"", "foobar", nil},
	TestCase{A, "expr_and", "3 & 6", 3 & 6, nil},
	TestCase{A, "expr_or", "7 | 8", 7 | 8, nil},
	TestCase{A, "expr_xor", "0x1f ^ 0xf1", 0x1f ^ 0xf1, nil},
	TestCase{A, "expr_arith", "((1+2)*3^4|99)%112", ((1+2)*3 ^ 4 | 99) % 112, nil},
	TestCase{A, "expr_shift", "7<<(10>>1)", 7 << (10 >> 1), nil},

	TestCase{A, "complex_1", "7i", 7i, nil},
	TestCase{A, "complex_2", "0.5+1.75i", 0.5 + 1.75i, nil},
	TestCase{A, "complex_3", "1i * 2i", 1i * 2i, nil},
	TestCase{A, "const_1", "const c1 = 11; c1", 11, nil},
	TestCase{A, "const_2", "const c2 = 0xff&555+23/12.2; c2", 0xff&555 + 23/12.2, nil},

	// the classic interpreter is not accurate in this cases... missing exact arithmetic on constants
	TestCase{I, "const_3", "const c3 = 0.1+0.2; c3", float64(0.1) + float64(0.2), nil},
	TestCase{I, "const_4", "const c4 = c3/3; c4", (float64(0.1) + float64(0.2)) / 3, nil},

	// the fast_interpreter instead *IS* accurate, thanks to exact arithmetic on untyped constants
	TestCase{F, "const_3", "const c3 = 0.1+0.2; c3", 0.1 + 0.2, nil},
	TestCase{F, "const_4", "const c4 = c3/3; c4", (0.1 + 0.2) / 3, nil},
	TestCase{F, "untyped_1", "2.0 >> 1", 1, nil},
	TestCase{A, "untyped_2", "1/2", 0, nil},
	TestCase{A, "untyped_unary", "-+^6", -+^6, nil},

	TestCase{A, "iota_1", "const c5 = iota^7; c5", 7, nil},
	TestCase{A, "iota_2", "const ( c6 = iota+6; c7=iota+6 ); c6", 6, nil},
	TestCase{A, "iota_3", "c7", 7, nil},
	TestCase{A, "iota_implicit_1", "const ( c8 uint = iota+8; c9 ); c8", uint(8), nil},
	TestCase{A, "iota_implicit_2", "c9", uint(9), nil},

	TestCase{A, "var_1", "var v1 bool; v1", false, nil},
	TestCase{A, "var_2", "var v2 uint8 = 7; v2", uint8(7), nil},
	TestCase{A, "var_3", "var v3 uint16 = 12; v3", uint16(12), nil},
	TestCase{A, "var_4", "var v uint32 = 99; v", uint32(99), nil},
	TestCase{A, "var_5", "var v5 string; v5", "", nil},
	TestCase{A, "var_6", "var v6 float32; v6", float32(0), nil},
	TestCase{A, "var_7", "var v7 complex64; v7", complex64(0), nil},
	TestCase{A, "var_8", "var err error; err", nil, nil},
	TestCase{A, "var_9", `var ve string = ""; ve`, "", nil},
	TestCase{A, "var_pointer", "var vp *string; vp", (*string)(nil), nil},
	TestCase{A, "var_map", "var vm *map[error]bool; vm", (*map[error]bool)(nil), nil},
	TestCase{A, "var_slice", "var vs []byte; vs", ([]byte)(nil), nil},
	TestCase{A, "var_array", "var va [2][]rune; va", [2][]rune{}, nil},
	TestCase{A, "var_interface_1", "var vi interface{} = 1; vi", 1, nil},
	TestCase{A, "var_interface_2", "var vnil interface{}; vnil", nil, nil},
	TestCase{A, "var_shift_1", "7 << 8", 7 << 8, nil},
	TestCase{A, "var_shift_2", "8 >> 2", 8 >> 2, nil},
	TestCase{A, "var_shift_3", "v2 << 3", uint8(7) << 3, nil},
	TestCase{A, "var_shift_4", "v2 >> 1", uint8(7) >> 1, nil},
	TestCase{A, "var_shift_5", "0xff << v2", 0xff << 7, nil},
	TestCase{A, "var_shift_6", "0x12345678 >> v2", 0x12345678 >> uint8(7), nil},
	TestCase{A, "var_shift_7", "v << v2", uint32(99) << uint8(7), nil},
	TestCase{A, "var_shift_8", "v3 << v3 >> v2", uint16(12) << 12 >> uint8(7), nil},
	TestCase{A, "var_shift_9", "v3 << 0", uint16(12), nil},
	TestCase{A, "var_shift_overflow", "v3 << 13", uint16(32768), nil},
	TestCase{A, "eql_nil_1", "err == nil", true, nil},
	TestCase{A, "eql_nil_2", "vp == nil", true, nil},
	TestCase{A, "eql_nil_3", "vm == nil", true, nil},
	TestCase{A, "eql_nil_4", "vs == nil", true, nil},
	TestCase{A, "eql_nil_5", "vi == nil", false, nil},
	TestCase{A, "eql_nil_6", "vnil == nil", true, nil},
	TestCase{A, "eql_halfnil", "var vhalfnil interface{} = vm; vhalfnil == nil", false, nil},
	TestCase{A, "eql_interface", "vi == 1", true, nil},
	TestCase{A, "typed_unary_1", "!!!v1", true, nil},
	TestCase{A, "typed_unary_2", "+-^v2", uint8(8), nil},
	TestCase{A, "typed_unary_3", "+^-v3", uint16(11), nil},
	TestCase{A, "typed_unary_4", "v7 = 2.5i; -v7", complex64(-2.5i), nil},

	TestCase{A, "type_int8", "type t8 int8; var v8 t8; v8", int8(0), nil},
	TestCase{A, "type_complicated", "type tfff func(int,int) func(error, func(bool)) string; var vfff tfff; vfff", (func(int, int) func(error, func(bool)) string)(nil), nil},
	TestCase{A, "type_struct", "type Pair struct { A, B int }; var pair Pair; pair", struct{ A, B int }{}, nil},
	TestCase{I, "struct", "pair.A, pair.B = 1, 2; pair", struct{ A, B int }{1, 2}, nil},
	TestCase{I, "pointer", "var p = 1.25; if *&p != p { p = -1 }; p", 1.25, nil},
	TestCase{I, "defer_1", "v = 0; func testdefer(x uint32) { if x != 0 { defer func() { v = x }() } }; testdefer(29); v", uint32(29), nil},
	TestCase{I, "defer_2", "v = 12; testdefer(0); v", uint32(12), nil},
	TestCase{I, "make_chan", "cx := make(chan interface{}, 2)", make(chan interface{}, 2), nil},
	TestCase{I, "make_map", "m := make(map[rune]bool)", make(map[rune]bool), nil},
	TestCase{I, "make_slice", "y := make([]uint8, 7); y[0] = 100; y[3] = 103; y", []uint8{100, 0, 0, 103, 0, 0, 0}, nil},
	TestCase{I, "expr_slice", "y = y[:4]", []uint8{100, 0, 0, 103}, nil},
	TestCase{I, "expr_slice3", "y = y[:3:4]", []uint8{100, 0, 0}, nil},

	TestCase{A, "set_const_1", "v1 = true;    v1", true, nil},
	TestCase{A, "set_const_2", "v2 = 9;       v2", uint8(9), nil},
	TestCase{A, "set_const_3", "v3 = 60000;   v3", uint16(60000), nil},
	TestCase{A, "set_const_4", "v  = 987;      v", uint32(987), nil},
	TestCase{A, "set_const_5", `v5 = "8y57r"; v5`, "8y57r", nil},
	TestCase{A, "set_const_6", "v6 = 0.12345678901234; v6", float32(0.12345678901234), nil},        // v6 is declared float32
	TestCase{A, "set_const_7", "v7 = 0.98765432109i; v7", complex(0, float32(0.98765432109)), nil}, // v7 is declared complex64

	TestCase{A, "set_expr_1", "v1 = v1 == v1;    v1", true, nil},
	TestCase{A, "set_expr_2", "v2 = v2 - 7;      v2", uint8(2), nil},
	TestCase{A, "set_expr_3", "v3 = v3 % 7;      v3", uint16(60000) % 7, nil},
	TestCase{A, "set_expr_4", "v  = v * 10;      v", uint32(9870), nil},
	TestCase{A, "set_expr_5", `v5 = v5 + "iuh";  v5`, "8y57riuh", nil},
	TestCase{A, "set_expr_6", "v6 = 1/v6;        v6", 1 / float32(0.12345678901234), nil},                              // v6 is declared float32
	TestCase{A, "set_expr_7", "v7 = v7 * v7;     v7", complex(-float32(0.98765432109)*float32(0.98765432109), 0), nil}, // v7 is declared complex64

	TestCase{A, "add_2", "v2 += 255;    v2", uint8(1), nil}, // overflow
	TestCase{A, "add_3", "v3 += 536;    v3", uint16(60000)%7 + 536, nil},
	TestCase{A, "add_4", "v  += 111;     v", uint32(9870 + 111), nil},
	TestCase{A, "add_5", `v5 += "@#$";  v5`, "8y57riuh@#$", nil},
	TestCase{A, "add_6", "v6 += 0.975319; v6", 1/float32(0.12345678901234) + float32(0.975319), nil}, // v6 is declared float32
	TestCase{A, "add_7", "v7 = 1; v7 += 0.999999i; v7", complex(float32(1), float32(0.999999)), nil}, // v7 is declared complex64

	TestCase{A, "if_1", "if v2 < 1 { v2 = v2-1 } else { v2 = v2+1 }; v2", uint8(2), nil},
	TestCase{A, "if_2", "if v2 < 5 { v2 = v2+2 } else { v2 = v2-2 }; v2", uint8(4), nil},

	TestCase{A, "for_1", "var i, j, k int; for i=1; i<=2; i=i+1 { if i<2 {j=i} else {k=i} }; i", 3, nil},
	TestCase{A, "for_2", "j", 1, nil},
	TestCase{A, "for_3", "k", 2, nil},

	TestCase{A, "continue_1", "j=0; k=0; for i=1; i<=7; i=i+1 { if i==3 {j=i; continue}; k=k+i }; j", 3, nil},
	TestCase{A, "continue_2", "k", 25, nil},
	TestCase{A, "continue_3", "j=0; k=0; for i=1; i<=7; i=i+1 { var ii = i; if ii==3 {j=ii; continue}; k=k+ii }; j", 3, nil},
	TestCase{A, "continue_4", "k", 25, nil},

	TestCase{I, "for_range_chan", "i := 0; c := make(chan int, 2); c <- 1; c <- 2; close(c); for e := range c { i += e }; i", 3, nil},
	TestCase{I, "function", "func ident(x uint) uint { return x }; ident(42)", uint(42), nil},
	TestCase{I, "function_variadic", "func list_args(args ...interface{}) []interface{} { args }; list_args('x', 'y', 'z')", []interface{}{'x', 'y', 'z'}, nil},
	TestCase{I, "fibonacci", fib_s + "; fibonacci(13)", uint(233), nil},
	TestCase{I, "import", "import \"fmt\"", "fmt", nil},
	TestCase{I, "literal_struct", "Pair{A: 73, B: 94}", struct{ A, B int }{A: 73, B: 94}, nil},
	TestCase{I, "literal_array", "[3]int{1,2:3}", [3]int{1, 0, 3}, nil},
	TestCase{I, "literal_map", "map[int]string{1: \"foo\", 2: \"bar\"}", map[int]string{1: "foo", 2: "bar"}, nil},
	TestCase{I, "literal_slice", "[]rune{'a','b','c'}", []rune{'a', 'b', 'c'}, nil},
	TestCase{I, "method_on_ptr", "func (p *Pair) SetLhs(a int) { p.A = a }; pair.SetLhs(8); pair.A", 8, nil},
	TestCase{I, "method_on_value", "func (p Pair) SetLhs(a int) { p.A = a }; pair.SetLhs(11); pair.A", 8, nil}, // method on value gets a copy of the receiver - changes to not propagate
	TestCase{I, "multiple_values_1", "func twins(x float32) (float32,float32) { return x, x+1 }; twins(17.0)", nil, []interface{}{float32(17.0), float32(18.0)}},
	TestCase{I, "multiple_values_2", "func twins2(x float32) (float32,float32) { return twins(x) }; twins2(19.0)", nil, []interface{}{float32(19.0), float32(20.0)}},
	TestCase{A, "pred_bool_1", "false==false && true==true && true!=false", true, nil},
	TestCase{A, "pred_bool_2", "false!=false || true!=true || true==false", false, nil},
	TestCase{A, "pred_int", "1==1 && 1<=1 && 1>=1 && 1!=2 && 1<2 && 2>1 || 0==1", true, nil},
	TestCase{A, "pred_string_1", `""=="" && "">="" && ""<="" && ""<"a" && ""<="b" && "a">"" && "b">=""`, true, nil},
	TestCase{A, "pred_string_2", `ve=="" && ve>="" && ve<="" && ve<"a" && ve<="b" && "a">ve && "b">=ve`, true, nil},
	TestCase{A, "pred_string_3", `"x"=="x" && "x"<="x" && "x">="x" && "x"!="y" && "x"<"y" && "y">"x"`, true, nil},
	TestCase{A, "pred_string_4", `"x"!="x" || "y"!="y" || "x">="y" || "y"<="x"`, false, nil},
	TestCase{I, "recover", `var vpanic interface{}
		func test_recover(rec bool, panick interface{}) {
			defer func() {
				if rec {
					vpanic = recover()
				}
			}()
			panic(panick)
		}
		test_recover(true, -3)
		vpanic`, -3, nil},
	TestCase{I, "recover_nested_1", `var vpanic2, vpanic3 interface{}
		func test_nested_recover(repanic bool, panick interface{}) {
			defer func() {
				vpanic = recover()
			}()
			defer func() {
				func() {
					vpanic3 = recover()
				}()
				vpanic2 = recover()
				if repanic {
					panic(vpanic2)
				}
			}()
			panic(panick)
		}
		test_nested_recover(false, -4)
		Values(vpanic, vpanic2, vpanic3)
		`, nil, []interface{}{nil, -4, nil}},
	TestCase{I, "recover_nested_2", `vpanic, vpanic2, vpanic3 = nil, nil, nil
		test_nested_recover(true, -5)
		Values(vpanic, vpanic2, vpanic3)
		`, nil, []interface{}{-5, -5, nil}},
	TestCase{I, "send_recv", "cx <- \"x\"; <-cx", nil, []interface{}{"x", true}},
	TestCase{I, "sum", sum_s + "; sum(100)", 5050, nil},

	TestCase{I, "select_1", "cx <- 1; { var x interface{}; select { case x=<-cx: x; default: } }", 1, nil},
	TestCase{I, "select_2", "cx <- m; select { case x:=<-cx: x; default: }", make(map[rune]bool), nil},
	TestCase{I, "select_3", "select { case cx<-1: 1; default: 0 }", 1, nil},
	TestCase{I, "select_4", "select { case cx<-2: 2; default: 0 }", 2, nil},
	TestCase{I, "select_5", "select { case cx<-3: 3; default: 0 }", 0, nil},
	TestCase{I, "select_6", "select { case cx<-4: 4; case x:=<-cx: x; default: 0 }", 1, nil},

	TestCase{I, "switch_1", "switch { case false: 0; default: 1 }", 1, nil},
	TestCase{I, "switch_2", "switch v:=20; v { case 20: '@' }", '@', nil},
	TestCase{I, "switch_fallthrough", "switch 0 { default: fallthrough; case 1: 10; fallthrough; case 2: 20 }", 20, nil},

	TestCase{I, "typeswitch_1", "var x interface{} = \"abc\"; switch y := x.(type) { default: 0; case string: 1 }", 1, nil},
	TestCase{I, "typeswitch_2", "switch x.(type) { default: 0; case interface{}: 2 }", 2, nil},
	TestCase{I, "typeswitch_3", "switch x.(type) { default: 0; case int: 3 }", 0, nil},
	TestCase{I, "typeswitch_4", "switch nil.(type) { default: 0; case nil: 4 }", 4, nil},

	TestCase{A, "quote_1", `~quote{7}`, &ast.BasicLit{Kind: token.INT, Value: "7"}, nil},
	TestCase{A, "quote_2", `~quote{x}`, &ast.Ident{Name: "x"}, nil},
	TestCase{A, "quote_3", `var ab = ~quote{a;b}; ab`, &ast.BlockStmt{List: []ast.Stmt{
		&ast.ExprStmt{X: &ast.Ident{Name: "a"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "b"}},
	}}, nil},
	TestCase{A, "quote_4", `~'{"foo"+"bar"}`, &ast.BinaryExpr{
		Op: token.ADD,
		X:  &ast.BasicLit{Kind: token.STRING, Value: `"foo"`},
		Y:  &ast.BasicLit{Kind: token.STRING, Value: `"bar"`},
	}, nil},
	TestCase{I, "quasiquote_1", `~quasiquote{1 + ~unquote{2+3}}`, &ast.BinaryExpr{
		Op: token.ADD,
		X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
		Y:  &ast.BasicLit{Kind: token.INT, Value: "5"},
	}, nil},
	TestCase{I, "quasiquote_2", `~"{2 * ~,{3<<1}}`, &ast.BinaryExpr{
		Op: token.MUL,
		X:  &ast.BasicLit{Kind: token.INT, Value: "2"},
		Y:  &ast.BasicLit{Kind: token.INT, Value: "6"},
	}, nil},
	TestCase{I, "unquote_splice_1", `~quasiquote{~unquote_splice ab ; c}`, &ast.BlockStmt{List: []ast.Stmt{
		&ast.ExprStmt{X: &ast.Ident{Name: "a"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "b"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "c"}},
	}}, nil},
	TestCase{I, "unquote_splice_2", `~"{zero ; ~,@ab ; one}`, &ast.BlockStmt{List: []ast.Stmt{
		&ast.ExprStmt{X: &ast.Ident{Name: "zero"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "a"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "b"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "one"}},
	}}, nil},
	TestCase{I, "macro", "~macro second_arg(a,b,c interface{}) interface{} { return b }; 0", 0, nil},
	TestCase{I, "macro_call", "v = 98; second_arg;1;v;3", uint32(98), nil},
	TestCase{I, "macro_nested", "second_arg;1;{second_arg;2;3;4};5", 3, nil},
	TestCase{I, "values", "Values(3,4,5)", nil, []interface{}{3, 4, 5}},
	TestCase{I, "eval", "Eval(Values(3,4,5))", 3, nil},
	TestCase{I, "eval_quote", "Eval(~quote{Values(3,4,5)})", nil, []interface{}{3, 4, 5}},
}

func (c *TestCase) compareResults(t *testing.T, actual []r.Value) {
	expected := c.results
	if len(expected) == 0 {
		expected = []interface{}{c.result0}
	}
	if len(actual) != len(expected) {
		c.fail(t, actual, expected)
		return
	}
	for i := range actual {
		c.compareResult(t, actual[i], expected[i])
	}
}

func (c *TestCase) compareResult(t *testing.T, actualv r.Value, expected interface{}) {
	if actualv == Nil || actualv == None {
		if expected != nil {
			c.fail(t, nil, expected)
		}
		return
	}
	actual := actualv.Interface()
	if !r.DeepEqual(actual, expected) {
		if r.TypeOf(actual) == r.TypeOf(expected) {
			if actualNode, ok := actual.(ast.Node); ok {
				if expectedNode, ok := expected.(ast.Node); ok {
					c.compareAst(t, ToAst(actualNode), ToAst(expectedNode))
					return
				}
			} else if actualv.Kind() == r.Chan {
				// for channels just check the type, length and capacity
				expectedv := r.ValueOf(expected)
				if actualv.Len() == expectedv.Len() && actualv.Cap() == expectedv.Cap() {
					return
				}
			}
		}
		c.fail(t, actual, expected)
	}
}

func (c *TestCase) compareAst(t *testing.T, actual Ast, expected Ast) {
	if r.TypeOf(actual) == r.TypeOf(expected) {
		switch actual := actual.(type) {
		case BadDecl, BadExpr, BadStmt:
			return
		case Ident:
			if actual.X.Name == expected.(Ident).X.Name {
				return
			}
		case BasicLit:
			actualp := actual.X
			expectedp := expected.(BasicLit).X
			if actualp.Kind == expectedp.Kind && actualp.Value == expectedp.Value {
				return
			}
		default:
			na := actual.Size()
			ne := expected.Size()
			if actual.Op() == expected.Op() && na == ne {
				for i := 0; i < na; i++ {
					c.compareAst(t, actual.Get(i), expected.Get(i))
				}
				return
			}
		}
	}
	c.fail(t, actual, expected)
}

func (c *TestCase) fail(t *testing.T, actual interface{}, expected interface{}) {
	t.Errorf("expected %v <%T>, found %v <%T>\n", expected, expected, actual, actual)
}