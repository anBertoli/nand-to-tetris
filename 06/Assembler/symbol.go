package Assembler


type SymbolTable struct {
	table     map[string]int
	VarMemIdx int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		VarMemIdx: 16,
		table: map[string]int{
			"R0": 0,
			"R1": 1,
			"R2": 2,
			"R3": 3,
			"R4": 4,
			"R5": 5,
			"R6": 6,
			"R7": 7,
			"R8": 8,
			"R9": 9,
			"R10": 10,
			"R11": 11,
			"R12": 12,
			"R13": 13,
			"R14": 14,
			"R15": 15,
			"SCREEN": 16384,
			"KBD": 24576,
			"SP": 0,
			"LCL": 1,
			"ARG": 2,
			"THIS": 3,
			"THAT": 4,
		},
	}
}


func (s *SymbolTable) AddEntry(symbol string, address int) {
	s.table[symbol] = address
}
func (s *SymbolTable) Contains(symbol string) bool {
	_, ok := s.table[symbol]
	return ok
}
func (s *SymbolTable) GetAddress(symbol string) int {
	if addr, ok := s.table[symbol]; ok {
		return addr
	}
	return -1
}