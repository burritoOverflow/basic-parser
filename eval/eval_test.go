package eval

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/parser"
	"basic-arithmetic-parser/token"
	"fmt"
	"testing"
)

func TestEvalOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected float64
		hasError bool
	}{
		{
			name: "Addition before Multiplication",
			// 2 + 3 * 4 = 14
			node: &ast.BinaryOpNode{
				Left: &ast.NumberNode{Value: 2},
				Op:   token.Token{Type: token.PLUS, Value: "+"},
				Right: &ast.BinaryOpNode{
					Left:  &ast.NumberNode{Value: 3},
					Op:    token.Token{Type: token.MULTIPLY, Value: "*"},
					Right: &ast.NumberNode{Value: 4},
				},
			},
			expected: 14,
			hasError: false,
		},
		{
			name: "Parentheses override Multiplication",
			// (2 + 3) * 4 = 20
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left:  &ast.NumberNode{Value: 2},
					Op:    token.Token{Type: token.PLUS, Value: "+"},
					Right: &ast.NumberNode{Value: 3},
				},
				Op:    token.Token{Type: token.MULTIPLY, Value: "*"},
				Right: &ast.NumberNode{Value: 4},
			},
			expected: 20,
			hasError: false,
		},
		{
			name: "Division before Subtraction",
			// 10 / 2 - 1 = 4
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left:  &ast.NumberNode{Value: 10},
					Op:    token.Token{Type: token.DIVIDE, Value: "/"},
					Right: &ast.NumberNode{Value: 2},
				},
				Op:    token.Token{Type: token.MINUS, Value: "-"},
				Right: &ast.NumberNode{Value: 1},
			},
			expected: 4,
			hasError: false,
		},
		{
			name: "Unary minus precedence",
			// -2 + 5 = 3
			node: &ast.BinaryOpNode{
				Left: &ast.UnaryOpNode{
					Op:   token.Token{Type: token.MINUS, Value: "-"},
					Expr: &ast.NumberNode{Value: 2},
				},
				Op:    token.Token{Type: token.PLUS, Value: "+"},
				Right: &ast.NumberNode{Value: 5},
			},
			expected: 3,
			hasError: false,
		},
		{
			name: "Unary minus with multiplication",
			// 5 * -2 = -10
			node: &ast.BinaryOpNode{
				Left: &ast.NumberNode{Value: 5},
				Op:   token.Token{Type: token.MULTIPLY, Value: "*"},
				Right: &ast.UnaryOpNode{
					Op:   token.Token{Type: token.MINUS, Value: "-"},
					Expr: &ast.NumberNode{Value: 2},
				},
			},
			expected: -10,
			hasError: false,
		},
		{
			name: "Unary plus (identity)",
			// +3 * 4 = 12
			node: &ast.BinaryOpNode{
				Left: &ast.UnaryOpNode{
					Op:   token.Token{Type: token.PLUS, Value: "+"},
					Expr: &ast.NumberNode{Value: 3},
				},
				Op:    token.Token{Type: token.MULTIPLY, Value: "*"},
				Right: &ast.NumberNode{Value: 4},
			},
			expected: 12,
			hasError: false,
		},
		{
			name: "Complex precedence",
			// 2 + 3 * 4 / 2 - 1 = 7
			// AST: (2 + ((3 * 4) / 2)) - 1
			node: &ast.BinaryOpNode{
				Left: &ast.BinaryOpNode{
					Left: &ast.NumberNode{Value: 2},
					Op:   token.Token{Type: token.PLUS, Value: "+"},
					Right: &ast.BinaryOpNode{
						Left: &ast.BinaryOpNode{
							Left:  &ast.NumberNode{Value: 3},
							Op:    token.Token{Type: token.MULTIPLY, Value: "*"},
							Right: &ast.NumberNode{Value: 4},
						},
						Op:    token.Token{Type: token.DIVIDE, Value: "/"},
						Right: &ast.NumberNode{Value: 2},
					},
				},
				Op:    token.Token{Type: token.MINUS, Value: "-"},
				Right: &ast.NumberNode{Value: 1},
			},
			expected: 7,
			hasError: false,
		},
		{
			name: "Division by zero",
			// 5 / 0
			node: &ast.BinaryOpNode{
				Left:  &ast.NumberNode{Value: 5},
				Op:    token.Token{Type: token.DIVIDE, Value: "/"},
				Right: &ast.NumberNode{Value: 0},
			},
			expected: 0, // Expected value doesn't matter when error occurs
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Eval(tt.node)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %g, but got %g", tt.expected, result)
				}
			}
		})
	}
}

