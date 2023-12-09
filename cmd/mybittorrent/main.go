package main

import (
	// Uncomment this line to pass the first stage
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeInt(bencodedString string) (interface{}, error){
	var lastIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			lastIndex = i
			break
		}
	}

	num, err := strconv.Atoi(bencodedString[1:lastIndex])
	return num, err
	

}


func decodeString(bencodedString string) (interface{}, error){
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}

func decodeList(bencodedString string) (interface{}, error){
	 slice:=  make([]interface{},0,4)
	i := 1
	for i < len(bencodedString)-1{
		if unicode.IsDigit(rune(bencodedString[i])) {

			
			decoded, err := decodeString(bencodedString[i:])
			
			if err != nil {
				slice = append(slice, decoded)
			}
			var length string
			var decodedLength int
			if str, ok := decoded.(string); ok {
				length = string(len(str))
				decodedLength = len(str)
			 }

			
			i += decodedLength+len(length)+1

			
		}else if (bencodedString[i]) == 'i' { 
			decoded, err := decodeInt(bencodedString[i:])
			if err != nil {
				slice = append(slice, decoded)
			}
			for bencodedString[i] != 'e'{
				i++
			}
			i++


	
		}else {
			return "", fmt.Errorf("Only strings are supported at the moment")
		}

	


		
	}
	fmt.Println(len(bencodedString),slice)
	return slice,nil




}


func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	}else if (bencodedString[0]) == 'i' { 
		return decodeInt(bencodedString)
	
	}else if (bencodedString[0]) == 'l' {
	
		
		slice,err := decodeList(bencodedString)
		print("hi",slice)
		return slice,err
		

	}else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage
		//
		bencodedValue := os.Args[2]
		
		decoded, err := decodeBencode(bencodedValue)
		print("main",decoded)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
