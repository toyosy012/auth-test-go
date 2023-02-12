package main

import (
	"log"

	"auth-test/infra"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and Token.
// @BasePath /v1
func main() {
	if err := infra.Run(); err != nil {
		log.Fatal(err)
	}
}
