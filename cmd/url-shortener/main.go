package main

import (
	"URL_SHORTENER/internal/config"
	"fmt"
)

func main(){
	config := config.MustLoad()

	fmt.Print(config)

	//TODO: init logger

	//TODO: init storage

	//TODO: init router

	//TODO: run server
}