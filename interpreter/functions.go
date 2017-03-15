package interpreter

import (
	"github.com/ckirkendall/glisp/parser"
	"fmt"
)


type Sexp interface {
	Apply(env Environment, args ...interface{}) (interface{}, error)
}

type Def struct {}
type FnBuilder struct {}
type Fn struct {
	Env Environment
	Ident string
	Args []parser.Ident
	Body []interface{}
}

type MacroBuilder struct {}
type Macro struct {
	Env Environment
	Ident string
	Args []parser.Ident
	Body []interface{}
}

type Add struct {}
type Minus struct {}
type Div struct {}
type Mult struct {}
type Cons struct {}
type First struct {}
type Rest struct {}
type Print struct {}

func fnError(args []interface{}) error {
	return GLispError{"Invalid function def:" + fmt.Sprint(args)}
}

func wrongNumArgsError(call string) error {
	return GLispError{"Wrong number of args for call " + call}
}

func invalidArgError(call string) error {
	return GLispError{"Invalid arg error for call " + call}
}

func decomposeFn(args []interface{}) (*parser.Ident, []parser.Ident, []interface{}, error) {
	if len(args) < 2 {
		return nil, nil, nil, fnError(args)
	}
	sym, sok := args[0].(parser.Ident)
	if !sok {
		return nil, nil, nil, fnError(args)
	}
	argLst, aok := args[1].(parser.List)
	if !aok  {
		return nil, nil, nil, fnError(args)
	}
	identArgs := make([]parser.Ident, 0, 0)
	for i := 0; i < len(argLst.Val); i++ {
		sid, ok := argLst.Val[0].(parser.Ident)
		if !ok {
			return nil, nil, nil, fnError(args)
		}
		identArgs = append(identArgs, sid)
	}
	return &sym, identArgs, args[2:], nil
}

func (f Def) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, wrongNumArgsError("Def")
	}
	sym, sok := args[0].(parser.Ident)
	if !sok {
		return nil, invalidArgError("Def")
	}
	body, err := Eval(args[1], callerEnv)
	if err != nil {
		return nil, GLispError{"Unexpected error in def:" + sym.Val}
	}
	PutEnv(&callerEnv, sym.Val, body)
	return body, nil
}

func (f FnBuilder) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	sym, identArgs, body, err := decomposeFn(args)
	if err != nil {
		return nil, err
	}
	return Fn{callerEnv, sym.Val, identArgs, body}, nil

}

func (f Fn) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != len(f.Args) {
		err := wrongNumArgsError("Fn")
		return parser.Nill{}, err
	}
	nenv := Environment{ make(map[string]interface{}), &f.Env }
	for i := 0; i < len(args); i++ {
		val, e := Eval(args[i], nenv)
		if e != nil {
			err := GLispError{"Error evaluating args of " + f.Ident + " with " + fmt.Sprint(args) }
			return parser.Nill{}, err
		}
		PutEnv(&nenv, f.Args[i].Val, val)
	}
	var tail interface{}
	for _, el := range f.Body {
		res, e := Eval(el, nenv)
		if e != nil {
			err := GLispError{"Error evaluating body of " + f.Ident + " with " + fmt.Sprint(args) }
			return parser.Nill{}, err
		}
		tail = res
	}
	return tail, nil
}

func (f MacroBuilder) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	sym, identArgs, body, err := decomposeFn(args)
	if err != nil {
		return nil, err
	}
	return Macro{callerEnv, sym.Val, identArgs, body}, nil

}

func (f Macro) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != len(f.Args) {
		err := wrongNumArgsError("Macro")
		return parser.Nill{}, err
	}
	nenv := Environment{make(map[string]interface{}), &f.Env }
	for i := 0; i < len(args); i++ {
		PutEnv(&nenv, f.Args[i].Val, args[i])
	}
	var tail interface{}
	for _, el := range f.Body {
		res, e := Eval(el, nenv)
		if e != nil {
			err := GLispError{"Error evaluating body of " + f.Ident + " with " + fmt.Sprint(args) }
			return parser.Nill{}, err
		}
		tail = res
	}
	return Eval(tail, callerEnv)
}

func numArgs(args []interface{}) ([]parser.Number, error){
	nums := make([]parser.Number, 0, 0)
	for _, arg := range args {
		num, ok := arg.(parser.Number)
		if !ok {
			return nil, invalidArgError("Add")
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func (f Add) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	nums, err := numArgs(args)
	if err != nil {
		return nil, err
	}
	var res float64 = 0
	for i := 0; i < len(nums); i++ {
		res += nums[i].Val
	}
	return parser.Number{ res }, nil
}

func (f Minus) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	nums, err := numArgs(args)
	if err != nil {
		return nil, err
	}
	var res float64 = 0
	for i := 0; i < len(nums); i++ {
		res -= nums[i].Val
	}
	return parser.Number{ res }, nil
}

func (f Mult) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	nums, err := numArgs(args)
	if err != nil {
		return nil, err
	}
	if len(nums) == 1 {
		return nums[0], nil
	}
	res := nums[0].Val
	for i := 1; i < len(nums); i++ {
		res *= nums[i].Val
	}
	return parser.Number{ res }, nil
}

func (f Div) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	nums, err := numArgs(args)
	if err != nil {
		return nil, invalidArgError("Div")
	}
	if len(nums) == 1 {
		return nums[0], nil
	}
	res := nums[0].Val
	for i := 1; i < len(nums); i++ {
		res /= nums[i].Val
	}
	return parser.Number{ res }, nil
}

func (f Cons) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, wrongNumArgsError("Cons")
	}
	el := args[0]
	lst, ok := args[1].(parser.List)
	if !ok {
		return nil, invalidArgError("Cons")
	}

	return parser.List{ append([]interface{}{el}, lst.Val)}, nil
}

func (f First) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, wrongNumArgsError("First")
	}
	lst, ok := args[0].(parser.List)
	if ok {
		return nil, invalidArgError("First")
	}
	if len(lst.Val) == 0 {
		return parser.Nill{}, nil
	}
	return lst.Val[0], nil
}

func (f Rest) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, wrongNumArgsError("Rest")
	}
	lst, ok := args[0].(parser.List)
	if !ok {
		return nil, invalidArgError("Rest")
	}
	if len(lst.Val) == 0 {
		return lst, nil
	}
	return parser.List{ lst.Val[1:] }, nil
}

func (f Print) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	var first = true
	for _, arg := range args {
		if !first {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	return nil, nil
}