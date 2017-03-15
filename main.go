package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/ckirkendall/glisp/parser"
	"github.com/ckirkendall/glisp/interpreter"
)

func main() {
	const input = "(+ 1 2)"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(parser.Tokenize)
	// Count the words.
	toks := make([]string,0, 0)
	for scanner.Scan() {
		toks = append(toks, scanner.Text())
	}
	metaTokens := parser.BuildMetaTokens(toks)
	fmt.Println(metaTokens)
	ast := parser.BuildAst(metaTokens)
	fmt.Println(ast)
	env := interpreter.DefaultEnv()
	val, err := interpreter.Eval(ast[0],env)
	fmt.Println(val, err)

}
