package main

import (
	"log"

	_ "github.com/homearchbishop/downgram/internal/cache"
	"github.com/homearchbishop/downgram/internal/gui"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	gui.Start()
}
