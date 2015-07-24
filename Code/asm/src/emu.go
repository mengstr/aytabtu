package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"time"
)

/*

00: pro pro pro pro pro pro pro pro   Reg A : 00 - 0000 0000           In 0 : 00 - 0000 0000
08: pro pro pro pro pro pro pro pro   Reg B : 00 - 0000 0000           In 1 : 00 - 0000 0000
10: pro pro pro pro pro pro pro pro   Reg C : 00 - 0000 0000           In 2 : 00 - 0000 0000
18: pro pro pro pro pro pro pro pro   Reg D : 00 - 0000 0000           In 3 : 00 - 0000 0000
20: pro pro pro pro pro pro pro pro   Reg X : 00 - 0000 0000           In 4 : 00 - 0000 0000
28: pro pro pro pro pro pro pro pro   Reg Y : 00 - 0000 0000           In 5 : 00 - 0000 0000
30: pro pro pro pro pro pro pro pro   Reg JH: 00 - 0000 0000           In 6 : 00 - 0000 0000
38: pro pro pro pro pro pro pro pro   Reg JL: 00 - 0000 0000           In 7 : 00 - 0000 0000
40: pro pro pro pro pro pro pro pro                                    In 8 : 00 - 0000 0000
48: pro pro pro pro pro pro pro pro   PC    : 0000                     In 9 : 00 - 0000 0000
50: pro pro pro pro pro pro pro pro   PChold: 0000                     In A : 00 - 0000 0000
58: pro pro pro pro pro pro pro pro   Flag C: 0                        In B : 00 - 0000 0000
60: pro pro pro pro pro pro pro pro   Flag Z: 0                        In C : 00 - 0000 0000
68: pro pro pro pro pro pro pro pro   Cycles: 000000000                In D : 00 - 0000 0000
70: pro pro pro pro pro pro pro pro                                    In E : 00 - 0000 0000
78: pro pro pro pro pro pro pro pro                                    In F : 00 - 0000 0000
80: pro pro pro pro pro pro pro pro   Dat 00: da da da da da da da da  Out 0: 00 - 0000 0000
88: pro pro pro pro pro pro pro pro   Dat 08: da da da da da da da da  Out 1: 00 - 0000 0000
90: pro pro pro pro pro pro pro pro   Dat 10: da da da da da da da da  Out 2: 00 - 0000 0000
98: pro pro pro pro pro pro pro pro   Dat 18: da da da da da da da da  Out 3: 00 - 0000 0000
A0: pro pro pro pro pro pro pro pro   Dat 20: da da da da da da da da  Out 4: 00 - 0000 0000
A8: pro pro pro pro pro pro pro pro   Dat 28: da da da da da da da da  Out 5: 00 - 0000 0000
B0: pro pro pro pro pro pro pro pro   Dat 30: da da da da da da da da  Out 6: 00 - 0000 0000
B8: pro pro pro pro pro pro pro pro   Dat 38: da da da da da da da da  Out 7: 00 - 0000 0000
C0: pro pro pro pro pro pro pro pro   Dat 40: da da da da da da da da  Out 8: 00 - 0000 0000
C8: pro pro pro pro pro pro pro pro   Dat 48: da da da da da da da da  Out 9: 00 - 0000 0000
D0: pro pro pro pro pro pro pro pro   Dat 50: da da da da da da da da  Out A: 00 - 0000 0000
D8: pro pro pro pro pro pro pro pro   Dat 58: da da da da da da da da  Out B: 00 - 0000 0000
E0: pro pro pro pro pro pro pro pro   Dat 60: da da da da da da da da  Out C: 00 - 0000 0000
E8: pro pro pro pro pro pro pro pro   Dat 68: da da da da da da da da  Out D: 00 - 0000 0000
F0: pro pro pro pro pro pro pro pro   Dat 70: da da da da da da da da  Out E: 00 - 0000 0000
F8: pro pro pro pro pro pro pro pro   Dat 78: da da da da da da da da  Out F: 00 - 0000 0000


*/

