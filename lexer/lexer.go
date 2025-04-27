package lexer

import (
	"basic-arithmetic-parser/token"
	"fmt"
	"unicode"
)

type Lexer struct {
	input       string
	position    int
	currentChar byte
}

func New(input string) *Lexer {
	lexer := &Lexer{
		input:    input,
		position: 0,
	}
	if len(input) > 0 {
		lexer.currentChar = input[0]
	} else {
		lexer.currentChar = 0
	}
	return lexer
}

func (l *Lexer) advance() {
	l.position++
	if l.position < len(l.input) {
		l.currentChar = l.input[l.position]
	} else {
		l.currentChar = 0 // End of input
	}
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar != 0 && unicode.IsSpace(rune(l.currentChar)) {
		l.advance()
	}
}

func (l *Lexer) number() string {
	result := ""
	for l.currentChar != 0 && (unicode.IsDigit(rune(l.currentChar)) || l.currentChar == '.') {
		result += string(l.currentChar)
		l.advance()
	}
	return result
}

func (l *Lexer) GetNextToken() token.Token {
	for l.currentChar != 0 {
		if unicode.IsSpace(rune(l.currentChar)) {
			l.skipWhitespace()
			continue
		}

		if unicode.IsDigit(rune(l.currentChar)) {
			return token.Token{Type: token.NUMBER, Value: l.number()}
		}

		switch l.currentChar {
		case '+':
			l.advance()
			return token.Token{Type: token.PLUS, Value: "+"}
		case '-':
			l.advance()
			return token.Token{Type: token.MINUS, Value: "-"}
		case '*':
			l.advance()
			return token.Token{Type: token.MULTIPLY, Value: "*"}
		case '/':
			l.advance()
			return token.Token{Type: token.DIVIDE, Value: "/"}
		case '(':
			l.advance()
			return token.Token{Type: token.LPAREN, Value: "("}
		case ')':
			l.advance()
			return token.Token{Type: token.RPAREN, Value: ")"}
		default:
			panic(fmt.Sprintf("Invalid character: %c", l.currentChar))
		}
	}

	return token.Token{Type: token.EOF, Value: ""}
}
