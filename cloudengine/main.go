package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"

	"github.com/shima-park/asr/cloudengine/yt"
	"github.com/shima-park/asr/cloudengine/yt/log"
)

func ServerGetPort() int {
	var port int = 0
	flag.IntVar(&port, "port", 16800, "server port")
	flag.Parse()
	return port
}

func ServerEchoSystemInfo(port int) {
	fmt.Println("-------------------------------------")
	fmt.Println("")
	fmt.Println("Welcome to use Youngtone CloudEngine.")
	fmt.Println("  Version V1.00")
	fmt.Println("  Write By:")
	fmt.Println("    fanghua.nie")
	fmt.Println("    xinwang.liu")
	fmt.Println("")
	fmt.Println("-------------------------------------")

}

func ServerEchoListeningInfo(port int) {
	// out put Listening infomation
	fmt.Println("")
	fmt.Printf("  Listening on at port: %d \r\n", port)
	fmt.Println("")
	fmt.Println("-------------------------------------")
}

func ServerStart(port int) {

	// set up max cpu process unit
	runtime.GOMAXPROCS(runtime.NumCPU())

	// make a tcp address and port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("  err:" + err.Error())
		return
	}

	// Listening the port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("  err:" + err.Error())
		return
	}

	// Output listening info
	ServerEchoListeningInfo(port)

	// we define a loop, accept all use request
	for {
		// If there has a new user, call the client to process
		conn, err := listener.Accept()
		if err != nil {
			// if
			continue
		}
		go yt.NewClient(conn).Run()
	}
}

func serverWaitCommand() {
	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()

		command := string(data)

		if command == "stop" {
			running = false
		} else if command == "logoff" {
			log.LogOut = 0
		} else if command == "logon" {
			log.LogOut = 1
		}

		fmt.Println("command:>", command)
	}
}

func main() {
	var port int = 0
	port = ServerGetPort()
	ServerEchoSystemInfo(port)
	go ServerStart(port)
	serverWaitCommand()
}
