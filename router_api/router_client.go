package main

import (
	"fmt"
	"net"
	"log"

	"github.com/golang/protobuf/proto"

	"router/proto"
)

func GetRoute(module_id string)(code int32, ip string){
	get_route_request := &router.GetRouteRequest{
		ModuleId:proto.String(module_id),
	}
	data, err := proto.Marshal(get_route_request)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}

	request := &router.Request{
		Cmd:  proto.Int32(1),
		Body: proto.String(string(data)),
	}

	// 进行编码
	data, err = proto.Marshal(request)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return
	}

	client, err := net.Dial("tcp", "127.0.0.1:1988")
	if err != nil {
		fmt.Printf("Failure to connet:%s\n", err.Error())
		return
	}
	client.Write(data)

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		return
	}

	response := &router.Response{}
	err = proto.Unmarshal(buf[0:n], response)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}

	code = response.GetCode()
	data1 := response.GetResponse()
	log.Println("get Ip: ", data1)

	response1 := &router.GetRouteResponse{}
	err = proto.Unmarshal([]byte(data1), response1)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}
	ip = response1.GetIp()

	log.Println("get Ip: ", ip)
	return

}


func UpdateRoute(module_id string, ip string, status int32)(code int32, error string){
	route_request := &router.UpdateRouteRequest{
		ModuleId:proto.String(module_id),
		Ip:proto.String(ip),
		Status:proto.Int32(status),
		Timeout:proto.Int32(1),
	}
	data, err := proto.Marshal(route_request)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}

	request := &router.Request{
		Cmd:  proto.Int32(2),
		Body: proto.String(string(data)),
	}

	// 进行编码
	data, err = proto.Marshal(request)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return
	}

	client, err := net.Dial("tcp", "127.0.0.1:1988")
	if err != nil {
		fmt.Printf("Failure to connet:%s\n", err.Error())
		return
	}
	client.Write(data)

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		return
	}

	response := &router.Response{}
	err = proto.Unmarshal(buf[0:n], response)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}

	code = response.GetCode()
	data1 := response.GetResponse()

	response1 := &router.UpdateRouteResponse{}
	err = proto.Unmarshal([]byte(data1), response1)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}
	error = response1.GetError()

	return

}


func NotifyReloadRoute(module_id string)(code int32, error string){
	route_request := &router.ReLoadRouteNotify{
		ModuleId:proto.String(module_id),
	}
	data, err := proto.Marshal(route_request)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}

	request := &router.Request{
		Cmd:  proto.Int32(102),
		Body: proto.String(string(data)),
	}

	// 进行编码
	data, err = proto.Marshal(request)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return
	}

	client, err := net.Dial("tcp", "127.0.0.1:2188")
	if err != nil {
		fmt.Printf("Failure to connet:%s\n", err.Error())
		return
	}
	client.Write(data)

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		return
	}

	response := &router.Response{}
	err = proto.Unmarshal(buf[0:n], response)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}

	code = response.GetCode()
	error = ""
	return

}

func main() {
	code, ip := GetRoute("order")
	fmt.Println("get route of order", code, ip)

	//code, error := UpdateRoute("order", "172.16.0.92", 1)
	//fmt.Println("update route of order ", code, error)
	//
	//code, error := NotifyReloadRoute("order")
	//fmt.Println("notify route of order ", code, error)

}
