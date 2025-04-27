package eval

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/token"
	"fmt"
)

// Eval evaluates the given AST node and returns the result as a float64.
// It returns an error for invalid operations like division by zero.
func Eval(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return n.Value, nil
	case *ast.BinaryOpNode:
		leftVal, err := Eval(n.Left)
		if err != nil {
			return 0, err
		}
		rightVal, err := Eval(n.Right)
		if err != nil {
			return 0, err
		}

		switch n.Op.Type {
		case token.PLUS:
			return leftVal + rightVal, nil
		case token.MINUS:
			return leftVal - rightVal, nil
		case token.MULTIPLY:
			return leftVal * rightVal, nil
		case token.DIVIDE:
			if rightVal == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return leftVal / rightVal, nil
		default:
			return 0, fmt.Errorf("unknown binary operator: %s", n.Op.Value)
		}
	case *ast.UnaryOpNode:
		exprVal, err := Eval(n.Expr)
		if err != nil {
			return 0, err
		}

		switch n.Op.Type {
		case token.PLUS: // Unary plus (identity)
			return exprVal, nil
		case token.MINUS: // Unary minus (negation)
			return -exprVal, nil
		default:
			return 0, fmt.Errorf("unknown unary operator: %s", n.Op.Value)
		}
	default:
		return 0, fmt.Errorf("unknown node type: %T", node)
	}
}
