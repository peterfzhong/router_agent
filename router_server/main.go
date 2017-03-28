package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"strings"

	"router/common"
	"github.com/golang/protobuf/proto"
	"router/proto"
	"log"
)


const (
	msg_length = 1024
)


var g_router_info = make(map[string][]*common.RouterMachineInfo)

var g_lock sync.RWMutex

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getRemoteIp(remote_addr string)(remoteIp string){
	index := strings.Index(remote_addr, ":")
	remoteIp = remote_addr[0:index]
	return
}

func handle_process(remote_addr string, request *router.Request, response *router.Response) {
	//response = &router.Response{
	//	Code:     proto.Int32(0),
	//	Response: proto.String("hello world"),
	//}
	cmd := request.GetCmd()
	code := int32(0)
	res := ""
	switch cmd {
	case 101:
		code, res = HandleGetRouteInfoRequest(remote_addr, request.GetBody())
	case 102:
		code, res = HandleNotifyUpdateRouteInfoRequest(request.GetBody())
	}
	response.Code = new(int32)
	*response.Code = code
	response.Response = &res

	return
}

func handle_conn(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, msg_length)

	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("read message from lotus failed")
		return
	}
	fmt.Println("recv byte ", n, data)

	newRequest := &router.Request{}
	err = proto.Unmarshal(data[0:n], newRequest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}

	newResponse := &router.Response{}

	remote_addr := getRemoteIp(conn.RemoteAddr().String())
	handle_process(remote_addr, newRequest, newResponse)
	//response := &router.Response{
	//	Code:     proto.Int32(0),
	//	Response: proto.String("hello"),
	//}

	data, err = proto.Marshal(newResponse)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}

	conn.Write([]byte(data))

	return
}

func init_server() {
	l, err := net.Listen("tcp", ":2188")
	if err != nil {
		fmt.Printf("Failure to listen: %s\n", err.Error())
		return
	}

	defer l.Close()

	for {
		if c, err := l.Accept(); err == nil {
			go handle_conn(c) //new thread
		}
	}
}

func main() {
	init_server()
}
