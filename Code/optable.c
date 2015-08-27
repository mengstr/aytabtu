
//INST     BASEOP    P1  P2  P1start  P2start

{'HALT   ', 0x000, P_none, P_none, 0, 0 },  //  0   0   0   0   0   x   x   x   x   x   x
{'NOP    ', 0x080, P_none, P_none, 0, 0 },  //  0   0   0     1   0   x   x     x   x   x   x
{'SJMP   ', 0x0C0, P_of6 , P_none, 5, 0 },  //  0   0   0     1   1   o5  o4    o3  o2  o1  o0
{'BRAZ   ', 0x100, P_of6 , P_none, 5, 0 },  //  0   0   1     0   0   o5  o4    o3  o2  o1  o0
{'BRANZ  ', 0x140, P_of6 , P_none, 5, 0 },  //  0   0   1     0   1   o5  o4    o3  o2  o1  o0
{'BRAC   ', 0x180, P_of6 , P_none, 5, 0 },  //  0   0   1     1   0   o5  o4    o3  o2  o1  o0
{'BRANC  ', 0x1C0, P_of6 , P_none, 5, 0 },  //  0   0   1     1   1   o5  o4    o3  o2  o1  o0
{'LDAI   ', 0x200, P_va8,  P_none, 7, 0 },  //  0   1   0     id7 id6 id5 id4   id3 id2 id1 id0
{'LDAZP  ', 0x300, P_mem,  P_none, 7, 0 },  //  0   1   1     a7  a6  a5  a4    a3  a2  a1  a0
{'STAZP  ', 0x400, P_mem,  P_none, 7, 0 },  //  1   0   0     a7  a6  a5  a4    a3  a2  a1  a0
{'CLR    ', 0x500, P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   0   0     0   r2  r1  r0
{'SETFF  ', 0x508  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   0   0     1   r2  r1  r0
{'NOT    ', 0x510  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   0   1     0   r2  r1  r0
{'OR     ', 0x518  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   0   1     1   r2  r1  r0
{'AND    ', 0x520  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   1   0     0   r2  r1  r0
{'XOR    ', 0x528  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   1   0     1   r2  r1  r0
{'INC    ', 0x530  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   1   1     0   r2  r1  r0
{'DEC    ', 0x538  P_reg,  P_none, 2, 0 },  //  1   0   1     0   0   1   1     1   r2  r1  r0
{'ADD    ', 0x540  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   0   0     0   r2  r1  r0
{'SUB    ', 0x548  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   0   0     1   r2  r1  r0
{'ADDC   ', 0x550  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   0   1     0   r2  r1  r0
{'SUBC   ', 0x558  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   0   1     1   r2  r1  r0
{'LSHIFT ', 0x560  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   1   0     0   r2  r1  r0
{'RSHIFT ', 0x568  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   1   0     1   r2  r1  r0
{'LSHIFTC', 0x570  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   1   1     0   r2  r1  r0
{'RSHIFTC', 0x578  P_reg,  P_none, 2, 0 },  //  1   0   1     0   1   1   1     1   r2  r1  r0
{'MOVE   ', 0x580, P_reg,  P_reg,  5, 2 },  //  1   0   1     1   0   rs2 rs1   rs0 rd2 rd1 rd0
{'TEST   ', 0x580, P_reg,  P_bit,  2, 5 },  //  1   0   1     1   1   b2  b1    b0  r2  r1  r0
{'PEEK   ', 0x600, P_io ,  P_none, 3, 0 }, //  1   1   0     0   0   0   0     p3  p2  p1  p0
{'POKE   ', 0x610, P_io ,  P_none, 3, 0 }, //  1   1   0     0   0   0   1     p3  p2  p1  p0
{'LDAXY  ', 0x620, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   0   0   0
{'STAXY  ', 0x621, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   0   0   1
{'JUMP   ', 0x622, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   0   1   0
{'CALL   ', 0x623, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   0   1   1
{'RET    ', 0x624, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   1   0   0
{'CLRC   ', 0x626, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   1   1   0
{'SETC   ', 0x627, P_none, P_none, 0, 0 },  //  1   1   0     0   0   1   0     0   1   1   1

