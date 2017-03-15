package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/ckirkendall/glisp/parser"
	"github.com/ckirkendall/glisp/interpreter"
)

func main() {
	const input = "(+ 1 2)(println \"test\" (* 2 3)) "
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
	for _, sexpr := range ast {
		val, err := interpreter.Eval(sexpr,env)
		fmt.Println(val, err)
	}



}
