package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/ckirkendall/glisp/parser"
)

func main() {
	const input = "(fn (test test2) (println \"fun\") (+ 1 1.02))"
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

}
