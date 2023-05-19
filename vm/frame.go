package vm

import (
	"iscript/code"
	"iscript/object"
)

type Frame struct {
	fn      *object.CompiledFunc
	ip      int
	basePtr int
}

func NewFrame(fn *object.CompiledFunc, basePtr int) *Frame {
	return &Frame{
		fn:      fn,
		ip:      -1,
		basePtr: basePtr,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
