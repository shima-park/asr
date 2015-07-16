package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
)

func IntToBytes(i int) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	return buf
}

func Rand(begin int, end int) int {
	if end > 0 && end > begin {
		return rand.Intn(end-begin) + begin
	}
	return 0
}

func test() {

	conn, err := net.Dial("tcp", "127.0.0.1:16800")
	if err != nil {
		fmt.Println("connect failed!")
		return
	}

	go func() {
		read_buf := make([]byte, 1024)
		n, err := conn.Read(read_buf)
		for err == nil && n > 0 {
			fmt.Println("result", string(read_buf), n)
		}
	}()

	//fileName := strconv.Itoa(Rand(1, 9))
	fileName := `1.amr`

	fl, err := os.Open(fileName)
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	defer fl.Close()

	buffer := make([]byte, 12)
	buffer[0] = 'Y'
	buffer[1] = 'T'

	buf := make([]byte, 1024)
	for {
		fn, err := fl.Read(buf)
		if 0 == fn || err == io.EOF {
			break
		}

		length := IntToBytes(fn + 12)
		fmt.Println("-----------length---------", fn+12)
		buffer[2] = length[0]
		buffer[3] = length[1]
		buffer[4] = 0
		buffer[5] = 0
		buffer[6] = 0
		buffer[7] = 0
		buffer[8] = length[0]
		buffer[9] = length[1]
		buffer[10] = length[2]
		buffer[11] = length[3]

		_, err = conn.Write(buffer)

		if err != nil {
			fmt.Printf("Write head error!" + err.Error() + "\r\n")
			return
		}
		var data []byte
		data = append(data, buf[:fn]...)

		data = append(data, 0)
		data = append(data, 0)
		data = append(data, 'C')
		data = append(data, 'N')

		_, err = conn.Write(data)

		if err != nil {
			fmt.Printf("Write data error!" + err.Error() + "\r\n")
			return
		}

	}

	conn.Close()
}

func Run(num int) {
	if num == 0 {
		num = 1
	}

	for i := 0; i < num; i++ {
		go test()
	}

}

func main() {
	var num int = 0
	flag.IntVar(&num, "num", 20, "server thread")
	flag.Parse()

	Run(num)

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}

	}
}
