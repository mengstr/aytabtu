package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type opcode struct {
	base    int
	arg1    int
	arg1pos int
	arg2    int
	arg2pos int
}

//
// P_none - Not allowed
// P_of6  - 6 bits - label or value -31..+32
// P_va8  - 8 bits - const or value 0..255
// P_mem  - 8 bits - const or value 0..255
// P_reg  - 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
// P_bit  - 3 bits - const or value 0..7
// P_io   - 4 bits - const or value 0..15
//
//
// Labels starts in the first column
// Directives have at least one leading space or tab and starts with a dot
// Opcodes have at least one leading space or tab
// Numbers can be specified as decimal, hex with leading $ or binary with leading %
//

const (
	P_none = iota // Not allowed
	P_of6  = iota // 6 bits - label or value -31..+32
	P_va8  = iota // 8 bits - const or value 0..255
	P_mem  = iota // 8 bits - const or value 0..255
	P_reg  = iota // 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
	P_bit  = iota // 3 bits - const or value 0..7
	P_io   = iota // 4 bits - const or value 0..15
)

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
	"TEST":    {base: 0x580, arg1: P_reg, arg1pos: 3, arg2: P_bit, arg2pos: 0},   //  1   0   1     1   1   b2  b1    b0  r2  r1  r0
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

var debug int = 1
var labels = make(map[string]int)
var pc int = 0
var lineNo int = 0

//
//
//
func addLabel(name string, pc int) bool {
	if string(name[0]) == ";" {
		return false
	}
	name = strings.ToUpper(name)
	// label must start with A..Z and then only A..Z0..9
	ok, _ := regexp.MatchString("^[A-Z]+[A-Z0-9]*$", name)
	if !ok {
		fmt.Printf("Invalid label at line %d: %s\n", lineNo, name)
		os.Exit(1)
	}

	// Don't redefine a label
	_, exists := labels[name]
	if exists {
		fmt.Printf("Redefinition of label at line %d: %s\n", lineNo, name)
		os.Exit(1)
	}

	labels[name] = pc
	return true
}

//
//
//
func convertValue(arg string) (int, bool) {
	if arg == "" {
		fmt.Printf("Missing value at line %d\n", lineNo)
		os.Exit(1)
	}

	// Is decimal value?
	if string(arg[0]) >= "0" && string(arg[0]) <= "9" {
		if val, err := strconv.ParseInt(arg, 10, 64); err != nil {
			fmt.Printf("Invalid decimal value at line %d: %s\n", lineNo, arg)
			os.Exit(1)
		} else {
			return int(val), true
		}
	}

	// Is hex value?
	if string(arg[0]) == "$" {
		if val, err := strconv.ParseInt(arg[1:], 16, 64); err != nil {
			fmt.Printf("Invalid hex value at line %d: %s\n", lineNo, arg)
			os.Exit(1)
		} else {
			return int(val), true
		}
	}

	// Is binary value?
	if string(arg[0]) == "%" {
		if val, err := strconv.ParseInt(arg[1:], 2, 64); err != nil {
			fmt.Printf("Invalid binary value at line %d: %s\n", lineNo, arg)
			os.Exit(1)
		} else {
			return int(val), true
		}
	}

	val, exists := labels[arg]

	return val, exists
}

//
//
//
func handleDirectives(tokens []string, tp int) {
	if tokens[tp] == ".DEF" {
		val, ok := convertValue(tokens[tp+2])
		if !ok {
			fmt.Printf("Invalid value at line %d: %s\n", lineNo, tokens[tp+2])
			os.Exit(1)
		}
		_ = addLabel(tokens[tp+1], val)
		if debug > 0 {
			fmt.Printf("%s defined as %04x\n", tokens[tp+1], int(val))
		}
		return
	}

	if tokens[tp] == ".ORG" {
		val, ok := convertValue(tokens[tp+1])
		if !ok {
			fmt.Printf("Invalid value at line %d: %s\n", lineNo, tokens[tp+1])
			os.Exit(1)
		}
		pc = val
		if debug > 0 {
			fmt.Printf("PC set to %04x\n", pc)
		}
		return
	}
	fmt.Printf("Invalid directive at line %d: %s\n", lineNo, tokens[tp])
	os.Exit(1)

}

//
//
//
func main() {
	file, err := os.Open("test.asm")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawline := scanner.Text()
		// Add some comments at end of line so we don't run out of tokens to process, it's an ugly patch
		// but it simplifies the code
		line := rawline + " ; ; ; ;"
		lineNo++

		// Must check for a label starting in the first column before tokenizing the line
		// since leading spaces/tabs are removed
		haslabel := (line[0] != 9 && line[0] != 32)

		tokens := strings.Fields(line)
		tp := 0
		if debug > 2 {
			fmt.Printf("Label is %t Fields are: %q\n", haslabel, tokens)
		}

		// Process the label if we had one
		if haslabel {
			if addLabel(tokens[tp], pc) {
				tp++
			}
		}

		// Check for and process compiler directives starting with a dot
		if string(tokens[tp][0]) == "." {
			handleDirectives(tokens, tp)
			continue
		}

		// Check if this is a comment
		if string(tokens[tp][0]) == ";" {
			continue
		}

		// By now the current token should be an opcode
		_, exists := opcodes[tokens[tp]]
		if !exists {
			fmt.Printf("Invalid opcode at line %d: %s\n", lineNo, tokens[tp])
			os.Exit(1)
		}
		// Start by retreiving the base opcode value
		op := opcodes[tokens[tp]].base
		arg1val := 0

		// Check for different requred types of the first argument, retreive the argument value
		// and merge it into the base opcode value

		if opcodes[tokens[tp]].arg1 == P_va8 {
			arg1val, _ := convertValue(tokens[tp+1])
			if arg1val < 0 || arg1val > 255 {
				fmt.Printf("Invalid value at line %d: %s\n", lineNo, tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_mem {
			arg1val, _ := convertValue(tokens[tp+1])
			if arg1val < 0 || arg1val > 255 {
				fmt.Printf("Invalid memory address at line %d: %s\n", lineNo, tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_bit {
			arg1val, _ := convertValue(tokens[tp+1])
			if arg1val < 0 || arg1val > 7 {
				fmt.Printf("Invalid bit value at line %d: %s\n", lineNo, tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_io {
			arg1val, _ := convertValue(tokens[tp+1])
			if arg1val < 0 || arg1val > 15 {
				fmt.Printf("Invalid io-port %d: %s\n", lineNo, tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		// P_reg  - 3 bits - A=0 B=1 C=2 D=3 X=4 Y=5 JH=6 JL=7
		if opcodes[tokens[tp]].arg1 == P_reg {
			switch tokens[tp+1] {
			case "A":
				arg1val = 0
			case "B":
				arg1val = 1
			case "C":
				arg1val = 2
			case "D":
				arg1val = 3
			case "X":
				arg1val = 4
			case "Y":
				arg1val = 5
			case "JH":
				arg1val = 6
			case "JL":
				arg1val = 7
			default:
				fmt.Printf("Invalid register name at line %d: %s\n", lineNo, tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		fmt.Printf("%04x: $%03x %%%011b   %s \n", pc, op, op, rawline)

		pc++
	}

}
