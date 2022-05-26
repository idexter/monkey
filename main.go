package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/idexter/monkey/repl"
)

func main() {
	entrypoint := flag.String("in", "", "Runs script from file.\nUsage: monkeyc -in ./example.monkey")
	runREPL := flag.Bool("repl", false, "Runs REPL")
	flag.Parse()

	if *runREPL {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", usr.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.StartREPL(os.Stdin, os.Stdout)
		return
	}

	if *entrypoint != "" {
		f, err := os.OpenFile(*entrypoint, os.O_RDONLY, os.ModePerm)
		defer f.Close()
		if err != nil {
			fmt.Printf("Unable to read script: %v\n", err)
			return
		}

		repl.RunScript(f, os.Stdout)
		return
	}
}
