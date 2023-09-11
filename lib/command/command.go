package command

import (
	"fmt"
	"goup/lib/app"
	"os"

	_ "github.com/spf13/cobra"
)

var command string = "help"

var fnMap = map[string]func(){
	"help":    app.App.Print,
	"update":  app.CheckUpdate,
	"version": app.App.GetVersion,
}

func init() {
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if fn, ok := fnMap[command]; ok {
		fn()
	} else {
		noSupport()
	}
}

func noSupport() {
	fmt.Println("no support this command", os.Args)
}
