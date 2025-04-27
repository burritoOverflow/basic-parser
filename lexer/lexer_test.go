package lexer

import (
	"basic-arithmetic-parser/token"
	"testing"
)

func TestGetNextToken(t *testing.T) {
	input := `10 + 2 * (5 - 3) / 4.5`

	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.NUMBER, "10"},
		{token.PLUS, "+"},
		{token.NUMBER, "2"},
		{token.MULTIPLY, "*"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.MINUS, "-"},
		{token.NUMBER, "3"},
		{token.RPAREN, ")"},
		{token.DIVIDE, "/"},
		{token.NUMBER, "4.5"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.GetNextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - token value wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}

func TestWhitespaceAndNumbers(t *testing.T) {
	input := `  123   45.67  `

	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.NUMBER, "123"},
		{token.NUMBER, "45.67"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.GetNextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("tests[%d] - token value wrong. expected=%q, got=%q",
				i, tt.expectedValue, tok.Value)
		}
	}
}
