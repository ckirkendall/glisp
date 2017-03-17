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
	Args []parser.Ident
	Body []interface{}
}

func (l Fn) String() string {
	return fmt.Sprint("(fn ...)")
}

type MacroBuilder struct {}
type Macro struct {
	Env Environment
	Args []parser.Ident
	Body []interface{}
}

type If struct {}
type Equal struct {}
type Not struct {}
type Quote struct {}

type Add struct {}
type Minus struct {}
type Div struct {}
type Mult struct {}

type List struct {}
type Cons struct {}
type First struct {}
type Rest struct {}
type Empty struct {}

type Print struct {}
type PrintLn struct {}

func fnError(args []interface{}) error {
	return GLispError{"Invalid function def:" + fmt.Sprint(args)}
}

func wrongNumArgsError(call string) error {
	return GLispError{"Wrong number of args for call " + call}
}

func invalidArgError(call string) error {
	return GLispError{"Invalid arg error for call " + call}
}

func evalArgs(args []interface{}, env Environment) ([]interface{}, error){
	vals := make([]interface{}, 0, 0)
	for _, arg := range args {
		val, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func decomposeFn(args []interface{}) ([]parser.Ident, []interface{}, error) {
	if len(args) < 2 {
		return nil, nil, fnError(args)
	}
	argLst, aok := args[0].(parser.List)
	if !aok  {
		return nil, nil, fnError(args)
	}
	identArgs := make([]parser.Ident, 0, 0)
	for i := 0; i < len(argLst.Val); i++ {
		sid, ok := argLst.Val[0].(parser.Ident)
		if !ok {
			return nil, nil, fnError(args)
		}
		identArgs = append(identArgs, sid)
	}
	return identArgs, args[1:], nil
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
		return nil, GLispError{"Unexpected error in def:" + sym.Val + " : " + err.Error()}
	}
	PutEnv(&callerEnv, sym.Val, body)
	return body, nil
}

func (f FnBuilder) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	identArgs, body, err := decomposeFn(args)
	if err != nil {
		return nil, err
	}
	return Fn{callerEnv, identArgs, body}, nil

}

func (f Fn) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	if len(args) != len(f.Args) {
		err := wrongNumArgsError("Fn")
		return parser.Nill{}, err
	}
	nenv := Environment{ make(map[string]interface{}), &f.Env }
	for i := 0; i < len(args); i++ {
		val, e := Eval(args[i], callerEnv)
		if e != nil {
			err := GLispError{"Error evaluating args " + fmt.Sprint(args) + ":" + e.Error()}
			return parser.Nill{}, err
		}
		PutEnv(&nenv, f.Args[i].Val, val)
	}
	for i := 0; i<len(f.Body); i++ {
		if i == (len(f.Body) - 1) {
			return Thunk{nenv,f.Body[i]}, nil
		}
		_, e := Eval(f.Body[i], nenv)
		if e != nil {
			err := GLispError{"Error evaluating body of " + fmt.Sprint(args) }
			return parser.Nill{}, err
		}
	}
	return nil, GLispError{"EEEK how did we get here!"}
}

func (f MacroBuilder) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	identArgs, body, err := decomposeFn(args)
	if err != nil {
		return nil, err
	}
	return Macro{callerEnv, identArgs, body}, nil

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
			err := GLispError{"Error evaluating body of " + fmt.Sprint(args) }
			return parser.Nill{}, err
		}
		tail = res
	}
	return Eval(tail, callerEnv)
}

func (f If) Apply(callerEnv Environment, args...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, wrongNumArgsError("IF")
	}
	test := false
	testSexp, err := Eval(args[0], callerEnv)
	if err != nil {
		return nil, err
	}
	tbool, ok := testSexp.(parser.Bool)
	if !ok {
		_, nok := testSexp.(parser.Nill)
		if !nok {
			test = true
		}
	}else{
		test = tbool.Val
	}
	if test {
		return Thunk{callerEnv, args[1]}, nil
	}
	if len(args) > 2 {
		return Thunk{callerEnv, args[2]}, nil
	}
	return parser.Nill{}, nil
}

func (f Quote) Apply(callerEnv Environment, args...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, wrongNumArgsError("Quote")
	}
	return args[0], nil
}

func (f Equal) Apply(callerEnv Environment, args...interface{}) (interface{}, error) {
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
	if len(args) <= 1 {
		return nil, wrongNumArgsError("Equal")
	}
	for i := 1; i< len(args); i++ {
		if args[i-1] != args[i] {
			return parser.Bool{false}, nil
		}
	}
	return parser.Bool{true}, nil
}

func (f Not) Apply(callerEnv Environment, args...interface{}) (interface{}, error) {
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
	if len(args) != 1 {
		return nil, wrongNumArgsError("Not")
	}
	val, ok := args[0].(parser.Bool)
	if ok {
		return parser.Bool{!val.Val}, nil
	}
	_, nok := args[0].(parser.Nill)
	return parser.Bool{nok}, nil
}

func numArgs(args []interface{}, callerEnv Environment) ([]parser.Number, error){
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
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
	nums, err := numArgs(args, callerEnv)
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
	nums, err := numArgs(args, callerEnv)
	if err != nil {
		return nil, err
	}
	if len(args) < 1 {
		return nil, wrongNumArgsError("Minus")
	}
	var res float64 = nums[0].Val
	for i := 1; i < len(nums); i++ {
		res -= nums[i].Val
	}
	return parser.Number{ res }, nil
}

func (f Mult) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	nums, err := numArgs(args, callerEnv)
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
	nums, err := numArgs(args, callerEnv)
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
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
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
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
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
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
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

func (f Empty) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
	if len(args) != 1 {
		return nil, wrongNumArgsError("Empty")
	}
	lst, ok := args[0].(parser.List)
	if !ok {
		return nil, invalidArgError("Empty")
	}
	return parser.Bool{len(lst.Val) == 0}, nil
}

func (f List) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	return parser.List{args}, nil
}

func (f Print) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
	fmt.Print(args...)
	return parser.Nill{}, nil
}

func (f PrintLn) Apply(callerEnv Environment, args ...interface{}) (interface{}, error) {
	args, aerr := evalArgs(args, callerEnv)
	if aerr != nil {
		return nil, aerr
	}
	fmt.Println(args...)
	return parser.Nill{}, nil
}