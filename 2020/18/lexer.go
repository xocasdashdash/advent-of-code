package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

//Position as
type Position struct {
	Line int
	Char int
}

//TokenType The type we're gonna be using instead of an int
type TokenType int

const (
	//TokenError a
	TokenError TokenType = iota
	//TokenEOF Emitted when finished parsing
	TokenEOF
	//TokenNumber Emitted when parsed a number
	TokenNumber // 12345
	//TokenLParenthesis Emitted when found a (
	TokenLParenthesis // (
	//TokenRParenthesis Emitted when found a )
	TokenRParenthesis // )
	//TokenMultiplication Emitted when found a *
	TokenMultiplication // *
	//TokenAddition Emitted when found a +
	TokenAddition //+
)

func (tt TokenType) String() string {
	switch tt {
	case TokenMultiplication:
		return "*"
	case TokenAddition:
		return "+"
	case TokenLParenthesis:
		return "("
	case TokenRParenthesis:
		return ")"
	default:
		return fmt.Sprintf("<token %d >", int(tt))
	}
}

// Token defines a single token which can be obtained via the Scanner
type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}
type lexer struct {
	input string // the string being lexed

	pos    int        // the current position of the input
	start  int        // the start of the current token
	width  int        // the width of the last read rune
	line   int        // the line number of the current token
	char   int        // the character number of the current token
	tokens chan Token // channel on which to emit tokens
}

func newLexer(input string) *lexer {
	return &lexer{
		input:  input,
		pos:    0,
		start:  0,
		width:  0,
		line:   1,
		char:   1,
		tokens: make(chan Token),
	}
}

//Lex Lexes an input and returns a channel where tokens are going to be emitted
func Lex(input string) <-chan Token {
	l := newLexer(input)
	go func() {
		defer close(l.tokens)
		for state := lexToken; state != nil; {
			state = state(l)
		}
	}()
	return l.tokens
}

type stateFn func(l *lexer) stateFn

const eof = -1

func (l *lexer) emit(t TokenType) {
	value := l.current()
	l.tokens <- Token{
		Pos:   l.position(),
		Type:  t,
		Value: value,
	}
	l.updatePosCounters()
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.updatePosCounters()
}

func (l *lexer) updatePosCounters() {
	value := l.current()
	// Update position counters
	l.start = l.pos

	// Count lines
	lastLine := 0
	for {
		i := strings.IndexRune(value[lastLine:], '\n')
		if i == -1 {
			break
		}
		lastLine += i + 1
		l.line++
		l.char = 1
	}
	l.char += len(value) - lastLine
}

func (l *lexer) position() Position {
	return Position{
		Line: l.line,
		Char: l.char,
	}
}
func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

//Backup the lexer to the previous rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// error emits an error token with the err and returns the terminal state.
func (l *lexer) error(err error) stateFn {
	l.tokens <- Token{Pos: l.position(), Type: TokenError, Value: err.Error()}
	return nil
}

// errorf emits an error token with the formatted arguments and returns the terminal state.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{Pos: l.position(), Type: TokenError, Value: fmt.Sprintf(format, args...)}
	return nil
}

// ignore a contiguous block of spaces.
func (l *lexer) ignoreSpace() {
	for unicode.IsSpace(l.next()) {
		l.ignore()
	}
	l.backup()
}

// lexToken is the top level state
func lexToken(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			l.backup()
			return lexNumberDigits
		case unicode.IsSpace(r):
			l.ignore()
		case r == '(':
			l.emit(TokenLParenthesis)
		case r == ')':
			l.emit(TokenRParenthesis)
		case r == '*':
			l.emit(TokenMultiplication)
		case r == '+':
			l.emit(TokenAddition)
		case r == eof:
			l.emit(TokenEOF)
			return nil
		default:
			return l.errorf("unexpected token %v", r)
		}
	}
}
func lexNumberDigits(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
		default:
			l.backup()
			l.emit(TokenNumber)
			return lexToken
		}
	}
}
