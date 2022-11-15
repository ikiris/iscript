package main

import (
	"fmt"
	"iscript/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is some bullshit!\n", user)
	fmt.Printf("Feel free to type in something or whatever\n")
	repl.Start(os.Stdin, os.Stdout)
}
