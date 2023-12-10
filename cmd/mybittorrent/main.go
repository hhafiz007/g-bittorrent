package main

import (
	// Uncomment this line to pass the first stage
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	"bytes"
	"io/ioutil"
	bencode "github.com/jackpal/bencode-go" // Available if you need it!
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

	num, err := strconv.Atoi(bencodedString[idx+1:lastIndex])
	
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

	lengthStr := bencodedString[idx:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 1,err
	}
	lastIndex := firstColonIndex+length
	

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length],lastIndex, nil
}

func decodeList(bencodedString string,idx int) (interface{}, int,error){
	slice:=  make([]interface{},0,4)

	i := idx+1

	for bencodedString[i] != 'e' {
		decoded, newIdx,_ := decodeBencode(bencodedString,i)
		
		slice = append(slice,decoded)
		
		i = newIdx+1
		//fmt.Println(newIdx)
	
		idx = i
		
		}
	
	return slice,idx,nil




}

func decodeDict(bencodedString string,idx int) (interface{}, int,error){
	myMap := make(map[string]interface{})


	i := idx+1


	for bencodedString[i] != 'e' {
		key,newIdx,_ := decodeBencode(bencodedString,i)
		i = newIdx+1
		value,newIdx,_ := decodeBencode(bencodedString,i)
		i = newIdx+1
		var strKey string
		if str, ok := key.(string); ok {
			strKey = string(str)
	
			// Now you can use 'str' as a string
		}
		myMap[strKey] = value
		idx = i

	}

	return myMap,idx,nil

}


func decodeBencode(bencodedString string,idx int) (interface{},int, error) {
	if unicode.IsDigit(rune(bencodedString[idx])) {
		return decodeString(bencodedString,idx)
	}else if (bencodedString[idx]) == 'i' { 
		return decodeInt(bencodedString,idx)
	
	}else if (bencodedString[idx]) == 'l' {
	
		
		slice,newIdx,err := decodeList(bencodedString,idx)
		
		return slice,newIdx,err
		

	}else if (bencodedString[idx]) == 'd' {
	
		
		return decodeDict(bencodedString,idx)

		

	}else {
		return "", idx,fmt.Errorf("Only strings are supported at the moment")
	}
}

type TorrentFile struct {
    Announce string
    Info     struct {
        Length      int
        Name        string
        PieceLength int `bencode:"piece length"`
        Pieces      string
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
	} else if command == "info"{

		torrentFilePath := os.Args[2]
		torrentData, err := ioutil.ReadFile(torrentFilePath)
		if err != nil {
			fmt.Println(err) }
		var torrent  TorrentFile
		reader := bytes.NewReader(torrentData)
		err = bencode.Unmarshal(reader, &torrent)
	if err != nil {
		fmt.Println(err)
	}
		fmt.Println("Tracker URL:",torrent.Announce)
		fmt.Println("Length:",torrent.Info.Length)





	} else{
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
