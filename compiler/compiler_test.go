package compiler

import (
	"code"
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{int64(1), int64(2)},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
			},
		},
	}

	runCompilerTests(t, tests)
}

func parse(input string) ast.Node {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}
		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func concatInstructions(input []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, instr := range input {
		out = append(out, instr...)
	}
	return out
}

func testInstructions(want []code.Instructions, got code.Instructions) error {
	concatted := concatInstructions(want)

	if len(concatted) != len(got) {
		return fmt.Errorf("wrong length of instructions.  Wanted %q, got %q", want, got)
	}
	for i, val := range concatted {
		if val != got[i] {
			return fmt.Errorf("wrong instruction value at %d.  Wanted %q, got %q", i, val, got[i])
		}
	}
	return nil
}

func testConstants(expected []any, got []object.Object) error {
	if len(expected) != len(got) {
		return fmt.Errorf("wrong number of constants.  Wanted %d, got %d", len(expected), len(got))
	}

	for i, expectedVal := range expected {
		switch expectedVal.(type) {
		case int:
			err := testIntegerConstant(expectedVal.(int64), got[i])
			if err != nil {
				return fmt.Errorf("integer value failed at %d. %s", i, err)
			}

		}
	}
	return nil
}

func testIntegerConstant(expected int64, got object.Object) error {
	result, ok := got.(*object.Integer)
	if !ok {
		return fmt.Errorf("Object was not an integer.  Got %T (%v)", got, got)
	}
	if result.Value != expected {
		return fmt.Errorf("incorrect value.  Wanted %d, got %d", expected, result.Value)
	}
	return nil
}
