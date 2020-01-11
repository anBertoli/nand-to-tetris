package writer

type symbolType string
const (
	INT     	symbolType = "int"
	CHAR    	symbolType = "char"
	BOOLEAN 	symbolType = "boolean"
	CLASS   	symbolType = "class"
)

type symbolKind string
const (
	FIELD    	symbolKind = "this"
	STATIC   	symbolKind = "static"
	LOCAL    	symbolKind = "local"
	ARGUMENT 	symbolKind = "argument"
)

type symbolScope string
const (
	CLASS_SCOPE			symbolScope = "class"
	SUBROUTINE_SCOPE 	symbolScope = "subroutine"
)


type symbol struct {
	name  string
	index int
	typ   symbolType
	kind  symbolKind
	scope symbolScope
}
type SymbolTable struct {
	table []symbol
}

func (s* SymbolTable) cleanSubroutineScope() {
	var newSymbols []symbol
	for _, sym := range s.table {
		if sym.scope == CLASS_SCOPE {
			newSymbols = append(newSymbols, sym)
		}
	}
	s.table = newSymbols
}

func (s* SymbolTable) lastOfKind(kind symbolKind) int {
	lastOfKind := -1
	for _, s := range s.table {
		if s.kind == kind && s.index > lastOfKind {
			lastOfKind = s.index
		}
	}
	return lastOfKind
}

func (s* SymbolTable) exists(name string) bool {
	for _, s := range s.table {
		if s.name == name {
			return true
		}
	}
	return false
}

func (s* SymbolTable) define(name string, typ symbolType, kind symbolKind, scope symbolScope) {
	index := s.lastOfKind(kind)
	index++

	s.table = append(s.table, symbol{
		name: name,
		index: index,
		typ: typ,
		kind: kind,
		scope: scope,
	})
}



// following methods check first for subroutine-level
// symbols, then for class-level ones
func (s* SymbolTable) kindOf(name string) symbolKind {
	for _, s := range s.table {
		if s.name == name && s.scope == SUBROUTINE_SCOPE {
			return s.kind
		}
	}
	for _, s := range s.table {
		if s.name == name && s.scope == CLASS_SCOPE{
			return s.kind
		}
	}
	return ""
}


func (s* SymbolTable) typeOf(name string) symbolType {
	for _, s := range s.table {
		if s.name == name && s.scope == CLASS_SCOPE {
			return s.typ
		}
	}
	for _, s := range s.table {
		if s.name == name && s.scope == SUBROUTINE_SCOPE {
			return s.typ
		}
	}
	return ""
}


func (s* SymbolTable) indexOf(name string) int {
	for _, s := range s.table {
		if s.name == name && s.scope == CLASS_SCOPE {
			return s.index
		}
	}
	for _, s := range s.table {
		if s.name == name && s.scope == SUBROUTINE_SCOPE {
			return s.index
		}
	}
	return -1
}

