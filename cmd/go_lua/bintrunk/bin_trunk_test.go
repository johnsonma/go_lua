package bintrunk

import (
	"testing"
)

func TestUndump(t *testing.T) {
	// Create a minimal valid Lua bytecode chunk
	data := []byte{
		// Header
		// LUA_SIGNATURE (4 bytes) - "x1bLua" as bytes
		0x1B, 0x4C, 0x75, 0x61,
		// LUAC_VERSION
		LUAC_VERSION,
		// LUAC_FORMAT
		LUAC_FORMAT,
		// LUAC_DATA (6 bytes)
		0x19, 0x93, 0x0d, 0x0a, 0x1a, 0x0a,
		// CINT_SIZE
		CINT_SIZE,
		// CSIZET_SIZE
		CSIZET_SIZE,
		// INSTRUCTION_SIZE
		INSTRUCTION_SIZE,
		// LUA_INTEGER_SIZE
		LUA_INTEGER_SIZE,
		// LUA_NUMBER_SIZE
		LUA_NUMBER_SIZE,
		// LUAC_INT (8 bytes little endian)
		0x78, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// LUAC_NUM (8 bytes little endian) - 370.5 as float64
		0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x77, 0x40,
		// sizeUpvalues
		0x00,
		// Prototype
		// Source (empty string)
		0x00,
		// LineDefined
		0x00, 0x00, 0x00, 0x00,
		// LastLineDefined
		0x00, 0x00, 0x00, 0x00,
		// NumParams
		0x00,
		// IsVarag
		0x00,
		// MaxStackSize
		0x00,
		// Code (empty)
		0x00, 0x00, 0x00, 0x00,
		// Constants (empty)
		0x00, 0x00, 0x00, 0x00,
		// Upvalues (empty)
		0x00, 0x00, 0x00, 0x00,
		// Protos (empty)
		0x00, 0x00, 0x00, 0x00,
		// LineInfo (empty)
		0x00, 0x00, 0x00, 0x00,
		// LocVars (empty)
		0x00, 0x00, 0x00, 0x00,
		// UpvalueNames (empty)
		0x00, 0x00, 0x00, 0x00,
	}

	proto := Undump(data)
	if proto == nil {
		t.Error("Expected non-nil prototype, got nil")
	}

	// Verify basic prototype structure
	if proto.Source != "" {
		t.Errorf("Expected empty source, got '%v'", proto.Source)
	}
	if proto.LineDefined != 0 {
		t.Errorf("Expected LineDefined = 0, got %v", proto.LineDefined)
	}
	if proto.LastLineDefined != 0 {
		t.Errorf("Expected LastLineDefined = 0, got %v", proto.LastLineDefined)
	}
	if proto.NumParams != 0 {
		t.Errorf("Expected NumParams = 0, got %v", proto.NumParams)
	}
	if proto.IsVarag != 0 {
		t.Errorf("Expected IsVarag = 0, got %v", proto.IsVarag)
	}
	if proto.MaxStackSize != 0 {
		t.Errorf("Expected MaxStackSize = 0, got %v", proto.MaxStackSize)
	}
	if len(proto.Code) != 0 {
		t.Errorf("Expected empty code, got %v instructions", len(proto.Code))
	}
	if len(proto.Constants) != 0 {
		t.Errorf("Expected empty constants, got %v constants", len(proto.Constants))
	}
	if len(proto.Upvalues) != 0 {
		t.Errorf("Expected empty upvalues, got %v upvalues", len(proto.Upvalues))
	}
	if len(proto.Protos) != 0 {
		t.Errorf("Expected empty protos, got %v protos", len(proto.Protos))
	}
	if len(proto.LineInfo) != 0 {
		t.Errorf("Expected empty lineInfo, got %v entries", len(proto.LineInfo))
	}
	if len(proto.LocVars) != 0 {
		t.Errorf("Expected empty locVars, got %v entries", len(proto.LocVars))
	}
	if len(proto.UpvalueNames) != 0 {
		t.Errorf("Expected empty upvalueNames, got %v entries", len(proto.UpvalueNames))
	}
}

