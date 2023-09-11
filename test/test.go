package main

import (
	"fmt"
	"wechat/utils/secure"
)

func main() {

	fmt.Println(secure.RandomStr(10))
	fmt.Println(secure.RandomStr(10))
	fmt.Println(secure.RandomStr(10))
}
