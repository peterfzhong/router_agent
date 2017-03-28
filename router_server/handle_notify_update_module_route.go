package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"router/proto"

	"database/sql"
)

import _ "github.com/go-sql-driver/mysql"

func query_module_interest_ip_list(module_id string, ip_list []string)(code int32){
	code = 0
	db, err := sql.Open("mysql", "peterzhong:12345678@/db_router")

	if err != nil {
		log.Fatal("fatal error happend")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT ip FROM db_router.tb_module_machine_interest where module_id = '%s'", module_id)

	log.Println(sql)
	rows, err := db.Query(sql)
	//checkErr(err)
	//
	for rows.Next() {
		var ip string
		err = rows.Scan(&ip)

		ip_list = append(ip_list, ip)

		fmt.Println(ip)

	}

	log.Println("ip_list length is ", len(ip_list))
	return
}

func notify_agent_reload(ip string, request_str string)(code int32){
	code = 0

	remote_addr := fmt.Sprintf("%s:1988", ip)
	client, err := net.Dial("tcp", remote_addr)
	if err != nil {
		fmt.Printf("Failure to connet:%s\n", err.Error())
		return
	}
	request := &router.Request{
		Cmd:  proto.Int32(102),
		Body: proto.String(string(request_str)),
	}

	// 进行编码
	data, err := proto.Marshal(request)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return
	}

	client.Write([]byte(data))

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

	return
}

func HandleNotifyUpdateRouteInfoRequest(request_str string)(code int32, response string){
	code = 0

	request := &router.ReLoadRouteNotify{}
	err := proto.Unmarshal([]byte(request_str), request)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}
	module_id := request.GetModuleId()
	ip_list := make([]string, 10)

	code = query_module_interest_ip_list(module_id, ip_list)
	if code != 0 {
		log.Fatal("query_module_info_db error: ", code, module_id)
		return
	}

	for _, ip := range ip_list{
		notify_agent_reload(ip, request_str)
	}

	fmt.Println("get request module id: ", module_id)

	response = ""

	return
}

