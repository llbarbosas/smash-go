package main

import (
	"fmt"
	"log"
	"os"

	"github.com/llbarbosas/smash-tictactoe"
)

func printHelp() {
	fmt.Println("c: Load console_renderer module")
	fmt.Println("i: Load input_manager module")
	fmt.Println("s: Load simulator module")
	fmt.Println("r: Register modules")
}

func main() {
	client, err := smash.NewManagmentAPIClient("http://127.0.0.1:3000")

	if err != nil {
		log.Fatal(err)
	}

	printHelp()

	b := make([]byte, 1)

	for {
		os.Stdin.Read(b)
		key := string(b)

		switch key {
		case "c":
			res, err := client.LoadModule(smash.ModuleLoadRequest{
				Path: "./bin/modules/console_renderer.so",
			})

			if err != nil {
				log.Println(err)
			}

			fmt.Println(res)
		case "s":
			res, err := client.LoadModule(smash.ModuleLoadRequest{
				Path: "./bin/modules/simulator.so",
			})

			if err != nil {
				log.Println(err)
			}

			fmt.Println(res)
		case "r":
			if err := client.RegisterModules(); err != nil {
				log.Println(err)
			}
		default:
			printHelp()
		}
	}

}
