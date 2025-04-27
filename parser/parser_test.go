package parser

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/token"
	"testing"
)

func checkNumberNode(t *testing.T, node ast.Node, expected float64) bool {
	t.Helper()
	numNode, ok := node.(*ast.NumberNode)
	if !ok {
		t.Errorf("node not *ast.NumberNode. got=%T", node)
		return false
	}
	if numNode.Value != expected {
		t.Errorf("node.Value not %f. got=%f", expected, numNode.Value)
		return false
	}
	return true
}

func checkBinaryOpNode(t *testing.T, node ast.Node, expectedOp token.TokenType) (*ast.BinaryOpNode, bool) {
	t.Helper()
	binOpNode, ok := node.(*ast.BinaryOpNode)
	if !ok {
		t.Errorf("node not *ast.BinaryOpNode. got=%T", node)
		return nil, false
	}
	if binOpNode.Op.Type != expectedOp {
		t.Errorf("node.Op.Type not %v. got=%v", expectedOp, binOpNode.Op.Type)
		return binOpNode, false
	}
	return binOpNode, true
}

func checkUnaryOpNode(t *testing.T, node ast.Node, expectedOp token.TokenType) (*ast.UnaryOpNode, bool) {
	t.Helper()
	unOpNode, ok := node.(*ast.UnaryOpNode)
	if !ok {
		t.Errorf("node not *ast.UnaryOpNode. got=%T", node)
		return nil, false
	}
	if unOpNode.Op.Type != expectedOp {
		t.Errorf("node.Op.Type not %v. got=%v", expectedOp, unOpNode.Op.Type)
		return unOpNode, false
	}
	return unOpNode, true
}

func TestSimpleAddition(t *testing.T) {
	input := "3 + 5"
	l := lexer.New(input)
	p := New(l)
	rootNode := p.Parse()

	binOp, ok := checkBinaryOpNode(t, rootNode, token.PLUS)
	if !ok {
		t.Fatalf("Root node is not a BinaryOpNode with PLUS operator")
	}

	if !checkNumberNode(t, binOp.Left, 3) {
		t.Errorf("Left operand check failed")
	}
	if !checkNumberNode(t, binOp.Right, 5) {
		t.Errorf("Right operand check failed")
	}
}

func TestOperatorPrecedence(t *testing.T) {
	// Expected: 3 + (5 * 2)
	input := "3 + 5 * 2"
	l := lexer.New(input)
	p := New(l)
	rootNode := p.Parse()

	// Root should be PLUS
	rootBinOp, ok := checkBinaryOpNode(t, rootNode, token.PLUS)
	if !ok {
		t.Fatalf("Root node is not a BinaryOpNode with PLUS operator")
	}

	// Left of PLUS should be 3
	if !checkNumberNode(t, rootBinOp.Left, 3) {
		t.Errorf("Left operand of PLUS check failed")
	}

	// Right of PLUS should be MULTIPLY
	rightBinOp, ok := checkBinaryOpNode(t, rootBinOp.Right, token.MULTIPLY)
	if !ok {
		t.Fatalf("Right operand of PLUS is not a BinaryOpNode with MULTIPLY operator")
	}

	// Left of MULTIPLY should be 5
	if !checkNumberNode(t, rightBinOp.Left, 5) {
		t.Errorf("Left operand of MULTIPLY check failed")
	}
	// Right of MULTIPLY should be 2
	if !checkNumberNode(t, rightBinOp.Right, 2) {
		t.Errorf("Right operand of MULTIPLY check failed")
	}
}

func TestParentheses(t *testing.T) {
	// Expected: (3 + 5) * 2
	input := "(3 + 5) * 2"
	l := lexer.New(input)
	p := New(l)
	rootNode := p.Parse()

	// Root should be MULTIPLY
	rootBinOp, ok := checkBinaryOpNode(t, rootNode, token.MULTIPLY)
	if !ok {
		t.Fatalf("Root node is not a BinaryOpNode with MULTIPLY operator")
	}

	// Right of MULTIPLY should be 2
	if !checkNumberNode(t, rootBinOp.Right, 2) {
		t.Errorf("Right operand of MULTIPLY check failed")
	}

	// Left of MULTIPLY should be PLUS
	leftBinOp, ok := checkBinaryOpNode(t, rootBinOp.Left, token.PLUS)
	if !ok {
		t.Fatalf("Left operand of MULTIPLY is not a BinaryOpNode with PLUS operator")
	}

	// Left of PLUS should be 3
	if !checkNumberNode(t, leftBinOp.Left, 3) {
		t.Errorf("Left operand of PLUS check failed")
	}
	// Right of PLUS should be 5
	if !checkNumberNode(t, leftBinOp.Right, 5) {
		t.Errorf("Right operand of PLUS check failed")
	}
}

func TestUnaryMinus(t *testing.T) {
	input := "-5"
	l := lexer.New(input)
	p := New(l)
	rootNode := p.Parse()

	unOp, ok := checkUnaryOpNode(t, rootNode, token.MINUS)
	if !ok {
		t.Fatalf("Root node is not a UnaryOpNode with MINUS operator")
	}

	if !checkNumberNode(t, unOp.Expr, 5) {
		t.Errorf("Unary operand check failed")
	}
}

func TestUnaryPlus(t *testing.T) {
	input := "+5"
	l := lexer.New(input)
	p := New(l)
	rootNode := p.Parse()

	unOp, ok := checkUnaryOpNode(t, rootNode, token.PLUS)
	if !ok {
		t.Fatalf("Root node is not a UnaryOpNode with PLUS operator")
	}

	if !checkNumberNode(t, unOp.Expr, 5) {
		t.Errorf("Unary operand check failed")
	}
}
