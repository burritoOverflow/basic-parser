package parser

import (
	"basic-arithmetic-parser/ast"
	"basic-arithmetic-parser/lexer"
	"basic-arithmetic-parser/token"
	"fmt"
	"strconv"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
	}
	p.currentToken = p.lexer.GetNextToken()
	return p
}

func (p *Parser) eat(tokenType token.TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.GetNextToken()
	} else {
		// parse error -- TODO: collect errors instead of panic
		panic(fmt.Sprintf("Syntax error: expected %v, got %v", tokenType, p.currentToken.Type))
	}
}

// expr → term ((PLUS | MINUS) term)*
func (p *Parser) expr() ast.Node {
	node := p.term()

	for p.currentToken.Type == token.PLUS || p.currentToken.Type == token.MINUS {
		currTok := p.currentToken
		if currTok.Type == token.PLUS {
			p.eat(token.PLUS)
			node = &ast.BinaryOpNode{
				Left:  node,
				Op:    currTok,
				Right: p.term(),
			}
		} else if currTok.Type == token.MINUS {
			p.eat(token.MINUS)
			node = &ast.BinaryOpNode{
				Left:  node,
				Op:    currTok,
				Right: p.term(),
			}
		}
	}

	return node
}

// term → factor ((MUL | DIV) factor)*
func (p *Parser) term() ast.Node {
	node := p.factor()

	for p.currentToken.Type == token.MULTIPLY || p.currentToken.Type == token.DIVIDE {
		currTok := p.currentToken
		if currTok.Type == token.MULTIPLY {
			p.eat(token.MULTIPLY)
			node = &ast.BinaryOpNode{
				Left:  node,
				Op:    currTok,
				Right: p.factor(),
			}
		} else if currTok.Type == token.DIVIDE {
			p.eat(token.DIVIDE)
			node = &ast.BinaryOpNode{
				Left:  node,
				Op:    currTok,
				Right: p.factor(),
			}
		}
	}

	return node
}

// factor → NUMBER | LPAREN expr RPAREN | (PLUS | MINUS) factor
func (p *Parser) factor() ast.Node {
	currTok := p.currentToken

	switch currTok.Type {
	case token.NUMBER:
		p.eat(token.NUMBER)
		val, err := strconv.ParseFloat(currTok.Value, 64)
		if err != nil {
			panic(err)
		}
		return &ast.NumberNode{Value: val}
	case token.LPAREN:
		p.eat(token.LPAREN)
		node := p.expr()
		p.eat(token.RPAREN)
		return node
	case token.MINUS:
		p.eat(token.MINUS)
		return &ast.UnaryOpNode{
			Op:   token.Token{Type: token.MINUS, Value: "-"},
			Expr: p.factor(),
		}
	case token.PLUS:
		p.eat(token.PLUS)
		return &ast.UnaryOpNode{
			Op:   token.Token{Type: token.PLUS, Value: "+"},
			Expr: p.factor(),
		}
	default:
		// TODO as below, collect errors
		panic(fmt.Sprintf("Syntax error: unexpected token %v", currTok.Type))
	}
}

// Parse the input and return the AST
func (p *Parser) Parse() ast.Node {
	node := p.expr()
	// Check for trailing tokens--after a valid expression, we should only have EOF
	if p.currentToken.Type != token.EOF {
		// for now, failures result in a panic
		// TODO collect errors
		panic(fmt.Sprintf("Syntax error: unexpected token %v", p.currentToken.Type))
	}
	return node
}
