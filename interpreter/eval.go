package interpreter

import (
	"github.com/ckirkendall/glisp/parser"
	"fmt"
)

type Thunk struct {
	Env Environment
	Exp interface{}
}

type GLispError struct {
	text string
}

func (e GLispError) Error() string {
	return e.text
}

func unrollThunk(th Thunk) (interface{}, error){
	for {
		lth, err := EvalTh(th.Exp, th.Env, true)
		if err != nil {
			return nil, err
		}
		nth, ok := lth.(Thunk)
		if !ok {
			return lth, nil
		}
		th = nth
	}
}

func EvalTh(el interface{}, env Environment, returnThunk bool) (interface{}, error) {
	switch el.(type){
	case parser.Number, parser.String:
		return el, nil
	case parser.Bool:
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
			return nil, GLispError{"Problem calling:" + fmt.Sprint(lst) + ":" + err.Error()}
		}
		fn := val.(Sexp)
		res, terr := fn.Apply(env, lst.Val[1:]...)
		if terr != nil {
			return nil, terr
		}
		if returnThunk {
			return res, nil
		}
		th, ok := res.(Thunk)
		if ok {
			return unrollThunk(th)
		}else{
			return res, nil
		}
	}
	err := GLispError{"Unknown element: " + fmt.Sprint(el)}
	return nil, err
}

func Eval(el interface{}, env Environment) (interface{}, error) {
	return EvalTh(el,env,false)
}


