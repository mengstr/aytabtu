package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

	012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789

	 PROGRAM MEMORY                         REGISTERS                      I/O INPUTS
0	00: 000 000 000 000 000 000 000 000     A: 00 - 00000000              In  0: 00 - 00000000
1	08: 000 000 000 000 000 000 000 000     B: 00 - 00000000              In  1: 00 - 00000000
2	10: 000 000 000 000 000 000 000 000     C: 00 - 00000000              In  2: 00 - 00000000
3	18: 000 000 000 000 000 000 000 000     D: 00 - 00000000              In  3: 00 - 00000000
4	20: 000 000 000 000 000 000 000 000     X: 00 - 00000000              In  4: 00 - 00000000
5	28: 000 000 000 000 000 000 000 000     Y: 00 - 00000000              In  5: 00 - 00000000
6	30: 000 000 000 000 000 000 000 000    JH: 00 - 00000000              In  6: 00 - 00000000
7	38: 000 000 000 000 000 000 000 000    JL: 00 - 00000000              In  7: 00 - 00000000
8	40: 000 000 000 000 000 000 000 000                                   In  8: 00 - 00000000
9	48: 000 000 000 000 000 000 000 000     CPU INTERNAL                  In  9: 00 - 00000000
0	50: 000 000 000 000 000 000 000 000    PC    : 0000                   In  A: 00 - 00000000
1	58: 000 000 000 000 000 000 000 000    PChold: 0000                   In  B: 00 - 00000000
2	60: 000 000 000 000 000 000 000 000    Flag C: false                  In  C: 00 - 00000000
3	68: 000 000 000 000 000 000 000 000    Flag Z: true                   In  D: 00 - 00000000
4	70: 000 000 000 000 000 000 000 000    Cycles: 0                      In  E: 00 - 00000000
5	78: 000 000 000 000 000 000 000 000    Instr : 000 HALT               In  F: 00 - 00000000
6	80: 000 000 000 000 000 000 000 000
7	88: 000 000 000 000 000 000 000 000     DATA MEMORY                    I/O OUTPUTS
8	90: 000 000 000 000 000 000 000 000    00: 00 00 00 00 00 00 00 00    Out 0: 00 - 00000000
9	98: 000 000 000 000 000 000 000 000    08: 00 00 00 00 00 00 00 00    Out 1: 00 - 00000000
0	a0: 000 000 000 000 000 000 000 000    10: 00 00 00 00 00 00 00 00    Out 2: 00 - 00000000
1	a8: 000 000 000 000 000 000 000 000    18: 00 00 00 00 00 00 00 00    Out 3: 00 - 00000000
2	b0: 000 000 000 000 000 000 000 000    20: 00 00 00 00 00 00 00 00    Out 4: 00 - 00000000
3	b8: 000 000 000 000 000 000 000 000    28: 00 00 00 00 00 00 00 00    Out 5: 00 - 00000000
4	c0: 000 000 000 000 000 000 000 000    30: 00 00 00 00 00 00 00 00    Out 6: 00 - 00000000
5	c8: 000 000 000 000 000 000 000 000    38: 00 00 00 00 00 00 00 00    Out 7: 00 - 00000000
6	d0: 000 000 000 000 000 000 000 000    40: 00 00 00 00 00 00 00 00    Out 8: 00 - 00000000
7	d8: 000 000 000 000 000 000 000 000    48: 00 00 00 00 00 00 00 00    Out 9: 00 - 00000000
8	e0: 000 000 000 000 000 000 000 000    50: 00 00 00 00 00 00 00 00    Out A: 00 - 00000000
9	e8: 000 000 000 000 000 000 000 000    58: 00 00 00 00 00 00 00 00    Out B: 00 - 00000000
0	f0: 000 000 000 000 000 000 000 000    60: 00 00 00 00 00 00 00 00    Out C: 00 - 00000000
1	f8: 000 000 000 000 000 000 000 000    68: 00 00 00 00 00 00 00 00    Out D: 00 - 00000000
2	                                       70: 00 00 00 00 00 00 00 00    Out E: 00 - 00000000
3	                                       78: 00 00 00 00 00 00 00 00    Out F: 00 - 00000000


