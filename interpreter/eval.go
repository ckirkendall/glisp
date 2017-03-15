package interpreter

import (
	"github.com/ckirkendall/glisp/parser"
	"fmt"
)


type GLispError struct {
	text string
}

func (e GLispError) Error() string {
	return e.text
}

func Eval(el interface{}, env Environment) (interface{}, error) {
	switch el.(type){
	case parser.Number, parser.String:
		return el, nil
	case parser.Ident:
		sym, _ := el.(parser.Ident)
		return LookUp(sym.Val, &env)
	case parser.List:
		lst, _ := el.(parser.List)
		if len(lst.Val) == 0 {
			return lst, nil
		}
		val, err := Eval(lst.Val[0],env)
		if err != nil {
			return nil, GLispError{"Problem calling:" + fmt.Sprint(val) }
		}
		fn := val.(Sexp)
		return fn.Apply(env, lst.Val[1:]...)
	}
	err := GLispError{"Unknown element: " + fmt.Sprint(el)}
	return nil, err
}