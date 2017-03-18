package main

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/ckirkendall/glisp/parser"
	"github.com/ckirkendall/glisp/interpreter"
)

func main() {
	const let = "(def let* (macro (plst & body)" +
		"(cons (cons (quote fn)" +
		"(cons (cons (first plst) (list))" +
		"body))" +
		"(rest plst))))" +
		"(def let (macro (plst & body)" +
		"(if (empty? plst)" +
		"(cons (cons (quote fn) (cons (list) body)) (list))" +
		"(cons (quote let*)" +
		"(cons (cons (first plst)" +
		"(cons (first (rest plst)) (list)))" +
		"(cons (cons (quote let)" +
		"(cons (rest (rest plst))" +
		"body))" +
		"(list)))))))"

	const _map = "(def map (fn (f c)" +
		"(if (empty? c)" +
		"(list)" +
		"(cons (f (first c)) (map f (rest c))))))"

	const reduce = "(def reduce (fn (f a b)" +
		"(if (empty? b)" +
		" a " +
		"(reduce f (f a (first b)) (rest b)))))"

	const input = let + _map + reduce + "(+ 1 2)(println \"test\" (* 2 3)) " +
		"(if (empty? (list)) (println \"empty\") (println \"not empty\"))" +
		"(println (quote (list 1 2 3 4)))" +
		"(- 2 3)" +
		"(def inc (fn (num) (+ 1 num)))" +
		"(inc (inc (- 3 1)))" +
		"(def tmp (fn (x & y) (println \"T1:\" x) (println \"T2:\" y)))" +
		"(def even? (fn (num) (if (= num 0) true (odd? (- num 1)))))" +
		"(def odd? (fn (num) (if (= num 0) false (even? (- num 1)))))" +
		"(tmp 1 2 3 4 5 6)" +
		"(let (x 3 y 4) (+ x y))" +
		"(let (x 1 y (quote (1 2 3)) w (map (fn (z) (+ z x)) y))" +
		"(println x)" +
		"(println y)" +
		"(println (reduce (fn (a b) (+ a b)) 1 w)))" +
		"(odd? 9998)"
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
