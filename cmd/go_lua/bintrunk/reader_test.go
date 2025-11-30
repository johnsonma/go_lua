package bintrunk

import (
	"testing"
)

func TestReaderReadByte(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	reader := &reader{data: data}

	// Test reading first byte
	b := reader.readByte()
	if b != 0x01 {
		t.Errorf("Expected 0x01, got %v", b)
	}

	// Test reading second byte
	b = reader.readByte()
	if b != 0x02 {
		t.Errorf("Expected 0x02, got %v", b)
	}

	// Test reading third byte
	b = reader.readByte()
	if b != 0x03 {
		t.Errorf("Expected 0x03, got %v", b)
	}

	// Verify data slice is empty
	if len(reader.data) != 0 {
		t.Errorf("Expected empty data slice, got %v bytes remaining", len(reader.data))
	}
}

func TestReaderReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	reader := &reader{data: data}

	// Test reading 3 bytes
	bytes := reader.readBytes(3)
	expected := []byte{0x01, 0x02, 0x03}
	if !equalBytes(bytes, expected) {
		t.Errorf("Expected %v, got %v", expected, bytes)
	}

	// Verify remaining data
	if len(reader.data) != 2 {
		t.Errorf("Expected 2 bytes remaining, got %v", len(reader.data))
	}
}

func TestReaderReadUint32(t *testing.T) {
	// Test with little endian data: 0x04030201
	data := []byte{0x01, 0x02, 0x03, 0x04}
	reader := &reader{data: data}

	result := reader.readUint32()
	expected := uint32(0x04030201)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Verify data slice is empty
	if len(reader.data) != 0 {
		t.Errorf("Expected empty data slice, got %v bytes remaining", len(reader.data))
	}
}

func TestReaderReadUint64(t *testing.T) {
	// Test with little endian data: 0x0807060504030201
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	reader := &reader{data: data}

	result := reader.readUint64()
	expected := uint64(0x0807060504030201)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Verify data slice is empty
	if len(reader.data) != 0 {
		t.Errorf("Expected empty data slice, got %v bytes remaining", len(reader.data))
	}
}

func TestReaderReadLuaInteger(t *testing.T) {
	// Test with little endian data: 0x0807060504030201
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	reader := &reader{data: data}

	result := reader.readLuaInteger()
	expected := int64(0x0807060504030201)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReaderReadLuaNumber(t *testing.T) {
	// Test with float64 value 1.5
	// 1.5 in IEEE 754 binary64: 0x3FF8000000000000
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF8, 0x3F}
	reader := &reader{data: data}

	result := reader.readLuaNumber()
	expected := 1.5
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReaderReadString(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "empty string",
			data:     []byte{0x00},
			expected: "",
		},
		{
			name:     "short string",
			data:     []byte{0x05, 'H', 'e', 'l', 'l', 'o'},
			expected: "Hello",
		},
		{
			name:     "long string with 0xFF size",
			data:     []byte{0xFF, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 'W', 'o', 'r', 'l', 'd'},
			expected: "World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &reader{data: tt.data}
			result := reader.readString()
			if result != tt.expected {
				t.Errorf("Expected '%v', got '%v'", tt.expected, result)
			}
		})
	}
}

func TestReaderCheckHeader(t *testing.T) {
	// Create valid header data
	data := []byte{
		// LUA_SIGNATURE (4 bytes)
		0x1b, 'L', 'u', 'a',
		// LUAC_VERSION
		0x54,
		// LUAC_FORMAT
		0x00,
		// LUAC_DATA (6 bytes)
		0x19, 0x93, 0x0d, 0x0a, 0x1a, 0x0a,
		// CINT_SIZE
		0x04,
		// CSIZET_SIZE
		0x08,
		// INSTRUCTION_SIZE
		0x04,
		// LUA_INTEGER_SIZE
		0x08,
		// LUA_NUMBER_SIZE
		0x08,
		// LUAC_INT (8 bytes little endian)
		0x78, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// LUAC_NUM (8 bytes little endian)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x73, 0x40,
	}

	reader := &reader{data: data}

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("checkHeader panicked unexpectedly: %v", r)
		}
	}()

	reader.checkHeader()
}

