package main

import (
	"errors"
	"io"
)

func init() {
	prefixParselets = map[TokenType]prefixParselet{
		TokLeftParen: {0, plGroup},
		TokName:      {0, plName},

		TokPlus:  {PPrefix, plPrefix},
		TokMinus: {PPrefix, plPrefix},
		TokTilde: {PPrefix, plPrefix},
		TokBang:  {PPrefix, plPrefix},
	}

	infixParselets = map[TokenType]infixParselet{
		TokBang:      {PPostfix, plPostfix},
		TokPlus:      {PSum, plBinary},
		TokMinus:     {PSum, plBinary},
		TokAsterisk:  {PMul, plBinary},
		TokSlash:     {PMul, plBinary},
		TokCaret:     {PExp, plBinary},
		TokQuestion:  {PCond, plCond},
		TokLeftParen: {PCall, plCall},
		TokAssign:    {PAssign, plAssign},
	}
}

const (
	PAssign = iota + 1
	PCond
	PSum
	PMul
	PExp
	PPrefix
	PPostfix
	PCall
)

type prefixParselet struct {
	Prec int
	Func func(int, *parser, Token) Expr
}

func (pl prefixParselet) Call(p *parser, tok Token) Expr {
	return pl.Func(pl.Prec, p, tok)
}

type infixParselet struct {
	Prec int
	Func func(int, *parser, Token, Expr) Expr
}

func (pl infixParselet) Call(p *parser, tok Token, left Expr) Expr {
	return pl.Func(pl.Prec, p, tok, left)
}

var prefixParselets map[TokenType]prefixParselet
var infixParselets map[TokenType]infixParselet

func Parse(r io.Reader) (expr Expr, err error) {
	toks := make(chan Token)
	go Tokenize(r, toks)
	p := parser{<-toks, toks}
	defer func() {
		switch e := recover().(type) {
		case string:
			err = errors.New(e)
		case error:
			err = e
		}
	}()
	expr = p.parseExpr(0)
	p.require(TokEOF)
	return
}

type parser struct {
	tok  Token
	toks <-chan Token
}

func (p *parser) next() Token {
	tok := p.tok
	p.tok = <-p.toks // Zero value of Token is EOF, so this works nicely with closed channels
	return tok
}
func (p *parser) peek() Token {
	return p.tok
}
func (p *parser) require(ty TokenType) Token {
	if tok := p.next(); tok.Ty == ty {
		return tok
	} else {
		panic("Expected " + ty.String() + " got " + tok.Ty.String())
	}
}

func (p *parser) parseExpr(prec int) Expr {
	tok := p.next()
	if tok.Ty == TokEOF {
		panic("Unexpected EOF")
	}
	prefixPl, ok := prefixParselets[tok.Ty]
	if !ok {
		panic(tok.Ty.String() + " is not a prefix operator")
	}
	left := prefixPl.Call(p, tok)

	for {
		tok = p.peek()
		infixPl := infixParselets[tok.Ty]
		if infixPl.Prec <= prec {
			return left
		}
		p.next()
		left = infixPl.Call(p, tok, left)
	}
	return left
}

func plGroup(_ int, p *parser, tok Token) Expr {
	expr := p.parseExpr(0)
	p.require(TokRightParen)
	return expr
}
func plName(_ int, p *parser, tok Token) Expr {
	return NameExpr(tok.S)
}
func plPrefix(prec int, p *parser, tok Token) Expr {
	return PrefixExpr{tok.Ty, p.parseExpr(prec)}
}
func plPostfix(_ int, p *parser, tok Token, left Expr) Expr {
	return PostfixExpr{tok.Ty, left}
}
func plBinary(prec int, p *parser, tok Token, left Expr) Expr {
	return OperatorExpr{tok.Ty, left, p.parseExpr(prec)}
}
func plCond(prec int, p *parser, tok Token, left Expr) Expr {
	then := p.parseExpr(prec)
	p.require(TokColon)
	else_ := p.parseExpr(prec)
	return CondExpr{left, then, else_}
}
func plCall(_ int, p *parser, tok Token, fun Expr) Expr {
	call := CallExpr{fun, nil}
	for {
		expr := p.parseExpr(0)
		call.Args = append(call.Args, expr)
		if p.peek().Ty == TokRightParen {
			p.next()
			break
		} else {
			p.require(TokComma)
		}
	}
	return call
}
func plAssign(prec int, p *parser, tok Token, left Expr) Expr {
	// prec-1 makes this right associative
	return AssignExpr{left.(NameExpr), p.parseExpr(prec - 1)}
}
