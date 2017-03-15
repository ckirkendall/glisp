package parser

import (
	"strconv"
	"fmt"
)

type Nill struct {}

type List struct {
	Val []interface{}
}

func (l List) String() string {
	return fmt.Sprint(l.Val)
}

type String struct {
	Val string
}

func (l String) String() string {
	return fmt.Sprint("\"", l.Val, "\"")
}

type Number struct {
	Val float64
}

func (l Number) String() string {
	return fmt.Sprint(l.Val)
}

type Ident struct {
	Val string
}

func (l Ident) String() string {
	return fmt.Sprint(l.Val)
}

func buildList(start int, tokens []MetaToken) (adv int, lst []interface{}){
	listToks := make([]MetaToken,0,0)
	numLeftParens := 0
	idx := start;
	for ; idx < len(tokens); idx++ {
		if tokens[idx].Tok == RIGHT_PAREN {
			if numLeftParens == 0 {
				break;
			}
			numLeftParens--
		}
		if tokens[idx].Tok == LEFT_PAREN {
			numLeftParens++
		}
		listToks = append(listToks, tokens[idx])
	}
	return idx, BuildAst(listToks)
}

func BuildAst(tokens []MetaToken) []interface{} {
	ast := make([]interface{}, 0, 0)
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Tok {
		case LEFT_PAREN:
			adv, lst := buildList(i + 1, tokens)
			i = adv
			ast = append(ast, List{ lst })
		case STRING:
			ast = append(ast, String{tokens[i].Lit})
		case NUMBER:
			nval, _ := strconv.ParseFloat(tokens[i].Lit, 64)
			ast = append(ast, Number{ nval })
		case IDENT:
			ast = append(ast, Ident{ tokens[i].Lit })

		}
	}
	return ast
}