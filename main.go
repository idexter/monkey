package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/idexter/monkey/rppl"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", usr.Username)
	fmt.Printf("Feel free to type in commands\n")
	rppl.Start(os.Stdin, os.Stdout)
}
