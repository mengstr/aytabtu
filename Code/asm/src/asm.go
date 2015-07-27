package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
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

const (
	LABEL_lbl = iota // Standard label from code
	LABEL_def = iota // Label created by .DEF
)

type Label struct {
	address   int
	labelType int
}

var labels = make(map[string]Label)
var pass int = 1
var pc int = 0
var lineSource string = ""
var lineNo int = 0

//
// Prints an error message explaining the problem and then exits the program
//
func dieError(msg string, extra string) {
	fmt.Printf("%04x: $??? %%???????????   %s \n", pc, lineSource)
	fmt.Printf("ERROR: %s at line %d: %s\n", msg, lineNo, extra)
	os.Exit(1)
}

//
//
//
func addLabel(name string, pc int, labelType int) bool {
	if string(name[0]) == ";" {
		return false
	}
	name = strings.ToUpper(name)
	// label must start with A..Z and then only A..Z0..9
	ok, _ := regexp.MatchString("^[A-Z]+[A-Z0-9]*$", name)
	if !ok {
		dieError("Invalid label name", name)
	}

	// If in pass 1 don't allow redefinition of a label
	_, exists := labels[name]
	if exists && pass == 1 {
		dieError("Redefinition of label", name)
	}

	labels[name] = Label{pc, labelType}
	return true
}

//
//
//
func convertValue(arg string, minVal int, maxVal int) (int, bool) {
	if arg == "" {
		dieError("Missing value", "")
	}

	// Is decimal value?
	if string(arg[0]) >= "0" && string(arg[0]) <= "9" {
		if val, err := strconv.ParseInt(arg, 10, 64); err != nil {
			dieError("Invalid decimal value", arg)
		} else {
			if int(val) < minVal {
				dieError("Value too small", arg)
			}
			if int(val) > maxVal {
				dieError("Value too large", arg)
			}
			return int(val), true
		}
	}

	// Is hex value?
	if string(arg[0]) == "$" {
		if val, err := strconv.ParseInt(arg[1:], 16, 64); err != nil {
			dieError("Invalid hex value", arg)
		} else {
			if int(val) < minVal {
				dieError("Value too small", arg)
			}
			if int(val) > maxVal {
				dieError("Value too large", arg)
			}
			return int(val), true
		}
	}

	// Is binary value?
	if string(arg[0]) == "%" {
		if val, err := strconv.ParseInt(arg[1:], 2, 64); err != nil {
			dieError("Invalid binary value", arg)
		} else {
			if int(val) < minVal {
				dieError("Value too small", arg)
			}
			if int(val) > maxVal {
				dieError("Value too large", arg)
			}
			return int(val), true
		}
	}

	tmpval, exists := labels[arg]
	if !exists && pass == 2 {
		dieError("Symbol not found", arg)
	}

	val := tmpval.address
	if int(val) < minVal {
		dieError("Value too small", arg)
	}
	if int(val) > maxVal {
		dieError("Value too large", arg)
	}

	return val, exists
}

//
// Returns true if .END directive found
//
func handleDirectives(tokens []string, tp int) bool {
	if tokens[tp] == ".END" {
		return true
	}

	if tokens[tp] == ".DEF" {
		val, ok := convertValue(tokens[tp+2], 0, 65535)
		if !ok {
			dieError("Invalid value", tokens[tp+2])
		}
		_ = addLabel(tokens[tp+1], val, LABEL_def)
		return false
	}

	if tokens[tp] == ".ORG" {
		val, ok := convertValue(tokens[tp+1], 0, 65535)
		if !ok {
			dieError("Invalid value", tokens[tp+1])
		}
		pc = val
		return false
	}
	dieError("Invalid directive", tokens[tp])
	return false
}

