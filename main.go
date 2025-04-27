package main

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/parser"
	"fmt"
	"strings"
)

func ParseExpression(input string) (ast.Node, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	input = strings.TrimSpace(input)
	l := lexer.New(input)
	p := parser.New(l)
	return p.Parse(), nil
}

func main() {
	expressions := []string{
		"3 + 4 * 2",
		"(3 + 4) * 2",
		"7 - 3 + 2",
		"7 - (3 + 2)",
		"10 / 2 - 3",
		"3 * (4 + 2) / 3",
		"-5 + 3",
		"+2 * 3",
	}

	for _, expr := range expressions {
		fmt.Printf("\nExpression: %s\n", expr)
		fmt.Println("AST:")
		if exprAst, err := ParseExpression(expr); err == nil {
			fmt.Println(ast.PrettyPrintAST(exprAst, ""))
			fmt.Printf("Infix notation: %s\n", exprAst.String())
		} else {
			fmt.Printf("Error parsing '%s': %v\n", expr, err)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}
