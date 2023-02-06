package main

import (
	"log"

	"auth-test/infra"
)

func main() {
	if err := infra.Run(); err != nil {
		log.Fatal(err)
	}
}