//
// Print the symbol table in sorted order grouped in
// .defined values and "in-code" labels
//
func showSymbolTable() {
	// Sort the keys into a new array
	sortedLabels := make([]string, 0, len(labels))
	for k := range labels {
		sortedLabels = append(sortedLabels, k)
	}
	sort.Strings(sortedLabels)

	// Print the labels map indexed by the sorted array
	fmt.Println()
	fmt.Println("Symbol Table")
	fmt.Println("============")

	// Start with the .defined values
	for _, k := range sortedLabels {
		if labels[k].labelType == LABEL_def {
			fmt.Printf(".DEF %-16s $%04x %5d\n", k, labels[k].address, labels[k].address)
		}
	}

	fmt.Println()

	// Then print the regular label defined for jumps and branches
	for _, k := range sortedLabels {
		if labels[k].labelType == LABEL_lbl {
			fmt.Printf("     %-16s $%04x\n", k, labels[k].address)
		}
	}
}

//
//
//
func processFile(f *os.File, fOut *os.File) {
	fmt.Printf("Processing pass %d\n", pass)
	pc = 0
	lineNo = 0

	f.Seek(0, 0)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lineSource = scanner.Text()
		// Add some comments at end of line so we don't run out of tokens to process, it's an ugly patch
		// but it simplifies the code
		line := lineSource + " ; ; ; ;"
		lineNo++

		// Must check for a label starting in the first column before tokenizing the line
		// since leading spaces/tabs are removed
		haslabel := (line[0] != 9 && line[0] != 32)

		tokens := strings.Fields(line)
		tp := 0

		// Add the label (if we have one) into the symbol table
		if haslabel {
			if addLabel(tokens[tp], pc, LABEL_lbl) {
				// And continue processing the remainder of the line
				tp++
			}
		}

		// Check for and process compiler directives starting with a dot
		if string(tokens[tp][0]) == "." {
			if handleDirectives(tokens, tp) {
				// .END found, so exit the loop
				break
			}
			// The command was processed, go fetch the next line
			continue
		}

		// If this is a comment go fetch the next line
		if string(tokens[tp][0]) == ";" {
			continue
		}

		// By now the current token should be an opcode
		_, exists := opcodes[tokens[tp]]
		if !exists {
			dieError("Invalid opcode", tokens[tp])
		}
		// Start by retreiving the base opcode value
		op := opcodes[tokens[tp]].base
		arg1val := 0

		// Check for different requred types of the first argument, retreive the argument value
		// and merge it into the opcode base value

		if opcodes[tokens[tp]].arg1 == P_of6 {
			loc, _ := convertValue(tokens[tp+1], 0, 255)
			if pass == 2 {
				ofs := loc - pc
				fmt.Printf("Offset is %d\n", ofs)
				if ofs < 0 {
					ofs = 63 + ofs
				}
				fmt.Printf("Offset is %d\n", ofs)
				//uofs := uint(ofs)
				op |= (ofs << uint(opcodes[tokens[tp]].arg1pos))
			}
		}

		if opcodes[tokens[tp]].arg1 == P_va8 {
			arg1val, _ := convertValue(tokens[tp+1], 0, 255)
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_mem {
			arg1val, _ := convertValue(tokens[tp+1], 0, 255)
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_bit {
			arg1val, _ := convertValue(tokens[tp+1], 0, 7)
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if opcodes[tokens[tp]].arg1 == P_io {
			arg1val, _ := convertValue(tokens[tp+1], 0, 15)
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
				dieError("Invalid register name", tokens[tp+1])
				os.Exit(1)
			}
			op |= (arg1val << uint(opcodes[tokens[tp]].arg1pos))
		}

		if pass == 2 {
			fmt.Printf("%04x : $%03x %%%011b   %s \n", pc, op, op, lineSource)
			s := fmt.Sprintf("%04x:%03x   %s \n", pc, op, lineSource)
			fOut.WriteString(s)
		}

		pc++
	}

}

//
//
//
func main() {
	infile := os.Args[1]
	fileparts := strings.Split(infile, ".")
	outfile := fileparts[0] + ".hex"

	fIn, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer fIn.Close()

	fOut, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer fOut.Close()

	pass = 1
	processFile(fIn, nil)
	pass = 2
	processFile(fIn, fOut)

	showSymbolTable()

}
