### Grammar

Basic parser that respects opertor precedence:
```
expr → term ((PLUS | MINUS) term)*
term → factor ((MUL | DIV) factor)*
factor → NUMBER | LPAREN expr RPAREN | (PLUS | MINUS) factor
```