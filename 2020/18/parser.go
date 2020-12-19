package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

//Parse parses a string
func Parse(input string, version string) (int, error) {
	parser := newParser(Lex(input), version)
	return parser.Parse()
}

type parser struct {
	tokens    <-chan Token
	version   string
	lookahead [2]Token
	peekCount int
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (p *parser) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*errp = e.(error)
	}
	return
}
func newParser(tokens <-chan Token, version string) *parser {
	return &parser{
		tokens:  tokens,
		version: version,
	}
}
func (p *parser) Parse() (expression int, err error) {
	// Parsing uses panics to bubble up errors
	defer p.recover(&err)
	if p.version == "v2" {
		expression = p.advancedMath()
	} else {
		expression = p.simpleMath()
	}

	return
}

func (p *parser) nextToken() Token {
	return <-p.tokens
}

// next returns the next token.
func (p *parser) next() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.lookahead[0] = p.nextToken()
	}
	return p.lookahead[p.peekCount]
}

// backup backs the input stream up one token.
func (p *parser) backup() {
	p.peekCount++
}

// peek returns but does not consume the next token.
func (p *parser) peek() Token {
	if p.peekCount > 0 {
		return p.lookahead[p.peekCount-1]
	}
	p.peekCount = 1
	p.lookahead[1] = p.lookahead[0]
	p.lookahead[0] = p.nextToken()
	return p.lookahead[0]
}

// errorf formats the error and terminates processing.
func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("parser: %s", format)
	panic(fmt.Errorf(format, args...))
}

// expect consumes the next token and guarantees it has the required type.
func (p *parser) expect(expected TokenType) Token {
	t := p.next()
	if t.Type != expected {
		debug.PrintStack()
		p.unexpected(t, expected)
	}
	return t
}

// unexpected complains about the token and terminates processing.
func (p *parser) unexpected(tok Token, expected ...TokenType) {
	expectedStrs := make([]string, len(expected))
	for i := range expected {
		expectedStrs[i] = fmt.Sprintf("%q", expected[i])
	}
	expectedStr := strings.Join(expectedStrs, ",")
	debug.PrintStack()
	p.errorf("unexpected token %q with value %q at line %d char %d, expected: %s", tok.Type, tok.Value, tok.Pos.Line, tok.Pos.Char, expectedStr)
}

// error terminates processing.
func (p *parser) error(err error) {
	p.errorf("%s", err)
}

func (p *parser) advancedMath() int {
	var result int
	var currentOp string
	for {
		switch p.peek().Type {
		case TokenEOF:
			return result
		case TokenRParenthesis:
			return result
		case TokenLParenthesis:
			v := p.parenthesisSentence()
			if currentOp != "+" {
				if result == 0 {
					result = 1
				}
				result *= v
			} else {
				result += v
			}
		case TokenMultiplication:
			p.next()
			currentOp = "*"
			if result == 0 {
				result = 1
			}
			r := p.advancedMath()
			//fmt.Printf("%d * %d\n", result, r)
			result *= r
		case TokenAddition:
			p.next()
			currentOp = "+"

		case TokenNumber:
			tok := p.next()
			value, err := strconv.Atoi(tok.Value)
			if err != nil {
				panic(err)
			}
			if result == 0 {
				result = value
			} else if currentOp != "" {
				//fmt.Printf("%d %s %d\n", result, currentOp, value)
				if currentOp == "+" {
					result += value
				}
			}
		default:
			p.unexpected(p.next(), TokenEOF, TokenLParenthesis)
		}
	}
}
func (p *parser) parenthesisSentence() int {
	p.expect(TokenLParenthesis)
	v := p.advancedMath()
	p.expect(TokenRParenthesis)
	return v
}
func (p *parser) simpleMath() int {
	var result int
	var currentOp string
	for {
		//fmt.Printf("Next token: %+v\n", p.peek())
		switch p.next().Type {
		case TokenEOF:
			return result
		case TokenRParenthesis:
			return result
		case TokenLParenthesis:
			value := p.simpleMath()
			if currentOp == "*" {
				result *= value
			} else {
				result += value
			}
		case TokenMultiplication:
			currentOp = "*"

		case TokenAddition:
			currentOp = "+"

		case TokenNumber:
			p.backup()
			tok := p.next()
			value, err := strconv.Atoi(tok.Value)
			if err != nil {
				panic(err)
			}
			if result == 0 {
				result = value
			} else {
				if currentOp == "*" {
					result *= value
				} else {
					result += value
				}
			}
		default:
			p.unexpected(p.next(), TokenEOF, TokenLParenthesis)
		}
	}
}
