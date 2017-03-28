package main

import (
	"log"

	"router/common"

	"math/rand"
	"time"
)


func select_router(router_machine_list []*common.RouterMachineInfo)(code int32, router_machine_info *common.RouterMachineInfo){
	code = 0
	if len(router_machine_list) == 0{
		log.Fatal("router_machine_list length is 0")
		code = -100
		return
	}else if len(router_machine_list) == 1{
		router_machine_info = router_machine_list[0]
		return
	}

	code, router_machine_info = select_by_weight(router_machine_list)
	return
}

func select_by_weight(router_machine_list []*common.RouterMachineInfo)(code int32, router_machine_info *common.RouterMachineInfo)  {
	code = 0

	weight_total := int32(0)
	for _, router_machine := range router_machine_list{
		weight_total = weight_total + router_machine.Weight
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	found_weight := int32(rand.Int31()%weight_total)
	log.Println("found_weight: ", found_weight)
	log.Println(router_machine_list)

	weight_total = 0
	found := false
	for index, router_machine := range router_machine_list{
		weight_total = weight_total + router_machine.Weight
		log.Println("found_weight: ", found_weight, "weight_total: ", weight_total)
		if weight_total >= found_weight{
			router_machine_info = router_machine_list[index]
			found = true
			return
		}
	}

	if found == false{
		router_machine_info = router_machine_list[len(router_machine_list) - 1]
		return
	}

	return

}
