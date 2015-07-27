	; Try setting the origin/program counter to some 
	; dec/hex/binary values

	.ORG $0000

	CLR	A
	CLR	B
	CLR	C
	CLR	D

LOOP
	INC	A
	BRANZ	LOOP
	INC	B
	BRANZ	LOOP
	INC	C
	BRANZ	LOOP
	INC	D
	SJMP	LOOP







	