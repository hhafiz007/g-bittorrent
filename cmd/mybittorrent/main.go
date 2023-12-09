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
func decodeInt(bencodedString string,idx int) (interface{},int, error){
	var lastIndex int

	for i := idx; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			lastIndex = i
			break
		}
	}

	num, err := strconv.Atoi(bencodedString[1:lastIndex])
	
	return num,lastIndex, err
	

}


func decodeString(bencodedString string,idx int) (interface{},int, error){
	var firstColonIndex int

	for i := idx; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0,err
	}
	lastIndex := firstColonIndex+length
	

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length],lastIndex, nil
}

func decodeList(bencodedString string,idx int) (interface{}, int,error){
	slice:=  make([]interface{},0,4)

	i := idx+1

	for bencodedString[i] != 'e' {
		decoded, newIdx,_ := decodeBencode(bencodedString,i)
		fmt.Println(i)
		
		slice = append(slice,decoded)
		
		i = newIdx+1
		idx = i
		
		}
	
	return slice,idx,nil




}


func decodeBencode(bencodedString string,idx int) (interface{},int, error) {
	if unicode.IsDigit(rune(bencodedString[idx])) {
		return decodeString(bencodedString,idx)
	}else if (bencodedString[idx]) == 'i' { 
		return decodeInt(bencodedString,idx)
	
	}else if (bencodedString[idx]) == 'l' {
	
		
		slice,newIdx,err := decodeList(bencodedString,idx)
		
		return slice,newIdx,err
		

	}else {
		return "", idx,fmt.Errorf("Only strings are supported at the moment")
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
		
		decoded,_,err := decodeBencode(bencodedValue,0)
	
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
