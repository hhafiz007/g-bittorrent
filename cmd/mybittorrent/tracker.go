package main


import (
	// Uncomment this line to pass the first 

	//"encoding/json"
	"fmt"
	"os"
	//"strconv"
	//"unicode"
	"bytes"
	"io/ioutil"
	"crypto/sha1"
	//"encoding/hex"
  	"io"
  	"net/http"
	"net/url"
	"net"
	bencode "github.com/jackpal/bencode-go" // Available if you need it!
)


func getTracker(){

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
	bencodedInfo,_:=encodeToBencode(torrent.Info)
	h := sha1.New()
	io.WriteString(h, bencodedInfo)
	infoHash := h.Sum(nil)
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

	url := fmt.Sprintf("%s?%s",torrent.Announce,query.Encode())

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Oops! Something went wrong:", err)
		return
	}
	
	body, _ := ioutil.ReadAll(res.Body)
	decoded,_,err := decodeBencode(string(body),0)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	
	
	peers := decoded.(map[string]interface{})["peers"]
	strPeers :=  []byte(peers.(string))
	
	for i := 0; i < len(strPeers); i += 6 {
		
		ip := net.IP(( strPeers[i : i+4]))
		port := int(( strPeers[ i+4]))<<8 + int(( strPeers[i+5]))
		fmt.Printf("%s:%d\n",ip.String(),port)
	}






}