func TestReaderReadCode(t *testing.T) {
	// Test reading code array
	data := []byte{
		// Array length (4 bytes little endian)
		0x02, 0x00, 0x00, 0x00,
		// First instruction
		0x01, 0x00, 0x00, 0x00,
		// Second instruction
		0x02, 0x00, 0x00, 0x00,
	}
	reader := &reader{data: data}

	code := reader.readCode()
	expected := []uint32{1, 2}
	if len(code) != len(expected) {
		t.Errorf("Expected %v instructions, got %v", len(expected), len(code))
	}
	for i, v := range expected {
		if code[i] != v {
			t.Errorf("Expected code[%d] = %v, got %v", i, v, code[i])
		}
	}
}

func TestReaderReadConstants(t *testing.T) {
	// Test reading various constants
	data := []byte{
		// Array length (2 constants)
		0x02, 0x00, 0x00, 0x00,
		// First constant: nil
		TAG_NIL,
		// Second constant: boolean true
		TAG_BOOLEAN, 0x01,
	}
	reader := &reader{data: data}

	constants := reader.readConstants()
	if len(constants) != 2 {
		t.Errorf("Expected 2 constants, got %v", len(constants))
	}
	if constants[0] != nil {
		t.Errorf("Expected nil, got %v", constants[0])
	}
	if constants[1] != true {
		t.Errorf("Expected true, got %v", constants[1])
	}
}

func TestReaderReadUpvalues(t *testing.T) {
	// Test reading upvalues
	data := []byte{
		// Array length (2 upvalues)
		0x02, 0x00, 0x00, 0x00,
		// First upvalue
		0x01, 0x02,
		// Second upvalue
		0x03, 0x04,
	}
	reader := &reader{data: data}

	upvalues := reader.readUpvalues()
	if len(upvalues) != 2 {
		t.Errorf("Expected 2 upvalues, got %v", len(upvalues))
	}
	expected := []Upvalue{
		{Instack: 1, Idx: 2},
		{Instack: 3, Idx: 4},
	}
	for i, uv := range expected {
		if upvalues[i] != uv {
			t.Errorf("Expected upvalues[%d] = %v, got %v", i, uv, upvalues[i])
		}
	}
}

func TestReaderReadLineInfo(t *testing.T) {
	// Test reading line info
	data := []byte{
		// Array length (2 entries)
		0x02, 0x00, 0x00, 0x00,
		// First line
		0x0A, 0x00, 0x00, 0x00,
		// Second line
		0x0B, 0x00, 0x00, 0x00,
	}
	reader := &reader{data: data}

	lineInfo := reader.readLineInfo()
	expected := []uint32{10, 11}
	if len(lineInfo) != len(expected) {
		t.Errorf("Expected %v line info entries, got %v", len(expected), len(lineInfo))
	}
	for i, v := range expected {
		if lineInfo[i] != v {
			t.Errorf("Expected lineInfo[%d] = %v, got %v", i, v, lineInfo[i])
		}
	}
}

func TestReaderReadLocVars(t *testing.T) {
	// Test reading local variables
	data := []byte{
		// Array length (1 locvar)
		0x01, 0x00, 0x00, 0x00,
		// VarName (length 3 + "foo")
		0x03, 'f', 'o', 'o',
		// StartPC
		0x00, 0x00, 0x00, 0x00,
		// EndPC
		0x0A, 0x00, 0x00, 0x00,
	}
	reader := &reader{data: data}

	locVars := reader.readLocVars()
	if len(locVars) != 1 {
		t.Errorf("Expected 1 locvar, got %v", len(locVars))
	}
	expected := LocVar{
		VarName: "foo",
		StartPC: 0,
		EndPC:   10,
	}
	if locVars[0] != expected {
		t.Errorf("Expected %v, got %v", expected, locVars[0])
	}
}

func TestReaderReadUpvalueNames(t *testing.T) {
	// Test reading upvalue names
	data := []byte{
		// Array length (2 names)
		0x02, 0x00, 0x00, 0x00,
		// First name (length 3 + "up1")
		0x03, 'u', 'p', '1',
		// Second name (length 3 + "up2")
		0x03, 'u', 'p', '2',
	}
	reader := &reader{data: data}

	names := reader.readUpvalueNames()
	expected := []string{"up1", "up2"}
	if len(names) != len(expected) {
		t.Errorf("Expected %v names, got %v", len(expected), len(names))
	}
	for i, v := range expected {
		if names[i] != v {
			t.Errorf("Expected names[%d] = %v, got %v", i, v, names[i])
		}
	}
}

// Helper function to compare byte slices
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
