package command

import (
	"fmt"
	"goup/lib/app"
	"os"
)

var command string = "help"

var funcMap map[string]func() = map[string]func(){
	"help":    help,
	"update":  app.CheckUpdate,
	"version": app.GetVersion,
}

func init() {
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if fn, ok := funcMap[command]; ok {
		fn()
	} else {
		noSupport()
	}
}

func help() {
	fmt.Println("goup", app.Version)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("-", "goup.exe [COMMAND]")
	fmt.Println("")
	fmt.Println("COMMAND:")
	fmt.Println("-", "update", "Update golang stable version")
	fmt.Println("-", "version", "Print version information")
	fmt.Println("-", "help", "Print this information")
}

func noSupport() {
	fmt.Println("no support")
}
