package bintrayresource

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
)

func LogErrorAndExit(err error) {
	sayf(colorstring.Color("[red]error %s: %s\n"), "runningCommand", err)
	os.Exit(1)
}

func sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}