func TestUndumpInvalidHeader(t *testing.T) {
	// Test with invalid signature
	data := []byte{
		// Invalid signature
		'I', 'n', 'v', 'l',
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid header, but no panic occurred")
		}
	}()

	Undump(data)
}

func TestPrototypeEquality(t *testing.T) {
	// Test that Prototype structs can be compared
	p1 := &Prototype{
		Source:          "test",
		LineDefined:     1,
		LastLineDefined: 2,
		NumParams:       0,
		IsVarag:         0,
		MaxStackSize:    1,
		Code:            []uint32{1, 2, 3},
		Constants:       []interface{}{nil, true, false},
		Upvalues:        []Upvalue{{Instack: 1, Idx: 2}},
		Protos:          []*Prototype{},
		LineInfo:        []uint32{1, 2},
		LocVars:         []LocVar{{VarName: "x", StartPC: 0, EndPC: 10}},
		UpvalueNames:    []string{"up1"},
	}

	p2 := &Prototype{
		Source:          "test",
		LineDefined:     1,
		LastLineDefined: 2,
		NumParams:       0,
		IsVarag:         0,
		MaxStackSize:    1,
		Code:            []uint32{1, 2, 3},
		Constants:       []interface{}{nil, true, false},
		Upvalues:        []Upvalue{{Instack: 1, Idx: 2}},
		Protos:          []*Prototype{},
		LineInfo:        []uint32{1, 2},
		LocVars:         []LocVar{{VarName: "x", StartPC: 0, EndPC: 10}},
		UpvalueNames:    []string{"up1"},
	}

	// These should be equal in terms of their field values
	if p1.Source != p2.Source {
		t.Error("Prototypes should have equal Source fields")
	}
	if p1.LineDefined != p2.LineDefined {
		t.Error("Prototypes should have equal LineDefined fields")
	}
	if p1.LastLineDefined != p2.LastLineDefined {
		t.Error("Prototypes should have equal LastLineDefined fields")
	}
	if p1.NumParams != p2.NumParams {
		t.Error("Prototypes should have equal NumParams fields")
	}
	if p1.IsVarag != p2.IsVarag {
		t.Error("Prototypes should have equal IsVarag fields")
	}
	if p1.MaxStackSize != p2.MaxStackSize {
		t.Error("Prototypes should have equal MaxStackSize fields")
	}
	if len(p1.Code) != len(p2.Code) {
		t.Error("Prototypes should have equal Code lengths")
	}
	if len(p1.Constants) != len(p2.Constants) {
		t.Error("Prototypes should have equal Constants lengths")
	}
	if len(p1.Upvalues) != len(p2.Upvalues) {
		t.Error("Prototypes should have equal Upvalues lengths")
	}
	if len(p1.Protos) != len(p2.Protos) {
		t.Error("Prototypes should have equal Protos lengths")
	}
	if len(p1.LineInfo) != len(p2.LineInfo) {
		t.Error("Prototypes should have equal LineInfo lengths")
	}
	if len(p1.LocVars) != len(p2.LocVars) {
		t.Error("Prototypes should have equal LocVars lengths")
	}
	if len(p1.UpvalueNames) != len(p2.UpvalueNames) {
		t.Error("Prototypes should have equal UpvalueNames lengths")
	}
}

func TestUpvalueEquality(t *testing.T) {
	uv1 := Upvalue{Instack: 1, Idx: 2}
	uv2 := Upvalue{Instack: 1, Idx: 2}
	uv3 := Upvalue{Instack: 3, Idx: 4}

	if uv1 != uv2 {
		t.Error("Upvalues with same values should be equal")
	}
	if uv1 == uv3 {
		t.Error("Upvalues with different values should not be equal")
	}
}

func TestLocVarEquality(t *testing.T) {
	lv1 := LocVar{VarName: "x", StartPC: 0, EndPC: 10}
	lv2 := LocVar{VarName: "x", StartPC: 0, EndPC: 10}
	lv3 := LocVar{VarName: "y", StartPC: 0, EndPC: 10}

	if lv1 != lv2 {
		t.Error("LocVars with same values should be equal")
	}
	if lv1 == lv3 {
		t.Error("LocVars with different values should not be equal")
	}
}
