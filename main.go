package main

import (
	"fmt"
	"os"

	"lachut.net/gogs/dslachut/go-irleak/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
