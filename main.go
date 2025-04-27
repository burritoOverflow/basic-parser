package main

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/eval"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/parser"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var printAST = flag.Bool("ast", false, "Print the Abstract Syntax Tree")

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
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Basic Arithmetic Parser REPL")
	fmt.Println("Enter expressions to evaluate or type 'exit' to quit.")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		exprAst, parseErr := ParseExpression(input)

		if parseErr != nil {
			fmt.Printf("  Error: %v\n", parseErr)
			continue
		}

		if *printAST {
			fmt.Println("  AST:")
			fmt.Print(ast.PrettyPrintAST(exprAst, "    "))
			fmt.Printf("  Infix notation: %s\n", exprAst.String())
		}

		// Evaluate the AST
		result, evalErr := eval.Eval(exprAst)
		if evalErr != nil {
			fmt.Printf("  Evaluation error: %v\n", evalErr)
		} else {
			fmt.Printf("  Result: %g\n", result)
		}
	}
}
