package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type packet struct {
	FileName string
	Mode     int
	Data     []byte
}

func StartServer() {
	router := gin.Default()
	router.POST("/file", postFile)

	router.Run("localhost:8080")
}

func postFile(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}

	//unpack incoming data
	var incoming packet
	json.Unmarshal(data, &incoming)
	// fmt.Println(incoming)

	file := openFile(incoming)

	//write data to file
	if _, err := file.Write([]byte(incoming.Data)); err != nil {
		panic(err)
	}

	file.Close()
}

// Evaluates incoming packet file operation and returns the neccisary file pointer
func openFile(incoming packet) *os.File {
	var (
		file *os.File
		fErr error
	)

	//evaluate file operation mode (create/append)
	if incoming.Mode == 0 {
		file, fErr = os.Create(incoming.FileName)

	} else {
		file, fErr = os.OpenFile(incoming.FileName, os.O_APPEND, 0644)
	}
	if fErr != nil {
		fmt.Println(fErr)
	}

	return file
}

/*
Request body formatting

filePath : the path of the file relative to the parent track location (./ would be in the track folder)
fileName : string - The name of the file
dataSlice : number	- Contains sequence numbers in the event the file is too big for one send
package : bytes - Contains the data from the file


type packet struct {
	FileName string	//name of file
	Mode     int	//mode (new/append)
	Data     []byte	//actual file contents
}
*/