func TestEvalWithParser(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expected        float64
		expectPanic     bool
		expectEvalError bool
	}{
		{"Simple Addition", "2 + 3", 5, false, false},
		{"Simple Subtraction", "10 - 4", 6, false, false},
		{"Simple Multiplication", "5 * 6", 30, false, false},
		{"Simple Division", "20 / 4", 5, false, false},
		{"Operator Precedence 1", "2 + 3 * 4", 14, false, false},
		{"Operator Precedence 2", "10 - 2 * 3", 4, false, false},
		{"Operator Precedence 3", "12 / 3 + 1", 5, false, false},
		{"Parentheses Override 1", "(2 + 3) * 4", 20, false, false},
		{"Parentheses Override 2", "10 - (2 + 3)", 5, false, false},
		{"Unary Minus 1", "-5 + 10", 5, false, false},
		{"Unary Minus 2", "5 * -2", -10, false, false},
		{"Unary Plus", "+3 + 5", 8, false, false},
		{"Complex Expression", "-(2 + 3) * 4 / 10 - 1", -3, false, false},
		{"Division by Zero", "10 / 0", 0, false, true},  // Parse OK, Eval Error
		{"Invalid Syntax 1", "2 + * 3", 0, true, false}, // Expect Parser Panic
		{"Invalid Syntax 2", "1 + 2 )", 0, true, false}, // Expect Parser Panic (unmatched parenthesis)
		{"Invalid Syntax 3", "( 1 + 2", 0, true, false}, // Expect Parser Panic (missing closing parenthesis)
		{"Empty Input", "", 0, true, false},             // Expect Parser Panic
		{"Just Operator", "+", 0, true, false},          // Expect Parser Panic
		{"Number Only", "42", 42, false, false},
		{"Unary Only", "-10", -10, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var panicked bool
			var panicValue interface{}
			var program ast.Node

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
						panicValue = r
						fmt.Printf("Caught panic: %v\n", r)
					}
				}()

				// Attempt to parse inside the function with the defer/recover
				l := lexer.New(tt.input)
				p := parser.New(l)
				program = p.Parse() // This might panic (when parsing fails).
			}()

			if panicked != tt.expectPanic {
				t.Fatalf("Expected panic: %v, but got panicked=%v. Panic value: %v", tt.expectPanic, panicked, panicValue)
			}

			if tt.expectPanic {
				return
			}

			// Ensure the parsed program/node is not nil if no panic occurred
			if program == nil {
				// This case might indicate an issue where Parse returns nil without panicking
				t.Fatalf("Parsing did not panic but returned a nil AST node for input: %s", tt.input)
			}

			result, evalErr := Eval(program)
			hasEvalError := evalErr != nil
			if hasEvalError != tt.expectEvalError {
				t.Fatalf("Expected eval error: %v, got: %v. Error: %v", tt.expectEvalError, hasEvalError, evalErr)
			}

			// Stop if evaluation failed and we expected it to succeed
			if !tt.expectEvalError && hasEvalError {
				t.Fatalf("Did not expect eval error, but got: %v", evalErr)
			}

			// Stop if evaluation succeeded but we expected it to fail
			if tt.expectEvalError && !hasEvalError {
				t.Fatalf("Expected eval error, but got none for input: %s", tt.input)
			}

			// If evaluation was expected to fail, we don't compare the result
			if tt.expectEvalError {
				return
			}

			// Compare the result if no errors were expected
			if result != tt.expected {
				t.Errorf("Expected result %g, but got %g for input: %s", tt.expected, result, tt.input)
			}
		})
	}
}
