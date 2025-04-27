### Grammar

Basic arithmetic parser that respects operator precedence:
```
expr → term ((PLUS | MINUS) term)*
term → factor ((MUL | DIV) factor)*
factor → NUMBER | LPAREN expr RPAREN | (PLUS | MINUS) factor
```