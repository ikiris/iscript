package compiler

type SymScope string

const (
	GlobalScope SymScope = "GLOBAL"
)

type Sym struct {
	Name  string
	Scope SymScope
	Index int
}

type SymTable struct {
	store          map[string]Sym
	numDefinitions int
}

func NewSymTable() *SymTable {
	s := make(map[string]Sym)
	return &SymTable{store: s}
}

func (s *SymTable) Define(name string) Sym {
	sym := Sym{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = sym
	s.numDefinitions++
	return sym
}

func (s *SymTable) Resolve(name string) (Sym, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
