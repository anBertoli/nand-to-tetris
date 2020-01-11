
// push argument 1
@1
D=A
@ARG
A=M
A=A+D
D=M
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

// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop that 0              
@0
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

// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop that 1              
@1
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

// push argument 0
@0
D=A
@ARG
A=M
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1

// push constant 2
@2
D=A
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

// pop argument 0          
@0
D=A
@ARG
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

// label MAIN_LOOP_START
(UNDEFINED$MAIN_LOOP_START)

// push argument 0
@0
D=A
@ARG
A=M
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1

// if-goto COMPUTE_ELEMENT 
@SP
AM=M-1
D=M
@UNDEFINED$COMPUTE_ELEMENT
D;JGT
D;JLT

// goto END_PROGRAM        
@UNDEFINED$END_PROGRAM
0;JMP

// label COMPUTE_ELEMENT
(UNDEFINED$COMPUTE_ELEMENT)

// push that 0
@0
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

// push that 1
@1
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

// pop that 2              
@2
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

// push pointer 1
@4
D=M
@SP
A=M
M=D
@SP
M=M+1

// push constant 1
@1
D=A
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

// pop pointer 1           
@SP
M=M-1
A=M
D=M
@4
M=D

// push argument 0
@0
D=A
@ARG
A=M
A=A+D
D=M
@SP
A=M
M=D
@SP
M=M+1

// push constant 1
@1
D=A
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

// pop argument 0          
@0
D=A
@ARG
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

// goto MAIN_LOOP_START
@UNDEFINED$MAIN_LOOP_START
0;JMP

// label END_PROGRAM
(UNDEFINED$END_PROGRAM)

(ENDLOOP)
@ENDLOOP
0;JMP
