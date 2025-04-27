package ast

import (
	"basic-arithmetic-parser/token"
	"fmt"
	"strings"
	"testing"
)

func TestNodeType(t *testing.T) {
	numNode := &NumberNode{Value: 123}
	if numNode.Type() != NUMBER_NODE {
		t.Errorf("NumberNode.Type() failed. Expected %v, got %v", NUMBER_NODE, numNode.Type())
	}

	binOpNode := &BinaryOpNode{
		Left:  &NumberNode{Value: 1},
		Op:    token.Token{Type: token.PLUS, Value: "+"},
		Right: &NumberNode{Value: 2},
	}
	if binOpNode.Type() != BINARY_OP_NODE {
		t.Errorf("BinaryOpNode.Type() failed. Expected %v, got %v", BINARY_OP_NODE, binOpNode.Type())
	}
	if binOpNode.String() != "(1 + 2)" {
		t.Errorf("BinaryOpNode.String() failed. Expected %q, got %q", "(1 + 2)", binOpNode.String())
	}

	unOpNode := &UnaryOpNode{
		Op:   token.Token{Type: token.MINUS, Value: "-"},
		Expr: &NumberNode{Value: 5},
	}
	if unOpNode.Type() != UNARY_OP_NODE {
		t.Errorf("UnaryOpNode.Type() failed. Expected %v, got %v", UNARY_OP_NODE, unOpNode.Type())
	}
	if unOpNode.String() != "-5" {
		t.Errorf("UnaryOpNode.String() failed. Expected %q, got %q", "-5", unOpNode.String())
	}
}

func TestNodeString(t *testing.T) {
	tests := []struct {
		node     Node
		expected string
	}{
		{
			&NumberNode{Value: 123.45},
			"123.45",
		},
		{
			&NumberNode{Value: 10}, // Test integer formatting
			"10",
		},
		{
			&BinaryOpNode{
				Left:  &NumberNode{Value: 1},
				Op:    token.Token{Type: token.PLUS, Value: "+"},
				Right: &NumberNode{Value: 2},
			},
			"(1 + 2)",
		},
		{
			&BinaryOpNode{
				Left: &BinaryOpNode{
					Left:  &NumberNode{Value: 1},
					Op:    token.Token{Type: token.MULTIPLY, Value: "*"},
					Right: &NumberNode{Value: 2},
				},
				Op:    token.Token{Type: token.MINUS, Value: "-"},
				Right: &NumberNode{Value: 3},
			},
			"((1 * 2) - 3)",
		},
		{
			&UnaryOpNode{
				Op:   token.Token{Type: token.MINUS, Value: "-"},
				Expr: &NumberNode{Value: 5},
			},
			"-5",
		},
		{
			&UnaryOpNode{
				Op: token.Token{Type: token.PLUS, Value: "+"},
				Expr: &BinaryOpNode{
					Left:  &NumberNode{Value: 1},
					Op:    token.Token{Type: token.DIVIDE, Value: "/"},
					Right: &NumberNode{Value: 2},
				},
			},
			"+(1 / 2)",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			if tt.node.String() != tt.expected {
				t.Errorf("String() mismatch. Expected %q, got %q", tt.expected, tt.node.String())
			}
		})
	}
}

func TestPrettyPrintAST(t *testing.T) {
	// Build a sample AST: (1 + 2) * -3
	ast := &BinaryOpNode{
		Left: &BinaryOpNode{
			Left:  &NumberNode{Value: 1},
			Op:    token.Token{Type: token.PLUS, Value: "+"},
			Right: &NumberNode{Value: 2},
		},
		Op: token.Token{Type: token.MULTIPLY, Value: "*"},
		Right: &UnaryOpNode{
			Op:   token.Token{Type: token.MINUS, Value: "-"},
			Expr: &NumberNode{Value: 3},
		},
	}

	expectedOutput := `
BinaryOp(*)
  Left: 
    BinaryOp(+)
      Left:
        Number(1)
      Right: 
        Number(2)
  Right: 
    UnaryOp(-)
      Expr: 
        Number(3)
`
	// Trim leading/trailing whitespace and normalize line endings for comparison
	normalize := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\r\n", "\n")
		return s
	}

	actualOutput := PrettyPrintAST(ast, "")
	if normalize(actualOutput) != normalize(expectedOutput) {
		t.Errorf("PrettyPrintAST mismatch.\nExpected:\n%s\nGot:\n%s", normalize(expectedOutput), normalize(actualOutput))
	}

	// Test with a simple number node
	numNode := &NumberNode{Value: 42}
	expectedNumOutput := "Number(42)\n"
	actualNumOutput := PrettyPrintAST(numNode, "")
	if strings.TrimSpace(actualNumOutput) != strings.TrimSpace(expectedNumOutput) {
		t.Errorf("PrettyPrintAST for NumberNode mismatch.\nExpected:\n%s\nGot:\n%s", expectedNumOutput, actualNumOutput)
	}
}
