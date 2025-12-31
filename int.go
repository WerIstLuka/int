package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"slices"
	"strconv"
	"math/big"
)

var Version string = "2.0.4"

func Help(){
	fmt.Println(`Convert any base to any other
Usage:
	int [Options] [Integers]
Options:
	-b	add binary prefix to all integers
	-x	add hexadecimal prefix to all integers
	-o	add octal prefix to all integers
	-Bx	where x is the base of the integers
	-Ox	set the output to any base, also works with b, x and o
 	-l	use the same characters for bases below and above 36
	-h	show this Help text`)
	os.Exit(0)
}

func HasPipeInput()bool{
	FileInfo, _ := os.Stdin.Stat()
	return FileInfo.Mode() & os.ModeCharDevice == 0
}

func GetArguments() ([]string, []string){
	Options := []string{}
	Numbers := []string{}
	//read input from pipe if it exists
	if HasPipeInput(){
 		bytes, _ := io.ReadAll(os.Stdin)
		lines := strings.Split((string(bytes)), "\n")
		for i:=0; i<len(lines)-1; i++{
			line := strings.Split(lines[i], " ")
			for j:=0; j<len(line); j++{
				if len(line[j]) == 0{
					continue
				}
				if line[j][0] == '-'{
					Options = append(Options, strings.TrimSpace(line[j]))
				}else{
					Numbers = append(Numbers, strings.TrimSpace(line[j]))
				}
			}
		}
	}
	//get arguments
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] == '-'{
			Options = append(Options, os.Args[i])
		}else{
			Numbers = append(Numbers, os.Args[i])
		}
	}
	if len(Options) == 0 && len(Numbers) == 0{
		Help()
	}
	return Options, Numbers
}

func GetInt(char byte)int64{
	switch char{
		case 'b':
			return 2
		case 'o':
			return 8
		case 'x':
			return 16
		default:
			return 0
	}
}

func GetBase(Option string)int64{
	var Base int64
	var IsValid bool = true
	for i:=0;i<len(Option);i++{
		if !strings.Contains("0123456789", string(Option[i])){
			IsValid = false
		}
	}
	if IsValid{
		Base, _ = strconv.ParseInt(Option, 10, 64)
	} else {
		if len(Option) != 1{
			Base = 0
		} else {
			Base = GetInt(Option[0])
		}
	}
	if Base == 0 || Base > 62{
		fmt.Println("Error: invalid base:", Option)
		os.Exit(1)
	}
	return Base
}


func Parser(Options []string)(int, int64, bool){
	var GotInputBase bool = false
	var GotOutputBase bool = false
	var InputBase int64 = 0
	var OutputBase int64 = 0
	var ForceLong bool = false
	for i:=0; i<len(Options); i++{
		Option := Options[i]
		if len(Option) < 2{
			fmt.Println("Error: invalid option:", Option)
			os.Exit(1)
		}
		if GotInputBase && slices.Contains([]byte{'B', 'b', 'o', 'x'}, Option[1]){
			fmt.Println("Error: Input base was given twice")
			os.Exit(1)
		} else if GotOutputBase && Option[1] == 'O'{
			fmt.Println("Error: Output base was given twice")
			os.Exit(1)
		} else if len(Option) == 2 && slices.Contains([]byte{'b', 'o', 'x'}, Option[1]){
			InputBase = GetInt(Option[1])
			GotInputBase = true
		} else if Option == "-h" || Option == "--help"{
			Help()
		} else if Option == "-v" || Option == "--version"{
			fmt.Println(Version)
			os.Exit(0)
		} else if Option == "-l" || Option == "--long"{
			ForceLong = true
		} else if slices.Contains([]byte{'B', 'O'}, Option[1]) == false{
			fmt.Println("Error: invalid option:", Option)
			os.Exit(1)
		} else if len(Option) < 3{
			fmt.Println("Error: too few arguments:", Option)
			os.Exit(1)
		} else if Option[1] == 'B'{
			InputBase = GetBase(Option[2:len(Option)])
			GotInputBase = true
		} else if Option[1] == 'O'{
			OutputBase = GetBase(Option[2:len(Option)])
			GotOutputBase = true
		}
	}
	if GotInputBase && InputBase <= 1{
		fmt.Println("Error: Input base can't be 1 or less")
		os.Exit(1)
	} else if GotOutputBase && OutputBase <= 1{
		fmt.Println("Error: Output base can't be 1 or less")
		os.Exit(1)
	}
	return int(InputBase), OutputBase, ForceLong
}

func ErrorOperationNotPossible(Operation string){
	fmt.Println("Error: Operation not possible:", Operation)
	os.Exit(1)
}

func ConvertNumbers(Num string, InputBase int, OutputBase int64, ForceLong bool) string{
	if InputBase == 0 && len(Num) >= 2{
		InputBase = int(GetInt(Num[1]))
		if InputBase != 0{
			if len(Num) == 2{
				fmt.Println("Error: No number given:", Num)
				os.Exit(1)
			}
			Num = Num[2:len(Num)]
		}
	}
	if InputBase == 0{
		InputBase = 10
	}
	if OutputBase == 0{
		OutputBase = 10
	}



	var DigitsShort string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var DigitsLong string  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var DigitList []string
	
	if InputBase <= 35{
		DigitList = strings.Split(DigitsShort, "")
	} else {
		DigitList = strings.Split(DigitsLong, "")
	}
	var Index int
	for i:=0;i<len(Num);i++{
		Index = slices.Index(DigitList, string(Num[i]))
		if Index == -1 || ForceLong{
			Index = slices.Index(strings.Split(DigitsLong, ""), string(Num[i]))
			if Index == -1{
				fmt.Println("Error: Illegal Character:", string(Num[i]))
				os.Exit(1)
			}
		}
		if Index >= InputBase{
			fmt.Println("Error: Number", Num, "is not valid for base", InputBase)
			os.Exit(1)
		}
	}
	var OutputDigits *string
	if OutputBase <= 35 && !ForceLong{
		OutputDigits = &DigitsShort
	} else {
		OutputDigits = &DigitsLong
	}

	
	BigNum := new(big.Int)
	BigNum.SetString(Num, InputBase)
	
	BigOutputBase := big.NewInt(OutputBase)
	var Output = []string{}
	IndexBig := new(big.Int)
	if BigNum.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	for BigNum.Cmp(big.NewInt(0)) != 0 {
		Output = append(Output, string(string(*OutputDigits)[IndexBig.Mod(BigNum, BigOutputBase).Int64()]))
		BigNum.Div(BigNum, BigOutputBase)
	}
	slices.Reverse(Output)
	return strings.Join(Output, "")
}

func main() {
	Options, Numbers := GetArguments()
	InputBase, OutputBase, ForceLong := Parser(Options)
	if len(Numbers) == 0{
		fmt.Println("no input")
		os.Exit(0)
	}
	ConvertedNumbers := []string{}
	for i:=0;i<len(Numbers);i++{
		ConvertedNumbers = append(ConvertedNumbers, ConvertNumbers(Numbers[i], InputBase, OutputBase, ForceLong))
	}
	for i:=0;i<len(ConvertedNumbers);i++{
		fmt.Println(ConvertedNumbers[i])
	}
}
