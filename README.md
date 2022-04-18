# I made a Lisp!

And you can too!

I caught a cold one weekend so I couldn't leave the house. I used that time to
build a Lisp interpreter! This interpreter implements lambdas, mutable
variables, mutable pairs, closures, macros, and the quote and quasiquote reader
macros. Please see the [examples folder](examples/) or
[`prelude.lisp`](prelude.lisp) for demonstrations of this Lisp's features.

Below, I detail the pieces that go into creating a Lisp, with simplified code
samples for various parts of the interpreter. I also detail which
[resources](#resources) I used while making this Lisp.

## The Pieces

Writing an interpreter consists of three parts:

* A reader that takes source code as input and produces an AST (abstract syntax
tree) as output.
* An evaluator that takes an AST as input and produces an expression as output,
computing side-effects as well.
* A printer that takes an expression as input and prints a human-readable
representation.

These three parts are present in every interpreter that I'm aware of. Composing
these three functions in a loop is known as a REPL, or Read-Eval-Print Loop. The
Go code might look something like the following:

```go
for {
    Print(Eval(Read()))
}
```

While the printing step is pretty trivial, I'll cover every other step in this
section.

### Parser

Here's an example of a Lisp expression. This expression evaluates to the
number 30.

```lisp
(multiply 10 (add 1 2))
```

In this expression, `multiply` and `add` are symbols, `10`, `1`, and `2` are
numbers, and everything between a pair of parentheses is a list.

Lisp's syntax makes it trivially easy to parse. The name Lisp (sometimes
stylized LISP) is an abbreviation for "LISt Processor", so every expression is
written as a list. As such, the parsing strategy is trivial.

Our `Read` function is very straightforward. It looks something like this:

```go
func Read() Expr {
    ConsumeWhitespace()
    if e := ReadList(); e != nil {
        return e
    }
    if e := ReadNumber(); e != nil {
        return e
    }
    if e := ReadSymbol(); e != nil {
        return e
    }
    panic("failed to parse")
}
```

Reading a number or symbol is trivial, and the list parsing function looks
something like this:

``` go
// peekRune gets the next character without consuming it
// readRune gets the next character and consumes it
func ReadList() Expr {
    if peekRune() != '(' {
        return nil
    }
    readRune()
    result := makeList()
    for {
        if readRune() == ')' {
            return result
        }
        result.Append(Read())
    }
}
```

This is basically all that you need to parse Lisp. The `Read` function reads any
expression, calling functions like `ReadList`, `ReadSymbol`, and `ReadNumber`,
and the `ReadList` expression calls `Read` for each element of the list.

This parsing strategy is called a "recursive descent" parser, and it works
because the parser doesn't need to backtrack. Functions like `peekRune` and
`readRune` can operate on a stream with a fixed input buffer, since you don't
need to peek an arbitrary number of characters. This parsing strategy allows for
a very simple parser to be constructed without using any sort of library.

### Evaluator

You'll notice that the return type of `Read` was an `Expr`. This is because, in
Lisp, there is no distinction between an AST and data. Since the AST is written
as a list, we can internally treat it just like Lisp's list data structure. This
will allow us to do some cool things like macros (explained below).

Our evaluator function take two parameters, and expression and an environment,
and will return an expression.

An environment is the "scope" of our evaluation. The environment is simply a
mapping between variable names and values, and will be used to look up
variables.

Our `Eval` function looks something like this:

```go
func Eval(expr Expr, env *Env) Expr {
    switch expr := expr.(type) {
    case *Number:
        // return primitive data as itself
        return expr
    case *Symbol:
        // look up variable and return its value
        return env.Resolve(expr)
    case *List:
        // apply the function to its arguments
        return Apply(Eval(expr.First(), env), expr.Rest(), env)
    default:
        panic("unknown expression type evaluated")
    }
}
```

So we have three cases:

1. If it's primitive data, it evaluates to itself.
2. If it's a variable, look up the corresponding value. In our language, we
represent variables with a `Symbol` type that stores names (like the function
`add` or the number `n`).
3. If it's a list, apply the function to the arguments. The function is either
an anonymous function literal or a variable, so just call `Eval` on it to get
the resulting value.

### Apply

But what does `Apply` do? I lied when I said it only takes functions (which I'll
call prodecures from here on out). It actually takes one of three things:
procedures, macros, and special forms.

Our `Apply` function looks something like this:

```go
func Apply(f Expr, args *List, env *Env) Expr {
    switch f := f.(type) {
    case *Procedure:
        ApplyProcedure(f, args)
    case *Macro:
        ApplyMacro(f, args)
    case SpecialForm:
        f(args, env)
    }
}
```

I'll cover each of the three cases below.

#### Procedures

If you've programmed before, procedures (functions) should be familiar to you. A
procedure is made up of a few things:

* The environment (scope) that it was defined in (i.e. defined at the top level
or defined within some function as a closure)
* Its arguments (as a list)
* Its body (as an expression with the code as an AST)

Our `ApplyProcedure` function looks something like this:

