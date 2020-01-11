// push constant 111
@111
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 333
@333
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 888
@888
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop static 8
@SP
M=M-1
A=M
D=M
@24
M=D
// pop static 3
@SP
M=M-1
A=M
D=M
@19
M=D
// pop static 1
@SP
M=M-1
A=M
D=M
@17
M=D
// push static 3
@19
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@17
D=M
@SP
A=M
M=D
@SP
M=M+1
// sub
@SP
AM=M-1
D=M
A=A-1
M=M-D
D=A+1
@SP
M=D
// push static 8
@24
D=M
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
AM=M-1
D=M
A=A-1
M=D+M
D=A+1
@SP
M=D
(ENDLOOP)
@ENDLOOP
0;JMP
