package compiler

type SymScope string

const (
	GlobalScope   SymScope = "GLOBAL"
	LocalScope    SymScope = "LOCAL"
	BuiltinScope  SymScope = "BUILTIN"
	FreeScope     SymScope = "FREE"
	FunctionScope SymScope = "FUNCTION"
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
	FreeSyms       []Sym
}

func NewSymTable() *SymTable {
	s := make(map[string]Sym)
	free := []Sym{}
	return &SymTable{store: s, FreeSyms: free}
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
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}

	return obj, ok
}

func NewEnclosedSymTable(outer *SymTable) *SymTable {
	s := NewSymTable()
	s.Outer = outer
	return s
}

func (s *SymTable) DefineBuiltin(index int, name string) Sym {
	sym := Sym{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = sym
	return sym
}

func (s *SymTable) defineFree(orig Sym) Sym {
	s.FreeSyms = append(s.FreeSyms, orig)

	sym := Sym{Name: orig.Name, Index: len(s.FreeSyms) - 1}
	sym.Scope = FreeScope

	s.store[orig.Name] = sym
	return sym
}

func (s *SymTable) DefineFunctionName(name string) Sym {
	sym := Sym{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = sym
	return sym
}
