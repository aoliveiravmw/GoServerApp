package main

import (
	"simpleApp/api"
)

func main() {
	a := api.App{}
	a.Port = ":8080"
	a.Initialize()
	a.Run()
}
