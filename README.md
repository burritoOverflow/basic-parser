### Grammar

Basic arithmetic parser that respects operator precedence:

```
expr → term ((PLUS | MINUS) term)*
term → factor ((MUL | DIV) factor)*
factor → NUMBER | LPAREN expr RPAREN | (PLUS | MINUS) factor
```

Use like so:

```bash
./basic-arithmetic-parser -input example.txt -output output
```

Where `input` is a file containing basic, valid arithmetic expressions, like the following:

```
12 * 12
333 + 9
12 * 12
122 * 12
403 + 3
12 / 4
4 / 1
```

`output` is optional, and only used if a different name is desired for output artifacts (assembly, object file, executable).

Compilation has hard dependencies on `nasm` and `gcc`; failures to find these in `PATH` results in exiting with `1`.
