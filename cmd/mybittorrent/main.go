package main

import (
	// Uncomment this line to pass the first
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"

	//	"net/http"
	//		"net/url"
	bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeInt(bencodedString string, idx int) (interface{}, int, error) {
	var lastIndex int

	for i := idx; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			lastIndex = i
			break
		}
	}

	num, err := strconv.Atoi(bencodedString[idx+1 : lastIndex])

	return num, lastIndex, err

}

func decodeString(bencodedString string, idx int) (interface{}, int, error) {
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
		return "", 1, err
	}
	lastIndex := firstColonIndex + length

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], lastIndex, nil
}

func decodeList(bencodedString string, idx int) (interface{}, int, error) {
	slice := make([]interface{}, 0, 4)

	i := idx + 1

	for bencodedString[i] != 'e' {
		decoded, newIdx, _ := decodeBencode(bencodedString, i)

		slice = append(slice, decoded)

		i = newIdx + 1
		//fmt.Println(newIdx)

		idx = i

	}

	return slice, idx, nil

}

func decodeDict(bencodedString string, idx int) (interface{}, int, error) {
	myMap := make(map[string]interface{})

	i := idx + 1

	for bencodedString[i] != 'e' {
		key, newIdx, _ := decodeBencode(bencodedString, i)
		i = newIdx + 1
		value, newIdx, _ := decodeBencode(bencodedString, i)
		i = newIdx + 1
		var strKey string
		if str, ok := key.(string); ok {
			strKey = string(str)

			// Now you can use 'str' as a string
		}
		myMap[strKey] = value
		idx = i

	}

	return myMap, idx, nil

}

func decodeBencode(bencodedString string, idx int) (interface{}, int, error) {
	if unicode.IsDigit(rune(bencodedString[idx])) {
		return decodeString(bencodedString, idx)
	} else if (bencodedString[idx]) == 'i' {
		return decodeInt(bencodedString, idx)

	} else if (bencodedString[idx]) == 'l' {

		slice, newIdx, err := decodeList(bencodedString, idx)

		return slice, newIdx, err

	} else if (bencodedString[idx]) == 'd' {

		return decodeDict(bencodedString, idx)

	} else {
		return "", idx, fmt.Errorf("Only strings are supported at the moment")
	}
}

func encodeToBencode(info Info) (string, error) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "d%s:lengthi%se", strconv.Itoa(6), strconv.Itoa(info.Length))
	fmt.Fprintf(&builder, "%s:name%s:%s", strconv.Itoa(4), strconv.Itoa(len(info.Name)), info.Name)
	fmt.Fprintf(&builder, "%s:piece lengthi%se", strconv.Itoa(12), strconv.Itoa(info.PieceLength))
	fmt.Fprintf(&builder, "%s:pieces%s:%se", strconv.Itoa(6), strconv.Itoa(len(info.Pieces)), info.Pieces)
	return builder.String(), nil

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

type Info struct {
	Length      int
	Name        string
	PieceLength int `bencode:"piece length"`
	Pieces      string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage
		//
		bencodedValue := os.Args[2]

		decoded, _, err := decodeBencode(bencodedValue, 0)

		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)

		fmt.Println(string(jsonOutput))
	} else if command == "info" {

		torrentFilePath := os.Args[2]
		torrentData, err := ioutil.ReadFile(torrentFilePath)
		//testOutput, _ := json.Marshal(string(torrentData))

		//fmt.Println("Info Hash:",string(testOutput))

		if err != nil {
			fmt.Println(err)
		}
		var torrent TorrentFile
		reader := bytes.NewReader(torrentData)
		err = bencode.Unmarshal(reader, &torrent)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Tracker URL:", torrent.Announce)
		fmt.Println("Length:", torrent.Info.Length)

		bencodedInfo, _ := encodeToBencode(torrent.Info)
		h := sha1.New()
		io.WriteString(h, bencodedInfo)
		infoHash := h.Sum(nil)
		//jsonOutput, _ := json.Marshal(infoHash)
		fmt.Println("Info Hash:", fmt.Sprintf("%x", infoHash))
		fmt.Println("Piece Length:", torrent.Info.PieceLength)
		fmt.Println("Piece Length:")
		pieces := torrent.Info.Pieces
		for i := 0; i < len(pieces); i += 20 {
			hash := []byte(pieces[i : i+20])
			fmt.Println(hex.EncodeToString(hash))
		}

	} else if command == "peers" {
		getTracker(string(os.Args[2]))

	} else if command == "handshake" {
		mP := []byte{}
		getHandshake(os.Args[3], 0, &mP, os.Args[2], 0)

	} else if command == "download_piece" {
		downloadPiece()

	} else if command == "download" {
		finalDowwnload()

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
