package main

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/codegen"
	"basic-arithmetic-parser/eval"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/parser"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var printAST = flag.Bool("ast", false, "Print the Abstract Syntax Tree")
var inputFile = flag.String("input", "", "Input file to read expressions from")

// NOTE: this is also used for the object file name from nasm and the final executable
// optional, only if one desires to have a different name for output artifacts
var outputFile = flag.String("output", "", "Output file for generated assembly")

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

func getInfix(exprAst *ast.Node) string {
	return fmt.Sprintf("  Infix notation: %s\n", (*exprAst).String())
}

func showAST(exprAst *ast.Node) {
	// global flag--show only if set
	if *printAST {
		fmt.Println("  AST:")
		fmt.Print(ast.PrettyPrintAST(*exprAst, "    "))
		fmt.Printf("%s", getInfix(exprAst))
	}
}

func doEval(exprAst *ast.Node, prefix *string) {
	result, evalErr := eval.Eval(*exprAst)
	if evalErr != nil {
		// no need to propagate the error; each use will continue
		fmt.Printf("  %sEvaluation error: %v\n", *prefix, evalErr)
	} else {
		fmt.Printf("Result =  '%g'\n", result)
	}
}

func processLine(line string, lineNum int) ast.Node {
	input := strings.TrimSpace(line)
	if input == "" {
		return nil // Skip empty lines
	}
	prefix := fmt.Sprintf("L%d: ", lineNum)
	fmt.Printf("%s'%s'\n", prefix, input)

	exprAst, _ := parseExpression(input)

	// Handle potential nil node from parsing (e.g., empty input or panic)
	if exprAst == nil {
		// parseErr check might be useful if parseExpression changes
		fmt.Printf("  %sSkipping line due to parsing issue or empty input.\n", prefix)
		return nil
	}
	showAST(&exprAst)
	return exprAst
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

		showAST(&exprAst)
		doEval(&exprAst, nil)
	}

}

// get the path of the executable in the system PATH, or return an error if not found
// logs the status of the search in either cond
func findInPath(executable string) (string, error) {
	path, err := exec.LookPath(executable)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH: %v", executable, err)
	}
	fmt.Printf("Found %s at: '%s'\n", executable, path)
	return path, nil
}

// Attempts to compile the contents of the input file and generate assembly code.
// Exits with 1 if any error occurs at any point in the process
func compileFileContents(filePath string, outFilePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening input file '%s': %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Processing input file: %s\n", filePath)

	var validAsts []ast.Node
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		exprAst := processLine(scanner.Text(), lineNum)
		if exprAst != nil {
			validAsts = append(validAsts, exprAst)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file '%s': %v\n", filePath, err)
	}

	if len(validAsts) == 0 {
		fmt.Println("No valid expressions found in the input file. No assembly generated.")
		return
	}

	fmt.Printf("Generating assembly for %d expression(s)...\n", len(validAsts))
	compiler := codegen.New()
	fullAssembly, err := compiler.GenerateNasm(validAsts)
	if err != nil {
		fmt.Printf("Error during code generation: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Printf("Error creating output file '%s': %v\n", outFilePath, err)
		os.Exit(1)
	}
	_, writeErr := outFile.WriteString(fullAssembly)
	closeErr := outFile.Close()
	if writeErr != nil {
		fmt.Printf("Error writing assembly to file '%s': %v\n", outFilePath, writeErr)
		os.Exit(1)
	}
	if closeErr != nil {
		fmt.Printf("Error closing assembly file '%s': %v\n", outFilePath, closeErr)
	}
	fmt.Printf("Assembly written to '%s'.\n", outFilePath)

	baseName := strings.TrimSuffix(outFilePath, filepath.Ext(outFilePath))
	objectFilePath := baseName + ".o"
	executableFilePath := baseName

	nasmPath, err := findInPath("nasm")
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Running NASM\n  Command: '%s -f elf64 %s -o %s'\n", nasmPath, outFilePath, objectFilePath)
	nasmCmd := exec.Command(nasmPath, "-f", "elf64", outFilePath, "-o", objectFilePath)
	// gets both stdout and stderr
	nasmOutput, err := nasmCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("NASM Error:\n%s\n", string(nasmOutput))
		fmt.Printf("Failed to assemble %s: %v\n", outFilePath, err)
		os.Exit(1)
	}
	if len(nasmOutput) > 0 {
		fmt.Printf("NASM Output:\n%s\n", string(nasmOutput))
	}
	fmt.Println("NASM completed successfully.")

	gccPath, err := findInPath("gcc")
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Running GCC\n  Command: '%s %s -no-pie -o %s'\n", gccPath, objectFilePath, executableFilePath)
	gccCmd := exec.Command(gccPath, objectFilePath, "-no-pie", "-o", executableFilePath)
	gccOutput, err := gccCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("GCC Error:\n%s\n", string(gccOutput))
		fmt.Printf("Failed to link %s: %v\n", objectFilePath, err)
		os.Exit(1)
	}

	if len(gccOutput) > 0 {
		fmt.Printf("GCC Output:\n%s\n", string(gccOutput))
	}

	fmt.Println("GCC completed successfully.")
}

func main() {
	flag.Parse()
	fmt.Println("Basic Arithmetic Parser REPL")
	fmt.Println("Enter expressions to evaluate or type 'exit' to quit.")

	if *inputFile != "" {
		if *outputFile == "" {
			// if no output file is specified, use the input file name
			*outputFile = strings.TrimSuffix(*inputFile, filepath.Ext(*inputFile)) + ".asm"
		} else {
			// append ".asm" if not already present
			if !strings.HasSuffix(*outputFile, ".asm") {
				*outputFile = *outputFile + ".asm"
			}
		}
		compileFileContents(*inputFile, *outputFile)
	} else {
		repl()
	}

}
