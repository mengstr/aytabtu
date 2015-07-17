	.ORG 128	; $0080
	.ORG $100	; $0100
	.ORG %0000001000000000	;$0200

	.DEF	MATS $65
	.DEF	ERIK 1


FOO	NOP
BAR
;Only a comment at start of line
				; A comment with a tab first
	NOP	; Command With comment 1
	LDAI	32
	LDAI	$10
	LDAI	%00000001
	LDAI	255		; Comment
	LDAI	$FF		; Comment
	LDAI	%11111111	;Comment
	LDAI	MATS
	LDAI 	OLLE
	HALT


	