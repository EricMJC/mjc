package main

import (
	"fmt"
	//"golang.org/x/tools/go/ssa"
)

func main(){
	fmt.Println("hello world my golang")
	str :=  "hello world"

	for i:=0 ;i<len(str); i++  {
		ch:=str[i]
		fmt.Println(ch)
	}

	for i,v:=range str{
	fmt.Println(i,v)
	}

}