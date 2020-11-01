package main

import (
	"fmt"
	"io"
)

type Expr interface {
	printTo(w astWriter)
}

type astWriter interface {
	io.Writer
	WriteRune(rune) (int, error)
	WriteString(string) (int, error)
}

func (e AssignExpr) printTo(w astWriter) {
	w.WriteRune('(')
	e.Name.printTo(w)
	w.WriteString(" = ")
	e.Value.printTo(w)
	w.WriteRune(')')
}

func (e CallExpr) printTo(w astWriter) {
	e.Func.printTo(w)
	w.WriteRune('(')
	for i, arg := range e.Args {
		if i > 0 {
			w.WriteString(", ")
		}
		arg.printTo(w)
	}
	w.WriteRune(')')
}

func (e CondExpr) printTo(w astWriter) {
	w.WriteRune('(')
	e.Cond.printTo(w)
	w.WriteString(" ? ")
	e.Then.printTo(w)
	w.WriteString(" : ")
	e.Else.printTo(w)
	w.WriteRune(')')
}

func (e NameExpr) printTo(w astWriter) {
	w.WriteString(string(e))
}

func (e OperatorExpr) printTo(w astWriter) {
	w.WriteRune('(')
	e.L.printTo(w)
	fmt.Fprintf(w, " %c ", e.Op.Rune())
	e.R.printTo(w)
	w.WriteRune(')')
}

func (e PostfixExpr) printTo(w astWriter) {
	w.WriteRune('(')
	e.L.printTo(w)
	w.WriteRune(e.Op.Rune())
	w.WriteRune(')')
}

func (e PrefixExpr) printTo(w astWriter) {
	w.WriteRune('(')
	w.WriteRune(e.Op.Rune())
	e.R.printTo(w)
	w.WriteRune(')')
}

type AssignExpr struct {
	Name  NameExpr
	Value Expr
}

type CallExpr struct {
	Func Expr
	Args []Expr
}

type CondExpr struct {
	Cond, Then, Else Expr
}

type NameExpr string

type OperatorExpr struct {
	Op   TokenType
	L, R Expr
}

type PostfixExpr struct {
	Op TokenType
	L  Expr
}

type PrefixExpr struct {
	Op TokenType
	R  Expr
}
