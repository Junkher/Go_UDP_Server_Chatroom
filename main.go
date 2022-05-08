package main

import (
	"UDP/backend"
	"fmt"
	"os"
)

// const(
// 	address string = ":8000"
// )
// 从命令行参数获取端口号
func main() {

	if len(os.Args) != 2 {
		fmt.Println("请输入正确的端口号")
	}
	s := backend.NewServer()
	s.Startup(os.Args[1])

}