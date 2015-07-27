	; Try setting the origin/program counter to some 
	; dec/hex/binary values

	.ORG $0000				// Let's start at the beginning of code area

	.DEF	UART	15		// This is the 300 baud UART port

	CLR	A 	`				// Zero out the registers we'll work with since 
	CLR	B 					// the CPU reset in real life doesn't do that for us
	CLR	C
	CLR	D

LOOP
	INC	A
	BRANZ	LOOP
	INC	B
	BRANZ	LOOP
	LDAI	64				// Every 256*256 LOOPs send a @ character to the UART
	POKE	UART
	INC	C
	BRANZ	LOOP
	INC	D
	SJMP	LOOP







	
