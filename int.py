#!/bin/python3

#developed at https://github.com/WerIstLuka/int

import sys
import numpy
def Help():
	print("""\
Convert any base to any other
Usage:
	int [options] [integers]
options:
	-b		add binary prefix to all integers
	-x		add hexadecimal prefix to all integers
	-o		add octal prefix to all integers
	-Bx		where x is the base of the integers
	-Ox		set the output to any base, also works with b, x and o
	-h		show this Help text\
""")
	quit()

def Exit_Too_Many_Prefixes():
	print("Conflicting Prefixes")
	quit()
def CheckOptions(Arguments):
	InputBase = 10
	OutputBase = 10
	for i in range(len(Arguments)):
		match Arguments[i][:2]:
			case "-b":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 2
			case "-x":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 16
			case "-o":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				InputBase = 8
			case "-B":
				if InputBase != 10:
					Exit_Too_Many_Prefixes()
				if Arguments[i][2:] == "":
					print(f"No Base given: {Arguments[i]}")
					quit()
				InputBase = int(Arguments[i][2:])
			case "-O":
				if OutputBase != 10:
					print("Output Base given twice")
					quit()
				OutputBase = int(Arguments[i][2:])
			case _:
				pass
	if "-h" in Arguments:
		Help()
	return(InputBase, OutputBase)

def GetIntegers(Arguments, InputBase, OutputBase):
	if OutputBase != 10:
		Arguments.remove(f"-O{OutputBase}")
	DecimalIntegers = []
	if InputBase != 10:
		for i in range(len(Arguments)):
			try:
				int(Arguments[i], InputBase)
			except Exception:
				print(f"Operation not possible: {Arguments[i]}")
				quit()
			DecimalIntegers.append(int(Arguments[i], InputBase))
	else:
		for i in range(len(Arguments)):
			match Arguments[i][:2]:
				case "0b":
					DecimalIntegers.append(int(Arguments[i][2:], 2))
				case "0x":
					DecimalIntegers.append(int(Arguments[i][2:], 16))
				case _:
					try:
						int(Arguments[i])
					except Exception:
						print(f"Operation not possible: {Arguments[i]}")
						quit()
					DecimalIntegers.append(int(Arguments[i]))
	return DecimalIntegers

def OutputIntegers(DecimalIntegers, OutputBase):
	for i in range(len(DecimalIntegers)):
		print(numpy.base_repr(DecimalIntegers[i], int(OutputBase)))

Arguments = sys.argv[1:]
if len(Arguments) == 0:
	Help()
Options = CheckOptions(Arguments)
DecimalOutput =  GetIntegers(Arguments, Options[0], Options[1])
OutputIntegers(DecimalOutput, Options[1])