const (
	P_none = iota // Not allowed
	P_of6  = iota // 6 bits - label or value -31..+32
	P_va8  = iota // 8 bits - const or value 0..255
	P_mem  = iota // 8 bits - const or value 0..255
	P_reg  = iota // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
	P_bit  = iota // 3 bits - const or value 0..7
	P_io   = iota // 4 bits - const or value 0..15
)

type opcode struct {
	base    uint
	arg1    uint
	arg1pos uint
	arg2    uint
	arg2pos uint
}

var opcodes = map[string]opcode{
	"HALT":    {base: 0x000, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  0   0   0     0   0   x   x     x   x   x   x
	"NOP":     {base: 0x080, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  0   0   0     1   0   x   x     x   x   x   x
	"SJMP":    {base: 0x0C0, arg1: P_of6, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   0   0     1   1   o5  o4    o3  o2  o1  o0
	"BRAZ":    {base: 0x100, arg1: P_of6, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   0   1     0   0   o5  o4    o3  o2  o1  o0
	"BRANZ":   {base: 0x140, arg1: P_of6, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   0   1     0   1   o5  o4    o3  o2  o1  o0
	"BRAC":    {base: 0x180, arg1: P_of6, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   0   1     1   0   o5  o4    o3  o2  o1  o0
	"BRANC":   {base: 0x1C0, arg1: P_of6, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   0   1     1   1   o5  o4    o3  o2  o1  o0
	"LDAI":    {base: 0x200, arg1: P_va8, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   1   0     id7 id6 id5 id4   id3 id2 id1 id0
	"LDAZP":   {base: 0x300, arg1: P_mem, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  0   1   1     a7  a6  a5  a4    a3  a2  a1  a0
	"STAZP":   {base: 0x400, arg1: P_mem, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   0     a7  a6  a5  a4    a3  a2  a1  a0
	"CLR":     {base: 0x500, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   0   0     0   r2  r1  r0
	"SETFF":   {base: 0x508, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   0   0     1   r2  r1  r0
	"NOT":     {base: 0x510, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   0   1     0   r2  r1  r0
	"OR":      {base: 0x518, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   0   1     1   r2  r1  r0
	"AND":     {base: 0x520, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   1   0     0   r2  r1  r0
	"XOR":     {base: 0x528, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   1   0     1   r2  r1  r0
	"INC":     {base: 0x530, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   1   1     0   r2  r1  r0
	"DEC":     {base: 0x538, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   0   1   1     1   r2  r1  r0
	"ADD":     {base: 0x540, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   0   0     0   r2  r1  r0
	"SUB":     {base: 0x548, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   0   0     1   r2  r1  r0
	"ADDC":    {base: 0x550, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   0   1     0   r2  r1  r0
	"SUBC":    {base: 0x558, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   0   1     1   r2  r1  r0
	"LSHIFT":  {base: 0x560, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   1   0     0   r2  r1  r0
	"RSHIFT":  {base: 0x568, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   1   0     1   r2  r1  r0
	"LSHIFTC": {base: 0x570, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   1   1     0   r2  r1  r0
	"RSHIFTC": {base: 0x578, arg1: P_reg, arg1pos: 0, arg2: P_none, arg2pos: 0},  //  1   0   1     0   1   1   1     1   r2  r1  r0
	"MOVE":    {base: 0x580, arg1: P_reg, arg1pos: 3, arg2: P_reg, arg2pos: 0},   //  1   0   1     1   0   rs2 rs1   rs0 rd2 rd1 rd0
	"TEST":    {base: 0x5C0, arg1: P_reg, arg1pos: 3, arg2: P_bit, arg2pos: 0},   //  1   0   1     1   1   b2  b1    b0  r2  r1  r0
	"PEEK":    {base: 0x600, arg1: P_io, arg1pos: 0, arg2: P_none, arg2pos: 0},   //  1   1   0     0   0   0   0     p3  p2  p1  p0
	"POKE":    {base: 0x610, arg1: P_io, arg1pos: 0, arg2: P_none, arg2pos: 0},   //  1   1   0     0   0   0   1     p3  p2  p1  p0
	"LDAXY":   {base: 0x620, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   0   0   0
	"STAXY":   {base: 0x621, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   0   0   1
	"JUMP":    {base: 0x622, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   0   1   0
	"CALL":    {base: 0x623, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   0   1   1
	"RET":     {base: 0x624, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   1   0   0
	"CLRC":    {base: 0x626, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   1   1   0
	"SETC":    {base: 0x627, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1   1   0     0   0   1   0     0   1   1   1
	"LDPML":   {base: 0x628, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1	  1	  0	    0	0	1	0	  1	  0   0	  0
	"LDPMH":   {base: 0x629, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1	  1	  0	    0	0	1	0	  1	  0   0	  1
	"STPML":   {base: 0x62A, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1	  1	  0	    0	0	1	0	  1	  0   1	  0
	"STPMH":   {base: 0x62B, arg1: P_none, arg1pos: 0, arg2: P_none, arg2pos: 0}, //  1	  1	  0	    0	0	1	0	  1	  0   1	  1

}

const screenMinWidth = 92
const screenMinHeight = 32

type CPU struct {
	code    [256]uint
	data    [256]uint
	regs    [8]uint
	pc      uint
	pchold  uint
	flagz   bool
	flagc   bool
	in      [16]uint
	out     [16]uint
	cycles  uint
	running bool
}

var regnames = []string{"A", "B", "C", "D", "X", "Y", "JH", "JL"}

const (
	A  = 0
	B  = 1
	C  = 2
	D  = 3
	X  = 4
	Y  = 5
	JH = 6
	JL = 7
)

var cpu CPU

//
//
//
func findOpcode(op uint) string {
	for opKey, opData := range opcodes {
		var mask uint = 0
		switch opData.arg1 {
		case P_none:
			mask |= (0x00 << opData.arg1pos)
		case P_of6: // 6 bits - label or value -31..+32
			mask |= (0x3F << opData.arg1pos)
		case P_va8: // 8 bits - const or value 0..255
			mask |= (0xFF << opData.arg1pos)
		case P_mem: // 8 bits - const or value 0..255
			mask |= (0xFF << opData.arg1pos)
		case P_reg: // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
			mask |= (0x07 << opData.arg1pos)
		case P_bit: // 3 bits - const or value 0..7
			mask |= (0x07 << opData.arg1pos)
		case P_io: // 4 bits - const or value 0..15
			mask |= (0x0F << opData.arg1pos)
		}
		switch opData.arg2 {
		case P_none:
			mask |= (0x00 << opData.arg2pos)
		case P_of6: // 6 bits - label or value -31..+32
			mask |= (0x3F << opData.arg2pos)
		case P_va8: // 8 bits - const or value 0..255
			mask |= (0xFF << opData.arg2pos)
		case P_mem: // 8 bits - const or value 0..255
			mask |= (0xFF << opData.arg2pos)
		case P_reg: // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
			mask |= (0x07 << opData.arg2pos)
		case P_bit: // 3 bits - const or value 0..7
			mask |= (0x07 << opData.arg2pos)
		case P_io: // 4 bits - const or value 0..15
			mask |= (0x0F << opData.arg2pos)
		}

		if op&(^mask) == opData.base {
			//			fmt.Printf("opcode %d is %s\n", op, opKey)
			return opKey
		}
	}
	return ""
}

//
//
//
func printXY(x int, y int, fgcolor termbox.Attribute, bgcolor termbox.Attribute, txt string) {
	for i, elem := range txt {
		termbox.SetCell(x+i, y, rune(elem), fgcolor, bgcolor)
	}
}

//
//
//
func getParam(opcode uint, arg uint, pos uint) (haveparam bool, paramvalue uint) {
	switch arg {
	case P_none:
		return false, 0
	case P_of6: // 6 bits - label or value -31..+32
		return true, (opcode >> pos) & 0x3F
	case P_va8: // 8 bits - const or value 0..255
		return true, (opcode >> pos) & 0xFF
	case P_mem: // 8 bits - const or value 0..255
		return true, (opcode >> pos) & 0xFF
	case P_reg: // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
		return true, (opcode >> pos) & 0x07
	case P_bit: // 3 bits - const or value 0..7
		return true, (opcode >> pos) & 0x07
	case P_io: // 4 bits - const or value 0..15
		return true, (opcode >> pos) & 0x0F
	}

	return false, 0
}

//
//
//
func disassembleGetParam(opcode uint, arg uint, pos uint) string {
	txt := ""
	ok, val := getParam(opcode, arg, pos)
	if ok {
		switch arg {
		case P_none:
			txt = ""
		case P_of6: // 6 bits - label or value -31..+32
			if val < 33 {
				txt = fmt.Sprintf("+%d", val)
			} else {
				txt = fmt.Sprintf("-%d", 64-val)
			}
		case P_va8: // 8 bits - const or value 0..255
			txt = fmt.Sprintf("$%02x", val)
		case P_mem: // 8 bits - const or value 0..255
			txt = fmt.Sprintf("$%02x", val)
		case P_reg: // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
			txt = regnames[val]
		case P_bit: // 3 bits - const or value 0..7
			txt = fmt.Sprintf("$%x", val)
		case P_io: // 4 bits - const or value 0..15
			txt = fmt.Sprintf("$%x", val)
		}
	}
	return txt
}

//
//
//
func disassemble(opcode uint) string {
	opname := findOpcode(opcode)
	txt := fmt.Sprintf("%03x %s", opcode, opname)

	param1 := disassembleGetParam(opcode, opcodes[opname].arg1, opcodes[opname].arg1pos)
	param2 := disassembleGetParam(opcode, opcodes[opname].arg2, opcodes[opname].arg2pos)
	if param1 != "" {
		txt += " " + param1
		if param2 != "" {
			txt += ", " + param2
		}
	}
	return txt
}

//
//
//
func redrawCpu() {
	var txt string

	// Code memory
	printXY(0, 0, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " PROGRAM MEMORY                    ")
	for i := 0; i < 256; i += 8 {
		txt = fmt.Sprintf("%02x: %03x %03x %03x %03x %03x %03x %03x %03x", i, cpu.code[i+0], cpu.code[i+1], cpu.code[i+2], cpu.code[i+3], cpu.code[i+4], cpu.code[i+5], cpu.code[i+6], cpu.code[i+7])
		printXY(0, 1+i/8, termbox.ColorBlack, termbox.ColorDefault, txt)
	}

	// Data memory
	printXY(39, 18, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " DATA MEMORY               ")
	for i := 0; i < 128; i += 8 {
		txt = fmt.Sprintf("%02x: %02x %02x %02x %02x %02x %02x %02x %02x", i, cpu.data[i+0], cpu.data[i+1], cpu.data[i+2], cpu.data[i+3], cpu.data[i+4], cpu.data[i+5], cpu.data[i+6], cpu.data[i+7])
		printXY(39, 19+i/8, termbox.ColorBlack, termbox.ColorDefault, txt)
	}

	// Registers
	printXY(39, 0, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " REGISTERS       ")
	for i := 0; i < 8; i++ {
		txt = fmt.Sprintf("%2s: %02x - %08b", regnames[i], cpu.regs[i], cpu.regs[i])
		printXY(39, 1+i, termbox.ColorBlack, termbox.ColorDefault, txt)
	}

	// Internal registers
	printXY(39, 10, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " CPU INTERNAL ")
	txt = fmt.Sprintf("PC    : %04x", cpu.pc)
	printXY(39, 11, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("PChold: %04x", cpu.pchold)
	printXY(39, 12, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("Flag C: %t", cpu.flagc)
	printXY(39, 13, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("Flag Z: %t", cpu.flagz)
	printXY(39, 14, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("Cycles: %d", cpu.cycles)
	printXY(39, 15, termbox.ColorBlack, termbox.ColorDefault, txt)
	// Current instruction
	txt = "Instr : " + disassemble(cpu.code[cpu.pc]) + "                    " // spaces removes old crud
	printXY(39, 16, termbox.ColorBlack, termbox.ColorDefault, txt[:30])       // max 30 chars, don't overrun

	// I/O In values
	printXY(70, 0, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " I/O INPUTS         ")
	for i := 0; i < 16; i++ {
		txt = fmt.Sprintf("In  %X: %02x - %08b", i, cpu.in[i], cpu.in[i])
		printXY(70, 1+i, termbox.ColorBlack, termbox.ColorDefault, txt)
	}

	// I/O Out values
	printXY(70, 18, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " I/O OUTPUTS        ")
	for i := 0; i < 16; i++ {
		txt = fmt.Sprintf("Out %X: %02x - %08b", i, cpu.out[i], cpu.out[i])
		printXY(70, 3+16+i, termbox.ColorBlack, termbox.ColorDefault, txt)
	}

}

//
// Initialize the CPU data, the memories are set to random (but for program mem valid) data
//
func initCpu() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < len(cpu.code); i++ {
		opname := ""
		var opcode uint = 0
		for opname == "" {
			opcode = uint(rand.Uint32() & 0x7FF)
			opname = findOpcode(opcode)
		}
		cpu.code[i] = opcode
	}
	for i := 0; i < len(cpu.data); i++ {
		cpu.data[i] = uint(rand.Uint32() & 0xFF)
	}
	for i := 0; i < len(cpu.regs); i++ {
		cpu.regs[i] = 0x00
	}
	for i := 0; i < len(cpu.in); i++ {
		cpu.in[i] = 0x00
	}
	for i := 0; i < len(cpu.out); i++ {
		cpu.out[i] = 0x00
	}
	cpu.pc = 0x0000
	cpu.pchold = 0x0000
	cpu.flagc = false
	cpu.flagz = true
	cpu.cycles = 0
	cpu.running = false
}

//
//
//
func updatePcOffset(ofs uint) {
	if ofs < 33 {
		cpu.pc += ofs
	} else {
		cpu.pc -= (64 - ofs)
	}
}

//
//
//
func executeOneOp() {
	opcode := cpu.code[cpu.pc]
	opname := findOpcode(opcode)
	_, p1 := getParam(opcode, opcodes[opname].arg1, opcodes[opname].arg1pos)
	_, p2 := getParam(opcode, opcodes[opname].arg2, opcodes[opname].arg2pos)

	switch opname {
	case "HALT":
		cpu.running = false
	case "NOP":
	case "SJMP":
		updatePcOffset(p1)
	case "BRAZ":
		if cpu.flagz {
			updatePcOffset(p1)
		}
	case "BRANZ":
		if !cpu.flagz {
			updatePcOffset(p1)
		}
	case "BRAC":
		if cpu.flagc {
			updatePcOffset(p1)
		}
	case "BRANC":
		if !cpu.flagc {
			updatePcOffset(p1)
		}
	case "LDAI":
		cpu.regs[A] = p1
	case "LDAZP":
		cpu.regs[A] = cpu.data[p1]
	case "STAZP":
		cpu.data[p1] = cpu.regs[A]
	case "CLR":
		cpu.regs[p1] = 0x00
	case "SETFF":
		cpu.regs[p1] = 0xFF
	case "NOT":
		cpu.regs[p1] ^= 0xFF
	case "OR":
		cpu.regs[p1] |= cpu.regs[A]
	case "AND":
		cpu.regs[p1] &= cpu.regs[A]
	case "XOR":
		cpu.regs[p1] ^= cpu.regs[A]
	case "INC":
		cpu.regs[p1] = (cpu.regs[p1] + 1) & 0xFF
	case "DEC":
		cpu.regs[p1] = (cpu.regs[p1] - 1) & 0xFF
	case "ADD":
		cpu.regs[p1] = (cpu.regs[p1] + cpu.regs[A]) & 0xFF
	case "SUB":
		cpu.regs[p1] = (cpu.regs[p1] - cpu.regs[A]) & 0xFF
	case "ADDC":
		if cpu.flagc {
			cpu.regs[p1] = (cpu.regs[p1] + cpu.regs[A] + 1) & 0xFF
		} else {
			cpu.regs[p1] = (cpu.regs[p1] + cpu.regs[A]) & 0xFF
		}
	case "SUBC":
		if cpu.flagc {
			cpu.regs[p1] = (cpu.regs[p1] - cpu.regs[A] - 1) & 0xFF
		} else {
			cpu.regs[p1] = (cpu.regs[p1] - cpu.regs[A]) & 0xFF
		}
	case "LSHIFT":
		cpu.regs[p1] = (cpu.regs[p1] * 2) & 0xFE
	case "RSHIFT":
		cpu.regs[p1] = (cpu.regs[p1] / 2) & 0xEF
	case "LSHIFTC":
		cpu.regs[p1] = (cpu.regs[p1] * 2) & 0xFE
		if cpu.flagc {
			cpu.regs[p1] |= 0x01
		}
	case "RSHIFTC":
		cpu.regs[p1] = (cpu.regs[p1] / 2) & 0xEF
		if cpu.flagc {
			cpu.regs[p1] |= 0x80
		}
	case "MOVE":
		cpu.regs[p2] = cpu.regs[p1]
	case "TEST":
		if (cpu.regs[p2] & (1 << cpu.regs[p1])) == (1 << cpu.regs[p1]) {
			cpu.flagz = false
		} else {
			cpu.flagz = true
		}
	case "PEEK":
		cpu.regs[A] = cpu.in[p1]
	case "POKE":
		cpu.out[p1] = cpu.regs[A]
	case "LDAXY":
		cpu.regs[A] = cpu.data[cpu.regs[Y]<<8+cpu.regs[X]]
	case "STAXY":
		cpu.data[cpu.regs[Y]<<8+cpu.regs[X]] = cpu.regs[A]
	case "JUMP":
		cpu.pc = cpu.regs[JH]<<8 + cpu.regs[JL] - 1
	case "CALL":
		cpu.pchold = cpu.pc
		cpu.pc = cpu.regs[JH]<<8 + cpu.regs[JL]
	case "RET":
		cpu.pc = cpu.pchold
	case "CLRC":
		cpu.flagc = false
	case "SETC":
		cpu.flagc = true
	case "LDPML":
		cpu.regs[A] = (cpu.code[cpu.regs[Y]<<8+cpu.code[X]]) & 0xFF
	case "LDPMH":
		cpu.regs[A] = (cpu.code[cpu.regs[Y]<<8+cpu.code[X]]) >> 8
	case "STPML":
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] &= 0x700
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] |= cpu.regs[A]
	case "STPMH":
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] &= 0x0FF
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] |= ((cpu.regs[A] & 0x07) << 8)
	}
	cpu.pc++
}

//
//
//
func main() {
	// Initialize termbox library
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Check that the screen is big enough for us
	width, height := termbox.Size()
	if width < screenMinWidth || height < screenMinHeight {
		termbox.Close()
		fmt.Printf("Error Screen W/H %d/%d is less than %d/%d\n", width, height, screenMinWidth, screenMinHeight)
		os.Exit(1)
	}
	termbox.SetOutputMode(termbox.OutputNormal)

	initCpu()
	redrawCpu()
	termbox.Flush()
	exit := false
	for !exit {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if string(ev.Ch) == "s" {
				executeOneOp()
				redrawCpu()
				termbox.Flush()
			}
			if ev.Key == termbox.KeyCtrlC {
				exit = true
			}
			if ev.Key == termbox.KeyEnter {
				initCpu()
				redrawCpu()
				termbox.Flush()
			}

		}
	}

}
