// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

@SCREEN
D=A
@ACTUALPIXELPOS	// pointer to actual screen bit
M=D

@BEGINSCREENPOS	// pointer to first valid screen bit
M=D
			
@BEGINSCREENPOS	// pointer to last valid screen bit
D=M
@8190
D=D+A
@ENDSCREENPOS		
M=D


(START)
// check if key is pressed
@KBD
D=M
@KEYPRESS
D;JGT
@NOKEYPRESS
0;JMP


// key pressed
(KEYPRESS)
@ACTUALPIXELPOS
A=M
M=-1	// darken pixel

D=A 	// check if end of screen is reached
@ENDSCREENPOS
A=M
D=A-D
@ENDSCREEN
A-D;JGE

(NOTENDSCREEN)
@ACTUALPIXELPOS
M=M+1
@START
0;JMP

(ENDSCREEN)
@START
0;JMP



// key not pressed
(NOKEYPRESS)
@ACTUALPIXELPOS
A=M
M=0		// clear pixel
D=A 	// check if begin of screen is reached
@BEGINSCREENPOS
A=M
D=A-D
@BEGINSCREEN
D;JGE

(NOTBEGINSCREEN)
@ACTUALPIXELPOS
M=M-1
@START
0;JMP

(BEGINSCREEN)
@START
0;JMP
