package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of args: got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got=%s", args[0].Type())
			}
		}},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		}},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of args: got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return nil
		}},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of args: got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got=%s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return nil
		}},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of args: got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length == 0 {
				return nil
			}
			newElm := make([]Object, length-1)
			copy(newElm, arr.Elements[1:length])
			return &Array{Elements: newElm}
		}},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of args: got=%d, want=2", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got=%s", args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)

			newElm := make([]Object, length+1)
			copy(newElm, arr.Elements)
			newElm[length] = args[1]

			return &Array{Elements: newElm}
		}},
	},
	{
		"updateHash",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError("wrong number of args: got=%d, want=2", len(args))
			}
			if args[0].Type() != HASH_OBJ {
				return newError("argument 1 to `updateHash` must be HASH, got=%s", args[0].Type())
			}
			i, ok := args[1].(Hashable)
			if !ok {
				return newError("argument 2 to `updateHash` must be HASHKEY, got=%T", args[1].Type())
			}

			hk := i.HashKey()
			h := args[0].(*Hash)
			h.Pairs[hk] = HashPair{Key: args[1], Value: args[2]}

			return nil
		}},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
