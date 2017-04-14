## gomacro - A Go interpreter with Lisp-like macros

The package `fast_interpreter` contains a faster reimplementation of gomacro core interpreter.

To learn about gomacro, download, compile and use it, please refer to the original implementation [README.md](../README.md)

If you want to help with the reimplementation, or you are simply curious, read on :)

## Current Status

ALPHA.

The fast intepreter supports:
* parsing, including parsing macro-related syntax - shared with the classic interpreter
* iota and untyped constants
* binary expressions on untyped constants, booleans, integers, floats, complex numbers, and strings
* unary operators + - ^ ! <- (other unary operators: deref * is unimplemented, address-of & is implemented only for variables)
* constant, variable and type declarations
* assignment to variables, i.e. 'variable = constant' and 'variable = expression'
* if
* ~quote

Everything else is still missing. You are welcome to contribute.

Limitations:
* no distinction between named and unnamed types created by interpreted code.
  For the interpreter, `struct { A, B int }` and `type Pair struct { A, B int }`
  are exactly the same type. This has subtle consequences, including the risk
  that two different packages define the same type and overwrite each other's methods.

  The reason for such limitation is simple: the interpreter uses `reflect.StructOf()`
  to define new types, which can only create unnamed types.
  The interpreter then defines named types as aliases for the underlying unnamed types.
* operators << and >> do not follow type deduction rules for untyped constants.
  The implemented behavior is:
  * an untyped constant shifted by a non-constant expression always returns an int
  * an untyped floating point constant shifted by a constant expression returns an untyped integer constant.
    the interpreter signals an error during the precompile phase
    if the left operand has a non-zero fractional or imaginary part,
    or it overflows both int64 and uint64.
  See [Go Language Specification](https://golang.org/ref/spec#Operators) for the correct behavior
