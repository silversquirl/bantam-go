//go:generate stringer -type TokenType

package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type Token struct {
	Ty TokenType
	S  string
}

type TokenType int

const (
	TokEOF TokenType = iota
	TokLeftParen
	TokRightParen
	TokComma
	TokAssign
	TokPlus
	TokMinus
	TokAsterisk
	TokSlash
	TokCaret
	TokTilde
	TokBang
	TokQuestion
	TokColon
	TokName
)

func (ty TokenType) Rune() rune {
	switch ty {
	case TokLeftParen:
		return '('
	case TokRightParen:
		return ')'
	case TokComma:
		return ','
	case TokAssign:
		return '='
	case TokPlus:
		return '+'
	case TokMinus:
		return '-'
	case TokAsterisk:
		return '*'
	case TokSlash:
		return '/'
	case TokCaret:
		return '^'
	case TokTilde:
		return '~'
	case TokBang:
		return '!'
	case TokQuestion:
		return '?'
	case TokColon:
		return ':'
	default:
		return 0
	}
}

var punct = map[rune]TokenType{}

func init() {
	for ty := TokLeftParen; ty < TokName; ty++ {
		punct[ty.Rune()] = ty
	}
}

func Tokenize(r io.Reader, c chan<- Token) {
	br := bufio.NewReader(r)
	name := &strings.Builder{}
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			break
		}

		if unicode.IsLetter(r) {
			name.WriteRune(r)
		} else if name.Len() > 0 {
			c <- Token{TokName, name.String()}
			name.Reset()
		}
		if ty, ok := punct[r]; ok {
			c <- Token{ty, string(r)}
		}
	}
	if name.Len() > 0 {
		c <- Token{TokName, name.String()}
	}
	close(c)
}
