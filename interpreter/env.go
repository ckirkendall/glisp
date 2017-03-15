package interpreter

import (
	"github.com/ckirkendall/glisp/parser"
)

type Environment struct {
	Vars map[string]interface{}
	Parent *Environment
}

func PutEnv(env *Environment, ident string, val interface{}){
	env.Vars[ident] = val
}

func LookUp(ident string, env *Environment) (interface{}, error) {
	if env == nil {
		err := GLispError{ "Invalid identifier: " + ident }
		return parser.Nill{}, err
	}

	val := env.Vars[ident]
	if val != nil {
		return val, nil
	}
	return LookUp(ident, env.Parent)
}

func DefaultEnv() Environment {
	vars := make(map[string]interface{})
	env := Environment{ vars, nil }
	PutEnv(&env, "def", Def{})
	PutEnv(&env, "fn", FnBuilder{})
	PutEnv(&env, "macro", MacroBuilder{})
	PutEnv(&env, "+", Add{})
	PutEnv(&env, "-", Minus{})
	PutEnv(&env, "*", Mult{})
	PutEnv(&env, "/", Div{})
	PutEnv(&env, "cons", Cons{})
	PutEnv(&env, "first", First{})
	PutEnv(&env, "rest", Rest{})
	PutEnv(&env, "print", Print{})
	PutEnv(&env, "println", PrintLn{})
	return env
}