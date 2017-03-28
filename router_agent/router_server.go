package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"router/proto"

	"router/common"
)

func get_router_machine_from_server(module_id string)(code int32, router_machine_list []*common.RouterMachineInfo){
	code = 0

	route_request := &router.GetModuleRouteListRequest{
		ModuleId:proto.String(module_id),
	}
	data, err := proto.Marshal(route_request)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}

	request := &router.Request{
		Cmd:  proto.Int32(101),
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
	if code != 0 {
		fmt.Println("get resposne error : ", code)
		return
	}

	module_info_response_content := response.GetResponse()

	module_info_response := &router.GetModuleRouteListResponse{}
	err = proto.Unmarshal([]byte(module_info_response_content), module_info_response)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
		return
	}

	code = module_info_response.GetCode()
	if code != 0 {
		fmt.Println("get resposne error : ", code)
		return
	}

	module_info := module_info_response.GetModule()


	//var router_machine_list [10]*common.RouterMachineInfo

	for  _, module_machine := range module_info.GetMachineList(){
		router_machine := &common.RouterMachineInfo{
			ModuleId: module_machine.GetModuleId(),
			Ip: module_machine.GetIp(),
			Weight:module_machine.GetWeight(),
		}

		router_machine_list = append(router_machine_list, router_machine)
	}

	g_lock.Lock()
	g_router_info[module_info.GetModuleId()] = router_machine_list
	defer g_lock.Unlock()

	return
}

func get_router_machine(module_id string)(code int32, router_machine_info *common.RouterMachineInfo){
	code = 0

	router_list, ret := g_router_info[module_id]
	if ret == false {
		code, router_list = get_router_machine_from_server(module_id)
		if code != 0{
			return
		}
	}

	code , router_machine_info = select_router(router_list)
	return

}


func HandleGetRouteRequest(request_str string)(code int32, response string){
	code = 0
	get_route_request := &router.GetRouteRequest{}
	err := proto.Unmarshal([]byte(request_str), get_route_request)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}
	module_id := get_route_request.GetModuleId()

	fmt.Println("get request module id: ", module_id)

	get_route_response := &router.GetRouteResponse{
		Code:  proto.Int32(0),
		Ip: proto.String("172.16.0.92"),
	}

	code, router_machine := get_router_machine(get_route_request.GetModuleId())
	if code != 0 {
		log.Fatal("get_router_machine error: ", code)
		return
	}
	get_route_response.Ip = &router_machine.Ip

	data, err := proto.Marshal(get_route_response)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}
	response = string(data)

	return
}


func HandleUpdateRouteRequest(request_str string)(code int32, response string){
	code = 0
	route_request := &router.UpdateRouteRequest{}
	err := proto.Unmarshal([]byte(request_str), route_request)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}
	module_id := route_request.GetModuleId()
	ip := route_request.GetIp()
	status := route_request.GetStatus()
	time_out := route_request.GetTimeout()

	log.Println("get request module id: ", module_id, ip, status, time_out)

	route_response := &router.UpdateRouteResponse{
		Code:  proto.Int32(0),
		Error: proto.String(""),
	}

	data, err := proto.Marshal(route_response)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}
	response = string(data)

	return
}

func HandleNotifyUpdateRouteInfoRequest(request_str string)(code int32, response string){
	code = 0
	response = ""

	request := &router.ReLoadRouteNotify{}
	err := proto.Unmarshal([]byte(request_str), request)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}
	module_id := request.GetModuleId()

	code,_ = get_router_machine_from_server(module_id)
	if code != 0{
		log.Fatal("get request module id error, ", module_id, ", code: ", code)
		return
	}


	return
}

