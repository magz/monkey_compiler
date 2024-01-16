package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	offset := 0
	for offset < len(ins) {
		def, err := Lookup(ins[offset])
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			fmt.Fprintf(&out, "ERROR: %s", err)
			return out.String()
		}
		operands, offsetIncrease := ReadOperands(def, ins[offset+1:])

		fmt.Fprintf(&out, "%04d %s%s\n", offset, def.Name, ins.fmtInstruction(def, operands))

		offset += 1 + offsetIncrease
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	switch def.Name {
	case "OpConstant":
		return fmt.Sprintf(" %d", operands[0])
	case "OpAdd":
		return ""
	default:
		return "UNKOWN OPERAND NAME"
	}
}

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

func ReadOperands(def *Definition, ins []byte) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
