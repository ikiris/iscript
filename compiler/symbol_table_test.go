package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Sym{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Scope: LocalScope, Index: 0},
		"d": {Name: "d", Scope: LocalScope, Index: 1},
		"e": {Name: "e", Scope: LocalScope, Index: 0},
		"f": {Name: "f", Scope: LocalScope, Index: 1},
	}

	global := NewSymTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}

	firstLocal := NewEnclosedSymTable(global)

	c := firstLocal.Define("c")
	if c != expected["c"] {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}

	d := firstLocal.Define("d")
	if d != expected["d"] {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}

	secondLocal := NewEnclosedSymTable(firstLocal)
	e := secondLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
	}

	f := secondLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("expected f=%+v, got=%+v", expected["f"], d)
	}

}

func TestResolveGlobal(t *testing.T) {
	global := NewSymTable()
	global.Define("a")
	global.Define("b")

	expected := map[string]Sym{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Sym{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s is not resolvable.", sym.Name)
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymTable(global)
	local.Define("c")
	local.Define("d")

	secondLocal := NewEnclosedSymTable(local)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table        *SymTable
		expectedSyms []Sym
	}{
		{local,
			[]Sym{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			}},
		{secondLocal,
			[]Sym{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			}},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSyms {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s is not resolvable.", sym.Name)
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymTable()
	firstLocal := NewEnclosedSymTable(global)
	secondLocal := NewEnclosedSymTable(firstLocal)

	expected := []Sym{
		{Name: "a", Scope: BuiltinScope, Index: 0},
		{Name: "c", Scope: BuiltinScope, Index: 1},
		{Name: "e", Scope: BuiltinScope, Index: 2},
		{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			res, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if res != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, res)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table           *SymTable
		expectedSyms    []Sym
		expectedFreeSym []Sym
	}{
		{
			firstLocal,
			[]Sym{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
			[]Sym{},
		},
		{
			secondLocal,
			[]Sym{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
			[]Sym{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSyms {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
		if len(tt.table.FreeSyms) != len(tt.expectedFreeSym) {
			t.Errorf("wrong number of free symbols: got=%d, want=%d", len(tt.table.FreeSyms), len(tt.expectedFreeSym))
			continue
		}
		for i, sym := range tt.expectedFreeSym {
			result := tt.table.FreeSyms[i]
			if result != sym {
				t.Errorf("wrong free symbol: got=%+v, want=%+v", result, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymTable()
	global.Define("a")

	firstLocal := NewEnclosedSymTable(global)
	firstLocal.Define("c")

	secondLocal := NewEnclosedSymTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	expected := []Sym{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 0},
		{Name: "f", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s is not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}

	expectedUnresolvable := []string{"b", "d"}

	for _, name := range expectedUnresolvable {
		_, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestDefineResolveFuncName(t *testing.T) {
	global := NewSymTable()
	global.DefineFunctionName("a")

	expected := Sym{Name: "a", Scope: FunctionScope, Index: 0}

	res, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("func name %s not resolvable.", expected.Name)
	}

	if res != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, res)
	}
}

func TestShadowingFuncName(t *testing.T) {
	global := NewSymTable()
	global.DefineFunctionName("a")
	global.Define("a")

	expected := Sym{Name: "a", Scope: GlobalScope, Index: 0}

	res, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("func name %s not resolvable.", expected.Name)
	}

	if res != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v", expected.Name, expected, res)
	}
}
