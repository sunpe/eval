# eval

This is an expression parser that can evaluate expression. It is simple in design, implemented using go/ast, and does not need other dependencies.

## Install

```
go get -u "github.com/sunpe/eval"
```

## Usage

```
expr := "a > b"
data := map[string]any{"a": 3, "b":2}
res, err := Eval(expr, data)
fmt.Println(res)
fmt.Println(err)
```