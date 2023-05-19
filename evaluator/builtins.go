package evaluator

import (
	"iscript/object"
)

var builtins = map[string]*object.Builtin{
	"len":  object.GetBuiltinByName("len"),
	"puts": object.GetBuiltinByName("puts"),
}
