package compiler

type SymScope string

const (
	GlobalScope SymScope = "GLOBAL"
	LocalScope  SymScope = "LOCAL"
)

type Sym struct {
	Name  string
	Scope SymScope
	Index int
}

type SymTable struct {
	Outer *SymTable

	store          map[string]Sym
	numDefinitions int
}

func NewSymTable() *SymTable {
	s := make(map[string]Sym)
	return &SymTable{store: s}
}

func (s *SymTable) Define(name string) Sym {
	sym := Sym{Name: name, Index: s.numDefinitions}
	sym.Scope = LocalScope
	if s.Outer == nil {
		sym.Scope = GlobalScope
	}
	s.store[name] = sym
	s.numDefinitions++
	return sym
}

func (s *SymTable) Resolve(name string) (Sym, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}

func NewEnclosedSymTable(outer *SymTable) *SymTable {
	s := NewSymTable()
	s.Outer = outer
	return s
}