*/

const (
	uartPort   = "1420"
	uartType   = "tcp"
	logType    = "udp"
	logPort    = "1421"
	logSrcPort = "1422" // Must have fixed source to please netcat
)

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

// Mappings for register Names and Numbers
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

// Enumerate I/O port special functions
const (
	IO_UART = 15
)

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
	cycles  uint64
	running bool
}

var cpu, holdCpu CPU

var logChannel chan string // Channel used for sending log entries over UDP
var uartChannel chan uint  // Channel used for simulation an UART at PORT 15

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
	var color termbox.Attribute

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
		txt = fmt.Sprintf("%2s: ", regnames[i])
		printXY(39, 1+i, termbox.ColorBlack, termbox.ColorDefault, txt)
		txt = fmt.Sprintf("%02x - %08b", cpu.regs[i], cpu.regs[i])
		color = termbox.ColorBlack
		if cpu.regs[i] != holdCpu.regs[i] {
			color = termbox.ColorRed
		}
		printXY(39, 1+i, color, termbox.ColorDefault, txt)

	}

	// Internal registers
	printXY(39, 10, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " CPU INTERNAL ")

	txt = fmt.Sprintf("PC    : ")
	printXY(39, 11, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("%04x", cpu.pc)
	color = termbox.ColorBlack
	printXY(39+8, 11, color, termbox.ColorDefault, txt)

	txt = fmt.Sprintf("PChold: ")
	printXY(39, 12, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("%04x", cpu.pchold)
	color = termbox.ColorBlack
	if cpu.pchold != holdCpu.pchold {
		color = termbox.ColorRed
	}
	printXY(39+8, 12, color, termbox.ColorDefault, txt)

	txt = fmt.Sprintf("Flag C: ")
	printXY(39, 13, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("%t", cpu.flagc)
	color = termbox.ColorBlack
	if cpu.flagc != holdCpu.flagc {
		color = termbox.ColorRed
	}
	printXY(39+8, 13, color, termbox.ColorDefault, txt)

	txt = fmt.Sprintf("Flag Z: ")
	printXY(39, 14, termbox.ColorBlack, termbox.ColorDefault, txt)
	txt = fmt.Sprintf("%t", cpu.flagz)
	color = termbox.ColorBlack
	if cpu.flagz != holdCpu.flagz {
		color = termbox.ColorRed
	}
	printXY(39+8, 14, color, termbox.ColorDefault, txt)

	txt = fmt.Sprintf("Cycles: %d", cpu.cycles)
	printXY(39, 15, termbox.ColorBlack, termbox.ColorDefault, txt)
	// Current instruction
	txt = "Instr : " + disassemble(cpu.code[cpu.pc]) + "                    " // spaces removes old crud
	printXY(39, 16, termbox.ColorBlack, termbox.ColorDefault, txt[:30])       // max 30 chars, don't overrun

	// I/O In values
	printXY(70, 0, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " I/O INPUTS         ")
	for i := 0; i < 16; i++ {
		txt = fmt.Sprintf("In  %X: ", i)
		printXY(70, 1+i, termbox.ColorBlack, termbox.ColorDefault, txt)
		txt = fmt.Sprintf("%02x - %08b", cpu.in[i], cpu.in[i])
		color = termbox.ColorBlack
		if cpu.in[i] != holdCpu.in[i] {
			color = termbox.ColorRed
		}
		printXY(70+7, 1+i, color, termbox.ColorDefault, txt)
	}

	// I/O Out values
	printXY(70, 18, termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault, " I/O OUTPUTS        ")
	for i := 0; i < 16; i++ {
		txt = fmt.Sprintf("Out %X: ", i)
		printXY(70, 3+16+i, termbox.ColorBlack, termbox.ColorDefault, txt)
		txt = fmt.Sprintf("%02x - %08b", cpu.out[i], cpu.out[i])
		color = termbox.ColorBlack
		if cpu.out[i] != holdCpu.out[i] {
			color = termbox.ColorRed
		}
		printXY(70+7, 3+16+i, color, termbox.ColorDefault, txt)
	}

	// Send updates in local buffer to the screen
	termbox.Flush()
	// Update changes to holding cpu
	holdCpu = cpu
}

//
// Initialize the CPU data, the memories are set to random (but for program mem valid) data
//
func initCpu() {
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < len(cpu.code); i++ {
		cpu.code[i] = 0x00
	}
	for i := 0; i < len(cpu.data); i++ {
		cpu.data[i] = 0x00
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

	holdCpu = cpu
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
		// Do nothing

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
		cpu.flagz = cpu.regs[A] == 0

	case "LDAZP":
		cpu.regs[A] = cpu.data[p1]
		cpu.flagz = cpu.regs[A] == 0

	case "STAZP":
		cpu.data[p1] = cpu.regs[A]

	case "CLR":
		cpu.regs[p1] = 0x00
		cpu.flagz = true

	case "SETFF":
		cpu.regs[p1] = 0xFF
		cpu.flagz = false

	case "NOT":
		cpu.regs[p1] ^= 0xFF
		cpu.flagz = cpu.regs[p1] == 0

	case "OR":
		cpu.regs[p1] |= cpu.regs[A]
		cpu.flagz = cpu.regs[p1] == 0

	case "AND":
		cpu.regs[p1] &= cpu.regs[A]
		cpu.flagz = cpu.regs[p1] == 0

	case "XOR":
		cpu.regs[p1] ^= cpu.regs[A]
		cpu.flagz = cpu.regs[p1] == 0

	case "INC":
		cpu.regs[p1] = (cpu.regs[p1] + 1) & 0xFF
		cpu.flagz = cpu.regs[p1] == 0

	case "DEC":
		cpu.regs[p1] = (cpu.regs[p1] - 1) & 0xFF
		cpu.flagz = cpu.regs[p1] == 0

	case "ADD":
		cpu.regs[p1] = cpu.regs[p1] + cpu.regs[A]
		if cpu.regs[p1] > 255 {
			cpu.regs[p1] -= 256
			cpu.flagc = true
		} else {
			cpu.flagc = false
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "SUB":
		cpu.regs[p1] = cpu.regs[p1] - cpu.regs[A]
		if cpu.regs[p1] < 0 {
			cpu.regs[p1] += 256
			cpu.flagc = true
		} else {
			cpu.flagc = false
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "ADDC":
		if cpu.flagc {
			cpu.regs[p1] = cpu.regs[p1] + cpu.regs[A] + 1
		} else {
			cpu.regs[p1] = cpu.regs[p1] + cpu.regs[A]
		}
		if cpu.regs[p1] > 255 {
			cpu.regs[p1] -= 256
			cpu.flagc = true
		} else {
			cpu.flagc = false
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "SUBC":
		if cpu.flagc {
			cpu.regs[p1] = cpu.regs[p1] - cpu.regs[A] - 1
		} else {
			cpu.regs[p1] = cpu.regs[p1] - cpu.regs[A]
		}
		if cpu.regs[p1] < 0 {
			cpu.regs[p1] += 256
			cpu.flagc = true
		} else {
			cpu.flagc = false
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "LSHIFT":
		cpu.regs[p1] = (cpu.regs[p1] * 2) & 0xFE
		cpu.flagz = cpu.regs[p1] == 0

	case "RSHIFT":
		cpu.regs[p1] = (cpu.regs[p1] / 2) & 0xEF
		cpu.flagz = cpu.regs[p1] == 0

	case "LSHIFTC":
		cpu.regs[p1] = (cpu.regs[p1] * 2) & 0xFE
		if cpu.flagc {
			cpu.regs[p1] |= 0x01
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "RSHIFTC":
		cpu.regs[p1] = (cpu.regs[p1] / 2) & 0xEF
		if cpu.flagc {
			cpu.regs[p1] |= 0x80
		}
		cpu.flagz = cpu.regs[p1] == 0

	case "MOVE":
		cpu.regs[p2] = cpu.regs[p1]
		cpu.flagz = cpu.regs[p1] == 0

	case "TEST":
		if (cpu.regs[p2] & (1 << cpu.regs[p1])) == (1 << cpu.regs[p1]) {
			cpu.flagz = false
		} else {
			cpu.flagz = true
		}

	case "PEEK":
		cpu.regs[A] = cpu.in[p1]
		cpu.flagz = cpu.regs[A] == 0

	case "POKE":
		cpu.out[p1] = cpu.regs[A]
		if p1 == IO_UART {
			uartChannel <- cpu.regs[A]
		}

	case "LDAXY":
		cpu.regs[A] = cpu.data[cpu.regs[Y]<<8+cpu.regs[X]]
		cpu.flagz = cpu.regs[A] == 0

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
		cpu.flagz = cpu.regs[p1] == 0

	case "LDPMH":
		cpu.regs[A] = (cpu.code[cpu.regs[Y]<<8+cpu.code[X]]) >> 8
		cpu.flagz = cpu.regs[p1] == 0

	case "STPML":
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] &= 0x700
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] |= cpu.regs[A]

	case "STPMH":
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] &= 0x0FF
		cpu.code[cpu.regs[Y]<<8+cpu.code[X]] |= ((cpu.regs[A] & 0x07) << 8)
	}
	cpu.pc++
	cpu.cycles++
}

//
// Receive log messages over the channel and blast them out
// to UDP port 1234 using a fixed source address so it's
// possible to receive the messages with netcat
//
func udpLog(theChannel chan string) {
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+logPort)
	localAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+logSrcPort)
	socket, _ := net.DialUDP("udp", localAddr, serverAddr)

	for {
		msg := <-theChannel
		t := time.Time(time.Now()).Format(time.StampMilli)
		msg = fmt.Sprintf("%s: %s\n", t, msg)
		_, _ = socket.Write([]byte(msg))
	}
}

//
// Read incomming data from the tcp connection and send them to the channel for
// consumption in the main loop
//
func uartRxComms(theChannel chan uint, conn net.Conn) {
	// Begin with turning on thecharacter-by-character mode at the
	// client and turn off echo as well
	conn.Write([]byte("\377\373\003\n"))    // send IAC WILL SUPPRESS-GOAHEAD
	conn.Write([]byte("\377\375\003\n"))    // send IAC DO SUPPRESS-GO-AHEAD
	conn.Write([]byte("\377\373\001\n"))    // send IAC WILL SUPPRESS-ECHO
	conn.Write([]byte("\377\375\001\n"))    // send IAC DO SUPPRESS-ECHO
	conn.Write([]byte("AYTABTU ready\r\n")) // Tell telnet user that we're good to go

	buf := make([]byte, 16)
	for {
		_, _ = conn.Read(buf)
		theChannel <- uint(buf[0])
	}
}

//
// Listen for messages on the usart channel and send them over the socket to
// be shown in the terminal
//
func uartTxComms(theChannel chan uint, conn net.Conn) {
	var buf = []byte{0}
	for {
		v := <-theChannel
		buf[0] = byte(v)
		_, _ = conn.Write([]byte(buf))
	}
}

//
//
//
func uartListener(theChannel chan uint) {
	// Start listeng on TCP port
	tcp, err := net.Listen(uartType, "127.0.0.1:"+uartPort)
	if err != nil {
		fmt.Println("Error listening on port:", err.Error())
		os.Exit(1)
	}
	defer tcp.Close()

	for {
		tcpconn, err := tcp.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		// Handle rx/tx comms in two new goroutines.
		go uartRxComms(theChannel, tcpconn)
		go uartTxComms(theChannel, tcpconn)
	}
}

//
// Continusly poll for keyboard events and post them to the channel for
// consumtion in the main loop
//
func keyboardPoller(c chan termbox.Event) {
	for {
		c <- termbox.PollEvent()
	}
}

func loadHex(filename string) {
	// Load entire file in one go
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logChannel <- fmt.Sprintf("Can't open file %s", filename)
		return
	}
	// Split it up to separate lines
	lines := strings.Split(string(content), "\n")
	logChannel <- fmt.Sprintf("Loaded %d lines from %s", len(lines), filename)

	// If we got this far then clear out the CPU to be ready to accept new info
	initCpu()
	// Scan line by line
	for i := 0; i < len(lines); i++ {
		// Don't bother with lines that is shorter than 8 characters and
		// doesn't have a : as the fifth character
		if len(lines[i]) >= 8 {
			if string(lines[i][4]) == ":" {
				// Parse the two fields (address/data) and store it into the cpu
				ad64, err := strconv.ParseUint(lines[i][0:4], 16, 64)
				ad := uint(ad64)
				if err != nil || ad > uint(len(cpu.code)) {
					logChannel <- fmt.Sprintf("Invalid address at line %d [%s]", i, lines[i])
					return
				}
				da64, err := strconv.ParseUint(lines[i][5:8], 16, 64)
				da := uint(da64)
				if err != nil || da > 2047 {
					logChannel <- fmt.Sprintf("Invalid data at line %d [%s]", i, lines[i])
					return
				}
				logChannel <- fmt.Sprintf("Seting cpu.code[%d]=%d", ad, da)
				cpu.code[ad] = da
			}
		}
	}
	// All done here, refresh the screen
	redrawCpu()
}

//
//
//
func handleRx(rxData uint) {
	cpu.in[15] = rxData
	redrawCpu()
}

//
//
//
func handleKeyboard(ev termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		if string(ev.Ch) == "s" {
			executeOneOp()
			redrawCpu()
		}
		if string(ev.Ch) == "b" {
			logChannel <- "Run stopped"
			cpu.running = false
			redrawCpu()
		}
		if string(ev.Ch) == "l" {
			loadHex("load1.hex")
		}
		if string(ev.Ch) == "r" {
			logChannel <- "Run started"
			cpu.running = true
		}
		if ev.Key == termbox.KeyCtrlC {
			return true
		}
	}

	return false
}

//
//
//
func main() {
	// Create UDP logger
	logChannel = make(chan string)
	go udpLog(logChannel)
	logChannel <- "AYTABTU started"

	// Start listening on TCP for UART simulator
	uartChannel = make(chan uint)
	go uartListener(uartChannel)

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

	// Run the polling of termbox key events inside a goroutine and
	// send them over a channel
	keyboardChannel := make(chan termbox.Event)
	go keyboardPoller(keyboardChannel)

	initCpu()
	redrawCpu()
	termbox.Flush()
	exit := false

	refreshTicker := time.NewTicker(time.Millisecond * 100)

	for !exit {
		if cpu.running {
			select {
			case <-refreshTicker.C:
				redrawCpu()
			case rxData := <-uartChannel:
				handleRx(rxData)
			case ev := <-keyboardChannel:
				exit = handleKeyboard(ev)
			default:
				executeOneOp()
			}
		} else {
			select {
			case <-refreshTicker.C:
				// In halted mode nothing needs to be done
			case rxData := <-uartChannel:
				handleRx(rxData)
			case ev := <-keyboardChannel:
				exit = handleKeyboard(ev)
			}
		}
	}

}
