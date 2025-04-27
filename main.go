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
var inputFile = flag.String("input", "", "Input file to read expressions from")

func parseExpression(input string) (ast.Node, error) {
	// TODO: parser needs changes to collect errors rather than this 'exception handling' hack
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

func processLine(line string, lineNum int) {
	input := strings.TrimSpace(line)
	if input == "" {
		return
	}
	prefix := fmt.Sprintf("Line %d: ", lineNum)
	// Print the line being processed from the file
	fmt.Printf("%s'%s'\n", prefix, input)

	exprAst, parseErr := parseExpression(input)
	if parseErr != nil {
		fmt.Printf("  %sError: %v\n", prefix, parseErr)
		return // Don't proceed if parsing failed
	}

	// If exprAst is nil after parseExpression returns (without error),
	// it might mean a parsing issue was recovered via panic.
	// We check again to avoid nil pointer dereference.
	if exprAst == nil {
		return
	}

	if *printAST {
		fmt.Println("  AST:")
		fmt.Print(ast.PrettyPrintAST(exprAst, "    "))
		fmt.Printf("  Infix notation: %s\n", exprAst.String())
	}

	result, evalErr := eval.Eval(exprAst)
	if evalErr != nil {
		fmt.Printf("  %sEvaluation error: %v\n", prefix, evalErr)
	} else {
		fmt.Printf("Result =  '%g'\n", result)
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
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
		// don't evaluate empty input
		if input == "" {
			continue
		}

		exprAst, parseErr := parseExpression(input)
		if parseErr != nil {
			fmt.Printf("  Error: %v\n", parseErr)
			continue
		}

		if *printAST {
			fmt.Println("  AST:")
			fmt.Print(ast.PrettyPrintAST(exprAst, "    "))
			fmt.Printf("  Infix notation: %s\n", exprAst.String())
		}

		result, evalErr := eval.Eval(exprAst)
		if evalErr != nil {
			fmt.Printf("  Evaluation error: %v\n", evalErr)
		} else {
			fmt.Printf("  Result: %g\n", result)
		}
	}

}

func main() {
	flag.Parse()
	fmt.Println("Basic Arithmetic Parser REPL")
	fmt.Println("Enter expressions to evaluate or type 'exit' to quit.")

	if *inputFile != "" {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		fmt.Printf("Using input file: %s\n", *inputFile)
		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			processLine(scanner.Text(), lineNum)
		}

	} else {
		repl()
	}

}
