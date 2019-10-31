package main

import (
	"github.com/tamalsaha/prober-demo/cmd"
	"log"
)

func main() {
	err:=cmd.NewRootCmd().Execute()
	if err!=nil{
		log.Fatal(err)
	}
}
