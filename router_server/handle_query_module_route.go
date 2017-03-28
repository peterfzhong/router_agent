package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"router/proto"

  	"database/sql"
)

import _ "github.com/go-sql-driver/mysql"

func add_module_ip_interest(module_id string, ip string)(code int32){

	code = 0
	db, err := sql.Open("mysql", "peterzhong:12345678@/db_router")

	fmt.Println(module_id, ", ", ip)
	if err != nil {
		log.Fatal("fatal error happend")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	sql := fmt.Sprintf("replace into db_router.tb_module_machine_interest(module_id, ip) values('%s', '%s')", module_id, ip)

	_, err = db.Exec(sql)
	//checkErr(err)
	if err !=nil {
		log.Fatal("error when exec sql: ", sql, err)
		return
	}

	return
}

func query_module_info_db(module_id string, module_info *router.ModuleInfo)(code int32){

	code = 0
	db, err := sql.Open("mysql", "peterzhong:12345678@/db_router")

	fmt.Println(module_info)
	if err != nil {
		log.Fatal("fatal error happend")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT module_id, module_name FROM db_router.tb_module where module_id = '%s'", module_id)

	rows, err := db.Query(sql)
	//checkErr(err)

	var moduleid string
	var module_name string
	for rows.Next() {
		err = rows.Scan(&moduleid, &module_name)

		fmt.Println(moduleid)
		fmt.Println(module_name)

		break
	}
	module_info.ModuleId = &moduleid
	module_info.ModuleName = proto.String(module_name)

	log.Println(sql)
	return
}

func query_module_ip_db(module_id string, module_info *router.ModuleInfo)(code int32){
	code = 0
	db, err := sql.Open("mysql", "peterzhong:12345678@/db_router")

	if err != nil {
		log.Fatal("fatal error happend")
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT ip, weight FROM db_router.tb_module_machine where module_id = '%s'", module_id)

	log.Println(sql)
	rows, err := db.Query(sql)
	//checkErr(err)
	//
	for rows.Next() {
		var ip string
		var weight int32
		err = rows.Scan(&ip, &weight)

		router_machine := &router.ModuleMachine{
			ModuleId: &module_id,
			Ip: &ip,
			Weight: &weight,
		}
		module_info.MachineList = append(module_info.MachineList, router_machine)

		fmt.Println(ip)
		fmt.Println(weight)

	}

	log.Println("module_info.MachineList length is ", len(module_info.MachineList))
	return
}

func HandleGetRouteInfoRequest(ip string, request_str string)(code int32, response string){
	code = 0
	get_module_request := &router.GetModuleRouteListRequest{}
	err := proto.Unmarshal([]byte(request_str), get_module_request)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return
	}
	module_id := get_module_request.GetModuleId()

	fmt.Println("get request module id: ", module_id)

	get_module_response := &router.GetModuleRouteListResponse{
		Code:  proto.Int32(0),
	}
	get_module_response.Module = &router.ModuleInfo{}

	get_module_response.Module.MachineList = []*router.ModuleMachine{}

	code = query_module_info_db(module_id, get_module_response.Module)
	if code != 0 {
		log.Fatal("query_module_info_db error: ", code, module_id)
		return
	}

	code = query_module_ip_db(module_id, get_module_response.GetModule())
	if code != 0 {
		log.Fatal("query_module_info_db error: ", code, module_id)
		return
	}

	data, err := proto.Marshal(get_module_response)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return
	}
	response = string(data)

	add_module_ip_interest(module_id, ip)

	return
}

