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

// number returns a string representation of a number in the input
func (l *Lexer) number() string {
	result := ""
	decimalPointSeen := false

	// collect digits and decimal point
	for l.currentChar != 0 {
		if unicode.IsDigit(rune(l.currentChar)) {
			result += string(l.currentChar)
			l.advance()
		} else if l.currentChar == '.' {
			if decimalPointSeen {
				// Found a second decimal point, stop consuming the number here.
				// The next call to GetNextToken will likely treat it as an invalid character or operator.
				// TODO: we also need to collect cases like this
				panic(fmt.Sprintf("Invalid number format: %s", result))
			}
			decimalPointSeen = true
			result += string(l.currentChar)
			l.advance()
		} else {
			// Not a digit or a decimal point, stop consuming the number.
			break
		}
	}

	return result
}

func (l *Lexer) GetNextToken() token.Token {
	for l.currentChar != 0 {
		if unicode.IsSpace(rune(l.currentChar)) {
			l.skipWhitespace()
			continue
		}

		// Check for numbers
		if unicode.IsDigit(rune(l.currentChar)) {
			return token.Token{Type: token.NUMBER, Value: l.number()}
		}

		// Check for operators
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

	// End of input
	return token.Token{Type: token.EOF, Value: ""}
}
