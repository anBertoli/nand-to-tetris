// push constant 3030
@3030
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop pointer 0
@SP
M=M-1
A=M
D=M
@3
M=D

// push constant 3040
@3040
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop pointer 1
@SP
M=M-1
A=M
D=M
@4
M=D

// push constant 32
@32
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop this 2
@2
D=A
@THIS
A=M
D=A+D
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
A=M
M=D

// push constant 46
@46
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop that 6
@6
D=A
@THAT
A=M
D=A+D
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
A=M
M=D

// push pointer 0
@3
D=M
@SP
A=M
M=D
@SP
M=M+1

// push pointer 1
@4
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

// push this 2
@2
D=A
@THIS
A=M
A=A+D
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

// push that 6
@6
D=A
@THAT
A=M
A=A+D
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
