package vm

import (
	"iscript/code"
	"iscript/object"
)

type Frame struct {
	fn *object.CompiledFunc
	ip int
}

func NewFrame(fn *object.CompiledFunc) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
