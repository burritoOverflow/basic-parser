package codegen

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/token"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Compiler holds the state for assembly generation for a complete program.
type Compiler struct {
	dataSection bytes.Buffer // Accumulates .data section content
	textSection bytes.Buffer // Accumulates .text section content
	labelCount  int          // To generate unique labels
	floatFmt    string       // Label for the printf format string
}

func New() *Compiler {
	return &Compiler{labelCount: 0}
}

// GenerateNasm generates a complete x86 assembly program (in NASM syntax)
// that evaluates each AST node in sequence and prints the result.
func (c *Compiler) GenerateNasm(nodes []ast.Node) (string, error) {
	c.dataSection.WriteString("section .data\n")
	c.floatFmt = fmt.Sprintf("LC%d", c.labelCount)
	c.labelCount++

	// format to two decimal places
	c.dataSection.WriteString(fmt.Sprintf("    %s: db \"%%.2f\", 10, 0\n", c.floatFmt))

	c.textSection.WriteString("section .text\n")
	c.textSection.WriteString("global main\n")
	c.textSection.WriteString("extern printf\n")
	c.textSection.WriteString("main:\n")

	// Initialize the FPU:  see: Vol. 2A 3-402
	c.textSection.WriteString("    finit\n\n")

	for i, node := range nodes {
		c.textSection.WriteString(fmt.Sprintf("    ; --- Expression %d ---\n", i+1))
		err := c.generateExpr(node) // This will now use 'f' format
		if err != nil {
			return "", fmt.Errorf("error compiling expression %d: %w", i+1, err)
		}
		c.generatePrintFloat()
		c.textSection.WriteString("\n")
	}

	c.textSection.WriteString("    ; --- Exit program ---\n")
	c.textSection.WriteString("    xor eax, eax\n")
	c.textSection.WriteString("    ret\n")
	c.textSection.WriteString("    section .note.GNU-stack\n")

	var finalCode strings.Builder
	finalCode.WriteString(c.dataSection.String())
	finalCode.WriteString("\n")
	finalCode.WriteString(c.textSection.String())

	return finalCode.String(), nil
}

// generateExpr recursively generates code for a given node.
// Appends instructions to c.textSection and constants to c.dataSection.
// Leaves the result on the FPU stack st0.
func (c *Compiler) generateExpr(node ast.Node) error {
	switch n := node.(type) {
	case *ast.NumberNode:
		// where `LC` is adopted naming convention for labels
		// for constants in the data section (Literal Constant)
		// i.e Label N = LC1, LC2, LC3...
		label := fmt.Sprintf("LC%d", c.labelCount)
		c.labelCount++

		// Use 'f' format with precision 1 to ensure at least one decimal place
		// (omissions of this, where numbers are JUST integers,
		// cause issues in NASM,
		// resulting in 0.0 for calculations in the manner we perform)
		floatStr := strconv.FormatFloat(n.Value, 'f', 1, 64)

		// Add constant definition to the .data section buffer
		c.dataSection.WriteString(fmt.Sprintf("    %s: dq %s\n", label, floatStr))

		// Append load instruction to the .text section buffer
		// Update comment to use the same standard representation
		// when a value is pushed via `fld` all registers are shifted (st0 -> st1)
		// the latest is now the first in the stack
		c.textSection.WriteString(fmt.Sprintf("    fld QWORD [%s]  ; Load %s\n", label, floatStr))
		return nil

	case *ast.BinaryOpNode:
		// Right operand first
		err := c.generateExpr(n.Right)
		if err != nil {
			return err
		}

		// Left operand next
		err = c.generateExpr(n.Left)
		if err != nil {
			return err
		}
		// Emit the appropriate operation
		switch n.Op.Type {
		case token.PLUS:
			// add st0 to st1, store result in st1, and pop the register stack
			c.textSection.WriteString("    faddp st1, st0 ; Add\n")
		case token.MINUS:
			// subtract st1 from st0, store result in st1, and pop the register stack
			c.textSection.WriteString("    fsubrp st1, st0 ; Subtract\n")
		case token.MULTIPLY:
			// multiply st1 by st0 store result in st1, and pop the register stack
			c.textSection.WriteString("    fmulp st1, st0 ; Multiply\n")
		case token.DIVIDE:
			// divide st0 by st1, store result in st1, and pop the register stack
			c.textSection.WriteString("    fdivrp st1, st0 ; Divide\n")
		default:
			return fmt.Errorf("unknown binary operator: %s", n.Op.Value)
		}
		return nil

	case *ast.UnaryOpNode:
		err := c.generateExpr(n.Expr)
		if err != nil {
			return err
		}
		switch n.Op.Type {
		case token.PLUS:
			// Technicaly valid, but nothing to do; just a no-op.
			// We'll be nice and leave a comment
			c.textSection.WriteString("    ; Unary plus (no-op)\n")
		case token.MINUS:
			// complements sign of st0: Vol. 2A 3-375B
			c.textSection.WriteString("    fchs             ; Negate\n")
		default:
			return fmt.Errorf("unknown unary operator: %s", n.Op.Value)
		}
		return nil

	default:
		return fmt.Errorf("unknown node type for code generation: %T", node)
	}
}

// generatePrintFloat appends assembly instructions to call printf
// to print the float currently in st(0).
func (c *Compiler) generatePrintFloat() {
	// x86-64 SysV ABI requires float arguments in XMM registers for variadic functions like printf.
	// We need to store st0 to memory, then load it into xmm0.
	// Create a temporary memory location in the data section for this transfer.
	tempLabel := fmt.Sprintf("temp%d", c.labelCount)
	c.labelCount++

	// Reserve 8 bytes
	c.dataSection.WriteString(fmt.Sprintf("    %s: dq 0.0\n", tempLabel))

	c.textSection.WriteString("    ; Print value currently in st(0)\n")
	// copies value from st(0) to memory operand and pops the FPU stack see: (Vol. 2A 3-443)
	c.textSection.WriteString(fmt.Sprintf("    fstp QWORD [%s]   ; Store st(0) to memory and pop\n", tempLabel))
	// move the value from memory (tempLabel) to xmm0
	c.textSection.WriteString(fmt.Sprintf("    movsd xmm0, QWORD [%s] ; Load float from memory into xmm0\n", tempLabel))
	c.textSection.WriteString(fmt.Sprintf("    mov rdi, %s       ; 1st arg (format string)\n", c.floatFmt))
	// Required for variadic functions
	c.textSection.WriteString("    mov rax, 1          ; Number of XMM registers used (xmm0)\n")
	c.textSection.WriteString("    sub rsp, 8          ; Align stack pointer (16-byte boundary before call)\n")
	c.textSection.WriteString("    call printf         ; Call C printf function\n")
	c.textSection.WriteString("    add rsp, 8          ; Restore stack pointer\n")
}
