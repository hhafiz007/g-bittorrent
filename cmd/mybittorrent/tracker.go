package main

import (
	// Uncomment this line to pass the first

	//"encoding/json"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"path"
	"strconv"

	//"strconv"
	//"unicode"
	"bytes"
	"crypto/sha1"
	"io/ioutil"
	"math"

	//"encoding/hex"

	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	// "time"

	bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func getTracker(torrentFilePath string) interface{} {

	torrentData, err := ioutil.ReadFile(torrentFilePath)
	if err != nil {
		fmt.Println(err)
	}
	var torrent TorrentFile
	reader := bytes.NewReader(torrentData)
	err = bencode.Unmarshal(reader, &torrent)
	if err != nil {
		fmt.Println(err)
	}
	bencodedInfo, _ := encodeToBencode(torrent.Info)
	h := sha1.New()
	io.WriteString(h, bencodedInfo)
	infoHash := h.Sum(nil)
	// fmt.Println("Info Hash:", fmt.Sprintf("%x", infoHash))
	peer_id := "00112233445566778899"
	port := "6881"
	uploaded := "0"
	downloaded := "0"
	left := "1"
	compact := "1"

	query := url.Values{}

	query.Add("peer_id", peer_id)
	query.Add("port", port)
	query.Add("uploaded", uploaded)
	query.Add("downloaded", downloaded)
	query.Add("left", left)
	query.Add("compact", compact)
	query.Add("info_hash", string(infoHash))

	url := fmt.Sprintf("%s?%s", torrent.Announce, query.Encode())

	res, err := http.Get(url)

	defer res.Body.Close()

	if err != nil {
		fmt.Println("Oops! Something went wrong:", err)

	}

	body, _ := ioutil.ReadAll(res.Body)
	decoded, _, err := decodeBencode(string(body), 0)
	// fmt.Println(decoded)

	if err != nil {
		fmt.Println(err)

	}

	peers := decoded.(map[string]interface{})["peers"]
	strPeers := []byte(peers.(string))
	var peerIps []string

	for i := 0; i < len(strPeers); i += 6 {

		ip := net.IP((strPeers[i : i+4]))
		port := int((strPeers[i+4]))<<8 + int((strPeers[i+5]))
		fmt.Println(fmt.Sprintf("%s:%d", ip.String(), port))
		peerIps = append(peerIps, fmt.Sprintf("%s:%d", ip.String(), port))
	}
	return peerIps

}

func getPiece(conn net.Conn, myPiece *[]byte, currBlock int, pieceLength int, pieceIndex int, totalLength int) {

	for {
		bitField := make([]byte, 5)
		_, err := conn.Read(bitField)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(int(bitField[0]), int(bitField[4]))
		if int(bitField[4]) == 5 {
			// Message ID for "interested" is 2
			messageID := byte(2)

			// Message length (excluding the length bytes itself)
			messageLength := uint32(1)

			// Create a buffer to hold the message
			messageBuffer := new(bytes.Buffer)

			// Write the message length in big-endian order (4 bytes)
			binary.Write(messageBuffer, binary.BigEndian, messageLength)

			// Write the message ID
			messageBuffer.WriteByte(messageID)

			// Get the final message as a byte slice
			interestedMessage := messageBuffer.Bytes()
			_, err = conn.Write(interestedMessage)
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	if (pieceIndex+1)*pieceLength > totalLength {
		pieceLength = totalLength - ((pieceIndex) * pieceLength)

	}
	fmt.Println(pieceLength)

	totalBlocks := int(math.Ceil(float64(pieceLength) / float64(16*1024)))

	fmt.Println("Total Blocks", totalBlocks)

	// start := 0

	for {
		unchoke := make([]byte, 5)
		_, err := conn.Read(unchoke)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("hi", int(unchoke[4]))
		if int(unchoke[4]) == 1 {

			for i := currBlock; i < totalBlocks; i++ {

				// Message ID for "interested" is 2
				fmt.Println("Welcome to request ", i)
				// duration := 0001 * time.Second
				// time.Sleep(duration)

				messageID := byte(6)

				length := 16 * 1024

				// Message length (excluding the length bytes itself)
				messageLength := uint32(13)
				messageIndex := uint32(pieceIndex)
				messageBegin := uint32((i * (16 * 1024)))

				if (messageBegin + uint32(length)) > uint32(pieceLength) {
					length = pieceLength - int(messageBegin)
				}

				messageBlock := uint32(length)

				// Create a buffer to hold the message
				messageBuffer := new(bytes.Buffer)

				// Write the message length in big-endian order (4 bytes)
				binary.Write(messageBuffer, binary.BigEndian, messageLength)

				// Write the message ID
				messageBuffer.WriteByte(messageID)
				binary.Write(messageBuffer, binary.BigEndian, messageIndex)

				// Use binary.Write to write the integer to the buffer
				binary.Write(messageBuffer, binary.BigEndian, messageBegin)
				binary.Write(messageBuffer, binary.BigEndian, messageBlock)

				// Get the final message as a byte slice
				requestMessage := messageBuffer.Bytes()
				_, err = conn.Write(requestMessage)

				if err != nil {
					fmt.Println(err)
				}
				for {

					reqBlock := make([]byte, 5)
					_, err = conn.Read(reqBlock)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(length)

					if int(reqBlock[4]) == 7 {

						reqPieceMessage := bytes.Repeat([]byte{0}, length+4+4)

						totalReads := 0

						for totalReads < length+4+4 {

							singleByte := make([]byte, 1)

							_, err = conn.Read(singleByte)
							if err != nil {
								fmt.Println(err)
							}
							// totalReads+=1
							reqPieceMessage[totalReads] = singleByte[0]
							totalReads += 1

						}

						*myPiece = append(*myPiece, reqPieceMessage[8:]...)

						// fmt.Println((reqPieceMessage))
						fmt.Println("Welcome to the end 7 for block", i)
						break
						// os.Exit(2)

					}
					// else if int(reqBlock[4]) == 0 {
					// 	fmt.Println("connection choked")
					// 	// os.Exit(2)
					// 	getPiece(conn, myPiece, i, pieceLength)

					// }

				}
				fmt.Println("Welcome to end of request ", i)

			}
			break
		}

	}

}

func getHandshake(peerIp string, downloadP int, myPiece *[]byte, torrentFilePath string, pieceIndex int) {

	torrentData, err := ioutil.ReadFile(torrentFilePath)
	if err != nil {
		fmt.Println(err)
	}
	var torrent TorrentFile
	reader := bytes.NewReader(torrentData)
	err = bencode.Unmarshal(reader, &torrent)
	if err != nil {
		fmt.Println(err)
	}
	bencodedInfo, _ := encodeToBencode(torrent.Info)
	h := sha1.New()
	io.WriteString(h, bencodedInfo)
	infoHash := h.Sum(nil)

	conn, err := net.Dial("tcp", peerIp)

	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	//fmt.Println(peerIp)

	// Build your handshake
	pstrlen := byte(19) // The length of the string "BitTorrent protocol"
	pstr := []byte("BitTorrent protocol")
	reserved := make([]byte, 8) // Eight zeros
	peer_id := "00112233445566778899"
	handshake := append([]byte{pstrlen}, pstr...)
	handshake = append(handshake, reserved...)
	handshake = append(handshake, infoHash...)
	handshake = append(handshake, peer_id...)

	buffer := make([]byte, 68)

	// Send Handshake
	_, err = conn.Write(handshake)
	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Peer ID:", hex.EncodeToString(buffer[48:]))

	if downloadP == 1 {
		getPiece(conn, myPiece, 0, torrent.Info.PieceLength, pieceIndex, torrent.Info.Length)
		h := sha1.New()
		io.WriteString(h, string(*myPiece))
		infoHash := h.Sum(nil)
		// fmt.Println(pieceIndex*20, )

		if bytes.Equal(infoHash, []byte(torrent.Info.Pieces[(pieceIndex*20):(pieceIndex*20)+20])) {
			fmt.Println("Hashes are equal")

		} else {
			fmt.Println("Hash not equal")
		}
	}

}

func downloadPiece() {

	peerIps := getTracker(string(os.Args[4]))
	fmt.Println(peerIps.([]string)[0])
	myPiece := make([]byte, 0, 1)
	pieceIndex, _ := strconv.Atoi(os.Args[5])
	// fmt.Println("Piece len", len(myPiece))
	getHandshake(string(peerIps.([]string)[1]), 1, &myPiece, string(os.Args[4]), int(pieceIndex))

	h := sha1.New()
	io.WriteString(h, string(myPiece))
	infoHash := h.Sum(nil)
	fmt.Println(hex.EncodeToString(infoHash))

	fmt.Println("Info Hash:", fmt.Sprintf("%x", infoHash))

	// Create the directory and any necessary parents

	var output string
	flag.StringVar(&output, "o", "/tmp/test-piece-0", "Torrent file destination")
	flag.Parse()

	dirPath := path.Dir(string(flag.Arg(2)))

	// fmt.Println(dirPath)
	err := os.MkdirAll(dirPath, os.ModePerm)

	filepath := string(flag.Arg(2))

	f, _ := os.Create(filepath)

	defer f.Close()

	n2, _ := f.Write(myPiece)

	fmt.Printf("wrote %d bytes\n", n2)
	_, err = ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	} else {
		fmt.Printf("Piece %d downloaded to %s\n", pieceIndex, string(flag.Arg(2)))
	}

	// fmt.Println(tor) gvhvnvbmbnmbhbb
}

func downloadPiece1(fileName string, pieceIndex int) []byte {

	peerIps := getTracker(fileName)
	fmt.Println(peerIps.([]string)[0])
	myPiece := make([]byte, 0, 1)

	// fmt.Println("Piece len", len(myPiece))
	getHandshake(string(peerIps.([]string)[1]), 1, &myPiece, fileName, int(pieceIndex))

	return myPiece

	// fmt.Println(tor) gvhvnvbmbnmbhbb
}

func finalDowwnload() {
	myPiece := make([]byte, 0, 1)

	torrentFilePath := os.Args[4]

	var torrent TorrentFile
	torrentData, _ := ioutil.ReadFile(torrentFilePath)
	reader := bytes.NewReader(torrentData)
	_ = bencode.Unmarshal(reader, &torrent)

	for i := 0; i < int(len(torrent.Info.Pieces)/20); i++ {
		fmt.Print(i)
		newPiece := downloadPiece1(string(os.Args[4]), i)
		myPiece = append(myPiece, newPiece...)
	}

	fmt.Println(len(myPiece))

	fPath := string(os.Args[3])

	dirPath := path.Dir(fPath)

	// fmt.Println(dirPath)
	_ = os.MkdirAll(dirPath, os.ModePerm)

	filepath := fPath

	f, _ := os.Create(filepath)

	defer f.Close()

	_, _ = f.Write(myPiece)

	fmt.Printf("Downloaded test.torrent to /tmp/test.txt.\n")

}
