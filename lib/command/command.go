package command

import (
	"fmt"
	"goup/lib/app"
	"os"

	_ "github.com/spf13/cobra"
)

var command string = "help"

var funcMap map[string]func() = map[string]func(){
	"help":    app.App.Print,
	"update":  app.CheckUpdate,
	"version": app.App.GetVersion,
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

type description struct {
	AppName    string
	AppDesc    string
	AppVersion string
	AppOptions []map[string]string
}

func noSupport() {
	fmt.Println("no support")
}
