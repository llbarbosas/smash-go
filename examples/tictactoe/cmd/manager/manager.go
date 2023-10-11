package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/llbarbosas/smash-go"
)

func printHelp() {
	fmt.Println("c: Load console_renderer module")
	fmt.Println("i: Load input_manager module")
	fmt.Println("s: Load simulator module")
	fmt.Println("r: Register modules")
}

func main() {
	addr := flag.String("addr", "127.0.0.1:6061", "link node to a remote node")
	flag.Parse()

	client, err := smash.NewManagmentAPIClient(fmt.Sprintf("http://%s", *addr))

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

		case "i":
			res, err := client.LoadModule(smash.ModuleLoadRequest{
				Path: "./bin/modules/input_manager.so",
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
