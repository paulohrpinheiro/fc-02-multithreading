package main

import (
	"fmt"
	"os"

	service "github.com/paulohrpinheiro/fc-02-multithreading/service"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please, give-me one CEP value to search")
		return
	}
	r, err := service.GetAddress(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Provider:", r.Provider)
	fmt.Println("Response:", r.Response)
}
