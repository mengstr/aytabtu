	; Try setting the origin/program counter to some 
	; dec/hex/binary values

	.ORG 128		; $0080
	NOP

	.ORG $100		; $0100
	NOP

	.ORG %0000001000000000	; $0200
	NOP

	.DEF	HEX65 		$65
	.DEF	DEC1 		1
	.DEF	BIN254		%11111110


LABEL1	NOP
LABEL2
;Comment at start of line
				; Indented comment
	NOP			; Opcode with comment

	; Try to do immediate loads from direct values and .DEF constants
	LDAI	1
	LDAI	$65
	LDAI	%11111110
	LDAI 	DEC1
	LDAI	HEX65
	LDAI	BIN254

	; Now try a basic variant of all currently allowed opcodes
	HALT	   		; No params
	NOP	    		; No params
	SJMP	 END		; Offset  
	BRAZ	 END  		; Offset
	BRANZ	 END 		; Offset
	BRAC	 END  		; Offset
	BRANC	 END 		; Offset
	LDAI	 0		; Immediate value
	LDAZP	 0		; Zeropage Memory locaion 
	STAZP	 0 		; Zeropage Memory location
	CLR	 A 		; Register 		 
	SETFF	 A  		; Register
	NOT	 A  		; Register  
	OR	 A  		; Register   
	AND	 A  		; Register  
	XOR	 A  		; Register  
	INC	 A  		; Register  
	DEC	 A  		; Register  
	ADD	 A  		; Register  
	SUB	 A  		; Register  
	ADDC	 A  		; Register 
	SUBC	 A  		; Register 
	LSHIFT	 A 		; Register
	RSHIFT	 A 		; Register
	LSHIFTC	 A 		; Register
	RSHIFTC	 A 		; Register
;	MOVE	 A,B		; Src Register, Dst Register  
;	TEST	 A,0		; Register, BitNo  
	PEEK	 0		; I/O location  
	POKE	 0 		; I/O location 
	LDAXY	  		; No params
	STAXY	  		; No params
	JUMP	   		; No params
	CALL	   		; No params
	RET	    		; No params
	CLRC	   		; No params
	SETC	   		; No params
	LDPML	  		; No params
	LDPMH	  		; No params
	STPML	  		; No params
	STPMH	  		; No params
END






	