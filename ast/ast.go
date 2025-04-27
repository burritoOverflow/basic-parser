package ast

import (
	"basic-arithmetic-parser/token"
	"fmt"
)

type NodeType int

const (
	NUMBER_NODE NodeType = iota
	BINARY_OP_NODE
	UNARY_OP_NODE
)

type Node interface {
	Type() NodeType
	String() string
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) Type() NodeType {
	return NUMBER_NODE
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%g", n.Value)
}

type BinaryOpNode struct {
	Left  Node
	Op    token.Token
	Right Node
}

func (n *BinaryOpNode) Type() NodeType {
	return BINARY_OP_NODE
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op.Value, n.Right.String())
}

// Unary operation node
type UnaryOpNode struct {
	Op   token.Token
	Expr Node
}

func (n *UnaryOpNode) Type() NodeType {
	return UNARY_OP_NODE
}

func (n *UnaryOpNode) String() string {
	return fmt.Sprintf("%s%s", n.Op.Value, n.Expr.String())
}

// Generate a visual representation of the AST with indentation
// TODO have this return an error when appropriate instead of 'Unknown node type'
func PrettyPrintAST(node Node, indent string) string {
	switch n := node.(type) {
	case *NumberNode:
		return fmt.Sprintf("%sNumber(%g)\n", indent, n.Value)
	case *BinaryOpNode:
		result := fmt.Sprintf("%sBinaryOp(%s)\n", indent, n.Op.Value)
		result += fmt.Sprintf("%s  Left:\n", indent)
		result += PrettyPrintAST(n.Left, indent+"    ")
		result += fmt.Sprintf("%s  Right:\n", indent)
		result += PrettyPrintAST(n.Right, indent+"    ")
		return result
	case *UnaryOpNode:
		result := fmt.Sprintf("%sUnaryOp(%s)\n", indent, n.Op.Value)
		result += fmt.Sprintf("%s  Expr:\n", indent)
		result += PrettyPrintAST(n.Expr, indent+"    ")
		return result
	default:
		return fmt.Sprintf("%sUnknown node type\n", indent)
	}
}
