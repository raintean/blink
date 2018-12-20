package blink

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "blink", log.LstdFlags)
