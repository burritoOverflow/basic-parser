package eval

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/token"
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
