package main

import (
	"fmt"
	"github.com/gobuffalo/packr"
)
func main(){
	box := packr.NewBox("./config")
	data,err := box.FindString("redis.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Contents of file:", data)
}