package main

import (
	"stash.appscode.dev/prober-demo/pkg/cmd"
	"log"
)

func main() {
	err:= cmd.NewRootCmd().Execute()
	if err!=nil{
		log.Fatal(err)
	}
}