```go
func ApplyProcedure(proc *Procedure, args *List, env *Env) {
    // evaluate the arguments
    evaluatedArgs := args.Map(func (expr Expr) Expr {
        return Eval(expr, env)
    })
    // create a new environment and bind the arguments to their values
    newEnv := makeEnv()
    newEnv.SetParent(proc.ParentEnv)
    for i, argSymbol := range proc.Args {
        newEnv.Bind(argSymbol, evaluatedArgs.Index(i))
    }
    // evaluate the body with the new environment
    return Eval(proc.Body, newEnv)
}
```

We use each of the three parts of a procedure:

* The environment was used as the parent of our new environment. The function
body will have access to the same variables as where it was defined. The new
environment allows the body to define new variables without polluting the
outer environment (we don't want all variables to be global).
* The arguments were used to bind the passed in argument values with the
corresponding variable names.
* The body was evaluated in this new environment.

#### Macros

Macro are very similar to functions but their purpose is different. While
functions take expressions as arguments and return an expression, macros take
ASTs as arguments and return a new AST. This allows you to extend the language
with new syntax.

Here's an example of what a macro might look like:

```lisp
;; Example:
;; The expression (subtract 10 3)
;;   Evaluates to: 7
;; The expression (reverse-args (subtract 10 3))
;;   Expands to: (subtract 3 10)
;;   Evaluates to: -7
(defmacro reverse-args (expr)
  (append ; append two lists
   (list (first expr)) ; the first element is the procedure/macro/special form
   (reverse (rest (expr))))) ; the rest of the elements are the reversed arguments
```

Notice how the argument didn't evaluate to 7 before being passed in, but was
instead passed in as a list. When the resulting list was returned, only then was
it evaluated. This is how macros work, and it allows us to define new syntax
such as
[`let`](https://docs.racket-lang.org/reference/let.html#%28form._%28%28lib._racket%2Fprivate%2Fletstx-scheme..rkt%29._let%29%29)
for variable binding.

Anyways, now that I'm done motivating macros, let's discuss how to implement
them.

Our `ApplyMacro` function looks something like this:

```go
func ApplyMacro(macro *Macro, env *Env) {
    // create a new environment and bind the arguments to their values
    newEnv := makeEnv()
    newEnv.SetParent(proc.ParentEnv)
    for i, argSymbol := range proc.Args {
        newEnv.Bind(argSymbol, args.Index(i))
    }
    // evaluate the body with the new environment
    newAST := Eval(proc.Body, newEnv)
    // finally evaluate the resulting AST in the current environment
    return Eval(newAST, env)
}
```

There are two differences from `ApplyMacro`:

* The arguments are not evaluated before they are bound
* The returned value is evaluate (since its an AST)

We just moved the evaluation from before application to after. While macros are
very simple, they are very powerful since they can be used to essentially create
new syntax for the language. The [n-queens algorithm
example](examples/nqueens.lisp) in this repo provides some examples of macros
used in a practical context, as well as [the prelude](prelude.lisp) where I
define some new syntax using Lisp rather than Go.

#### Special Forms

Finally, special forms are simply special syntax implemented in Go. Because of
this, they can simply be called as functions. We can define special forms like
so:

```go
type SpecialForm func (*List, *Env) Expr
```

Here's an example of a define special form:

```go
// type checking omitted and error handling omitted
func Define(args *List, env *Env) Expr {
	name := args.Index(0).(*Symbol)
	value := Eval(args.Index(1), env)
	env.Bind(name, value)
	return value
}
```

Which can be used like so:

```scheme
(define x 3)
;; => 3
(define y (+ 1 (* x x)))
;; => 10
(* 2 (+ x y))
;; => 26
```

By implementing special forms in Go, we get access to the internals of our
evaluator, for example by allowing `Define` to bind new values in the
environment. A Lisp programmer never sees the environment and likely never even
thinks about it if they already expect lexical scoping from their programming
languages.

## Resources

Here is a collection of resources that I used while making this Lisp:

* [Structure and Interpretation of Computer
  Programs](https://mitpress.mit.edu/sites/default/files/sicp/index.html): This
  book is where I was first introduced to Lisp and where I learned about its
  syntax, the Environment model, and how a Lisp interpreter is put together.
* [MiniLisp](https://github.com/rui314/minilisp) by Rui Ueyama. A Lisp
  interpreter in less than 1000 lines of C. This is the first Lisp interpreter
  that I read the source code for.
* [SectorLisp](https://justine.lol/sectorlisp2/) by Justine Tunney. A Lisp
  interpreter that assembles to less than 512 bytes. This project broke down
  Lisp into its absolutely smallest necessary pieces and allowed me to
  understand its fundamentals.
* [The R5RS Scheme Standard](https://schemers.org/Documents/Standards/R5RS/)
  which is the most widely implemented standard for Scheme. I used this as a
  reference whenever I was unsure how something should behave. For simplicity's
  sake my implementation doesn't follow this spec, but it did help guide my
  design decisions.
