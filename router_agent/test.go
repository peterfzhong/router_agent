package main

import (
	"fmt"
)

import "database/sql"

import _ "github.com/go-sql-driver/mysql"
import "github.com/go-redis/redis"
import "github.com/golang/protobuf/proto"

import "router/proto"

//import "example"

import "log"

func protocbuf_test() {
// 创建一个消息 Test
test := &router.Test{
// 使用辅助函数设置域的值
Label: proto.String("hello"),
Type:  proto.Int32(17),
Optionalgroup: &router.Test_OptionalGroup{
RequiredField: proto.String("good bye"),
},
}

// 进行编码
data, err := proto.Marshal(test)
if err != nil {
log.Fatal("marshaling error: ", err)
}
fmt.Println(data)
// 进行解码
newTest := &router.Test{}
err = proto.Unmarshal(data, newTest)
if err != nil {
log.Fatal("unmarshaling error: ", err)
}

// 测试结果
if test.GetLabel() != newTest.GetLabel() {
log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
}

fmt.Println(newTest.GetLabel())
}

func redis_test(i int, ch chan int) {

client := redis.NewClient(&redis.Options{
Addr:     "localhost:6379",
Password: "", // no password set
DB:       0,  // use default DB
})

pong, err := client.Ping().Result()
fmt.Println(pong, err)
// Output: PONG <nil>

key := fmt.Sprintf("key%d", i)
val := fmt.Sprintf("val%d", i)
err = client.Set(key, val, 0).Err()
if err != nil {
panic(err)
}

val, err = client.Get(key).Result()
if err != nil {
panic(err)
}
fmt.Println(key, val)

val2, err := client.Get("key2").Result()
if err == redis.Nil {
fmt.Println("key2 does not exists")
} else if err != nil {
panic(err)
} else {
fmt.Println("key2", val2)
}
// Output: key value
// key2 does not exists

ch <- i
return

}

func mysql_test() {
db, err := sql.Open("mysql", "peterzhong:12345678@/db_router")

if err != nil {
panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
}
defer db.Close()

// Prepare statement for inserting data
//	stmtIns, err := db.Prepare("INSERT INTO tb_router VALUES(2, '172.16.0.92')") // ? = placeholder
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

// Prepare statement for reading data
//	stmtOut, err := db.Prepare("SELECT * FROM tb_router where module_id = ?")
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//	defer stmtOut.Close()

//	// Insert square numbers for 0-24 in the database
//	for i := 0; i < 25; i++ {
//		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
//		if err != nil {
//			panic(err.Error()) // proper error handling instead of panic in your app
//		}
//	}

//	_, err = stmtIns.Exec() // Insert tuples (i, i^2)
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}

//	var squareNum int // we "scan" the result in here

//	// Query the square-number of 13
//	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//	fmt.Printf("The square number of 13 is: %d", squareNum)

//	// Query another number.. 1 maybe?
//	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//	fmt.Printf("The square number of 1 is: %d", squareNum)

rows, err := db.Query("SELECT * FROM tb_router")
if err != nil {
panic(err.Error()) // proper error handling instead of panic in your app
}

// Get column names
columns, err := rows.Columns()
if err != nil {
panic(err.Error()) // proper error handling instead of panic in your app
}

fmt.Println("rows number ", len(columns))
}

func test_go() {
fmt.Print("hello world\n")

//	mysql_test()

protocbuf_test()
//	chs := make([]chan int, 5000)
//	for i := 0; i < 5000; i++ {
//		go redis_test(i, chs[i])
//	}

//	for _, ch := range chs {
//		data := <-ch
//		fmt.Println(data)
//	}
}