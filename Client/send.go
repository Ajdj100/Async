package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type packet struct {
	FileName string
	Mode     int
	Data     []byte
}

func main() {
	var fileData packet

	// fileData.FileName = os.Args[1]
	fPath := os.Args[1]
	nameStart := strings.LastIndex(fPath, "\\")
	fileData.FileName = fPath[nameStart+1:]
	fileData.Mode = 0

	file, err := os.Open(fPath)

	if err != nil {
		log.Fatalln(err)
	}

	//loop send until entire file has been sent
	for readBytes := -1; readBytes != 0; {
		// buf := make([]byte, 5120)
		buf := make([]byte, 1024000)

		//read data into file
		readBytes, _ = file.Read(buf)

		//check for EOF
		if readBytes == 0 {
			break
		}
		//truncate buffer
		if readBytes < len(buf) {
			buf = buf[:readBytes]
		}
		// buf = bytes.Trim(buf, "\x00")
		fileData.Data = buf

		responseBody := jsonify(fileData)

		pushFile(responseBody)

		fileData.Mode = 1 //switch to append mode
	}
	log.Println("Finished sending file")

}

// convert packet into json for sending
func jsonify(data packet) *bytes.Buffer {
	postBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	return bytes.NewBuffer(postBody)
}

// push file to server
func pushFile(body *bytes.Buffer) {
	_, err := http.Post("http://localhost:8080/file", "application/json", body)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	// res, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// sb := string(res)
	//log.Printf(sb)
}
