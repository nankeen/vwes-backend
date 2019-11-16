package main

import "github.com/nankeen/vwes-backend/router"

func main() {
	r := router.SetupRouter()
	r.Run(":1337")
}
