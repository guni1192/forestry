package main

import (
	"github.com/guni1192/forestry/pkg/api"
)

func main() {
	e := api.NewRouter(true)
	e.Logger.Fatal(e.Start(":1192"))
}
