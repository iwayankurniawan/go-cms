package main

import (
	"fmt"
)

type Content struct{
	id int
	text string
}

func main() {
	c := Content{id:1, text:"hello"}
	fmt.Println(c.text)
}

