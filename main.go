package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/ckirkendall/glisp/parser"
	"github.com/ckirkendall/glisp/interpreter"
)

func main() {
	const input = "(+ 1 2)(println \"test\" (* 2 3)) " +
		"(if (empty? (list)) (println \"empty\") (println \"not empty\"))" +
		"(println (quote (list 1 2 3 4)))" +
		"(- 2 3)" +
		"(def inc (fn (num) (+ 1 num)))" +
		"(inc (inc (- 3 1)))" +
		"(def tmp (fn (x & y) (println \"T1:\" x) (println \"T2:\" y)))" +
		"(def even? (fn (num) (if (= num 0) true (odd? (- num 1)))))" +
		"(def odd? (fn (num) (if (= num 0) false (even? (- num 1)))))" +
		"(tmp 1 2 3 4 5 6)"
		//"(odd? 9999998)"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(parser.Tokenize)
	// Count the words.
	toks := make([]string, 0, 0)
	for scanner.Scan() {
		toks = append(toks, scanner.Text())
	}
	metaTokens := parser.BuildMetaTokens(toks)
	fmt.Println(metaTokens)
	ast := parser.BuildAst(metaTokens)
	fmt.Println(ast)
	env := interpreter.DefaultEnv()
	var res interface{}
	for _, sexpr := range ast {
		val, err := interpreter.Eval(sexpr, env)
		if err != nil {
			fmt.Println(err)
		}
		res = val
	}
	fmt.Println(res)

}
